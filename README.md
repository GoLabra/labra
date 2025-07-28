# Labra - Go Backend with Next.js Frontend

Labra is a modern web application built with Go backend and Next.js frontend, featuring GraphQL API, PostgreSQL database, and a comprehensive admin interface.

## 🚀 Quick Start

### Prerequisites

- **Go** (1.21 or higher)
- **Node.js** (18 or higher)
- **PostgreSQL** (14 or higher)
- **Git**

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd labra
   ```

2. **Set up the database**
   ```bash
   # Start PostgreSQL (if not already running)
   brew services start postgresql
   
   # Create database
   createdb labra_db
   ```

3. **Configure environment variables**
   ```bash
   # Backend environment
   cp src/app/.env.example src/app/.env
   # Edit src/app/.env with your database credentials
   
   # Frontend environment
   cp src/admin/.env.example src/admin/.env
   # Edit src/admin/.env with your API URLs
   ```

4. **Start the backend**
   ```bash
   cd src/app
   go mod tidy
   go run main.go start
   ```

5. **Start the frontend**
   ```bash
   cd src/admin
   yarn install
   yarn dev
   ```

## 🌐 Available Endpoints

- **Frontend Application**: http://localhost:3000
- **Backend API**: http://localhost:4001/query
- **GraphQL Playground**: http://localhost:4001/playground
- **Admin Playground**: http://localhost:4001/aplayground

## 📁 Project Structure

```
labra/
├── src/
│   ├── app/           # Go backend application
│   ├── admin/         # Next.js frontend application
│   └── api/           # Core API package
├── docs/              # Documentation
└── README.md
```

## 🔧 Configuration

### Backend Environment Variables

Create `src/app/.env`:
```env
DSN=postgres://username:password@localhost:5432/labra_db?sslmode=disable
DB_DIALECT=postgres
SERVER_PORT=4001
ENT_SCHEMA_PATH=./ent/schema
SECRET_KEY=your-secret-key
CENTRIFUGO_API_ADDRESS=http://localhost:8000
CENTRIFUGO_API_KEY=your-centrifugo-key
FILE_STORAGE_PROVIDER=local
FILE_STORAGE_PATH=./uploads
```

### Frontend Environment Variables

Create `src/admin/.env`:
```env
NEXT_PUBLIC_API_URL=http://localhost:4001
NEXT_PUBLIC_GRAPHQL_URL=http://localhost:4001/query
NEXT_PUBLIC_ADMIN_GRAPHQL_URL=http://localhost:4001/admin
```

## 🛠️ Development

### Code Generation

The project uses code generation for GraphQL schemas and database models:

```bash
cd src/app
go generate ./...
```

### Database Migrations

Database migrations are handled automatically by Ent:

```bash
cd src/app
go run main.go start
```

## 🐛 Troubleshooting

### Common Issues

1. **Database connection errors**
   - Ensure PostgreSQL is running
   - Check database credentials in `.env` file
   - Verify database exists

2. **Port conflicts**
   - Change `SERVER_PORT` in backend .env file
   - Change `PORT` in frontend .env file
   - Ensure your chosen ports are available

3. **Code generation errors**
   - Run `go mod tidy` to ensure dependencies are correct
   - Check that all required Go packages are installed
   - Verify the `entc.go` file has correct import paths

4. **Frontend build errors**
   - Run `yarn install` to install dependencies
   - Check Node.js version (requires 18+)
   - Verify environment variables are set correctly

## 📚 Documentation

- [API Reference](docs/api-reference.md)
- [Architecture Overview](docs/architecture.md)
- [Development Guide](docs/development.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- [GitHub Issues](https://github.com/GoLabra/labra/issues)
- [Discussions](https://github.com/GoLabra/labra/discussions)
