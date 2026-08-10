# Decode Order Number API

A RESTful API built with Go and Gin framework to encode/encrypt and decode/decrypt order numbers using URL-encoding and DES encryption.

## Features

- **Decrypt Endpoint**: Decrypt encrypted order numbers using DES encryption
- **Encrypt Endpoint**: Encrypt plain order numbers using DES encryption
- **Health Check**: Monitor API health and availability
- **Clean Architecture**: Well-structured codebase with separation of concerns
- **Comprehensive Testing**: Unit tests for handlers and services
- **API Documentation**: Swagger/OpenAPI documentation
- **Middleware**: CORS, logging, and recovery middleware

## API Endpoints

### Health Check
- **GET** `/health` - Check API health status

### Decryption
- **POST** `/api/v1/decrypt` - Decrypt an encrypted order number

### Encryption
- **POST** `/api/v1/encrypt` - Encrypt a plain order number

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd decode-service
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## Configuration

Create a `configs/config.yaml` file:

```yaml
server:
  port: "8080"
  read_timeout: 30
  write_timeout: 30

security:
  secret_key: "secretkey"
  enable_cors: true
  max_request_size: 1048576  # 1MB
```

## API Usage

### Decrypt Request

```bash
curl -X POST http://localhost:8080/api/v1/decrypt \
  -H "Content-Type: application/json" \
  -d '{
    "encrypted_order": "qZCC%252FTvzJYUpNd9VllcaYfYKCkAKeOR6",
    "secret_key": "secretkey"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "decrypted_order": "A-121226-ABCDEF"
  },
  "message": "Decryption successful",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Encrypt Request

```bash
curl -X POST http://localhost:8080/api/v1/encrypt \
  -H "Content-Type: application/json" \
  -d '{
    "plain_order": "OA-121226-ABCDEF",
    "secret_key": "secretkey"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "encrypted_order": "qZCC%252FTvzJYUpNd9VllcaYfYKCkAKeOR6"
  },
  "message": "Encryption successful",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## API Documentation

Swagger/OpenAPI documentation is available at:
- Development: `http://localhost:8080/swagger.html`
- Production: `https://api.example.com/swagger.html`

## Project Structure

```
decode-service/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── handlers/
│   │   ├── decrypt_handler.go      # Decrypt endpoint handler
│   │   ├── encrypt_handler.go      # Encrypt endpoint handler
│   │   ├── health_handler.go       # Health check endpoint
│   │   └── base_handler.go        # Common handler functionality
│   ├── middleware/
│   │   └── middleware.go           # Middleware functions
│   ├── models/
│   │   └── request.go              # Request/response models
│   ├── services/
│   │   ├── crypto_service.go       # Core crypto logic
│   │   └── validation_service.go   # Input validation
│   └── utils/
├── pkg/
│   └── crypto/
│       └── des.go                  # Crypto utilities
├── api/
│   ├── v1/
│   │   └── handlers.go            # Route definitions
│   └── docs/
│       └── swagger.yaml            # API documentation
├── configs/
│   └── config.yaml                 # Application configuration
├── tests/
│   ├── handlers/
│   └── services/
└── go.mod
```

## Testing

Run unit tests:
```bash
go test -v ./...
```

Run specific test packages:
```bash
go test -v ./tests/handlers
go test -v ./tests/services
```

## Security Considerations

- All requests are validated before processing
- CORS is enabled for cross-origin requests
- Request size limits prevent abuse
- Error messages are sanitized to prevent information leakage
- Secret keys are handled securely

## Development

### Adding New Endpoints

1. Add new models in `internal/models/`
2. Create service logic in `internal/services/`
3. Implement handlers in `internal/handlers/`
4. Add routes in `api/v1/handlers.go`
5. Update Swagger documentation in `api/docs/swagger.yaml`

### Running in Development

```bash
# Run with hot reload (requires air)
air

# Run specific tests
go test -v ./tests/handlers -run TestDecryptHandler

# Run with coverage
go test -v -cover ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License.