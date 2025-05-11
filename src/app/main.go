package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"app/config"
	"app/domain/repo"
	"app/domain/resolvers"
	"app/domain/svc"
	"app/ent"
	"app/ent/migrate"
	"app/generated"
	"app/handler"
	"app/hooks"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/GoLabra/labrago/src/api/cache"
	"github.com/GoLabra/labrago/src/api/constants"
	apiSvc "github.com/GoLabra/labrago/src/api/domain/svc"
	"github.com/GoLabra/labrago/src/api/entgql/annotations"
	"github.com/GoLabra/labrago/src/api/entgql/entity"
	entityGenerated "github.com/GoLabra/labrago/src/api/entgql/generated"
	"github.com/GoLabra/labrago/src/api/entgql/generator"
	entityResolvers "github.com/GoLabra/labrago/src/api/entgql/resolvers"
	"github.com/GoLabra/labrago/src/api/strcase"
	"github.com/GoLabra/labrago/src/api/subscription"
	"github.com/GoLabra/labrago/src/api/utils"
	"github.com/centrifugal/gocent/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/cors"
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
	); err != nil {
		panic(err)
	}

	graph, err := entc.LoadGraph(conf.EntSchemaPath, &gen.Config{})
	if err != nil {
		panic(err)
	}

	cache.NewEntityCache(1 * time.Hour)
	cache.NewEdgeCache(1 * time.Hour)
	cache.NewFieldCache(1 * time.Hour)

	LoadSchema(graph)

	repository := repo.New(client)

	osFileSystem := utils.NewOSFileSystem()

	graphqlSubscriptionClient := subscription.NewGraphqlSubscriptionClient()

	schemaManager := generator.NewSchemaManager(osFileSystem, "./schema", ".", graphqlSubscriptionClient)

	service := svc.New(repository, schemaManager)

	gocentClient := gocent.New(gocent.Config{
		Addr: conf.CentrifugoApiAddress,
		Key:  conf.CentrifugoKey,
	})

	resolver := &resolvers.Resolver{
		Service: service,
	}

	entityResolvers := &entityResolvers.Resolver{
		Service:            &struct{ Entity *apiSvc.Entity }{apiSvc.NewEntity(schemaManager)},
		SubscriptionClient: graphqlSubscriptionClient,
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
		entitySrv := gqlHandler.New(entityGenerated.NewExecutableSchema(
			entityGenerated.Config{
				Resolvers: entityResolvers,
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

		entitySrv.AddTransport(transport.POST{})
		entitySrv.AddTransport(transport.GET{})
		entitySrv.AddTransport(transport.Options{})
		// AV: Please review this
		entitySrv.AddTransport(&transport.Websocket{
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
			},
		})
		entitySrv.Use(extension.Introspection{})

		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(handler.Authenticator)

		router.Handle("/query", srv)
		router.Handle("/entity", entitySrv)
	})

	router.Handle("/playground", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/eplayground", handler.Playground("GraphQL playground", "/entity"))

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

func LoadSchema(graph *gen.Graph) {
	for _, node := range graph.Nodes {
		var entityAnnotations annotations.Entity
		err := mapstructure.Decode(node.Annotations["Entity"], &entityAnnotations)
		if err != nil {
			panic(err)
		}

		entityName := strcase.NodeNameToGraphqlName(node.Name)
		cache.Entity.Set(entityName, entity.Entity{
			Name:             entityName,
			EntName:          node.Name,
			Caption:          entityAnnotations.Caption,
			Owner:            entityAnnotations.Owner,
			DisplayFieldName: entityAnnotations.DisplayField,
		})

		fields := []entity.Field{
			{
				Caption: "Id",
				Name:    "id",
				Type:    string(entity.FieldTypeID),
			},
		}
		for _, nodeField := range node.Fields {
			var fieldAnnotations annotations.Field
			err := mapstructure.Decode(nodeField.Annotations["Field"], &fieldAnnotations)
			if err != nil {
				panic(err)
			}

			required := !nodeField.Optional
			unique := nodeField.Unique

			field := entity.Field{
				Name:           strcase.ToLowerCamel(nodeField.Name),
				EntName:        nodeField.Name,
				Caption:        fieldAnnotations.Caption,
				Type:           string(fieldAnnotations.Type),
				Required:       &required,
				Unique:         &unique,
				Nillable:       nodeField.Nillable,
				UpdateDefault:  nodeField.UpdateDefault,
				AcceptedValues: fieldAnnotations.AcceptedValues,
			}

			if nodeField.Default {
				defaultValue := fmt.Sprint(nodeField.DefaultValue())
				if fieldAnnotations.DefaultValue != "" {
					defaultValue = fieldAnnotations.DefaultValue
				}
				field.DefaultValue = &defaultValue
			}

			if fieldAnnotations.Min != "" {
				field.Min = &fieldAnnotations.Min
			}

			if fieldAnnotations.Max != "" {
				field.Max = &fieldAnnotations.Max
			}

			if fieldAnnotations.Private {
				field.Private = &fieldAnnotations.Private
			}

			fields = append(fields, field)
		}
		cache.Field.Set(entityName, fields)

		edges := []entity.Edge{}
		for _, edge := range node.Edges {
			var edgeAnnotations annotations.Edge
			err := mapstructure.Decode(edge.Annotations["Edge"], &edgeAnnotations)
			if err != nil {
				panic(err)
			}

			required := !edge.Optional
			ref := ""
			if edge.Ref != nil && edge.IsInverse() {
				ref = edge.Ref.Name
			}

			edges = append(edges, entity.Edge{
				Name:         strcase.ToLowerCamel(edge.Name),
				EntName:      edge.Name,
				Caption:      edgeAnnotations.Caption,
				Required:     &required,
				Type:         edge.Type.Name,
				Ref:          ref,
				RelationType: edgeAnnotations.RelationType,
			})
		}
		cache.Edge.Set(entityName, edges)
	}
}
