package subscription

import (
	"context"

	"github.com/GoLabra/labra/entgql/entity"
)

func NewGraphqlSubscriptionClient() *GraphqlSubscriptionClient {
	return &GraphqlSubscriptionClient{
		AppStatusSubscribers: make([]AppStatusSubscriber, 0),
	}
}

type GraphqlSubscriptionClient struct {
	AppStatusSubscribers []AppStatusSubscriber
	EntitySubscribers    []EntitySubscriber
}

type AppStatusSubscriber struct {
	Ctx  context.Context
	Chan chan AppStatus
}

type EntitySubscriber struct {
	Ctx  context.Context
	Chan chan []*entity.Entity

	// nil => allow all (e.g., SuperAdmin)
	AllowedEntities map[string]struct{}
}

type AppStatus string

const (
	AppStatusUp         AppStatus = "UP"
	AppStatusGenerating AppStatus = "GENERATING"
	AppStatusReverting  AppStatus = "REVERTING"
	AppStatusRestarting AppStatus = "RESTARTING"
	AppStatusFatal      AppStatus = "FATAL"
)

func (s *GraphqlSubscriptionClient) PublishAppStatusMessage(appStatus AppStatus) {
	ApplicationStatus = appStatus
	for i := 0; i < len(s.AppStatusSubscribers); i++ {
		cgs := s.AppStatusSubscribers[i]
		select {
		case <-cgs.Ctx.Done():
			s.AppStatusSubscribers = append(s.AppStatusSubscribers[:i], s.AppStatusSubscribers[i+1:]...)
			i--
		case cgs.Chan <- appStatus:
		}
	}
}

func (s *GraphqlSubscriptionClient) PublishEntities(entities []*entity.Entity) {
	for i := 0; i < len(s.EntitySubscribers); i++ {
		sub := s.EntitySubscribers[i]

		select {
		case <-sub.Ctx.Done():
			s.EntitySubscribers = append(s.EntitySubscribers[:i], s.EntitySubscribers[i+1:]...)
			i--
			continue
		default:
		}

		payload := entities

		// If not nil => filter
		if sub.AllowedEntities != nil {
			filtered := make([]*entity.Entity, 0, len(entities))
			for _, e := range entities {
				if _, ok := sub.AllowedEntities[e.Name]; ok {
					filtered = append(filtered, e)
				}
			}
			payload = filtered
		}

		select {
		case <-sub.Ctx.Done():
			s.EntitySubscribers = append(s.EntitySubscribers[:i], s.EntitySubscribers[i+1:]...)
			i--
		case sub.Chan <- payload:
		}
	}
}
