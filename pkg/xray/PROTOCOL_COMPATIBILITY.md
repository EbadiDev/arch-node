# Xray Protocol Compatibility Guide

This guide explains the transport and security feature compatibility across different Xray protocols.

## Protocol Support Matrix

### VMess Protocol
**Supported Features:**
- **Transports:** TCP, WebSocket (ws), gRPC, KCP, HTTP Upgrade
- **Security:** TLS only
- **Encryption:** auto, aes-128-gcm, chacha20-poly1305, none

**NOT Supported:**
- ❌ REALITY security (VLESS-only feature)
- ❌ XHTTP transport (VLESS-only feature)

### VLESS Protocol  
**Supported Features:**
- **Transports:** TCP, WebSocket (ws), gRPC, KCP, HTTP Upgrade, XHTTP
- **Security:** TLS, REALITY
- **Encryption:** none (VLESS is unencrypted by design)

**All Features Supported:**
- ✅ All transport types including XHTTP
- ✅ Both TLS and REALITY security

### Trojan Protocol
**Supported Features:**
- **Transports:** TCP (recommended)
- **Security:** TLS (required)

**Restrictions:**
- ⚠️ Typically only supports TCP transport
- ⚠️ Requires TLS security

### Shadowsocks Protocol
**Supported Features:**
- **Transports:** TCP, WebSocket (basic)
- **Security:** Built-in encryption methods

**NOT Supported:**
- ❌ REALITY security
- ❌ XHTTP transport

## API Design Philosophy

### Composable Transport Configuration

Instead of creating separate convenience methods for every transport combination (e.g., `MakeVmessWebSocketTlsInbound`, `MakeVmessGrpcTlsInbound`, etc.), this package uses a **composable approach**:

```go
// ✅ Composable - Build transport settings step by step
wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com")
wsSettings = config.AddTlsToStreamSettings(wsSettings, "server.com", false)
inbound := config.MakeVmessInbound("vmess-ws", 8080, "uuid", "auto", wsSettings)

// ❌ Monolithic - Would require separate methods for each combination
// config.MakeVmessWebSocketTlsInbound(...)  // Too many specific methods
// config.MakeVmessGrpcTlsInbound(...)       // API bloat
// config.MakeVmessKcpInbound(...)           // Hard to maintain
```

### Benefits of Composable Design

1. **Flexibility:** Mix and match any transport with any security setting
2. **Maintainability:** No API bloat from combinatorial explosion of methods
3. **Reusability:** Transport settings can be reused across different protocols
4. **Extensibility:** Adding new transports doesn't require new convenience methods

### Building Transport Configurations

```go
// Step 1: Create base transport settings
wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com")
grpcSettings := config.MakeGrpcStreamSettings("service", "authority")
kcpSettings := config.MakeKcpStreamSettings("none", "seed")
xhttpSettings := config.MakeXhttpStreamSettings("host.com", "/path", "auto")

// Step 2: Add security if needed
wsSettings = config.AddTlsToStreamSettings(wsSettings, "server.com", false)
grpcSettings = config.AddRealityToStreamSettings(grpcSettings, "example.com:443", 
    []string{"example.com"}, "privateKey", "publicKey")

// Step 3: Create inbound with configured transport
vmessInbound := config.MakeVmessInbound("tag", 8080, "uuid", "auto", wsSettings)
vlessInbound := config.MakeVlessInbound("tag", 8080, "uuid", "tcp", grpcSettings)
```

## Code Examples

### VMess Examples (Valid Configurations)

```go
// VMess with basic TCP
config.MakeVmessInbound("vmess-tcp", 8080, "uuid", "auto", nil)

// VMess with WebSocket + TLS
wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com")
wsSettings = config.AddTlsToStreamSettings(wsSettings, "server.com", false)
config.MakeVmessInbound("vmess-ws", 8080, "uuid", "auto", wsSettings)

// VMess with gRPC + TLS  
grpcSettings := config.MakeGrpcStreamSettings("service", "authority")
grpcSettings = config.AddTlsToStreamSettings(grpcSettings, "server.com", false)
config.MakeVmessInbound("vmess-grpc", 8080, "uuid", "auto", grpcSettings)

// VMess with KCP
kcpSettings := config.MakeKcpStreamSettings("none", "seed")
config.MakeVmessInbound("vmess-kcp", 8080, "uuid", "auto", kcpSettings)

// VMess with HTTP Upgrade + TLS
httpSettings := config.MakeHttpUpgradeStreamSettings("/upgrade", "host.com")
httpSettings = config.AddTlsToStreamSettings(httpSettings, "server.com", false)
config.MakeVmessInbound("vmess-http", 8080, "uuid", "auto", httpSettings)
```

### VLESS Examples (All Features Supported)

```go
// VLESS with basic TCP
config.MakeVlessInbound("vless-tcp", 8080, "uuid", "tcp", nil)

// VLESS with WebSocket + TLS
wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com")
wsSettings = config.AddTlsToStreamSettings(wsSettings, "server.com", false)
config.MakeVlessInbound("vless-ws", 8080, "uuid", "tcp", wsSettings)

// VLESS with gRPC + REALITY (only VLESS supports this)
grpcSettings := config.MakeGrpcStreamSettings("service", "")
grpcSettings = config.AddRealityToStreamSettings(grpcSettings, "example.com:443", 
    []string{"example.com"}, "privateKey", "publicKey")
config.MakeVlessInbound("vless-reality", 8080, "uuid", "tcp", grpcSettings)

// VLESS with XHTTP (only VLESS supports this)
xhttpSettings := config.MakeXhttpStreamSettings("host.com", "/path", "auto")
config.MakeVlessInbound("vless-xhttp", 8080, "uuid", "tcp", xhttpSettings)

// VLESS with XHTTP + REALITY (most advanced combination)
xhttpSettings := config.MakeXhttpStreamSettings("host.com", "/path", "auto")
xhttpSettings = config.AddRealityToStreamSettings(xhttpSettings, "example.com:443",
    []string{"example.com"}, "privateKey", "publicKey")
config.MakeVlessInbound("vless-advanced", 8080, "uuid", "tcp", xhttpSettings)
```

## Validation Functions

The package provides validation functions to check compatibility:

```go
// Check if a protocol supports specific settings
err := config.ValidateProtocolCompatibility("vmess", streamSettings)
if err != nil {
    // Handle incompatible combination
}

// Automatically fix incompatible settings
sanitized := config.SanitizeStreamSettingsForProtocol("vmess", streamSettings)
```

## Migration Guide

### From VMess to VLESS
If you need REALITY or XHTTP features, migrate from VMess to VLESS:

```go
// Old VMess (limited features)
wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com")
wsSettings = config.AddTlsToStreamSettings(wsSettings, "server.com", false)
vmessInbound := config.MakeVmessInbound("vmess-ws", 8080, "uuid", "auto", wsSettings)

// New VLESS (all features available)
wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com") 
wsSettings = config.AddTlsToStreamSettings(wsSettings, "server.com", false)
vlessInbound := config.MakeVlessInbound("vless-ws", 8080, "uuid", "tcp", wsSettings)

// Or use advanced VLESS features
grpcSettings := config.MakeGrpcStreamSettings("service", "")
grpcSettings = config.AddRealityToStreamSettings(grpcSettings, "example.com:443",
    []string{"example.com"}, "privateKey", "publicKey")
vlessRealityInbound := config.MakeVlessInbound("vless-reality", 8080, "uuid", "tcp", grpcSettings)
```

### Transport Fallbacks
The sanitization function provides automatic fallbacks:
- VMess + REALITY → VMess + TLS
- VMess + XHTTP → VMess + WebSocket
- Trojan + non-TLS → Trojan + TLS
- Trojan + non-TCP → Trojan + TCP

## Best Practices

1. **Use VLESS for modern features:** If you need REALITY or XHTTP, use VLESS protocol
2. **Validate configurations:** Always validate compatibility before deployment
3. **Test transport combinations:** Some combinations work better in certain network environments
4. **Security considerations:** REALITY provides better censorship resistance than TLS

## Common Mistakes to Avoid

❌ **Don't use VMess + REALITY**
```go
// This will fail validation
grpcSettings := config.MakeGrpcStreamSettings("service", "")
grpcSettings = config.AddRealityToStreamSettings(grpcSettings, "example.com:443",
    []string{"example.com"}, "privateKey", "publicKey")
config.MakeVmessInbound("tag", 8080, "uuid", "auto", grpcSettings) // Invalid!
```

✅ **Use VLESS + REALITY instead**
```go
// This is the correct way
grpcSettings := config.MakeGrpcStreamSettings("service", "")
grpcSettings = config.AddRealityToStreamSettings(grpcSettings, "example.com:443",
    []string{"example.com"}, "privateKey", "publicKey")
config.MakeVlessInbound("tag", 8080, "uuid", "tcp", grpcSettings) // Valid!
```

❌ **Don't use VMess + XHTTP**
```go
// This will fail validation
xhttpSettings := config.MakeXhttpStreamSettings("host.com", "/path", "auto")
config.MakeVmessInbound("tag", 8080, "uuid", "auto", xhttpSettings) // Invalid!
```

✅ **Use VLESS + XHTTP instead**
```go
// This is the correct way
xhttpSettings := config.MakeXhttpStreamSettings("host.com", "/path", "auto")
config.MakeVlessInbound("tag", 8080, "uuid", "tcp", xhttpSettings) // Valid!
```

❌ **Don't create API bloat with too many convenience methods**
```go
// Bad - creates combinatorial explosion
config.MakeVmessWebSocketTlsInbound(...)
config.MakeVmessGrpcTlsInbound(...)  
config.MakeVmessKcpInbound(...)
config.MakeVlessWebSocketTlsInbound(...)
config.MakeVlessGrpcRealityInbound(...)
// ... hundreds of combinations
```

✅ **Use composable transport configuration instead**
```go
// Good - flexible and maintainable
transport := config.MakeWebSocketStreamSettings("/path", "host.com")
transport = config.AddTlsToStreamSettings(transport, "server.com", false)
inbound := config.MakeVmessInbound("tag", 8080, "uuid", "auto", transport)
```
