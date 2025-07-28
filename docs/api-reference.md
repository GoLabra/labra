# API Reference

This document provides an overview of the GraphQL API endpoints and operations available in Labra.

## 🌐 Base URLs

- **Main API**: `http://localhost:4001/query`
- **Admin API**: `http://localhost:4001/admin`
- **GraphQL Playground**: `http://localhost:4001/playground`
- **Admin Playground**: `http://localhost:4001/aplayground`

## 🔐 Authentication

The API uses JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## 📝 Authentication Endpoints

### Login

```graphql
mutation Login($email: String!, $password: String!) {
  login(email: $email, password: $password) {
    token
    user {
      id
      email
      firstName
      lastName
    }
  }
}
```

### Signup

```graphql
mutation Signup($input: CreateUserInput!) {
  signup(input: $input) {
    token
    user {
      id
      email
      firstName
      lastName
    }
  }
}
```

## 👥 User Management

### Get Current User

```graphql
query GetCurrentUser {
  me {
    id
    email
    firstName
    lastName
    roles {
      id
      name
    }
  }
}
```

### Get Users

```graphql
query GetUsers($first: Int, $after: String) {
  users(first: $first, after: $after) {
    edges {
      node {
        id
        email
        firstName
        lastName
        createdAt
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

### Create User

```graphql
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
    email
    firstName
    lastName
  }
}
```

### Update User

```graphql
mutation UpdateUser($id: ID!, $input: UpdateUserInput!) {
  updateUser(id: $id, input: $input) {
    id
    email
    firstName
    lastName
  }
}
```

### Delete User

```graphql
mutation DeleteUser($id: ID!) {
  deleteUser(id: $id) {
    id
  }
}
```

## 🎭 Role Management

### Get Roles

```graphql
query GetRoles {
  roles {
    id
    name
    createdAt
  }
}
```

### Create Role

```graphql
mutation CreateRole($input: CreateRoleInput!) {
  createRole(input: $input) {
    id
    name
  }
}
```

### Update Role

```graphql
mutation UpdateRole($id: ID!, $input: UpdateRoleInput!) {
  updateRole(id: $id, input: $input) {
    id
    name
  }
}
```

### Delete Role

```graphql
mutation DeleteRole($id: ID!) {
  deleteRole(id: $id) {
    id
  }
}
```

## 📁 File Management

### Get Files

```graphql
query GetFiles($first: Int, $after: String) {
  files(first: $first, after: $after) {
    edges {
      node {
        id
        name
        path
        size
        mimeType
        createdAt
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

### Upload File

```graphql
mutation UploadFile($file: Upload!) {
  uploadFile(file: $file) {
    id
    name
    path
    size
    mimeType
  }
}
```

### Delete File

```graphql
mutation DeleteFile($id: ID!) {
  deleteFile(id: $id) {
    id
  }
}
```

## 🔐 Permission Management

### Get Permissions

```graphql
query GetPermissions {
  permissions {
    id
    name
    description
    createdAt
  }
}
```

### Create Permission

```graphql
mutation CreatePermission($input: CreatePermissionInput!) {
  createPermission(input: $input) {
    id
    name
    description
  }
}
```

### Update Permission

```graphql
mutation UpdatePermission($id: ID!, $input: UpdatePermissionInput!) {
  updatePermission(id: $id, input: $input) {
    id
    name
    description
  }
}
```

### Delete Permission

```graphql
mutation DeletePermission($id: ID!) {
  deletePermission(id: $id) {
    id
  }
}
```

## 📊 Pagination

The API uses cursor-based pagination. All list queries support:

- `first`: Number of items to return
- `after`: Cursor for pagination
- `before`: Cursor for reverse pagination
- `last`: Number of items to return in reverse

## 🔍 Filtering

Many queries support filtering using `where` arguments:

```graphql
query GetUsers($where: UserWhereInput) {
  users(where: $where) {
    edges {
      node {
        id
        email
        firstName
        lastName
      }
    }
  }
}
```

## 📝 Sorting

Queries support sorting using `orderBy` arguments:

```graphql
query GetUsers($orderBy: [UserOrderByInput!]) {
  users(orderBy: $orderBy) {
    edges {
      node {
        id
        email
        firstName
        lastName
        createdAt
      }
    }
  }
}
```

## 🔄 Subscriptions

The API supports GraphQL subscriptions for real-time updates:

```graphql
subscription OnUserCreated {
  userCreated {
    id
    email
    firstName
    lastName
  }
}
```

## 🛠️ Error Handling

The API returns structured errors with:

- `message`: Human-readable error message
- `code`: Error code for programmatic handling
- `path`: Field path where the error occurred

## 📚 Schema Introspection

You can explore the full schema using introspection:

```graphql
query IntrospectSchema {
  __schema {
    types {
      name
      description
      fields {
        name
        type {
          name
        }
      }
    }
  }
}
```

## 🔗 Related Documentation

- [Architecture Overview](architecture.md)
- [Development Guide](development.md)
- [Getting Started](../README.md) 