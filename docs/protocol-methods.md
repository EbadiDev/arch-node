# Arch-Node Protocol Support Documentation

This document describes the protocol methods available in the arch-node package for creating Xray configurations.

## Overview

The arch-node package provides factory methods for generating Xray inbound and outbound configurations for all major proxy protocols:

- **Shadowsocks** - AEAD cipher-based protocol
- **VLESS** - Lightweight protocol with minimal overhead  
- **VMess** - Feature-rich protocol with multiple encryption options
- **Trojan** - TLS-based protocol that mimics HTTPS traffic

## Protocol Methods

### Shadowsocks

#### Inbound
```go
func (c *Config) MakeShadowsocksInbound(tag, password, method, network string, port int, clients []*Client) *Inbound
```

**Parameters:**
- `tag`: Unique identifier for the inbound
- `password`: Shared secret for authentication
- `method`: Encryption method (e.g., "aes-256-gcm", "chacha20-poly1305")
- `network`: Network type ("tcp", "udp", "tcp,udp")
- `port`: Listen port
- `clients`: Array of client configurations

#### Outbound
```go
func (c *Config) MakeShadowsocksOutbound(tag, host, password, method string, port int) *Outbound
```

**Parameters:**
- `tag`: Unique identifier for the outbound
- `host`: Server address
- `password`: Shared secret for authentication
- `method`: Encryption method
- `port`: Server port

### VLESS

#### Inbound
```go
func (c *Config) MakeVlessInbound(tag string, port int, uuid string, network string, security interface{}) *Inbound
```

**Parameters:**
- `tag`: Unique identifier for the inbound
- `port`: Listen port
- `uuid`: Client UUID (UUID v4 format)
- `network`: Network type ("tcp", "ws", "http", "grpc")
- `security`: Security configuration (TLS, Reality, etc.)

**Features:**
- No encryption overhead (encryption handled by transport layer)
- Supports XTLS for better performance
- UUID-based client authentication
- Decryption set to "none" by default

#### Outbound
```go
func (c *Config) MakeVlessOutbound(tag, address string, port int, uuid, network string) *Outbound
```

**Parameters:**
- `tag`: Unique identifier for the outbound
- `address`: Server address
- `port`: Server port
- `uuid`: Client UUID
- `network`: Network type

### VMess

#### Inbound
```go
func (c *Config) MakeVmessInbound(tag string, port int, uuid, encryption, network string) *Inbound
```

**Parameters:**
- `tag`: Unique identifier for the inbound
- `port`: Listen port
- `uuid`: Client UUID (UUID v4 format)
- `encryption`: Encryption method ("auto", "aes-128-gcm", "chacha20-poly1305", "none")
- `network`: Network type ("tcp", "ws", "http", "grpc")

**Features:**
- UUID-based client authentication
- Multiple encryption options
- AlterId automatically set to 0 (recommended)
- Level set to 0 by default

#### Outbound
```go
func (c *Config) MakeVmessOutbound(tag, address string, port int, uuid, encryption, network string) *Outbound
```

**Parameters:**
- `tag`: Unique identifier for the outbound
- `address`: Server address
- `port`: Server port
- `uuid`: Client UUID
- `encryption`: Encryption method
- `network`: Network type

### Trojan

#### Inbound
```go
func (c *Config) MakeTrojanInbound(tag string, port int, password, network string, security interface{}) *Inbound
```

**Parameters:**
- `tag`: Unique identifier for the inbound
- `port`: Listen port
- `password`: Client password
- `network`: Network type ("tcp", "ws", "http", "grpc")
- `security`: Security configuration (typically TLS)

**Features:**
- Password-based authentication
- Designed to work with TLS
- Mimics HTTPS traffic patterns

#### Outbound
```go
func (c *Config) MakeTrojanOutbound(tag, address string, port int, password, network string) *Outbound
```

**Parameters:**
- `tag`: Unique identifier for the outbound
- `address`: Server address
- `port`: Server port
- `password`: Client password
- `network`: Network type

## Data Structures

### Client Structure
```go
type Client struct {
    Password string `json:"password" validate:"omitempty,min=1,max=64"`
    Method   string `json:"method" validate:"required"`
    Email    string `json:"email" validate:"required"`
    ID       string `json:"id,omitempty"`       // For VMess/VLESS UUID
    AlterId  int    `json:"alterId,omitempty"`  // For VMess
    Level    int    `json:"level,omitempty"`    // User level
}
```

### Inbound Settings
```go
type InboundSettings struct {
    Address    string    `json:"address,omitempty"`
    Clients    []*Client `json:"clients,omitempty" validate:"omitempty,dive"`
    Network    string    `json:"network,omitempty"`
    Method     string    `json:"method,omitempty"`
    Password   string    `json:"password,omitempty"`
    Decryption string    `json:"decryption,omitempty"` // For VLESS
}
```

### Outbound Server
```go
type OutboundServer struct {
    Address  string `json:"address" validate:"required"`
    Port     int    `json:"port" validate:"required,min=1,max=65536"`
    Method   string `json:"method" validate:"required"`
    Password string `json:"password" validate:"omitempty"`
    Uot      bool   `json:"uot"`
    ID       string `json:"id,omitempty"`      // For VMess/VLESS UUID  
    AlterId  int    `json:"alterId,omitempty"` // For VMess
    Level    int    `json:"level,omitempty"`   // User level
}
```

## Usage Examples

### Creating a Multi-Protocol Configuration

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/ebadidev/arch-node/pkg/xray"
)

func main() {
    config := xray.NewConfig("info")
    
    // Add VLESS inbound
    vlessInbound := config.MakeVlessInbound(
        "vless-in", 
        10001, 
        "550e8400-e29b-41d4-a716-446655440000", 
        "tcp", 
        nil,
    )
    config.Inbounds = append(config.Inbounds, vlessInbound)
    
    // Add VMess outbound
    vmessOutbound := config.MakeVmessOutbound(
        "vmess-out", 
        "server.example.com", 
        443, 
        "550e8400-e29b-41d4-a716-446655440001", 
        "auto", 
        "tcp",
    )
    config.Outbounds = append(config.Outbounds, vmessOutbound)
    
    // Serialize to JSON
    jsonData, _ := json.MarshalIndent(config, "", "  ")
    fmt.Println(string(jsonData))
}
```

### Shadowsocks with Custom Clients

```go
clients := []*xray.Client{
    {
        Password: "user1-password",
        Method:   "aes-256-gcm",
        Email:    "user1@example.com",
    },
    {
        Password: "user2-password", 
        Method:   "aes-256-gcm",
        Email:    "user2@example.com",
    },
}

ssInbound := config.MakeShadowsocksInbound(
    "ss-in",
    "",  // No global password when using clients
    "aes-256-gcm",
    "tcp",
    8388,
    clients,
)
```

## Network Types

All protocols support multiple network transports:

- **tcp**: Standard TCP transport
- **ws**: WebSocket transport (good for bypassing firewalls)
- **http**: HTTP/2 transport
- **grpc**: gRPC transport (experimental)

## Security Configurations

### TLS Settings
For protocols that support TLS (VLESS, VMess, Trojan):

```go
// TLS configuration would be passed as security parameter
// Implementation depends on your TLS certificate setup
```

### Reality Settings  
For VLESS with Reality:

```go
// Reality configuration for advanced traffic masking
// Implementation requires Reality key generation
```

## Best Practices

1. **Use appropriate encryption**: 
   - Shadowsocks: "aes-256-gcm" or "chacha20-poly1305"
   - VMess: "auto" for automatic selection
   - VLESS: Relies on transport layer encryption
   - Trojan: No application-layer encryption needed

2. **UUID Generation**: Use proper UUID v4 format for VMess/VLESS

3. **Port Management**: Ensure ports don't conflict with system services

4. **Network Selection**: Choose transport based on your environment:
   - TCP for general use
   - WebSocket for firewall traversal
   - gRPC for advanced scenarios

5. **Validation**: Always validate configurations before deployment:
   ```go
   if err := config.Validate(); err != nil {
       // Handle validation error
   }
   ```

## Testing

The package includes comprehensive tests for all protocols:

```bash
go test ./pkg/xray -v
```

Test coverage includes:
- Individual protocol method testing
- JSON serialization/deserialization
- Full configuration validation
- Protocol-specific feature verification
- Multi-protocol integration testing
