package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
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
	"github.com/GoLabra/labra/secrets"
	"github.com/GoLabra/labra/subscription"
	"github.com/GoLabra/labra/utils"
	"github.com/centrifugal/gocent/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"

	adminRepo "github.com/GoLabra/labra/entgql/domain/repo"
	adminResolvers "github.com/GoLabra/labra/entgql/domain/resolvers"
	adminSvc "github.com/GoLabra/labra/entgql/domain/svc"
	adminEnt "github.com/GoLabra/labra/entgql/ent"
	adminGenerated "github.com/GoLabra/labra/entgql/generated"
)

func main() {
	ctx := context.Background()

	// Load non-sensitive configuration from environment
	conf, err := config.New()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Load secrets from Infisical or environment variables
	log.Println("Loading secrets...")
	secretsData, err := secrets.LoadSecrets(ctx, conf.Environment)
	if err != nil {
		fmt.Printf("Failed to load secrets: %v\n", err)
		os.Exit(1)
	}

	// Combine config and secrets into AppConfig
	appConfig, err := config.NewAppConfig(conf, secretsData)
	if err != nil {
		fmt.Printf("Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	log.Println("Secrets loaded successfully")

	// Get the actual driver name for sql.Open()
	driverName, err := getDriverName(appConfig.DBDialect)
	if err != nil {
		fmt.Printf("Invalid DB_DIALECT: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open(driverName, appConfig.DSN)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = db.Ping()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Parse dialect from config for ent
	entDialect, err := parseDialect(appConfig.DBDialect)
	if err != nil {
		fmt.Printf("Invalid DB_DIALECT: %v\n", err)
		os.Exit(1)
	}

	// SQLite requires max open connections of 1 to prevent database locking
	if entDialect == dialect.SQLite {
		db.SetMaxOpenConns(1)
	} else {
		db.SetMaxOpenConns(100)
	}

	drv := entsql.OpenDB(entDialect, db)

	cache.NewEntityCache(1 * time.Hour)
	cache.NewEdgeCache(1 * time.Hour)
	cache.NewFieldCache(1 * time.Hour)

	utils.LoadSchema(appConfig)

	adminClient, adminRepository, adminService, adminResolver := InitAdmin(drv)
	appClient, repository, service, resolver := InitApp(drv, adminClient, adminRepository, adminService)

	graphqlSubscriptionClient := subscription.NewGraphqlSubscriptionClient()

	gocentClient := gocent.New(gocent.Config{
		Addr: appConfig.CentrifugoApiAddress,
		Key:  appConfig.CentrifugoKey,
	})

	tokenAuth := jwtauth.New("HS256", []byte(appConfig.SecretKey), nil)

	// Configure CORS
	allowedOrigins := appConfig.Config.CORSAllowedOriginsList()

	// Safety rail: never allow "*" when credentials are enabled
	for _, o := range allowedOrigins {
		if strings.TrimSpace(o) == "*" {
			log.Fatal("invalid CORS config: wildcard '*' is not allowed when AllowCredentials is true")
		}
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-CSRF-Token", "X-CSRF-TOKEN"},
	})
	router := chi.NewRouter()
	router.Use(corsMiddleware.Handler)
	router.Use(adminHandler.SecurityHeaders(&appConfig.Config))
	router.Use(adminHandler.CSRFMiddleware)

	// Rate limiters (per IP)
	loginLimiter := adminHandler.NewIPRateLimiter(appConfig.Config.AuthLoginRateLimitRPM)
	signupLimiter := adminHandler.NewIPRateLimiter(appConfig.Config.AuthSignupRateLimitRPM)
	apiLimiter := adminHandler.NewIPRateLimiter(appConfig.Config.AuthAPIRateLimitRPM)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			ctx = context.WithValue(ctx, constants.AdminServiceContextValue, adminService)
			ctx = context.WithValue(ctx, constants.AdminRepositoryContextValue, adminRepository)
			ctx = context.WithValue(ctx, constants.ServiceContextValue, service)
			ctx = context.WithValue(ctx, constants.RepositoryContextValue, repository)
			ctx = context.WithValue(ctx, constants.EntClientContextValue, appClient)
			ctx = context.WithValue(ctx, constants.AdminEntClientContextValue, adminClient)
			ctx = context.WithValue(ctx, constants.CentrifugeClientContextValue, gocentClient)
			ctx = context.WithValue(ctx, "config", appConfig)
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

		router.Use(apiLimiter.Middleware)
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(handler.Authenticator)

		router.Handle("/query", srv)
	})
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(handler.Authenticator)
		router.Use(apiLimiter.Middleware)

		router.Post("/change-session-role", handler.ChangeSessionRole)
		router.Post("/logout", handler.Logout)
	})
	router.Group(func(router chi.Router) {
		router.With(loginLimiter.Middleware).Post("/login", handler.Login)
		router.With(loginLimiter.Middleware).Post("/refresh", handler.Refresh)
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

		router.Use(apiLimiter.Middleware)
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(adminHandler.Authenticator)

		router.Handle("/admin/query", adminSrv)
	})
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(adminHandler.Authenticator)
		router.Use(apiLimiter.Middleware)

		router.Post("/admin/change-session-role", adminHandler.ChangeSessionRole)
		router.Post("/admin/logout", adminHandler.Logout)
	})
	router.Group(func(router chi.Router) {
		router.With(loginLimiter.Middleware).Post("/admin/login", adminHandler.Login)
		router.With(loginLimiter.Middleware).Post("/admin/refresh", adminHandler.Refresh)
		router.With(signupLimiter.Middleware).Post("/admin/signup", adminHandler.Signup)
		router.Mount("/labradmin", http.StripPrefix("/labradmin", adminHandler.ServeAdmin()))
		router.Handle("/admin/playground", adminHandler.Playground("GraphQL playground", "/admin/query"))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.ServerPort),
		Handler: router,
	}

	log.Printf("Server starting on port %s\n", appConfig.ServerPort)

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

// getDriverName converts a database dialect string to the actual driver name for sql.Open().
func getDriverName(dialectStr string) (string, error) {
	switch strings.ToLower(dialectStr) {
	case "sqlite", "sqlite3":
		return "sqlite3", nil
	case "postgres", "postgresql":
		return "postgres", nil
	case "mysql":
		return "mysql", nil
	default:
		return "", fmt.Errorf("unsupported database dialect: %s (supported: sqlite, postgres, mysql)", dialectStr)
	}
}

// parseDialect converts a database dialect string to an ent dialect constant.
func parseDialect(dialectStr string) (string, error) {
	switch strings.ToLower(dialectStr) {
	case "sqlite", "sqlite3":
		return dialect.SQLite, nil
	case "postgres", "postgresql":
		return dialect.Postgres, nil
	case "mysql":
		return dialect.MySQL, nil
	default:
		return "", fmt.Errorf("unsupported database dialect: %s (supported: sqlite, postgres, mysql)", dialectStr)
	}
}
