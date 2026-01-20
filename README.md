# E-Commerce API Gateway

A minimal API Gateway built with [Gin](https://github.com/gin-gonic/gin) that routes requests to backend microservices via gRPC.

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

## License

MIT License
