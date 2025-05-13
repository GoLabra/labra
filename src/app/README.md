# LabraGO – Backend/API

## Prerequisites

- Go (version 1.18 or later)
- PostgreSQL with an empty database

## Installation

1. Clone the repository
   ```bash
   git clone https://github.com/GoLabra/labra.git
   cd labra/src
   ```

2. Set Environment Variables

   Before running the app, you’ll need to configure your environment by copying the example file and filling in the values that match your local setup:
   ```bash
   cd app
   cp .env.example .env
   ```

   Here’s a sample .env. If you use these values, don’t forget to replace all secret keys first.

   ```bash
   DSN="host=localhost port=5432 user=ent dbname=ent password=123 sslmode=disable"
   DB_DIALECT=postgres
   SERVER_PORT=4000
   ENT_SCHEMA_PATH=./ent/schema
   SECRET_KEY=XyZ7WpPqY2VW3m1O9QkH1fLj8zT6sJgKAsDfGhJ7K0I=
   SUPER_ADMIN_EMAIL=admin@labrago.eu
   CENTRIFUGO_API_ADDRESS=http://localhost:8000/api
   CENTRIFUGO_API_KEY=m4Q9KvJNpY8Gh2LxU7sR0cVf3eZUw1PnKaYtXjBmOq0=
   ```

   To generate a strong secret key on Linux or macOS:
   ```bash
   head -c 32 /dev/urandom | base64
   ```


2. Development
Generate code and assets

```bash
cd app
go generate ./...
```

If you get errors, tidy first and re-generate:

```bash
go mod tidy
go generate ./...
```

## Running in Development Mode
The entrypoint lives in the cli subdirectory:

```bash
cd ../cli
go run main.go start
```

You should see startup logs indicating successful DB connection and GraphQL endpoint availability.

## API Development

 In `/app/go.mod`, replace the placeholder with your local path:

   ```go
   replace github.com/GoLabra/labra/src/api => ../api
   ```

   or run the command bellow

   ```bash
   cd app
   ```

   ```bash
   sed -i "/REPLACE_LABRAGO_DEVELOPMENT_API/c replace github.com\/GoLabra\/labra\/src\/api => ../api" go.mod
   ```

## Contributing
Check our [Contributing guide](https://github.com/GoLabra/labra/blob/feature/labra-module/CONTRIBUTING.md)
