# What is LabraGo?

## It's a headless CMS… but not *that* kind.

We're not here to help you build a blog or a basic website. LabraGo was born to handle **real apps**—the complex, data-heavy, API-driven kind.

Built with Go for performance and React for a clean, modern admin UI, it's fast, flexible, and actually enjoyable to work with.

Think of it as your backend brain—powerful enough to manage complex data models, events, and permissions, without drowning you in endless configurations.

<img height="96" src="https://github.com/GoLabra/labra/blob/develop/labra.jpeg" title="Labra Logo" width="96"/>

# Why does LabraGo exist?

## Because we exist 🙂

And because we've run into enough problems building data-heavy apps, we decided to do something about it.

This project is our way of making things simpler, cleaner, and more fun for anyone facing the same challenges.

---

# Quick Install

The easiest way to get started with LabraGo is using `labractl`, the official CLI tool. It automates project setup, dependency management, and development server orchestration.

**Note:** `labractl` v1.1.0+ supports Labra v0.1.3.

## Install labractl

```bash
go install github.com/GoLabra/labractl@latest
```

Make sure the Go bin directory is in your PATH. On macOS/Linux:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

## Create a new project

```bash
labractl create my-awesome-project
```

This command will:
- Clone the LabraGo repository
- Configure development environment
- Install all dependencies
- Set up PostgreSQL database
- Generate configuration files

## Start development servers

```bash
cd my-awesome-project
labractl start
```

This launches:
- **Backend server** on `http://localhost:4000`
- **Frontend admin** on `http://localhost:3000`
- **GraphQL Playground** on `http://localhost:4000/playground`

For more information about `labractl`, visit the [labractl repository](https://github.com/GoLabra/labractl).

---

# Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Go** (version 1.18 or later) - [Download](https://golang.org/dl/)
- **Database** - PostgreSQL, SQLite, or MySQL (with an empty database ready)
  - **PostgreSQL** - [Download](https://www.postgresql.org/download/)
  - **SQLite** - Usually pre-installed on most systems, or [Download](https://www.sqlite.org/download.html)
  - **MySQL** - [Download](https://dev.mysql.com/downloads/)
- **Node.js** (version 18 or later) - [Download](https://nodejs.org/)
- **Yarn** (version 1.22 or later) - [Installation Guide](https://yarnpkg.com/getting-started/install)
- **Centrifugo** - Real-time messaging server - [Installation Guide](https://centrifugo.dev/docs/getting-started/installation/)

---

# Getting Started (Manual Installation)

> **Prefer automated setup?** Use the [Quick Install](#quick-install) section above to get started with `labractl` in just a few commands.

This guide will walk you through manually setting up a new LabraGo project from scratch.

## Step 1: Copy the App Template

Start by copying the contents of `resources/app` to your new project directory:

```bash
# Create your new project directory
mkdir my-labrago-project
cd my-labrago-project

# Copy the app template
cp -r /path/to/labra/resources/app/* .
```

Or if you're cloning the LabraGo repository:

```bash
git clone https://github.com/GoLabra/labra.git
cd labra
mkdir my-project
cp -r resources/app/* my-project/
cd my-project
```

## Step 2: Configure Go Module

Update the `go.mod` file to replace the LabraGo dependency with your local path (for development) or your repository path (for production):

```bash
# For local development, replace with the path to the labra repository
sed -i "/REPLACE_LABRAGO_DEVELOPMENT_API/c replace github.com\/GoLabra\/labra => ..\/..\/." go.mod
```

Or manually edit `go.mod`:

```go
replace github.com/GoLabra/labra => /path/to/labra
```

## Step 3: Configure Environment Variables

### Backend Configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your configuration. See [Backend Environment Variables](#backend-environment-variables) for detailed documentation.

**Required variables:**
- `DSN` - Database connection string (PostgreSQL, SQLite, or MySQL)
- `DB_DIALECT` - Database dialect (`postgres`, `sqlite`, or `mysql`)
- `SECRET_KEY` - JWT secret key (generate with: `head -c 32 /dev/urandom | base64`)
- `CENTRIFUGO_API_ADDRESS` - Centrifugo API address
- `CENTRIFUGO_API_KEY` - Centrifugo API key (generate with: `head -c 32 /dev/urandom | base64`)

**Optional variables (with defaults):**
- `SERVER_PORT` - Defaults to `4000` if not set
- `ENT_SCHEMA_PATH` - Defaults to `./ent/schema` if not set
- `FILE_STORAGE_PROVIDER` - Defaults to `local` if not set
- `FILE_STORAGE_PATH` - Defaults to `./storage` if not set

For complete documentation, see [`resources/app/.env.example`](resources/app/.env.example).

### Admin Frontend Configuration

If you're setting up the admin frontend, navigate to the admin directory:

```bash
cd ../admin  # or wherever you have the admin frontend
cp .env.example .env.local
```

Edit `.env.local` with your configuration. See [Admin Frontend Environment Variables](#admin-frontend-environment-variables) for detailed documentation.

For complete documentation, see [`resources/admin/.env.example`](resources/admin/.env.example).

## Step 4: Set Up Database

### Option A: PostgreSQL

Create a PostgreSQL database for your project:

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE labrago;

# Exit psql
\q
```

Update your `.env` file with the PostgreSQL connection string:

```env
DSN=postgres://postgres:yourpassword@localhost:5432/labrago?sslmode=disable
DB_DIALECT=postgres
```

### Option B: SQLite

SQLite requires no setup - the database file will be created automatically. Update your `.env` file:

```env
DSN=file:labra.db?cache=shared&mode=rwc&_fk=1
DB_DIALECT=sqlite
```

**Important:** The `_fk=1` parameter is required to enable foreign key constraints in SQLite, which are necessary for ent to work correctly.

**Note:** SQLite is great for development and small applications. For production with high concurrency, consider PostgreSQL or MySQL.

### Option C: MySQL

Create a MySQL database for your project:

```bash
# Connect to MySQL
mysql -u root -p

# Create database
CREATE DATABASE labrago;

# Exit MySQL
exit
```

Update your `.env` file with the MySQL connection string:

```env
DSN=user:password@tcp(localhost:3306)/labrago?parseTime=true
DB_DIALECT=mysql
```

## Step 5: Set Up File Storage

Create the storage directory (if using the default local storage):

```bash
mkdir -p storage
chmod 755 storage
```

The storage directory will be created automatically if it doesn't exist, but it's good practice to create it beforehand with proper permissions.

See [File Storage Configuration](#file-storage-configuration) for more details.

## Step 6: Generate Code

Install dependencies and generate code:

```bash
go mod tidy
go generate ./...
```

This will:
- Install Go dependencies
- Generate Ent schema code
- Generate GraphQL resolvers and types

## Step 7: Run the Backend

Start the GraphQL API server:

```bash
go run main.go
```

The server will start on the port specified in your `SERVER_PORT` environment variable (default: 4000).

You should see startup logs indicating:
- Successful database connection
- GraphQL endpoint availability at `http://localhost:4000/query`
- Admin GraphQL endpoint at `http://localhost:4000/admin/query`
- GraphQL playgrounds available

## Step 8: Set Up and Run Admin Frontend (Optional)

If you want to use the admin UI:

```bash
# Navigate to admin directory
cd ../admin  # or wherever you have the admin frontend

# Install dependencies
yarn install

# Start development server
yarn dev
```

The admin UI will be available at `http://localhost:3000`.

---

# Environment Variables Reference

## Backend Environment Variables

All backend environment variables are configured in `.env` in your app directory. See [`resources/app/.env.example`](resources/app/.env.example) for a complete reference with inline documentation.

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DSN` | Database connection string | PostgreSQL: `postgres://user:pass@localhost:5432/dbname?sslmode=disable`<br>SQLite: `file:labra.db?cache=shared&mode=rwc&_fk=1`<br>MySQL: `user:pass@tcp(localhost:3306)/dbname?parseTime=true` |
| `DB_DIALECT` | Database dialect | `postgres`, `sqlite`, or `mysql` |
| `SECRET_KEY` | JWT secret key (generate securely) | Generated with `head -c 32 /dev/urandom \| base64` |
| `CENTRIFUGO_API_ADDRESS` | Centrifugo API server address | `http://localhost:8000/api` |
| `CENTRIFUGO_API_KEY` | Centrifugo API key (generate securely) | Generated with `head -c 32 /dev/urandom \| base64` |

### Optional Variables (with defaults)

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `SERVER_PORT` | API server port | `4000` | `4000` |
| `ENT_SCHEMA_PATH` | Path to Ent schema directory | `./ent/schema` | `./ent/schema` |
| `FILE_STORAGE_PROVIDER` | File storage provider | `local` | `local`, `S3`, `GCS` |
| `FILE_STORAGE_PATH` | File storage directory path | `./storage` | `./storage` or `/var/lib/labrago/storage` |

**Security Note:** Always generate new secret keys for production. Never use the example keys from `.env.example`.

## Admin Frontend Environment Variables

All admin frontend environment variables are configured in `.env.local` in the admin directory. See [`resources/admin/.env.example`](resources/admin/.env.example) for a complete reference with inline documentation.

All variables below are optional and have defaults if not specified.

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `NEXT_PUBLIC_BRAND_PRODUCT_NAME` | Product name displayed in UI | `Labra·GO` | `Labra·GO` |
| `NEXT_PUBLIC_BRAND_COLOR` | Brand color theme | `blue` | `blue` |
| `NEXT_PUBLIC_GRAPHQL_API_URL` | Base GraphQL API URL | `http://localhost:4000` | `http://localhost:4000` |
| `NEXT_PUBLIC_GRAPHQL_QUERY_API_URL` | App-level GraphQL query endpoint | `http://localhost:4000/query` | `http://localhost:4000/query` |
| `NEXT_PUBLIC_GRAPHQL_QUERY_PLAYGROUND_URL` | GraphQL query playground URL | `http://localhost:4000/playground` | `http://localhost:4000/playground` |
| `NEXT_PUBLIC_GRAPHQL_ADMIN_API_URL` | Admin-level GraphQL query endpoint | `http://localhost:4000/admin/query` | `http://localhost:4000/admin/query` |
| `NEXT_PUBLIC_GRAPHQL_ADMIN_PLAYGROUND_URL` | Admin GraphQL playground URL | `http://localhost:4000/admin/playground` | `http://localhost:4000/admin/playground` |
| `NEXT_PUBLIC_CENTRIFUGO_URL` | Centrifugo WebSocket URL | (no default) | `ws://localhost:8000/connection/websocket` |

---

# File Storage Configuration

LabraGo supports multiple file storage providers for handling file uploads. Currently, only the `local` provider is fully implemented.

## Storage Providers

### Local Storage (Default)

The `local` provider stores files on the local filesystem.

**Configuration:**
- `FILE_STORAGE_PROVIDER=local` (or omit for default)
- `FILE_STORAGE_PATH=./storage` (or omit for default)

**Requirements:**
- The storage directory must exist and be writable by the application
- Recommended permissions: `0755` for directory, `0644` for files
- Path can be relative (e.g., `./storage`) or absolute (e.g., `/var/lib/labrago/storage`)

**Example setup:**
```bash
mkdir -p storage
chmod 755 storage
```

### Future Providers

Support for the following providers is planned:
- **S3** - Amazon S3 storage
- **GCS** - Google Cloud Storage

When these providers are implemented, additional configuration will be required (credentials, bucket names, regions, etc.).

## Troubleshooting

### Permission Errors

If you encounter permission errors when uploading files:

1. Check directory permissions: `ls -ld storage`
2. Ensure the directory is writable: `chmod 755 storage`
3. Verify the application user has write access

### Path Issues

- Use absolute paths in production for better reliability
- Ensure relative paths are relative to where the application is executed
- The default `./storage` path is relative to the application's working directory

### Storage Directory Not Found

The application will attempt to create the storage directory if it doesn't exist, but it's recommended to create it manually with proper permissions.

---

# Development Workflow

## Code Generation

After modifying your Ent schemas or GraphQL definitions, regenerate code:

```bash
go generate ./...
```

This runs all `//go:generate` directives in your codebase, including:
- Ent schema generation
- GraphQL code generation

## Database Migrations

LabraGo uses Ent for database schema management. Schema changes are automatically applied when the application starts. The application will:

- Create new tables and columns
- Drop removed columns (if `WithDropColumn(true)` is enabled)
- Drop removed indexes (if `WithDropIndex(true)` is enabled)

**Important:** Always backup your database before running migrations in production.

## Local Development with LabraGo Source

If you're developing LabraGo itself or need to use a local version:

1. Update `go.mod` to point to your local LabraGo repository:
   ```go
   replace github.com/GoLabra/labra => /path/to/labra
   ```

2. Run `go mod tidy` to update dependencies

3. Regenerate code: `go generate ./...`

## Testing

Run tests:

```bash
go test ./...
```

---

# Contributing

We welcome contributions to LabraGo! Please check our [Contributing guide](CONTRIBUTING.md) for details on:

- Code style and standards
- How to submit pull requests
- Development setup
- Testing requirements

---

# Additional Resources

- [GraphQL Playground](http://localhost:4000/playground) - Test your GraphQL queries
- [Admin GraphQL Playground](http://localhost:4000/admin/playground) - Test admin GraphQL queries
- [Centrifugo Documentation](https://centrifugo.dev/docs/) - Real-time messaging setup
- [Ent Documentation](https://entgo.io/) - Entity framework documentation
- [gqlgen Documentation](https://gqlgen.com/) - GraphQL code generation
