# Development Guide

## Overview

This guide covers setting up a development environment, understanding the codebase structure, contributing guidelines, and debugging techniques for Arch-Node.

## Development Environment Setup

### Prerequisites

- **Go**: Version 1.24 or higher
- **Git**: For version control
- **Make**: For build automation
- **Docker**: For containerized testing (optional)
- **jq**: For JSON processing in scripts

### Installation

**Ubuntu/Debian:**
```bash
# Install Go
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install dependencies
sudo apt-get update
sudo apt-get install -y make git jq curl vim

# Verify installation
go version
```

**macOS:**
```bash
# Using Homebrew
brew install go make jq

# Verify installation
go version
```

### Clone and Setup

```bash
# Clone repository
git clone https://github.com/ebadidev/arch-node.git
cd arch-node

# Setup development environment
make local-setup

# Verify setup
make local-run
```

## Project Structure

### Directory Layout

```
arch-node/
├── cmd/                    # Command-line interface
│   ├── root.go            # Root command
│   └── start.go           # Start command
├── internal/              # Private application code
│   ├── app/               # Application orchestration
│   ├── config/            # Configuration management
│   ├── coordinator/       # Manager synchronization
│   ├── database/          # Data persistence
│   ├── http/              # HTTP server and handlers
│   │   ├── handlers/      # HTTP handlers
│   │   │   ├── home.go    # Health check handler
│   │   │   └── v1/        # API v1 handlers
│   │   └── server/        # HTTP server setup
│   └── utils/             # Internal utilities
├── pkg/                   # Public/reusable packages
│   ├── http/              # HTTP client and middleware
│   │   ├── client/        # HTTP client
│   │   ├── middleware/    # HTTP middleware
│   │   └── validator/     # Request validation
│   ├── logger/            # Structured logging
│   ├── worker/            # Background workers
│   └── xray/              # Xray integration
├── configs/               # Configuration files
├── scripts/               # Deployment and utility scripts
├── storage/               # Runtime data and logs
├── third_party/           # External binaries
└── docs/                  # Documentation
```

### Code Organization

**Internal Packages (`internal/`):**
- `app/`: Main application lifecycle management
- `config/`: Configuration loading and validation
- `coordinator/`: Background synchronization with manager
- `database/`: JSON-based data persistence
- `http/`: Web server and API handlers

**Public Packages (`pkg/`):**
- `http/`: Reusable HTTP utilities
- `logger/`: Structured logging system
- `worker/`: Background task processing
- `xray/`: Xray-core integration

## Development Workflow

### 1. Local Development

```bash
# Setup development environment
make local-setup

# Run locally (with hot reload)
make local-run

# Clean logs
make local-clean

# Fresh start (clear all data)
make local-fresh
```

### 2. Building

```bash
# Build for current platform
go build -o arch-node

# Build for Linux (production)
make build

# Cross-compile for different platforms
GOOS=darwin GOARCH=amd64 go build -o arch-node-darwin
GOOS=windows GOARCH=amd64 go build -o arch-node.exe
```

### 3. Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/database

# Verbose output
go test -v ./...
```

### 4. Code Quality

```bash
# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Check for vulnerabilities
go list -json -m all | nancy sleuth

# Static analysis
go vet ./...
```

## Key Components

### 1. Application Lifecycle (`internal/app/app.go`)

**Main Application Structure:**
```go
type App struct {
    Context    context.Context
    Cancel     context.CancelFunc
    Shutdown   chan struct{}
    Config     *config.Config
    Logger     *logger.Logger
    HttpServer *server.Server
    HttpClient *client.Client
    Xray       *xray.Xray
    Syncer     *coordinator.Coordinator
    Database   *database.Database
}
```

**Initialization Order:**
1. Configuration loading
2. Logger setup
3. Database initialization
4. Xray setup
5. HTTP server creation
6. Coordinator setup
7. Signal handling

### 2. Configuration System (`internal/config/`)

**Configuration Loading:**
```go
func (c *Config) Init() error {
    // 1. Load defaults
    content, err := os.ReadFile(defaultConfigPath)
    json.Unmarshal(content, &c)
    
    // 2. Load overrides
    if utils.FileExist(configPath) {
        content, err = os.ReadFile(configPath)
        json.Unmarshal(content, &c)
    }
    
    // 3. Validate
    return validator.New().Struct(c)
}
```

### 3. HTTP Server (`internal/http/server/`)

**Server Setup:**
```go
func (s *Server) Run() {
    s.engine.Use(echoMiddleware.CORS())
    s.engine.Use(middleware.Logger(s.l))
    s.engine.Use(middleware.General())

    s.engine.GET("/", handlers.HomeShow())

    g2 := s.engine.Group("/v1")
    g2.Use(middleware.Authorize(func() string {
        return s.database.Data.Settings.HttpToken
    }))

    g2.GET("/stats", v1.StatsShow(s.xray))
    g2.POST("/configs", v1.ConfigsStore(s.xray))
    g2.POST("/manager", v1.ManagerStore(s.database))
}
```

### 4. Database System (`internal/database/`)

**Thread-Safe Operations:**
```go
func (d *Database) Save() error {
    d.locker.Lock()
    defer d.locker.Unlock()
    
    content, err := json.Marshal(d.Data)
    if err != nil {
        return err
    }
    
    return os.WriteFile(Path, content, 0755)
}
```

## Debugging

### 1. Local Debugging

**Enable Debug Logging:**
```json
{
  "logger": {
    "level": "debug"
  },
  "xray": {
    "log_level": "debug"
  }
}
```

**Debug with Delve:**
```bash
# Install Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug

# Set breakpoint and run
(dlv) break main.main
(dlv) continue
```

### 2. Logging and Monitoring

**Log Files:**
```
storage/logs/
├── app-std.log        # Application stdout
├── app-err.log        # Application stderr
├── xray-access.log    # Xray access logs
└── xray-error.log     # Xray error logs
```

**Log Monitoring:**
```bash
# Monitor all logs
tail -f storage/logs/*.log

# Monitor specific component
tail -f storage/logs/app-std.log | grep coordinator

# Monitor with filtering
journalctl -f -u arch-node-1 | grep -E "(error|warn)"
```

### 3. Performance Profiling

**CPU Profiling:**
```go
import _ "net/http/pprof"

// Add to main function
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

**Memory Profiling:**
```bash
# Get heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Get CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

### 4. Network Debugging

**API Testing:**
```bash
# Test health endpoint
curl -v http://localhost:15888/

# Test with authentication
TOKEN=$(cat storage/database/app.json | jq -r '.settings.http_token')
curl -v -H "Authorization: Bearer $TOKEN" http://localhost:15888/v1/stats

# Test manager connectivity
curl -v -H "Authorization: Bearer manager-token" https://manager.com/configs
```

## Contributing

### 1. Code Style

**Go Conventions:**
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions small and focused

**Example:**
```go
// ProcessConfig validates and applies the given configuration
func (x *Xray) ProcessConfig(config *Config) error {
    if err := config.Validate(); err != nil {
        return errors.Wrap(err, "config validation failed")
    }
    
    x.config = config
    return x.saveConfig()
}
```

### 2. Error Handling

**Use structured errors:**
```go
import "github.com/cockroachdb/errors"

// Wrap errors with context
if err := someOperation(); err != nil {
    return errors.Wrap(err, "operation failed")
}

// Stack traces for debugging
if err := criticalOperation(); err != nil {
    return errors.WithStack(err)
}
```

### 3. Testing

**Unit Test Example:**
```go
func TestDatabaseSave(t *testing.T) {
    // Setup
    tempDir := t.TempDir()
    database := &Database{
        Data: &Data{
            Settings: &Settings{
                HttpPort:  8080,
                HttpToken: "test-token",
            },
        },
    }
    
    // Test
    err := database.Save()
    
    // Assert
    assert.NoError(t, err)
    assert.FileExists(t, filepath.Join(tempDir, "app.json"))
}
```

### 4. Pull Request Guidelines

**Before Submitting:**
1. Run all tests: `go test ./...`
2. Format code: `go fmt ./...`
3. Check for race conditions: `go test -race ./...`
4. Update documentation if needed
5. Test locally with `make local-run`

**PR Template:**
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style
- [ ] Self-review completed
- [ ] Documentation updated
```

## Advanced Development

### 1. Adding New API Endpoints

**Create Handler:**
```go
// internal/http/handlers/v1/new_endpoint.go
func NewEndpoint(deps *Dependencies) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Implementation
        return c.JSON(http.StatusOK, response)
    }
}
```

**Register Route:**
```go
// internal/http/server/server.go
g2.GET("/new-endpoint", v1.NewEndpoint(deps))
```

### 2. Adding New Configuration Options

**Update Config Struct:**
```go
// internal/config/config.go
type Config struct {
    Logger struct {
        Level  string `json:"level" validate:"required,oneof=debug info warn error"`
        Format string `json:"format" validate:"required"`
    } `json:"logger" validate:"required"`
    
    NewSection struct {
        Option string `json:"option" validate:"required"`
    } `json:"new_section" validate:"required"`
}
```

**Update Default Config:**
```json
{
  "logger": {
    "level": "warn",
    "format": "2006-01-02 15:04:05.000"
  },
  "new_section": {
    "option": "default_value"
  }
}
```

### 3. Adding New Components

**Create Component:**
```go
// pkg/newcomponent/newcomponent.go
type NewComponent struct {
    logger *logger.Logger
    config *Config
}

func (nc *NewComponent) Start() error {
    // Implementation
}

func (nc *NewComponent) Stop() error {
    // Implementation
}

func New(logger *logger.Logger, config *Config) *NewComponent {
    return &NewComponent{
        logger: logger,
        config: config,
    }
}
```

**Integrate with App:**
```go
// internal/app/app.go
type App struct {
    // ... existing fields
    NewComponent *newcomponent.NewComponent
}

func New() (*App, error) {
    // ... existing initialization
    a.NewComponent = newcomponent.New(a.Logger, a.Config)
    return a, nil
}
```

## Troubleshooting Development Issues

### Common Development Problems

**1. Import Cycle:**
```bash
package command-line-arguments
imports github.com/ebadidev/arch-node/internal/app
imports github.com/ebadidev/arch-node/internal/http/server
imports github.com/ebadidev/arch-node/internal/app
import cycle not allowed
```

**Solution:** Restructure imports, use interfaces, or move shared code to a common package.

**2. Module Issues:**
```bash
go: module github.com/ebadidev/arch-node requires Go 1.24
```

**Solution:** Update Go version or adjust `go.mod` requirements.

**3. Race Conditions:**
```bash
==================
WARNING: DATA RACE
Write at 0x... by goroutine 7:
```

**Solution:** Add proper synchronization (mutex, channels) or use atomic operations.

### Debug Commands

```bash
# Check Go environment
go env

# Verify module dependencies
go mod verify
go mod tidy

# Check for unused dependencies
go mod why -m <module>

# Build with detailed output
go build -v

# Show assembly output
go build -gcflags="-S"
```

This development guide provides a comprehensive foundation for contributing to and extending Arch-Node.
