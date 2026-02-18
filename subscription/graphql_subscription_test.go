package subscription

import (
	"context"
	"testing"
	"time"

	"github.com/GoLabra/labra/entgql/entity"
)

func TestPublishEntities_FiltersPerSubscriberr(t *testing.T) {
	client := NewGraphqlSubscriptionClient()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan []*entity.Entity, 1)

	client.EntitySubscribers = append(client.EntitySubscribers, EntitySubscriber{
		Ctx:  ctx,
		Chan: ch,
		AllowedEntities: map[string]struct{}{
			"Post": {},
		},
	})

	payload := []*entity.Entity{
		{Name: "Post"},
		{Name: "User"},
	}

	client.PublishEntities(payload)

	got := <-ch
	if len(got) != 1 || got[0].Name != "Post" {
		t.Fatalf("expected only Post, got %#v", got)
	}
}

func TestPublishEntities_FiltersPerSubscriber(t *testing.T) {
	client := NewGraphqlSubscriptionClient()

	ctx := context.Background()

	chA := make(chan []*entity.Entity, 1)
	chB := make(chan []*entity.Entity, 1)

	// Subscriber A can only see "Post"
	client.EntitySubscribers = append(client.EntitySubscribers, EntitySubscriber{
		Ctx:  ctx,
		Chan: chA,
		AllowedEntities: map[string]struct{}{
			"Post": {},
		},
	})

	// Subscriber B can only see "User"
	client.EntitySubscribers = append(client.EntitySubscribers, EntitySubscriber{
		Ctx:  ctx,
		Chan: chB,
		AllowedEntities: map[string]struct{}{
			"User": {},
		},
	})

	payload := []*entity.Entity{
		{Name: "Post"},
		{Name: "User"},
		{Name: "Permission"},
	}

	client.PublishEntities(payload)

	gotA := mustRecv(t, chA)
	if len(gotA) != 1 || gotA[0].Name != "Post" {
		t.Fatalf("subscriber A expected only [Post], got %#v", names(gotA))
	}

	gotB := mustRecv(t, chB)
	if len(gotB) != 1 || gotB[0].Name != "User" {
		t.Fatalf("subscriber B expected only [User], got %#v", names(gotB))
	}
}

func TestPublishEntities_AllowsAllWhenAllowedEntitiesNil(t *testing.T) {
	client := NewGraphqlSubscriptionClient()

	ctx := context.Background()
	ch := make(chan []*entity.Entity, 1)

	// AllowedEntities == nil means allow all (e.g. SuperAdmin)
	client.EntitySubscribers = append(client.EntitySubscribers, EntitySubscriber{
		Ctx:             ctx,
		Chan:            ch,
		AllowedEntities: nil,
	})

	payload := []*entity.Entity{
		{Name: "Post"},
		{Name: "User"},
	}

	client.PublishEntities(payload)

	got := mustRecv(t, ch)
	if len(got) != len(payload) {
		t.Fatalf("expected %d entities, got %d (%#v)", len(payload), len(got), names(got))
	}
}

func mustRecv(t *testing.T, ch <-chan []*entity.Entity) []*entity.Entity {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(250 * time.Millisecond):
		t.Fatalf("timed out waiting for subscription payload")
		return nil
	}
}

func names(list []*entity.Entity) []string {
	out := make([]string, 0, len(list))
	for _, e := range list {
		out = append(out, e.Name)
	}
	return out
}
