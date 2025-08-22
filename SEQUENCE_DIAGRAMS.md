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

## 3. Create Consignment Request Flow

This diagram shows a player creating a new consignment request with multiple items.

```mermaid
sequenceDiagram
    participant Client as Player
    participant GinRouter as Gin Router
    participant AuthMiddleware
    participant RoleMiddleware
    participant ConsignmentHandler
    participant ConsignmentService
    participant ConsignmentRepo
    participant Database

    Client->>GinRouter: POST /api/consignments (Authorization: Bearer <token>, body: {store_id, card_ids[]})
    GinRouter->>AuthMiddleware: Handle request
    AuthMiddleware-->>GinRouter: c.Next() (JWT valid, claims set)

    GinRouter->>RoleMiddleware: Handle request (requires 'PLAYER')
    RoleMiddleware-->>GinRouter: c.Next() (Role OK)

    GinRouter->>ConsignmentHandler: CreateConsignment(c)
    ConsignmentHandler->>ConsignmentHandler: Bind JSON to CreateConsignmentRequest
    ConsignmentHandler->>ConsignmentService: CreateConsignment(playerID, storeID, cardIDs)
    
    ConsignmentService->>ConsignmentRepo: CreateConsignment(consignment, items)
    Note over ConsignmentRepo: Begins DB Transaction
    ConsignmentRepo->>Database: INSERT INTO consignments (...)
    Database-->>ConsignmentRepo: (new consignment_id)
    loop for each card_id
        ConsignmentRepo->>Database: INSERT INTO consignment_items (consignment_id, card_id, ...)
    end
    Note over ConsignmentRepo: Commits DB Transaction
    ConsignmentRepo-->>ConsignmentService: (consignment object), nil
    
    ConsignmentService-->>ConsignmentHandler: (full consignment with items), nil
    ConsignmentHandler-->>GinRouter: 201 Created response
    GinRouter-->>Client: HTTP 201
```

## 4. Update Consignment Item Status Flow

This diagram shows a store owner approving or rejecting an individual consignment item.

```mermaid
sequenceDiagram
    participant Client as Store Owner
    participant GinRouter as Gin Router
    participant AuthMiddleware
    participant RoleMiddleware
    participant ConsignmentHandler
    participant ConsignmentService
    participant ConsignmentRepo
    participant Database

    Client->>GinRouter: PUT /api/consignments/items/{itemId} (Authorization: Bearer <token>, body: {status, reason})
    GinRouter->>AuthMiddleware: Handle request
    AuthMiddleware-->>GinRouter: c.Next() (JWT valid, claims set)

    GinRouter->>RoleMiddleware: Handle request (requires 'STORE')
    RoleMiddleware-->>GinRouter: c.Next() (Role OK)

    GinRouter->>ConsignmentHandler: UpdateConsignmentItemStatus(c)
    ConsignmentHandler->>ConsignmentHandler: Bind JSON to UpdateStatusRequest
    ConsignmentHandler->>ConsignmentService: UpdateConsignmentItemStatus(storeUserID, itemID, newStatus, reason)
    
    ConsignmentService->>ConsignmentRepo: GetConsignmentItemByID(itemID)
    ConsignmentRepo-->>ConsignmentService: (item object)
    
    ConsignmentService->>ConsignmentRepo: GetConsignmentByID(item.ConsignmentID)
    ConsignmentRepo-->>ConsignmentService: (consignment object with storeID)

    Note over ConsignmentService: Verifies storeUserID owns the consignment.storeID

    ConsignmentService->>ConsignmentRepo: UpdateConsignmentItemStatus(itemID, newStatus, reason)
    ConsignmentRepo->>Database: UPDATE consignment_items SET status = ?, ... WHERE id = ?
    Database-->>ConsignmentRepo: (success)
    ConsignmentRepo-->>ConsignmentService: nil
    
    ConsignmentService-->>ConsignmentHandler: (updated item object), nil
    ConsignmentHandler-->>GinRouter: 200 OK response
    GinRouter-->>Client: HTTP 200
```
