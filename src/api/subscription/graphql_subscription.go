// Package subscription defines pub/sub mechanisms for GraphQL notifications.
// This file implements an in-memory subscription client and is located in
// src/api/subscription.
package subscription

import (
	"context"

	"github.com/GoLabra/labra/src/api/entgql/entity"
)

// NewGraphqlSubscriptionClient creates a new subscription client instance.
func NewGraphqlSubscriptionClient() *GraphqlSubscriptionClient {
	return &GraphqlSubscriptionClient{
		AppStatusSubscribers: make([]AppStatusSubscriber, 0),
	}
}

// GraphqlSubscriptionClient manages subscribers for application status and entity changes.
type GraphqlSubscriptionClient struct {
	AppStatusSubscribers []AppStatusSubscriber
	EntitySubscribers    []EntitySubscriber
}

// AppStatusSubscriber represents a channel subscribed to application status updates.
type AppStatusSubscriber struct {
	Ctx  context.Context
	Chan chan AppStatus
}

// EntitySubscriber represents a channel subscribed to entity updates.
type EntitySubscriber struct {
	Ctx  context.Context
	Chan chan []*entity.Entity
}

// AppStatus enumerates possible status messages sent over subscriptions.
type AppStatus string

const (
	AppStatusUp         AppStatus = "UP"
	AppStatusGenerating AppStatus = "GENERATING"
	AppStatusReverting  AppStatus = "REVERTING"
	AppStatusRestarting AppStatus = "RESTARTING"
	AppStatusFatal      AppStatus = "FATAL"
)

// PublishAppStatusMessage notifies all subscribers about an application status change.
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

// PublishEntities notifies subscribers about updated entities.
func (s *GraphqlSubscriptionClient) PublishEntities(entities []*entity.Entity) {
	for i := 0; i < len(s.EntitySubscribers); i++ {
		cgs := s.EntitySubscribers[i]
		select {
		case <-cgs.Ctx.Done():
			s.EntitySubscribers = append(s.EntitySubscribers[:i], s.EntitySubscribers[i+1:]...)
			i--
		case cgs.Chan <- entities:
		}
	}
}
