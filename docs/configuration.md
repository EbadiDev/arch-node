# Configuration Guide

## Overview

Arch-Node uses a layered configuration system that allows for flexible deployment and customization. Configuration is handled through JSON files, environment variables, and runtime settings.

## Configuration Layers

### 1. Default Configuration

**File:** `configs/main.defaults.json`

This file contains the base configuration that ships with the application:

```json
{
  "logger": {
    "level": "warn",
    "format": "2006-01-02 15:04:05.000"
  },
  "xray": {
    "log_level": "info"
  }
}
```

**Purpose:** Provides sensible defaults for all environments.

### 2. Override Configuration

**File:** `configs/main.json` (optional)

Create this file to override default settings:

```json
{
  "logger": {
    "level": "debug"
  }
}
```

**Purpose:** Customize settings without modifying default configuration.

### 3. Environment Variables

Environment variables can override any configuration setting:

```bash
export LOG_LEVEL=debug
export HTTP_PORT=8080
```

### 4. Runtime Database

**File:** `storage/database/app.json`

Contains runtime-generated settings:

```json
{
  "settings": {
    "http_port": 15888,
    "http_token": "9CwH8bSQDR1nNtcO"
  },
  "manager": {
    "url": "https://manager.example.com/v1/nodes/1",
    "token": "manager-auth-token"
  }
}
```

## Configuration Schema

### Logger Configuration

```go
type Logger struct {
    Level  string `json:"level" validate:"required,oneof=debug info warn error"`
    Format string `json:"format" validate:"required"`
}
```

**Options:**
- `level`: Log verbosity (debug, info, warn, error)
- `format`: Timestamp format for logs

**Example:**
```json
{
  "logger": {
    "level": "info",
    "format": "2006-01-02 15:04:05.000"
  }
}
```

### Xray Configuration

```go
type Xray struct {
    LogLevel string `json:"log_level" validate:"required,oneof=debug info warn error"`
}
```

**Options:**
- `log_level`: Xray-core logging level

**Example:**
```json
{
  "xray": {
    "log_level": "warn"
  }
}
```

## Environment Variables

### Application Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `TAG` | Docker image tag | `latest` | `v25.8.21` |
| `HTTP_PORT` | API server port | Auto-generated | `8080` |
| `LOG_LEVEL` | Application log level | `warn` | `debug` |

### Docker Variables

When using Docker Compose, these variables are supported:

```yaml
services:
  app:
    image: ghcr.io/ebadidev/arch-node:${TAG:-latest}
    environment:
      - LOG_LEVEL=${LOG_LEVEL:-info}
      - HTTP_PORT=${HTTP_PORT:-8080}
```

## Configuration Loading Process

### 1. Initialization Sequence

```go
func (c *Config) Init() error {
    // 1. Load default configuration
    content, err := os.ReadFile("configs/main.defaults.json")
    if err != nil {
        return err
    }
    json.Unmarshal(content, &c)

    // 2. Load override configuration (if exists)
    if utils.FileExist("configs/main.json") {
        content, err = os.ReadFile("configs/main.json")
        if err != nil {
            return err
        }
        json.Unmarshal(content, &c)  // Overrides defaults
    }

    // 3. Apply environment variables
    // (handled by individual components)

    // 4. Validate final configuration
    return validator.New().Struct(c)
}
```

### 2. Configuration Precedence

1. **Environment Variables** (highest priority)
2. **Override Configuration** (`configs/main.json`)
3. **Default Configuration** (`configs/main.defaults.json`)
4. **Built-in Defaults** (lowest priority)

## Deployment Configurations

### 1. Development Environment

**`configs/main.json`:**
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

**Benefits:**
- Verbose logging for debugging
- Detailed Xray logs
- Easy troubleshooting

### 2. Production Environment

**`configs/main.json`:**
```json
{
  "logger": {
    "level": "warn"
  },
  "xray": {
    "log_level": "error"
  }
}
```

**Benefits:**
- Minimal logging overhead
- Focus on errors and warnings
- Better performance

### 3. Testing Environment

**`configs/main.json`:**
```json
{
  "logger": {
    "level": "info"
  },
  "xray": {
    "log_level": "warn"
  }
}
```

**Benefits:**
- Balanced logging
- Adequate debugging information
- Performance testing friendly

## Advanced Configuration

### 1. Custom Log Formats

The timestamp format follows Go's reference time format:

```json
{
  "logger": {
    "format": "2006-01-02T15:04:05.000Z07:00"  // ISO 8601
  }
}
```

**Common Formats:**
- `2006-01-02 15:04:05.000` - Default format
- `2006-01-02T15:04:05.000Z07:00` - ISO 8601
- `Jan 2 15:04:05` - Syslog format
- `15:04:05.000` - Time only

### 2. Xray Configuration Override

Xray configuration is managed separately in `storage/app/xray.json`. This file is:
- Generated automatically on startup
- Updated by the manager synchronization
- Can be manually modified (will be overwritten by sync)

**Sample Xray Configuration:**
```json
{
  "log": {
    "logLevel": "info",
    "access": "./storage/logs/xray-access.log",
    "error": "./storage/logs/xray-error.log"
  },
  "inbounds": [
    {
      "tag": "api",
      "protocol": "dokodemo-door",
      "listen": "127.0.0.1",
      "port": 3411,
      "settings": {
        "address": "127.0.0.1",
        "network": "tcp"
      }
    }
  ],
  "outbounds": [
    {
      "tag": "out",
      "protocol": "freedom"
    }
  ],
  "dns": {
    "servers": ["8.8.8.8", "8.8.4.4", "localhost"]
  },
  "stats": {},
  "api": {
    "tag": "api",
    "services": ["StatsService"]
  }
}
```

## Configuration Management

### 1. Validation

All configuration is validated using struct tags:

```go
type Config struct {
    Logger struct {
        Level  string `json:"level" validate:"required,oneof=debug info warn error"`
        Format string `json:"format" validate:"required"`
    } `json:"logger" validate:"required"`
}
```

**Validation Rules:**
- Required fields must be present
- String values must match allowed options
- Format strings must be valid

### 2. Configuration Display

On startup, the application displays the loaded configuration:

```bash
Config: {"logger":{"level":"warn","format":"2006-01-02 15:04:05.000"},"xray":{"log_level":"info"}}
```

### 3. Runtime Updates

Some configuration can be updated at runtime:

**Manager Configuration:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"url":"https://new-manager.com","token":"new-token"}' \
  http://localhost:15888/v1/manager
```

**Xray Configuration:**
- Updated automatically through manager synchronization
- Applied with automatic Xray restart
- Validated before application

## Configuration Files Location

### Directory Structure

```
configs/
├── main.defaults.json    # Default configuration (committed to git)
└── main.json            # Override configuration (optional, gitignored)

storage/
├── database/
│   └── app.json         # Runtime database
└── app/
    └── xray.json        # Xray configuration
```

### File Permissions

```bash
# Configuration files should be readable
chmod 644 configs/*.json

# Database files should be writable by application
chmod 644 storage/database/*.json
chmod 644 storage/app/*.json

# Ensure directories exist and are writable
chmod 755 storage/ storage/database/ storage/app/ storage/logs/
```

## Configuration Best Practices

### 1. Environment-Specific Settings

**Development:**
```json
{
  "logger": {
    "level": "debug"
  }
}
```

**Production:**
```json
{
  "logger": {
    "level": "error"
  }
}
```

### 2. Security Considerations

- Never commit sensitive tokens to version control
- Use environment variables for sensitive data
- Restrict file permissions on configuration files
- Regularly rotate authentication tokens

### 3. Configuration Management

- Use version control for default configurations
- Document configuration changes
- Test configuration changes in non-production environments
- Keep backups of working configurations

### 4. Monitoring Configuration

- Monitor configuration file changes
- Log configuration updates
- Validate configuration before applying
- Alert on configuration errors

## Troubleshooting

### Common Configuration Issues

1. **Invalid JSON Syntax**
   ```bash
   # Validate JSON syntax
   cat configs/main.json | jq .
   
   # Error: parse error: Invalid numeric literal at line 3, column 10
   ```

2. **Validation Errors**
   ```bash
   # Check logs for validation errors
   journalctl -u arch-node-1 | grep -i validation
   
   # Sample error: Validation error: Level must be one of [debug info warn error]
   ```

3. **File Permission Issues**
   ```bash
   # Check file permissions
   ls -la configs/
   
   # Fix permissions
   chmod 644 configs/*.json
   ```

4. **Missing Configuration Files**
   ```bash
   # Check if default config exists
   ls -la configs/main.defaults.json
   
   # If missing, the application will fail to start
   ```

### Debug Commands

```bash
# View current configuration
grep "Config:" /var/log/syslog | tail -1

# Check configuration file syntax
jq . configs/main.defaults.json
jq . configs/main.json

# Monitor configuration changes
inotifywait -m configs/ storage/database/

# Test configuration validation
./arch-node start --dry-run
```

### Configuration Recovery

```bash
# Reset to default configuration
rm -f configs/main.json

# Reset runtime database
rm -f storage/database/app.json

# Restart service
systemctl restart arch-node-1
```

This flexible configuration system allows for easy deployment across different environments while maintaining security and reliability.
