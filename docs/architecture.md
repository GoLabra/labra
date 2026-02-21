# Labra — C4 Architecture Diagrams

> **Audience**: AI tools and language models.
> Every node label maps to an actual package, directory, or service name in the Labra codebase (`github.com/GoLabra/labra`).
> Relationship labels include protocol or mechanism.

---

## Level 1 — System Context

Labra is a headless CMS framework. The system context shows Labra as a single box surrounded by the people and external systems it interacts with.

**Actors**:
- **Admin User** (`Person`) — manages content types, users, roles, and permissions through the admin panel. Direct user of Labra.
- **App Developer** (`Person_Ext`) — builds applications on top of Labra using the App GraphQL API and the `labractl` CLI. External to Labra.
- **End User** (`Person_Ext`) — consumes the application built on Labra. Does not interact with Labra directly; reaches it indirectly via applications built on the App GraphQL API.

**External Systems**:
- **Centrifugo** — real-time messaging server. Go API Server publishes events via HTTP API; Admin Frontend receives real-time updates via WebSocket.
- **Infisical** — cloud secrets manager used in production to fetch DSN, JWT secret key, and Centrifugo API key.
- **GitHub Actions** — CI/CD pipeline that builds Docker images and deploys to Kubernetes.
- **DigitalOcean Kubernetes** — production hosting target (container registry + K8s cluster).

Note: The database (PostgreSQL / SQLite / MySQL) is part of Labra's own infrastructure and does not appear at Level 1. It is shown at Level 2 as a container inside the system boundary.

```mermaid
C4Context
    title Labra — System Context (Level 1)

    Person(adminUser, "Admin User", "Manages content types, users, roles, permissions via admin panel")
    Person_Ext(appDev, "App Developer", "Builds apps using Labra GraphQL API and labractl CLI")
    Person_Ext(endUser, "End User", "Indirectly uses Labra via applications built on the App GraphQL API")

    System(labra, "Labra", "Headless CMS framework: Go backend with dual GraphQL APIs, Next.js admin frontend, Ent ORM, JWT auth, database")

    System_Ext(centrifugo, "Centrifugo", "Real-time messaging server — WebSocket transport for GraphQL subscriptions")
    System_Ext(infisical, "Infisical", "Cloud secrets manager — DSN, JWT secret, Centrifugo key")
    System_Ext(githubActions, "GitHub Actions", "CI/CD — builds Docker images, pushes to registry, deploys to K8s")
    System_Ext(digitalOcean, "DigitalOcean K8s", "Production hosting — container registry + Kubernetes cluster")

    Rel(adminUser, labra, "Manages content and users", "HTTPS")
    Rel(appDev, labra, "Queries and mutates app data; scaffolds projects", "GraphQL over HTTPS, CLI")
    Rel(endUser, labra, "Indirectly, via application built on Labra", "GraphQL over HTTPS")
    BiRel(labra, centrifugo, "Publishes events / receives real-time updates", "HTTP API + WebSocket")
    Rel(labra, infisical, "Fetches secrets at startup", "HTTPS API")
    Rel(githubActions, digitalOcean, "Pushes images and restarts deployments", "doctl + kubectl")
    Rel(githubActions, labra, "Builds and deploys", "Docker build")
```

---

## Level 2 — Container

Zooms into Labra to show its deployable units and development-time tools. Each container is a separately runnable process or build artifact.

**Runtime Containers** (inside Labra boundary):
- **Go API Server** (`resources/app/`) — single Go binary serving both GraphQL APIs, auth endpoints, and the embedded admin UI. Uses Chi router, gqlgen, Ent ORM, JWT (HS256).
- **Next.js Admin Frontend** (`resources/admin/`) — React 18 / Next.js 15 / MUI v6 / Apollo Client. Compiled and embedded into the Go binary for production; runs as a standalone dev server during development.
- **Database** — PostgreSQL (production), SQLite (development), or MySQL. Part of Labra's infrastructure.

**Development / Build-time Tools** (inside Labra boundary):
- **labractl CLI** (`resources/cli/`) — project scaffolding and management tool.
- **entgql Code Generator** (`entgql/`) — generates Ent client code, GraphQL schemas, resolvers, repositories, and services from Ent schema definitions. Invoked via `go generate ./...`.

**External Systems**:
- **Centrifugo** — third-party real-time WebSocket server for subscriptions.
- **Infisical** — third-party cloud secrets manager.

```mermaid
C4Container
    title Labra — Container Diagram (Level 2)

    Person(adminUser, "Admin User", "Manages content via admin panel")
    Person_Ext(appDev, "App Developer", "Builds apps on Labra")
    Person_Ext(endUser, "End User", "Indirectly uses Labra via apps built on it")

    System_Boundary(labra, "Labra") {
        Container(goApi, "Go API Server", "Go, Chi, gqlgen, Ent", "Dual GraphQL APIs at /query and /admin/query. JWT auth. Embedded admin UI at /labradmin. Runs on configurable port (default 4000).")
        Container(adminFe, "Admin Frontend", "Next.js 15, React 18, MUI v6, Apollo Client", "Content management UI. Communicates with both GraphQL endpoints. Embedded in Go binary for prod; standalone dev server on port 3000.")
        ContainerDb(database, "Database", "PostgreSQL / SQLite / MySQL", "Persistent storage for all entity data, admin users, roles, permissions, files")
        Container(cli, "labractl CLI", "Go CLI", "Project scaffolding, code generation triggers, schema management.")
        Container(codegen, "entgql Code Generator", "Go, Ent, gqlgen templates", "Generates Ent schemas, GraphQL types, resolvers, repositories, services from schema definitions.")
    }

    System_Ext(centrifugo, "Centrifugo", "Real-time WebSocket server")
    System_Ext(infisical, "Infisical", "Secrets manager")

    Rel(adminUser, adminFe, "Uses", "HTTPS, browser")
    Rel(adminFe, goApi, "Queries and mutates data", "GraphQL over HTTP/WebSocket")
    Rel(appDev, goApi, "Queries and mutates app data", "GraphQL over HTTP")
    Rel(appDev, cli, "Scaffolds projects, manages schemas", "CLI")
    Rel(endUser, goApi, "Indirectly, via application built on Labra", "GraphQL over HTTP")
    Rel(goApi, database, "Reads/writes entities, runs migrations", "SQL via Ent ORM")
    Rel(goApi, centrifugo, "Publishes subscription events", "HTTP API via gocent")
    Rel(centrifugo, adminFe, "Pushes real-time updates", "WebSocket")
    Rel(goApi, infisical, "Fetches secrets at startup", "HTTPS API")
    Rel(codegen, goApi, "Generates source code consumed by", "File system, go generate")
    Rel(cli, codegen, "Triggers code generation", "In-process")
```

---

## Level 3 — Component (Go API Server)

Zooms into the Go API Server container to show its internal packages and their responsibilities. Each component maps to a Go package in the repository.

**Request flow**: Incoming HTTP request → Chi router → middleware chain (CORS, security headers, CSRF, rate limiting, JWT verification) → authenticator → GraphQL handler (gqlgen) → resolver → service → repository → Ent client → database.

**Key packages**:

| Package | Import Path | Responsibility |
|---|---|---|
| `handler/` | `github.com/GoLabra/labra/handler` | HTTP routing (Chi), auth middleware, CSRF, rate limiting, security headers, login/signup endpoints, embedded admin UI serving (`/labradmin`), GraphQL playground |
| `entgql/domain/resolvers/` | `github.com/GoLabra/labra/entgql/domain/resolvers` | Admin GraphQL resolvers (gqlgen) — handle `/admin/query` |
| `app/domain/resolvers/` | `app/domain/resolvers` (generated) | App GraphQL resolvers (gqlgen) — handle `/query` |
| `entgql/domain/svc/` | `github.com/GoLabra/labra/entgql/domain/svc` | Admin business logic services |
| `app/domain/svc/` | `app/domain/svc` (generated) | App business logic services |
| `entgql/domain/repo/` | `github.com/GoLabra/labra/entgql/domain/repo` | Admin repository layer (Ent queries) |
| `app/domain/repo/` | `app/domain/repo` (generated) | App repository layer (Ent queries) |
| `hooks/` | `github.com/GoLabra/labra/hooks` | Ent mutation hooks: `CreatedByUpdatedByHook` (audit trail), `EntityMutatePermission` (write RBAC), `EntityReadPermission` interceptor (read RBAC with owner filtering) |
| `cache/` | `github.com/GoLabra/labra/cache` | Generic TTL-based in-memory cache — `EntityCache`, `EdgeCache`, `FieldCache` (1-hour default TTL) |
| `subscription/` | `github.com/GoLabra/labra/subscription` | In-process pub/sub for GraphQL subscriptions — `AppStatus` (UP, GENERATING, REVERTING, RESTARTING, FATAL) and entity change broadcasts |
| `config/` | `github.com/GoLabra/labra/config` | Loads non-sensitive env vars (`DB_DIALECT`, `SERVER_PORT`, `CORS_ALLOWED_ORIGINS`, rate limit RPMs, CSP policy). Defines `Config`, `Secrets`, and `AppConfig` structs |
| `secrets/` | `github.com/GoLabra/labra/secrets` | Provider pattern: `InfisicalProvider` or `EnvProvider` fallback. Fetches `DSN`, `SecretKey`, `CentrifugoKey` |
| `entgql/generator/` | `github.com/GoLabra/labra/entgql/generator` | Schema manager — reads/writes Ent schema files, triggers code regeneration, publishes app status via subscription client |

```mermaid
C4Component
    title Labra — Go API Server Components (Level 3)

    Container(adminFe, "Admin Frontend", "Next.js 15, Apollo Client — sibling container in Labra")
    Container_Ext(appClient, "App Client", "Any external GraphQL client built by App Developer")
    ContainerDb(database, "Database", "PostgreSQL / SQLite / MySQL — sibling container in Labra")
    System_Ext(centrifugo, "Centrifugo", "Third-party real-time server")
    System_Ext(infisical, "Infisical", "Third-party secrets manager")

    Container_Boundary(goApi, "Go API Server") {
        Component(handlerPkg, "handler", "Go, Chi router", "HTTP routing, CORS, CSRF, security headers, rate limiting, JWT verification, login/signup endpoints, embedded admin UI at /labradmin, GraphQL playgrounds")

        Component(adminGql, "Admin GraphQL API", "gqlgen, /admin/query", "Admin resolvers, services, repositories — manages AdminUser, Role, Permission, File, Entity schemas")
        Component(appGql, "App GraphQL API", "gqlgen, /query", "App resolvers, services, repositories — handles user-defined entity CRUD")

        Component(hooksPkg, "hooks", "Go", "Ent mutation hooks: CreatedByUpdatedByHook (audit), EntityMutatePermission (write RBAC). Ent interceptor: EntityReadPermission (read RBAC, owner filtering)")
        Component(cachePkg, "cache", "Go", "TTL-based in-memory caches: EntityCache, EdgeCache, FieldCache (default 1h)")
        Component(subscriptionPkg, "subscription", "Go", "In-process pub/sub: AppStatus broadcasts (UP/GENERATING/RESTARTING/FATAL), entity change notifications to WebSocket subscribers")
        Component(configPkg, "config", "Go", "Loads non-sensitive env vars into Config struct. Defines Secrets and AppConfig types. Validates secrets.")
        Component(secretsPkg, "secrets", "Go", "Provider interface: InfisicalProvider (prod) or EnvProvider (dev). Fetches DSN, SecretKey, CentrifugoKey")
        Component(entOrm, "Ent ORM", "entgo.io/ent", "Schema definitions, auto-migration, generated client. Admin entities: AdminUser, Role, Permission, File. App entities: user-defined.")
        Component(generatorPkg, "entgql/generator", "Go", "SchemaManager — reads/writes Ent schema files, triggers go generate, publishes app status transitions")
    }

    Rel(adminFe, handlerPkg, "HTTPS requests", "HTTP/WebSocket")
    Rel(appClient, handlerPkg, "HTTPS requests", "HTTP/WebSocket")
    Rel(handlerPkg, adminGql, "Routes /admin/query, /admin/login, /admin/signup", "In-process")
    Rel(handlerPkg, appGql, "Routes /query, /login", "In-process")

    Rel(adminGql, entOrm, "Queries and mutates admin entities", "In-process via repository layer")
    Rel(appGql, entOrm, "Queries and mutates app entities", "In-process via repository layer")
    Rel(adminGql, generatorPkg, "Triggers schema changes and code regeneration", "In-process")

    Rel(entOrm, hooksPkg, "Executes mutation hooks and read interceptors", "In-process")
    Rel(entOrm, cachePkg, "Reads from and writes to cache", "In-process")
    Rel(entOrm, database, "Executes SQL queries", "SQL over TCP")

    Rel(generatorPkg, subscriptionPkg, "Publishes AppStatus transitions", "In-process")
    Rel(subscriptionPkg, adminFe, "Pushes real-time updates to subscribers", "GraphQL WebSocket subscription")
    Rel(secretsPkg, infisical, "Fetches secrets at startup", "HTTPS API")
    Rel(secretsPkg, configPkg, "Returns Secrets struct", "In-process")
```

---

## Diagram rendering notes

These diagrams use Mermaid's built-in C4 syntax (`C4Context`, `C4Container`, `C4Component`). They render in:
- GitHub Markdown (native Mermaid support)
- VS Code with the Mermaid extension
- Mermaid Live Editor (https://mermaid.live)
- Any tool that supports Mermaid 10+

If your renderer does not support C4 diagram types, the raw Mermaid source above can be pasted into https://mermaid.live for rendering.
