# Sequence Diagrams

This document contains sequence diagrams illustrating the key interaction flows within the `card-manage` application. These diagrams are written in Mermaid syntax and help visualize how different components collaborate to fulfill a user request.

## 1. User Registration Flow

This diagram shows the process for a new user registering.

```mermaid
sequenceDiagram
    participant Client
    participant GinRouter as Gin Router
    participant UserHandler as User Handler
    participant UserService as User Service
    participant UserRepo as User Repository
    participant Database

    Client->>GinRouter: POST /register (email, password, role)
    GinRouter->>UserHandler: Register(c)
    UserHandler->>UserHandler: Bind JSON to RegisterRequest struct
    UserHandler->>UserService: Register(email, password, role)
    UserService->>UserRepo: GetUserByEmail(email)
    UserRepo->>Database: SELECT * FROM users WHERE email = ?
    Database-->>UserRepo: (user not found)
    UserRepo-->>UserService: nil, ErrNotFound
    UserService->>UserService: Hash password with bcrypt
    UserService->>UserRepo: CreateUser(user)
    UserRepo->>Database: INSERT INTO users (...)
    Database-->>UserRepo: (new user ID)
    UserRepo-->>UserService: (user object), nil
    UserService-->>UserHandler: (user object), nil
    UserHandler-->>GinRouter: 201 Created response
    GinRouter-->>Client: HTTP 201
```

## 2. User Login Flow

This diagram illustrates how a user logs in and receives a JWT.

```mermaid
sequenceDiagram
    participant Client
    participant GinRouter as Gin Router
    participant UserHandler as User Handler
    participant UserService as User Service
    participant UserRepo as User Repository
    participant JWTService as JWT Service
    participant Database

    Client->>GinRouter: POST /login (email, password)
    GinRouter->>UserHandler: Login(c)
    UserHandler->>UserHandler: Bind JSON to LoginRequest struct
    UserHandler->>UserService: Login(email, password)
    UserService->>UserRepo: GetUserByEmail(email)
    UserRepo->>Database: SELECT * FROM users WHERE email = ?
    Database-->>UserRepo: (user record with hashed password)
    UserRepo-->>UserService: (user object), nil
    UserService->>UserService: Compare provided password with hash
    alt Passwords Match
        UserService-->>UserHandler: (user object), nil
        UserHandler->>JWTService: GenerateToken(user)
        JWTService-->>UserHandler: (JWT string), nil
        UserHandler-->>GinRouter: 200 OK response with token
        GinRouter-->>Client: HTTP 200 ({"token": "..."})
    else Passwords Do Not Match
        UserService-->>UserHandler: nil, ErrInvalidCredentials
        UserHandler-->>GinRouter: 401 Unauthorized response
        GinRouter-->>Client: HTTP 401
    end
```

## 3. Create Consignment Flow (Protected Route)

This diagram shows a more complex flow for a protected endpoint. It includes authentication and role-based authorization middleware.

```mermaid
sequenceDiagram
    participant Client
    participant GinRouter as Gin Router
    participant AuthMiddleware as Auth Middleware
    participant RoleMiddleware as Role Middleware
    participant ConsignmentHandler as Consignment Handler
    participant ConsignmentService as Consignment Service
    participant ConsignmentRepo as Consignment Repository
    participant Database

    Client->>GinRouter: POST /api/consignments (Authorization: Bearer <token>)
    GinRouter->>AuthMiddleware: Handle request
    Note over AuthMiddleware: Validates JWT and extracts claims (userID, role)
    AuthMiddleware->>AuthMiddleware: Set claims into Gin context
    AuthMiddleware-->>GinRouter: c.Next()

    GinRouter->>RoleMiddleware: Handle request (requires 'PLAYER')
    Note over RoleMiddleware: Checks role from context claims
    RoleMiddleware-->>GinRouter: c.Next()

    GinRouter->>ConsignmentHandler: CreateConsignment(c)
    ConsignmentHandler->>ConsignmentHandler: Bind JSON to request struct
    ConsignmentHandler->>ConsignmentService: CreateConsignment(...)
    Note over ConsignmentService: Business logic (e.g., validate card, etc.)
    ConsignmentService->>ConsignmentRepo: CreateConsignment(consignment)
    ConsignmentRepo->>Database: INSERT INTO consignments (...)
    Database-->>ConsignmentRepo: (new consignment ID)
    ConsignmentRepo-->>ConsignmentService: (consignment object), nil
    ConsignmentService-->>ConsignmentHandler: (consignment object), nil
    ConsignmentHandler-->>GinRouter: 201 Created response
    GinRouter-->>Client: HTTP 201
```
