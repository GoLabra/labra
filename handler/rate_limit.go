package handler

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipVisitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*ipVisitor

	limit rate.Limit
	burst int
	ttl   time.Duration
}

func NewIPRateLimiter(rpm int) *IPRateLimiter {
	// rpm is "requests per minute"
	// convert to tokens per second for x/time/rate
	limit := rate.Limit(float64(rpm) / 60.0)
	if rpm <= 0 {
		// safety: effectively disables limiting (shouldn't happen if config defaults are set)
		limit = rate.Inf
	}

	burst := rpm
	if burst < 1 {
		burst = 1
	}

	return &IPRateLimiter{
		visitors: make(map[string]*ipVisitor),
		limit:    limit,
		burst:    burst,
		ttl:      10 * time.Minute, // cleanup idle IP buckets
	}
}

func (l *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// cleanup stale entries
	for k, v := range l.visitors {
		if now.Sub(v.lastSeen) > l.ttl {
			delete(l.visitors, k)
		}
	}

	v, ok := l.visitors[ip]
	if !ok {
		lim := rate.NewLimiter(l.limit, l.burst)
		l.visitors[ip] = &ipVisitor{limiter: lim, lastSeen: now}
		return lim
	}

	v.lastSeen = now
	return v.limiter
}

func clientIP(r *http.Request) string {
	// If behind a proxy, X-Forwarded-For often contains "client, proxy1, proxy2"
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}
	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
		return xrip
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}

	return r.RemoteAddr
}

func (l *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		lim := l.getLimiter(ip)

		now := time.Now()
		res := lim.ReserveN(now, 1)
		if !res.OK() {
			// Extremely unlikely unless limiter is misconfigured
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		delay := res.DelayFrom(now)
		if delay > 0 {
			// Don't consume a token when rejecting
			res.CancelAt(now)

			retryAfterSeconds := int(math.Ceil(delay.Seconds()))
			if retryAfterSeconds < 1 {
				retryAfterSeconds = 1
			}

			w.Header().Set("Retry-After", strconv.Itoa(retryAfterSeconds))
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
