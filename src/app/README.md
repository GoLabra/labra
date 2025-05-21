# Why does LabraGo exist?

## Because we exist 🙂

And because we’ve run into enough problems building data-heavy apps, we decided to do something about it.

This project is our way of making things simpler, cleaner, and more fun for anyone facing the same challenges.

---

# 📦 Installation Options

You can install and run LabraGo in two ways:

---

## 1. The Easy Way — Using `labractl` CLI 🚀

The official CLI tool lets you skip the manual process entirely.

### ✅ Install it with:

```bash
go install github.com/GoLabra/labractl@v.1.0.0
```

Make sure `$GOPATH/bin` is in your `$PATH`.

### ✅ Create and start a project:

```bash
labractl create myproject
cd myproject
labractl start
```

This will:

- Clone LabraGo and configure everything
- Create `.env` files for backend and frontend
- Ensure PostgreSQL user and database exist
- Install dependencies with Yarn
- Start both backend and frontend in parallel

---

## 2. Manual Setup (Advanced)

If you want full control over how things are configured and started, follow the steps below:

- [Running LabraGo Admin Panel](#running-labrago-admin-panel)
- [Running LabraGO – Backend/API](#running-labrago--backendapi)

---

## 🖼️ Running LabraGo Admin Panel

1. Make sure you have **Yarn** and **Node.js** installed.
2. Navigate to the admin panel directory:

```bash
cd src/admin
```

3. Install dependencies:

```bash
yarn install
```

4. Create a `.env` file based on `.env.example`, or directly:

```env
NEXT_PUBLIC_BRAND_PRODUCT_NAME="Labra·GO"
NEXT_PUBLIC_BRAND_COLOR="blue"
NEXT_PUBLIC_GRAPHQL_API_URL="http://localhost:4001"
NEXT_PUBLIC_GRAPHQL_QUERY_API_URL="http://localhost:4001/query"
NEXT_PUBLIC_GRAPHQL_QUERY_SUBSCRIPTION_URL="ws://localhost:4001/query"
NEXT_PUBLIC_GRAPHQL_QUERY_PLAYGROUND_URL="http://localhost:4001/playground"
NEXT_PUBLIC_GRAPHQL_ENTITY_API_URL="http://localhost:4001/entity"
NEXT_PUBLIC_GRAPHQL_ENTITY_PLAYGROUND_URL="http://localhost:4001/eplayground"
```

5. Start the frontend:

```bash
yarn dev
```

---

## ⚙️ Running LabraGO – Backend/API

1. Make sure **Go 1.20+** and **PostgreSQL** are installed.
2. Navigate to the backend directory:

```bash
cd src/app
```

3. Create a `.env` file like this:

```env
SERVER_PORT=4001
SECRET_KEY=supersecretdevkey
DSN=postgres://postgres:postgres@localhost:5432/<your-db>?sslmode=disable
DB_DIALECT=postgres
ENT_SCHEMA_PATH=absolute/path/to/src/app/ent/schema
CENTRIFUGO_API_ADDRESS=http://localhost:8000
CENTRIFUGO_API_KEY=secretkey
```

4. Run migrations if needed (if you're using Atlas or manual setup).

5. Run `go mod tidy`:

```bash
go mod tidy
```

6. Run `go generate`:

```bash
go generate ./...
```

7. Start the backend:

```bash
go run main.go start
```

---

## 🛠 Development

If you want to build or test the CLI locally:

```bash
go build -o labractl main.go
./labractl help
```

---

## 🤝 Contributing

PRs welcome. Open an issue or fork away if you want to improve the CLI.

---

## License

MIT © 2025 GoLabra
