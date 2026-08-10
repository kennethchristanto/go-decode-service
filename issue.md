;# Backend API Development Plan

## Overview
Convert the existing Go CLI tool into a RESTful backend API with clean architecture principles, providing two main endpoints for decrypt and encrypt operations.

## Architecture Overview

### Clean Architecture Structure
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
│   │   └── health_handler.go       # Health check endpoint
│   ├── middleware/
│   │   ├── logging.go              # Request logging middleware
│   │   ├── recovery.go             # Panic recovery middleware
│   │   └── cors.go                 # CORS middleware
│   ├── models/
│   │   ├── request.go              # Request DTOs
│   │   └── response.go            # Response DTOs
│   ├── services/
│   │   ├── crypto_service.go       # Core encryption/decryption logic
│   │   └── validation_service.go   # Input validation service
│   └── utils/
│       └── helpers.go              # Utility functions
├── pkg/
│   └── crypto/
│       └── des.go                  # DES encryption utilities
├── api/
│   ├── v1/
│   │   └── handlers.go              # Route definitions
│   └── docs/
│       └── swagger.yaml             # API documentation
├── configs/
│   └── config.yaml                 # Application configuration
├── tests/
│   ├── handlers/
│   ├── services/
│   └── integration/
├── go.mod
├── go.sum
└── README.md
```

## API Endpoints

### 1. Decrypt Endpoint
**Endpoint:** `POST /api/v1/decrypt`  
**Description:** Decrypts an encrypted order number

**Request Body:**
```json
{
  "encrypted_order": "qZCC%252FTvzJYUpNd9VllcaYfYKCkAKeOR6",
  "secret_key": "secretkey"
}
```

**Response Body:**
```json
{
  "success": true,
  "data": {
    "decrypted_order": "A-121226-ABCDEF",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "message": "Decryption successful"
}
```

### 2. Encrypt Endpoint
**Endpoint:** `POST /api/v1/encrypt`  
**Description:** Encrypts a plain order number

**Request Body:**
```json
{
  "plain_order": "A-121226-ABCDEF",
  "secret_key": "secretkey"
}
```

**Response Body:**
```json
{
  "success": true,
  "data": {
    "encrypted_order": "qZCC%252FTvzJYUpNd9VllcaYfYKCkAKeOR6",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "message": "Encryption successful"
}
```

### 3. Health Check Endpoint
**Endpoint:** `GET /health`  
**Description:** Health check endpoint

**Response Body:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Implementation Details

### 1. Configuration Management (`internal/config/config.go`)
```go
type Config struct {
    Server struct {
        Port         string `yaml:"port"`
        ReadTimeout  int    `yaml:"read_timeout"`
        WriteTimeout int    `yaml:"write_timeout"`
    } `yaml:"server"`
    
    Security struct {
        SecretKey     string `yaml:"secret_key"`
        EnableCORS    bool   `yaml:"enable_cors"`
        MaxRequestSize int64 `yaml:"max_request_size"`
    } `yaml:"security"`
}
```

### 2. Request/Response Models (`internal/models/`)

**Request Models:**
```go
// internal/models/request.go
type DecryptRequest struct {
    EncryptedOrder string `json:"encrypted_order" validate:"required"`
    SecretKey     string `json:"secret_key" validate:"required"`
}

type EncryptRequest struct {
    PlainOrder string `json:"plain_order" validate:"required"`
    SecretKey  string `json:"secret_key" validate:"required"`
}
```

**Response Models:**
```go
// internal/models/response.go
type APIResponse struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Message   string      `json:"message"`
    Timestamp time.Time   `json:"timestamp"`
}

type DecryptResponse struct {
    DecryptedOrder string `json:"decrypted_order"`
}

type EncryptResponse struct {
    EncryptedOrder string `json:"encrypted_order"`
}
```

### 3. Service Layer (`internal/services/`)

**Crypto Service:**
```go
// internal/services/crypto_service.go
type CryptoService interface {
    Decrypt(encryptedOrder, secretKey string) (string, error)
    Encrypt(plainOrder, secretKey string) (string, error)
}

type cryptoService struct{}

func NewCryptoService() CryptoService {
    return &cryptoService{}
}

func (s *cryptoService) Decrypt(encryptedOrder, secretKey string) (string, error) {
    // Implement existing decryption logic
    // URL decode twice → Base64 decode → DES decrypt
}

func (s *cryptoService) Encrypt(plainOrder, secretKey string) (string, error) {
    // Implement existing encryption logic  
    // PKCS5 pad → DES encrypt → Base64 encode → URL encode twice
}
```

**Validation Service:**
```go
// internal/services/validation_service.go
type ValidationService interface {
    ValidateDecryptRequest(req DecryptRequest) error
    ValidateEncryptRequest(req EncryptRequest) error
}

type validationService struct{}

func NewValidationService() ValidationService {
    return &validationService{}
}
```

### 4. Handler Layer (`internal/handlers/`)

**Base Handler:**
```go
// internal/handlers/base_handler.go
type BaseHandler struct {
    cryptoService  services.CryptoService
    validationService services.ValidationService
}

func NewBaseHandler(cryptoService services.CryptoService, validationService services.ValidationService) *BaseHandler {
    return &BaseHandler{
        cryptoService: cryptoService,
        validationService: validationService,
    }
}
```

**Decrypt Handler:**
```go
// internal/handlers/decrypt_handler.go
type DecryptHandler struct {
    BaseHandler
}

func NewDecryptHandler(cryptoService services.CryptoService, validationService services.ValidationService) *DecryptHandler {
    return &DecryptHandler{
        BaseHandler: BaseHandler{
            cryptoService: cryptoService,
            validationService: validationService,
        },
    }
}

func (h *DecryptHandler) HandleDecrypt(c *gin.Context) {
    var req models.DecryptRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid request format",
            Timestamp: time.Now(),
        })
        return
    }
    
    if err := h.validationService.ValidateDecryptRequest(req); err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: err.Error(),
            Timestamp: time.Now(),
        })
        return
    }
    
    decrypted, err := h.cryptoService.Decrypt(req.EncryptedOrder, req.SecretKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Decryption failed: " + err.Error(),
            Timestamp: time.Now(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Data: models.DecryptResponse{
            DecryptedOrder: decrypted,
        },
        Message: "Decryption successful",
        Timestamp: time.Now(),
    })
}
```

**Encrypt Handler:**
```go
// internal/handlers/encrypt_handler.go
type EncryptHandler struct {
    BaseHandler
}

func NewEncryptHandler(cryptoService services.CryptoService, validationService services.ValidationService) *EncryptHandler {
    return &EncryptHandler{
        BaseHandler: BaseHandler{
            cryptoService: cryptoService,
            validationService: validationService,
        },
    }
}

func (h *EncryptHandler) HandleEncrypt(c *gin.Context) {
    var req models.EncryptRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Invalid request format",
            Timestamp: time.Now(),
        })
        return
    }
    
    if err := h.validationService.ValidateEncryptRequest(req); err != nil {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: err.Error(),
            Timestamp: time.Now(),
        })
        return
    }
    
    encrypted, err := h.cryptoService.Encrypt(req.PlainOrder, req.SecretKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Encryption failed: " + err.Error(),
            Timestamp: time.Now(),
        })
        return
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Data: models.EncryptResponse{
            EncryptedOrder: encrypted,
        },
        Message: "Encryption successful",
        Timestamp: time.Now(),
    })
}
```

### 5. Middleware (`internal/middleware/`)

**Logging Middleware:**
```go
// internal/middleware/logging.go
func LoggingMiddleware() gin.HandlerFunc {
    return gin.LoggerWithConfig(gin.LoggerConfig{
        SkipPaths: []string{"/health"},
        Formatter: func(param gin.LogFormatterParams) string {
            return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
                param.ClientIP,
                param.TimeStamp.Format(time.RFC1123),
                param.Method,
                param.Path,
                param.Request.Proto,
                param.StatusCode,
                param.Latency,
                param.Request.UserAgent(),
                param.ErrorMessage,
            )
        },
    })
}
```

**CORS Middleware:**
```go
// internal/middleware/cors.go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

### 6. API Routes (`api/v1/handlers.go`)
```go
// api/v1/handlers.go
package v1

import (
    "github.com/gin-gonic/gin"
    "decode-service/internal/handlers"
)

func SetupRoutes(
    decryptHandler *handlers.DecryptHandler,
    encryptHandler *handlers.EncryptHandler,
) *gin.Engine {
    r := gin.Default()
    
    // Apply middleware
    r.Use(handlers.LoggingMiddleware())
    r.Use(handlers.CORSMiddleware())
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, map[string]interface{}{
            "status": "healthy",
            "timestamp": time.Now(),
        })
    })
    
    // API v1 routes
    v1 := r.Group("/api/v1")
    {
        v1.POST("/decrypt", decryptHandler.HandleDecrypt)
        v1.POST("/encrypt", encryptHandler.HandleEncrypt)
    }
    
    return r
}
```

### 7. Main Application (`cmd/server/main.go`)
```go
// cmd/server/main.go
package main

import (
    "log"
    "decode-service/internal/config"
    "decode-service/internal/handlers"
    "decode-service/internal/services"
    "decode-service/api/v1"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("configs/config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Initialize services
    cryptoService := services.NewCryptoService()
    validationService := services.NewValidationService()
    
    // Initialize handlers
    decryptHandler := handlers.NewDecryptHandler(cryptoService, validationService)
    encryptHandler := handlers.NewEncryptHandler(cryptoService, validationService)
    
    // Setup routes
    router := v1.SetupRoutes(decryptHandler, encryptHandler)
    
    // Start server
    log.Printf("Server starting on port %s", cfg.Server.Port)
    if err := router.Run(":" + cfg.Server.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

## Implementation Steps

### Phase 1: Setup and Configuration
1. Create project structure following clean architecture
2. Implement configuration management system
3. Set up basic Gin router with middleware

### Phase 2: Service Layer
1. Extract existing crypto logic into `CryptoService`
2. Implement `ValidationService` for input validation
3. Add unit tests for services

### Phase 3: Handler Layer
1. Create request/response DTOs
2. Implement `DecryptHandler` and `EncryptHandler`
3. Add proper error handling and validation

### Phase 4: API Integration
1. Set up routing with proper middleware
2. Implement health check endpoint
3. Add CORS and logging middleware

### Phase 5: Testing and Documentation
1. Write unit tests for handlers and services
2. Create integration tests
3. Generate API documentation (Swagger/OpenAPI)
4. Add configuration file support

### Phase 6: Deployment
1. Create Docker configuration
2. Set up environment variables
3. Add monitoring and logging
4. Performance optimization

## Technology Stack
- **Framework:** Gin (HTTP web framework)
- **Validation:** Go-validator
- **Configuration:** Viper
- **Documentation:** Swagger/OpenAPI
- **Testing:** Go testing with testify
- **Containerization:** Docker

## Security Considerations
1. Input validation and sanitization
2. Rate limiting middleware
3. Secure headers (CORS, Content-Type)
4. Environment-based configuration
5. Error message sanitization

## Performance Considerations
1. Connection pooling for database (if needed)
2. Request size limits
3. Timeout configurations
4. Proper logging without performance impact
5. Caching strategies for frequent operations