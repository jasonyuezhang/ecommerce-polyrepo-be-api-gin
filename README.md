# E-Commerce API Gateway

A minimal API Gateway built with [Gin](https://github.com/gin-gonic/gin) that routes requests to backend microservices via gRPC.

## 🎯 About This Repository

This repository is part of the **ecommerce-polyrepo** project - a polyrepo setup designed for testing the [Propel](https://propel.us) code review feature across multiple microservices.

### Role in Microservices Architecture

The API Gateway serves as the **REST API entry point** in the microservices architecture:

```
┌─────────────┐
│  Frontend   │
│  (Next.js)  │
└──────┬──────┘
       │
       ▼
┌─────────────────┐     gRPC      ┌──────────────┐
│  API Gateway    │◄──────────────►│ User Service │
│  (Go/Gin)       │                │  (Django)    │
│  [THIS REPO]    │                └──────────────┘
└─────────────────┘
       │
       │ gRPC
       ├─────────────────►┌──────────────────┐
       │                   │ Listing Service  │
       │                   │  (Spring Boot)   │
       │                   └──────────────────┘
       │
       └─────────────────►┌──────────────────┐
                           │Inventory Service │
                           │    (Rails)       │
                           └──────────────────┘
```

### Quick Start (Standalone Testing)

To test this service independently:

```bash
# 1. Install dependencies
go mod download

# 2. Set up environment
cp .env.example .env
# Edit .env with mock service addresses or local services

# 3. Run the service
go run main.go

# 4. Test health endpoint
curl http://localhost:8080/health

# 5. Test API endpoints (requires backend services)
curl http://localhost:8080/api/v1/products
```

**Note:** For full functionality, backend microservices (user, listing, inventory) must be running. See the [parent polyrepo](https://github.com/jasonyuezhang/ecommerce-polyrepo) for orchestrated setup with all services.

---

## Architecture

This API Gateway serves as the single entry point for all client requests and routes them to the appropriate microservices:

- **User Service** - Authentication and user management
- **Listing Service** - Product catalog and listings
- **Inventory Service** - Stock management and availability

## Project Structure

```
be-api-gin/
├── cmd/
│   └── server/
│       └── main.go          # Server initialization
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── handlers/
│   │   ├── product.go       # Product handlers
│   │   └── order.go         # Order handlers
│   ├── middleware/
│   │   ├── auth.go          # JWT authentication
│   │   └── cors.go          # CORS middleware
│   ├── models/
│   │   └── models.go        # Common models
│   └── routes/
│       └── routes.go        # Route definitions
├── pkg/
│   └── grpc/
│       └── client.go        # gRPC client connections
├── main.go                  # Entry point
├── Dockerfile
├── go.mod
├── .env.example
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker (optional)

### Configuration

Copy the example environment file and configure:

```bash
cp .env.example .env
```

### Running Locally

```bash
go mod download
go run main.go
```

### Running with Docker

```bash
docker build -t ecommerce-api-gateway .
docker run -p 8080:8080 --env-file .env ecommerce-api-gateway
```

## API Endpoints

### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/products | List all products |
| GET | /api/v1/products/:id | Get product by ID |
| POST | /api/v1/products | Create product (auth required) |
| PUT | /api/v1/products/:id | Update product (auth required) |
| DELETE | /api/v1/products/:id | Delete product (auth required) |

### Orders

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/orders | List user orders (auth required) |
| GET | /api/v1/orders/:id | Get order by ID (auth required) |
| POST | /api/v1/orders | Create order (auth required) |
| PUT | /api/v1/orders/:id/status | Update order status (auth required) |
| DELETE | /api/v1/orders/:id | Cancel order (auth required) |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check |
| GET | /ready | Readiness check |

## Authentication

The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## GraphQL Gateway

For complex data aggregations, real-time features, and efficient data fetching, the platform also provides a GraphQL endpoint via `be-graphql-go`:

- **GraphQL HTTP**: `http://localhost:30900/graphql`
- **GraphQL WebSocket**: `ws://localhost:30901/graphql` (for subscriptions)

### When to Use GraphQL vs REST

- **Use REST (this API)** for:
  - Simple CRUD operations
  - Internal service-to-service communication
  - Traditional backend integrations

- **Use GraphQL** for:
  - Complex data aggregations (e.g., homepage with products + categories + recommendations)
  - Product detail pages with related data (inventory, similar products, reviews)
  - Real-time updates via subscriptions
  - Mobile applications requiring flexible queries
  - Frontend applications needing precise data fetching

Both APIs coexist and share the same backend microservices via gRPC.

## License

MIT License
