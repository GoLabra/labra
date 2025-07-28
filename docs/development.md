# Development Guide

This guide covers the development workflow, best practices, and common tasks for working with the Labra project.

## 🚀 Getting Started

### Prerequisites

- **Go** (1.21 or higher)
- **Node.js** (18 or higher)
- **PostgreSQL** (14 or higher)
- **Git**

### Development Setup

1. **Clone and setup**
   ```bash
   git clone <repository-url>
   cd labra
   ```

2. **Backend setup**
   ```bash
   cd src/app
   cp .env.example .env
   # Edit .env with your database credentials
   go mod tidy
   ```

3. **Frontend setup**
   ```bash
   cd src/admin
   cp .env.example .env
   # Edit .env with your API URLs
   yarn install
   ```

4. **Database setup**
   ```bash
   createdb labra_db
   ```

## 🔄 Development Workflow

### Backend Development

#### **Code Generation**
The project uses code generation for database models and GraphQL schemas:

```bash
cd src/app
go generate ./...
```

This generates:
- Database client code
- GraphQL resolvers
- Type-safe database operations

#### **Database Schema Changes**
1. Modify schema files in `src/app/ent/schema/`
2. Run code generation: `go generate ./...`
3. Test the changes: `go run main.go start`

#### **Adding New Features**
1. **Define the data model** in Ent schema
2. **Generate code** with `go generate ./...`
3. **Implement business logic** in service layer
4. **Create GraphQL resolvers** for API endpoints
5. **Add tests** for new functionality

#### **Testing**
```bash
# Run all tests
go test ./...

# Run specific test
go test ./domain/svc

# Run with coverage
go test -cover ./...
```

### Frontend Development

#### **Component Development**
1. Create components in `src/admin/src/components/`
2. Add TypeScript types for GraphQL operations
3. Implement GraphQL queries/mutations
4. Add styling with Tailwind CSS

#### **Feature Development**
1. **Create feature module** in `src/admin/src/features/`
2. **Add GraphQL operations** for data fetching
3. **Build UI components** for the feature
4. **Add routing** if needed
5. **Test the feature** thoroughly

#### **Testing**
```bash
# Run tests
yarn test

# Run with coverage
yarn test --coverage

# Run specific test
yarn test --testNamePattern="UserManagement"
```

## 🏗️ Project Structure

### Backend Structure (`src/app/`)

```
src/app/
├── main.go                 # Application entry point
├── ent/                    # Database models
│   ├── schema/            # Schema definitions
│   └── client.go          # Generated client
├── domain/                # Business logic
│   ├── repo/              # Repository layer
│   ├── svc/               # Service layer
│   └── resolvers/         # GraphQL resolvers
├── generated/              # Generated GraphQL code
├── graphql/                # GraphQL schemas
└── interfaces/             # Interface definitions
```

### Frontend Structure (`src/admin/`)

```
src/admin/
├── src/
│   ├── app/               # Next.js app router
│   ├── components/        # Reusable components
│   ├── features/          # Feature modules
│   ├── lib/               # Utilities and config
│   └── types/             # TypeScript types
├── public/                # Static assets
└── package.json
```

## 🔧 Configuration

### Environment Variables

#### **Backend** (`src/app/.env`)
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

#### **Frontend** (`src/admin/.env`)
```env
NEXT_PUBLIC_API_URL=http://localhost:4001
NEXT_PUBLIC_GRAPHQL_URL=http://localhost:4001/query
NEXT_PUBLIC_ADMIN_GRAPHQL_URL=http://localhost:4001/admin
```

## 🛠️ Common Tasks

### Adding a New Entity

1. **Define the schema** in `src/app/ent/schema/`
   ```go
   // user.go
   package schema
   
   import (
       "entgo.io/ent"
       "entgo.io/ent/schema/field"
   )
   
   type User struct {
       ent.Schema
   }
   
   func (User) Fields() []ent.Field {
       return []ent.Field{
           field.String("email").Unique(),
           field.String("password"),
           field.String("first_name"),
           field.String("last_name"),
       }
   }
   ```

2. **Generate code**
   ```bash
   go generate ./...
   ```

3. **Implement repository** in `src/app/domain/repo/`
4. **Implement service** in `src/app/domain/svc/`
5. **Add GraphQL resolvers** in `src/app/domain/resolvers/`

### Adding a New API Endpoint

1. **Define GraphQL schema** in `src/app/graphql/`
2. **Implement resolver** in `src/app/domain/resolvers/`
3. **Add business logic** in service layer
4. **Test the endpoint**

### Adding Frontend Features

1. **Create feature module** in `src/admin/src/features/`
2. **Add GraphQL operations**
3. **Build UI components**
4. **Add routing** if needed
5. **Test the feature**

## 🔍 Debugging

### Backend Debugging

#### **Database Issues**
```bash
# Check database connection
psql -h localhost -U username -d labra_db

# View logs
tail -f logs/app.log
```

#### **GraphQL Issues**
- Use GraphQL Playground: `http://localhost:4001/playground`
- Check resolver logs
- Verify schema generation

#### **Code Generation Issues**
```bash
# Clean generated files
rm -rf src/app/generated/
rm -rf src/app/ent/

# Regenerate
go generate ./...
```

### Frontend Debugging

#### **Build Issues**
```bash
# Clear cache
rm -rf .next/
yarn dev
```

#### **GraphQL Issues**
- Use Apollo DevTools
- Check network tab for errors
- Verify environment variables

## 📝 Code Style

### Go Code Style

- Use `gofmt` for formatting
- Follow Go naming conventions
- Add comments for exported functions
- Use meaningful variable names

### TypeScript/React Code Style

- Use Prettier for formatting
- Follow ESLint rules
- Use TypeScript strict mode
- Add JSDoc comments for functions

## 🧪 Testing

### Backend Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./domain/svc -v
```

### Frontend Testing

```bash
# Run tests
yarn test

# Run with coverage
yarn test --coverage

# Run specific test
yarn test --testNamePattern="UserForm"
```

## 🚀 Deployment

### Development Deployment

1. **Backend**
   ```bash
   cd src/app
   go run main.go start
   ```

2. **Frontend**
   ```bash
   cd src/admin
   yarn dev
   ```

### Production Deployment

1. **Build backend**
   ```bash
   cd src/app
   go build -o main main.go
   ```

2. **Build frontend**
   ```bash
   cd src/admin
   yarn build
   ```

3. **Deploy with Docker**
   ```bash
   docker-compose up -d
   ```

## 🔗 Useful Commands

### Backend Commands
```bash
# Generate code
go generate ./...

# Run tests
go test ./...

# Build application
go build -o main main.go

# Run with hot reload (if using air)
air
```

### Frontend Commands
```bash
# Install dependencies
yarn install

# Start development server
yarn dev

# Build for production
yarn build

# Run tests
yarn test

# Lint code
yarn lint
```

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Next.js Documentation](https://nextjs.org/docs)
- [GraphQL Documentation](https://graphql.org/learn/)
- [Ent Documentation](https://entgo.io/docs/getting-started)

## 🔗 Related Documentation

- [API Reference](api-reference.md)
- [Architecture Overview](architecture.md)
- [Getting Started](../README.md) 