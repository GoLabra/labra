package tokenrevocation

import (
	"context"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	"github.com/GoLabra/labra/entgql/ent"
	_ "github.com/mattn/go-sqlite3"
)

func TestRevokeTokenAndCleanup(t *testing.T) {
	client := newTestEntClient(t)
	ctx := context.Background()

	now := time.Now().UTC()
	expiresAt := now.Add(30 * time.Minute)

	if err := RevokeToken(ctx, client, client.DialectName(), "access-token-1", "admin@example.com", "admin", expiresAt, now); err != nil {
		t.Fatalf("revoke token failed: %v", err)
	}

	revoked, err := IsTokenRevoked(ctx, client, client.DialectName(), "access-token-1")
	if err != nil {
		t.Fatalf("is token revoked failed: %v", err)
	}
	if !revoked {
		t.Fatalf("expected token to be revoked")
	}

	deleted, err := CleanupExpiredRevocations(ctx, client, client.DialectName(), now.Add(2*time.Hour))
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected one deleted row, got %d", deleted)
	}

	revoked, err = IsTokenRevoked(ctx, client, client.DialectName(), "access-token-1")
	if err != nil {
		t.Fatalf("is token revoked failed after cleanup: %v", err)
	}
	if revoked {
		t.Fatalf("expected token revocation row to be cleaned up")
	}
}

func TestRevokeSubjectTokens(t *testing.T) {
	client := newTestEntClient(t)
	ctx := context.Background()

	revokedAt := time.Now().UTC()
	if err := RevokeSubjectTokens(ctx, client, client.DialectName(), "user@example.com", "admin", revokedAt); err != nil {
		t.Fatalf("revoke subject tokens failed: %v", err)
	}

	revoked, err := IsSubjectTokenRevoked(ctx, client, client.DialectName(), "user@example.com", "admin", revokedAt.Add(-time.Minute))
	if err != nil {
		t.Fatalf("is subject token revoked failed: %v", err)
	}
	if !revoked {
		t.Fatalf("expected older token to be revoked")
	}

	revoked, err = IsSubjectTokenRevoked(ctx, client, client.DialectName(), "user@example.com", "admin", revokedAt)
	if err != nil {
		t.Fatalf("is subject token revoked failed: %v", err)
	}
	if !revoked {
		t.Fatalf("expected token issued in revocation second to be revoked")
	}

	revoked, err = IsSubjectTokenRevoked(ctx, client, client.DialectName(), "user@example.com", "admin", revokedAt.Add(time.Minute))
	if err != nil {
		t.Fatalf("is subject token revoked failed: %v", err)
	}
	if revoked {
		t.Fatalf("expected newer token to remain valid")
	}
}

func newTestEntClient(t *testing.T) *ent.Client {
	t.Helper()

	client, err := ent.Open(dialect.SQLite, "file:token-revocation-test?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed opening sqlite client: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Close()
	})

	return client
}
