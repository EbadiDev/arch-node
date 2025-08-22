# Architecture Overview

## System Architecture

Arch-Node is built with a modular, layered architecture designed for high performance and extensibility. The system consists of several key components that work together to provide proxy services and management capabilities.

## Core Components

### 1. Application Layer (`internal/app/`)

The main application orchestrator that initializes and coordinates all components:

- **Lifecycle Management**: Handles startup, shutdown, and graceful termination
- **Component Initialization**: Sets up all services in the correct order
- **Signal Handling**: Manages OS signals for graceful shutdown
- **Context Management**: Provides cancellation context for all components

### 2. HTTP Server (`internal/http/`)

RESTful API server for management and monitoring:

- **Echo Framework**: High-performance HTTP router
- **Middleware Stack**: Authentication, logging, CORS, and request validation
- **API Endpoints**: V1 API for configuration, statistics, and manager operations
- **Security**: Token-based authentication for all protected endpoints

### 3. Configuration Management (`internal/config/`)

Handles application configuration with layered approach:

- **Default Configuration**: Base settings from `configs/main.defaults.json`
- **Override Configuration**: Custom settings from `configs/main.json`
- **Environment Variables**: Runtime configuration overrides
- **Validation**: Comprehensive validation using struct tags

### 4. Database System (`internal/database/`)

JSON file-based persistence for runtime data:

- **Settings Storage**: HTTP port, authentication tokens
- **Manager Configuration**: Connection details for Arch-Manager
- **File-based**: Simple JSON files for easy debugging and backup
- **Thread-safe**: Mutex-based concurrency control

### 5. Coordinator (`internal/coordinator/`)

Manages synchronization with Arch-Manager:

- **Background Sync**: Periodic configuration fetching from manager
- **Configuration Updates**: Applies remote configuration changes
- **Worker Pool**: Background task processing
- **Error Handling**: Robust error handling and retry mechanisms

### 6. Xray Integration (`pkg/xray/`)

Core proxy functionality through Xray-core:

- **Process Management**: Manages Xray binary lifecycle
- **Configuration Management**: Dynamic configuration updates
- **API Communication**: gRPC connection to Xray API
- **Statistics Collection**: Traffic and performance metrics

### 7. HTTP Client (`pkg/http/client/`)

HTTP client for external communications:

- **Manager Communication**: Requests to Arch-Manager
- **Authentication**: Automatic token handling
- **Timeouts**: Configurable request timeouts
- **Error Handling**: Comprehensive error handling and retry logic

### 8. Logger (`pkg/logger/`)

Structured logging system:

- **Zap Integration**: High-performance structured logging
- **Log Levels**: Configurable log levels (debug, info, warn, error)
- **File Output**: Separate log files for different components
- **JSON Format**: Machine-readable log format

## Data Flow

### 1. Startup Flow

```
main.go → cmd/start.go → app.New() → Component Initialization
```

1. **Configuration Loading**: Load and validate configuration files
2. **Database Initialization**: Load or create database files
3. **Xray Setup**: Initialize Xray configuration and start process
4. **HTTP Server**: Start REST API server
5. **Coordinator**: Begin background synchronization

### 2. Configuration Sync Flow

```
Coordinator → HTTP Client → Arch-Manager → Configuration Update → Xray Restart
```

1. **Periodic Sync**: Coordinator runs every 30 seconds
2. **Fetch Configuration**: Request configuration from manager
3. **Compare**: Check if configuration differs from current
4. **Update**: Apply new configuration and restart Xray if needed

### 3. API Request Flow

```
Client → HTTP Server → Middleware → Handler → Database/Xray → Response
```

1. **Authentication**: Verify Bearer token
2. **Validation**: Validate request format and data
3. **Processing**: Execute business logic
4. **Response**: Return JSON response

## Security Architecture

### Authentication

- **Token-based**: Bearer token authentication for API access
- **Random Tokens**: Cryptographically secure random token generation
- **Per-node Tokens**: Each node has unique authentication token

### Network Security

- **Local Binding**: HTTP server binds to all interfaces but uses authentication
- **Xray API**: Xray API bound to localhost only for security
- **TLS Support**: Ready for TLS termination at proxy level

### Process Isolation

- **Binary Separation**: Xray runs as separate process
- **Resource Limits**: Configurable resource constraints
- **Signal Handling**: Clean shutdown procedures

## Scalability Design

### Horizontal Scaling

- **Multiple Instances**: Run multiple nodes on same server
- **Port Management**: Automatic free port detection
- **Service Isolation**: Each instance has separate systemd service

### Performance Optimization

- **Go Concurrency**: Goroutines for concurrent operations
- **Connection Pooling**: HTTP client connection reuse
- **Memory Management**: Efficient memory usage patterns
- **Background Processing**: Non-blocking background tasks

## Extensibility

### Plugin Architecture

- **Interface-based**: Abstract interfaces for core components
- **Proxy Cores**: Designed to support multiple proxy cores beyond Xray
- **Middleware**: Pluggable HTTP middleware system
- **Storage Backends**: Extensible storage layer

### Configuration System

- **Layered Configuration**: Default, override, and environment layers
- **Validation**: Comprehensive validation framework
- **Hot Reload**: Support for runtime configuration updates
- **Schema Evolution**: Forward-compatible configuration schema

## Monitoring and Observability

### Metrics Collection

- **Xray Statistics**: Traffic statistics via Xray API
- **Application Metrics**: Performance and error metrics
- **Health Checks**: Endpoint health status
- **Resource Usage**: Memory and CPU monitoring

### Logging Strategy

- **Structured Logs**: JSON-formatted logs for machine processing
- **Log Rotation**: Automatic log file rotation
- **Error Tracking**: Comprehensive error tracking and context
- **Audit Trail**: Security and configuration change logs

## Dependencies

### Core Dependencies

- **Go**: Programming language (1.24+)
- **Echo**: HTTP web framework
- **Zap**: Structured logging
- **Cobra**: CLI framework
- **Xray-core**: Proxy engine

### External Dependencies

- **Xray Binary**: Third-party proxy binary
- **System Dependencies**: make, curl, jq, systemd (for production)

This architecture provides a solid foundation for a distributed proxy system with excellent performance, security, and maintainability characteristics.
