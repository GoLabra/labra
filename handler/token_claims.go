package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func claimStringValue(claims map[string]interface{}, key string) string {
	raw, ok := claims[key]
	if !ok || raw == nil {
		return ""
	}

	switch v := raw.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func claimUnixTime(claims map[string]interface{}, key string) (time.Time, error) {
	raw, ok := claims[key]
	if !ok || raw == nil {
		return time.Time{}, fmt.Errorf("missing %s claim", key)
	}

	switch v := raw.(type) {
	case time.Time:
		return v.UTC(), nil
	case *time.Time:
		if v == nil {
			return time.Time{}, fmt.Errorf("invalid %s claim", key)
		}
		return v.UTC(), nil
	case int64:
		return time.Unix(v, 0).UTC(), nil
	case int:
		return time.Unix(int64(v), 0).UTC(), nil
	case int32:
		return time.Unix(int64(v), 0).UTC(), nil
	case float64:
		return time.Unix(int64(v), 0).UTC(), nil
	case float32:
		return time.Unix(int64(v), 0).UTC(), nil
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid %s claim: %w", key, err)
		}
		return time.Unix(n, 0).UTC(), nil
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid %s claim: %w", key, err)
		}
		return time.Unix(n, 0).UTC(), nil
	default:
		return time.Time{}, fmt.Errorf("invalid %s claim type %T", key, raw)
	}
}
