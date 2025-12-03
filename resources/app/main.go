package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"app/domain/repo"
	"app/domain/resolvers"
	"app/domain/svc"
	"app/ent"
	"app/ent/migrate"
	"app/generated"
	"app/handler"

	atlas "ariga.io/atlas/sql/schema"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/GoLabra/labra/cache"
	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/constants"
	"github.com/GoLabra/labra/entgql/generator"
	adminHandler "github.com/GoLabra/labra/handler"
	"github.com/GoLabra/labra/hooks"
	"github.com/GoLabra/labra/subscription"
	"github.com/GoLabra/labra/utils"
	"github.com/centrifugal/gocent/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	adminRepo "github.com/GoLabra/labra/entgql/domain/repo"
	adminResolvers "github.com/GoLabra/labra/entgql/domain/resolvers"
	adminSvc "github.com/GoLabra/labra/entgql/domain/svc"
	adminEnt "github.com/GoLabra/labra/entgql/ent"
	adminGenerated "github.com/GoLabra/labra/entgql/generated"
)

func main() {
	conf, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open(conf.DBDialect, conf.DSN)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = db.Ping()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(100)

	drv := entsql.OpenDB(dialect.Postgres, db)

	cache.NewEntityCache(1 * time.Hour)
	cache.NewEdgeCache(1 * time.Hour)
	cache.NewFieldCache(1 * time.Hour)

	utils.LoadSchema(conf)

	adminClient, adminRepository, adminService, adminResolver := InitAdmin(drv)
	_, repository, service, resolver := InitApp(drv, adminClient, adminRepository, adminService)

	graphqlSubscriptionClient := subscription.NewGraphqlSubscriptionClient()

	gocentClient := gocent.New(gocent.Config{
		Addr: conf.CentrifugoApiAddress,
		Key:  conf.CentrifugoKey,
	})

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	// Configure CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	router := chi.NewRouter()
	router.Use(corsMiddleware.Handler)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			ctx = context.WithValue(ctx, constants.AdminServiceContextValue, adminService)
			ctx = context.WithValue(ctx, constants.AdminRepositoryContextValue, adminRepository)
			ctx = context.WithValue(ctx, constants.ServiceContextValue, service)
			ctx = context.WithValue(ctx, constants.RepositoryContextValue, repository)
			ctx = context.WithValue(ctx, constants.CentrifugeClientContextValue, gocentClient)
			ctx = context.WithValue(ctx, "config", conf)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// APP routes
	router.Group(func(router chi.Router) {
		srv := gqlHandler.New(generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		))

		srv.AddTransport(transport.POST{})
		srv.AddTransport(transport.GET{})
		srv.AddTransport(transport.Options{})
		// AV: Please review this
		srv.AddTransport(&transport.Websocket{
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
			},
		})
		srv.Use(extension.Introspection{})

		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(handler.Authenticator)

		router.Handle("/query", srv)
	})
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(handler.Authenticator)

		router.Post("/change-session-role", handler.ChangeSessionRole)
	})
	router.Group(func(router chi.Router) {
		router.Post("/login", handler.Login)
		router.Mount("/labradmin", http.StripPrefix("/labradmin", adminHandler.ServeAdmin()))
		router.Handle("/playground", adminHandler.Playground("GraphQL playground", "/query"))
	})

	// ADMIN routes
	router.Group(func(router chi.Router) {
		adminSrv := gqlHandler.New(adminGenerated.NewExecutableSchema(
			adminGenerated.Config{
				Resolvers: adminResolver,
			},
		))

		adminSrv.AddTransport(transport.POST{})
		adminSrv.AddTransport(transport.GET{})
		adminSrv.AddTransport(transport.Options{})
		// AV: Please review this
		adminSrv.AddTransport(&transport.Websocket{
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
			},
		})
		adminSrv.Use(extension.Introspection{})

		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(adminHandler.Authenticator)

		router.Handle("/admin/query", adminSrv)
	})
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(adminHandler.Authenticator)

		router.Post("/admin/change-session-role", adminHandler.ChangeSessionRole)
	})
	router.Group(func(router chi.Router) {
		router.Post("/admin/login", adminHandler.Login)
		router.Post("/admin/signup", adminHandler.Signup)
		router.Handle("/admin/playground", adminHandler.Playground("GraphQL playground", "/admin/query"))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", conf.ServerPort),
		Handler: router,
	}

	log.Printf("Server starting on port %s\n", conf.ServerPort)

	graphqlSubscriptionClient.PublishAppStatusMessage(subscription.AppStatusUp)

	err = server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

func InitApp(drv *entsql.Driver, adminClient *adminEnt.Client, adminRepository *adminRepo.Repository, adminService *adminSvc.Service) (*ent.Client, *repo.Repository, *svc.Service, *resolvers.Resolver) {
	client := ent.NewClient(ent.Driver(drv))

	client = client.Debug()

	client.Use(
		hooks.CreatedByUpdatedByHook,
		hooks.EntityMutatePermission,
	)

	client.Intercept(hooks.EntityReadPermission())

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
		schema.WithDiffHook(skipDiffOnAdminEntities),
	); err != nil {
		panic(err)
	}

	repository := repo.New(client, adminClient)
	service := svc.New(repository, adminRepository)
	resolver := &resolvers.Resolver{
		Service: service,
	}

	return client, repository, service, resolver
}

func skipDiffOnAdminEntities(next schema.Differ) schema.Differ {
	return schema.DiffFunc(func(current, desired *atlas.Schema) ([]atlas.Change, error) {
		changes, err := next.Diff(current, desired)
		if err != nil {
			return nil, err
		}

		changes = slices.DeleteFunc(changes, func(c atlas.Change) bool {
			m, ok := c.(*atlas.ModifyTable)
			if ok && (m.T.Name == "admin_users" || m.T.Name == "files" || m.T.Name == "roles") {
				return true
			}
			return false
		})

		return changes, nil
	})
}

func InitAdmin(drv *entsql.Driver) (*adminEnt.Client, *adminRepo.Repository, *adminSvc.Service, *adminResolvers.Resolver) {
	client := adminEnt.NewClient(adminEnt.Driver(drv))

	client = client.Debug()

	client.Use(
		hooks.CreatedByUpdatedByHook,
		hooks.EntityMutatePermission,
	)

	client.Intercept(hooks.EntityReadPermission())

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
		schema.WithDiffHook(skipDiffOnUserEntities),
	); err != nil {
		panic(err)
	}

	repository := adminRepo.New(client)

	osFileSystem := utils.NewOSFileSystem()

	graphqlSubscriptionClient := subscription.NewGraphqlSubscriptionClient()

	schemaManager := generator.NewSchemaManager(osFileSystem, "./schema", ".", graphqlSubscriptionClient)

	service := adminSvc.New(repository, schemaManager)

	adminResolver := &adminResolvers.Resolver{
		Service:            service,
		SubscriptionClient: graphqlSubscriptionClient,
	}

	return client, repository, service, adminResolver
}

func skipDiffOnUserEntities(next schema.Differ) schema.Differ {
	return schema.DiffFunc(func(current, desired *atlas.Schema) ([]atlas.Change, error) {
		changes, err := next.Diff(current, desired)
		if err != nil {
			return nil, err
		}

		changes = slices.DeleteFunc(changes, func(c atlas.Change) bool {
			m, ok := c.(*atlas.ModifyTable)
			if ok && (m.T.Name == "users") {
				return true
			}
			return false
		})

		return changes, nil
	})
}
