# What is LabraGo?

## It’s a headless CMS… but not *that* kind.

We’re not here to help you build a blog or a basic website. LabraGo was born to handle **real apps**—the complex, data-heavy, API-driven kind.

Built with Go for performance and React for a clean, modern admin UI, it’s fast, flexible, and actually enjoyable to work with.

Think of it as your backend brain—powerful enough to manage complex data models, events, and permissions, without drowning you in endless configurations.

<img height="96" src="https://github.com/GoLabra/labra/blob/develop/labra.jpeg" title="Labra Logo" width="96"/>

# Why does LabraGo exist?

## Because we exist 🙂

And because we’ve run into enough problems building data-heavy apps, we decided to do something about it.

This project is our way of making things simpler, cleaner, and more fun for anyone facing the same challenges.

---

# 📦 How to Install LabraGo

You can install and run LabraGo in two ways:

---

## ✅ The Easy Way — Using `labractl` CLI

```bash
go install github.com/GoLabra/labractl@v.1.0.0
```

Then run:

```bash
labractl create myproject
cd myproject
labractl start
```

This will:

- Clone the repo
- Patch `go.mod`
- Set up `.env` files for frontend and backend
- Run `go mod tidy` and `go generate`
- Ensure PostgreSQL user and database
- Install frontend dependencies
- Start frontend and backend concurrently

🔥 Done in seconds.

---

## 🛠 The Manual Way (Advanced)

See below:

---

# Running LabraGo Admin Panel

## Prerequisites
To run this project, ensure you have the following installed on your system:

- [Node.js](https://nodejs.org/) (version 18 or later is recommended)
- [Yarn](https://yarnpkg.com/) (version 1.22 or later)

## Installation
Follow these steps to set up the project on your local machine:

```bash
git clone https://github.com/GoLabra/labra.git
cd labra/src/admin
yarn install
```

## Running in Development Mode

```bash
yarn dev
```

Go to `http://localhost:3000` to view the app.

## Building the Project

```bash
yarn build
```

## Starting the Production Server

```bash
yarn start
```

## Set Environment Variables

```bash
cp .env.example .env
```

```env
NEXT_PUBLIC_BRAND_PRODUCT_NAME="Labra·GO"
NEXT_PUBLIC_BRAND_COLOR="blue"
NEXT_PUBLIC_GRAPHQL_API_URL="http://localhost:4000"
NEXT_PUBLIC_GRAPHQL_QUERY_API_URL="http://localhost:4000/query"
NEXT_PUBLIC_GRAPHQL_QUERY_SUBSCRIPTION_URL="ws://localhost:4000/query"
NEXT_PUBLIC_GRAPHQL_QUERY_PLAYGROUND_URL="http://localhost:4000/playground"
NEXT_PUBLIC_GRAPHQL_ENTITY_API_URL="http://localhost:4000/entity"
NEXT_PUBLIC_GRAPHQL_ENTITY_SUBSCRIPTION_URL="ws://localhost:4000/entity"
NEXT_PUBLIC_GRAPHQL_ENTITY_PLAYGROUND_URL="http://localhost:4000/eplayground"
NEXT_PUBLIC_CENTRIFUGO_URL="ws://localhost:8000/connection/websocket"
```

---

# Running LabraGO – Backend/API

## Prerequisites

- Go 1.18 or later
- PostgreSQL with an empty database

## Setup

```bash
cd labra/src/app
cp .env.example .env
```

Sample `.env`:

```bash
DSN=postgres://postgres:postgres@localhost:5432/labrago?sslmode=disable
DB_DIALECT=postgres
SERVER_PORT=4000
ENT_SCHEMA_PATH=./ent/schema
SECRET_KEY=XyZ7WpPqY2VW3m1O9QkH1fLj8zT6sJgKAsDfGhJ7K0I=
SUPER_ADMIN_EMAIL=admin@labrago.eu
CENTRIFUGO_API_ADDRESS=http://localhost:8000/api
CENTRIFUGO_API_KEY=m4Q9KvJNpY8Gh2LxU7sR0cVf3eZUw1PnKaYtXjBmOq0=
```

Generate a key:

```bash
head -c 32 /dev/urandom | base64
```

## Development

```bash
go mod tidy
go generate ./...
```

## Run the API

```bash
cd ../cli
go run main.go start
```

## Local API Dev Shortcut

In `/app/go.mod`:

```go
replace github.com/GoLabra/labra/src/api => ../api
```

or run:

```bash
sed -i "/REPLACE_LABRAGO_DEVELOPMENT_API/c replace github.com\/GoLabra\/labra\/src\/api => ../api" go.mod
```

---

## Contributing
Check our [Contributing guide](https://github.com/GoLabra/labra/blob/feature/labra-module/CONTRIBUTING.md)
