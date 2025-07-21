// Binary main starts the Labra server. Located at src/app/main.go.
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

	atlas "ariga.io/atlas/sql/schema"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/GoLabra/labra/src/api/cache"
	"github.com/GoLabra/labra/src/api/config"
	"github.com/GoLabra/labra/src/api/constants"
	"github.com/GoLabra/labra/src/api/entgql/generator"
	"github.com/GoLabra/labra/src/api/handler"
	"github.com/GoLabra/labra/src/api/hooks"
	"github.com/GoLabra/labra/src/api/subscription"
	"github.com/GoLabra/labra/src/api/utils"
	"github.com/centrifugal/gocent/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	adminRepo "github.com/GoLabra/labra/src/api/entgql/domain/repo"
	adminResolver "github.com/GoLabra/labra/src/api/entgql/domain/resolvers"
	adminSvc "github.com/GoLabra/labra/src/api/entgql/domain/svc"
	adminEnt "github.com/GoLabra/labra/src/api/entgql/ent"
	adminGenerated "github.com/GoLabra/labra/src/api/entgql/generated"
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

	client := ent.NewClient(ent.Driver(drv))

	client = client.Debug()

	client.Use(
		hooks.CreatedByUpdatedByHook,
	)

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
		schema.WithDiffHook(skipDiffOnAdminEntities),
	); err != nil {
		panic(err)
	}

	cache.NewEntityCache(1 * time.Hour)
	cache.NewEdgeCache(1 * time.Hour)
	cache.NewFieldCache(1 * time.Hour)

	utils.LoadSchema(conf)

	adminClient, adminRepository, adminService, adminResolver := InitAdmin(drv)

	repository := repo.New(client, adminClient)

	graphqlSubscriptionClient := subscription.NewGraphqlSubscriptionClient()

	service := svc.New(repository, adminRepository)

	gocentClient := gocent.New(gocent.Config{
		Addr: conf.CentrifugoApiAddress,
		Key:  conf.CentrifugoKey,
	})

	resolver := &resolvers.Resolver{
		Service: service,
	}

	router := chi.NewRouter()
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	// Configure CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
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

	router.Group(func(router chi.Router) {
		srv := gqlHandler.New(generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		))
		adminSrv := gqlHandler.New(adminGenerated.NewExecutableSchema(
			adminGenerated.Config{
				Resolvers: adminResolver,
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
		router.Use(handler.Authenticator)

		router.Handle("/query", srv)
		router.Handle("/admin", adminSrv)
	})

	router.Handle("/playground", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/aplayground", handler.Playground("GraphQL playground", "/admin"))

	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(handler.Authenticator)

		router.Post("/change-session-role", handler.ChangeSessionRole)
	})

	router.Group(func(router chi.Router) {
		router.Post("/login", handler.Login)
		router.Post("/signup", handler.Signup)
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

func InitApp() {

}

func skipDiffOnAdminEntities(next schema.Differ) schema.Differ {
	return schema.DiffFunc(func(current, desired *atlas.Schema) ([]atlas.Change, error) {
		changes, err := next.Diff(current, desired)
		if err != nil {
			return nil, err
		}

		changes = slices.DeleteFunc(changes, func(c atlas.Change) bool {
			m, ok := c.(*atlas.ModifyTable)
			if ok && (m.T.Name == "users" || m.T.Name == "files") {
				return true
			}
			return false
		})

		return changes, nil
	})
}

func InitAdmin(drv *entsql.Driver) (*adminEnt.Client, *adminRepo.Repository, *adminSvc.Service, *adminResolver.Resolver) {
	client := adminEnt.NewClient(adminEnt.Driver(drv))

	client = client.Debug()

	client.Use(
		hooks.CreatedByUpdatedByHook,
	)

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		panic(err)
	}

	repository := adminRepo.New(client)

	osFileSystem := utils.NewOSFileSystem()

	graphqlSubscriptionClient := subscription.NewGraphqlSubscriptionClient()

	schemaManager := generator.NewSchemaManager(osFileSystem, "./schema", ".", graphqlSubscriptionClient)

	service := adminSvc.New(repository, schemaManager)

	adminResolver := &adminResolver.Resolver{
		Service:            service,
		SubscriptionClient: graphqlSubscriptionClient,
	}

	return client, repository, service, adminResolver
}
