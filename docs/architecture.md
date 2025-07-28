# Architecture Overview

This document provides an overview of the Labra application architecture, explaining the key components and how they work together.

## 🏗️ System Architecture

Labra is built as a modern web application with a clear separation between frontend and backend components:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   Database      │
│   (Next.js)     │◄──►│    (Go)         │◄──►│  (PostgreSQL)   │
│                 │    │                 │    │                 │
│ - Admin UI      │    │ - GraphQL API   │    │ - User data     │
│ - React         │    │ - Authentication │    │ - File storage  │
│ - TypeScript    │    │ - File handling  │    │ - Permissions   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📁 Project Structure

```
labra/
├── src/
│   ├── app/                    # Go backend application
│   │   ├── main.go            # Application entry point
│   │   ├── ent/               # Database models (Ent)
│   │   ├── domain/            # Business logic
│   │   │   ├── repo/          # Repository layer
│   │   │   ├── svc/           # Service layer
│   │   │   └── resolvers/     # GraphQL resolvers
│   │   ├── generated/         # Generated GraphQL code
│   │   └── graphql/           # GraphQL schemas
│   │
│   ├── admin/                  # Next.js frontend
│   │   ├── src/
│   │   │   ├── app/           # Next.js app router
│   │   │   ├── components/    # React components
│   │   │   ├── features/      # Feature modules
│   │   │   └── lib/           # Utilities and config
│   │   └── public/            # Static assets
│   │
│   └── api/                    # Core API package
│       ├── entgql/            # GraphQL generation
│       ├── config/            # Configuration
│       ├── handler/           # HTTP handlers
│       └── utils/             # Shared utilities
│
├── docs/                       # Documentation
└── README.md
```

## 🔧 Key Components

### Backend (Go)

#### **Main Application (`src/app/main.go`)**
- Application entry point
- Database connection setup
- GraphQL server configuration
- HTTP routing and middleware

#### **Database Layer (`src/app/ent/`)**
- **Ent ORM**: Type-safe database operations
- **Schema**: Database model definitions
- **Migrations**: Automatic database schema management
- **Generated Code**: Type-safe database clients

#### **Domain Layer**
- **Repository Pattern**: Data access abstraction
- **Service Layer**: Business logic implementation
- **Resolvers**: GraphQL query/mutation handlers

#### **GraphQL API**
- **Schema**: Type definitions and operations
- **Resolvers**: Query and mutation implementations
- **Authentication**: JWT-based security
- **Subscriptions**: Real-time updates

### Frontend (Next.js)

#### **Admin Interface (`src/admin/`)**
- **React Components**: Reusable UI components
- **Feature Modules**: Organized by functionality
- **State Management**: Client-side state handling
- **API Integration**: GraphQL client setup

#### **Key Features**
- **Dynamic Forms**: Auto-generated from schemas
- **Content Management**: CRUD operations
- **User Management**: Authentication and authorization
- **File Management**: Upload and storage

## 🔄 Data Flow

### 1. **Authentication Flow**
```
User Login → JWT Token → API Requests → Authorization
```

### 2. **Data Operations**
```
Frontend Request → GraphQL Resolver → Service Layer → Repository → Database
```

### 3. **File Upload**
```
File Upload → Backend Handler → File Storage → Database Record
```

## 🗄️ Database Schema

### Core Entities

#### **Users**
- Authentication and user management
- Role-based access control
- Profile information

#### **Roles**
- Permission groups
- Access control definitions

#### **Permissions**
- Granular access control
- Resource-specific permissions

#### **Files**
- File storage metadata
- Upload and download handling

## 🔐 Security Architecture

### **Authentication**
- JWT-based token authentication
- Secure password hashing
- Session management

### **Authorization**
- Role-based access control (RBAC)
- Permission-based authorization
- Resource-level security

### **Data Protection**
- Input validation and sanitization
- SQL injection prevention
- XSS protection

## 🚀 Performance Features

### **Backend**
- Connection pooling
- Query optimization
- Caching strategies
- Efficient database queries

### **Frontend**
- Code splitting
- Lazy loading
- Optimized builds
- Responsive design

## 🔄 Code Generation

The project uses extensive code generation to maintain consistency and reduce boilerplate:

### **Database Layer**
- Ent generates type-safe database clients
- Automatic migration generation
- Schema validation

### **GraphQL Layer**
- Schema-first development
- Type-safe resolvers
- Automatic client generation

### **Frontend**
- TypeScript types from GraphQL schema
- Component generation
- Form generation

## 🛠️ Development Workflow

### **Backend Development**
1. Define database schema in Ent
2. Generate database code
3. Implement business logic in services
4. Create GraphQL resolvers
5. Test API endpoints

### **Frontend Development**
1. Design UI components
2. Implement GraphQL queries
3. Build feature modules
4. Test user interactions

## 📊 Monitoring and Logging

### **Application Logs**
- Request/response logging
- Error tracking
- Performance metrics

### **Database Monitoring**
- Query performance
- Connection health
- Schema changes

## 🔗 External Dependencies

### **Backend Dependencies**
- **PostgreSQL**: Primary database
- **Ent**: ORM and code generation
- **GraphQL**: API layer
- **JWT**: Authentication

### **Frontend Dependencies**
- **Next.js**: React framework
- **Apollo Client**: GraphQL client
- **TypeScript**: Type safety
- **Tailwind CSS**: Styling

## 🔄 Deployment Architecture

### **Development**
- Local PostgreSQL database
- Hot reloading for both frontend and backend
- Development-specific configurations

### **Production**
- Containerized deployment
- Database clustering
- Load balancing
- CDN for static assets

## 📈 Scalability Considerations

### **Horizontal Scaling**
- Stateless backend design
- Database connection pooling
- Load balancer support

### **Performance Optimization**
- Query optimization
- Caching strategies
- Asset optimization

## 🔗 Related Documentation

- [API Reference](api-reference.md)
- [Development Guide](development.md)
- [Getting Started](../README.md) 