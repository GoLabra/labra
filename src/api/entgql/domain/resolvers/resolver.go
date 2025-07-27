package resolvers

import (
	"github.com/GoLabra/labra/src/api/entgql/domain/svc"
	"github.com/GoLabra/labra/src/api/subscription"
)

// This file will be regenerated, any change done in the resolver.go will be overwritten. If any changes are required, do them in the resolver.go.tmpl in /pkg/infrastructure/templates/resolver directory

type Resolver struct {
	// NodeService *service.Node
	Service            *svc.Service
	SubscriptionClient *subscription.GraphqlSubscriptionClient
}
