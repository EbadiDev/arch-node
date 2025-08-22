# Xray Integration

## Overview

Arch-Node integrates with Xray-core as its primary proxy engine. This integration provides high-performance proxy capabilities with support for multiple protocols and advanced traffic management features.

## Xray-Core Architecture

### What is Xray-Core?

Xray-core is a platform for building proxies to bypass network restrictions. It provides:

- **Multiple Protocols**: VMess, VLESS, Trojan, Shadowsocks, and more
- **Advanced Routing**: Complex routing rules based on domain, IP, and other criteria
- **High Performance**: Optimized for speed and low resource usage
- **Extensible**: Plugin architecture for custom protocols

### Integration Benefits

- **Process Isolation**: Xray runs as a separate process for stability
- **Dynamic Configuration**: Hot-reload configuration without service interruption
- **Statistics API**: Real-time traffic and performance metrics
- **Centralized Management**: Configuration sync from Arch-Manager

## Xray Process Management

### 1. Binary Management

**Binary Location:**
```
third_party/xray-linux-64/xray      # Linux binary
third_party/xray-macos-arm64/xray   # macOS binary (development)
```

**Binary Setup:**
```bash
# Automatic setup
make setup-xray

# Manual setup
./scripts/setup-xray.sh

# Download specific version
./scripts/setup-xray.sh v1.8.0
```

**Binary Verification:**
```bash
# Check binary exists and is executable
ls -la third_party/xray-linux-64/xray

# Test binary
./third_party/xray-linux-64/xray version
```

### 2. Process Lifecycle

**Startup Sequence:**
```go
func (x *Xray) Run() error {
    // 1. Save configuration to file
    if err := x.saveConfig(); err != nil {
        return err
    }
    
    // 2. Start Xray process
    go x.runCore()
    
    // 3. Connect to Xray API
    err := x.connect()
    return err
}
```

**Process Execution:**
```go
func (x *Xray) runCore() {
    x.command = exec.Command(x.binaryPath, "-c", x.configPath)
    x.command.Stderr = os.Stderr
    x.command.Stdout = os.Stdout
    
    if err := x.command.Run(); err != nil {
        x.l.Fatal("xray: cannot execute the binary", zap.Error(err))
    }
}
```

**Graceful Shutdown:**
```go
func (x *Xray) Close() error {
    // Close API connection
    if x.connection != nil {
        x.connection.Close()
    }
    
    // Kill process
    if x.command != nil && x.command.Process != nil {
        x.command.Process.Kill()
    }
    
    return nil
}
```

## Configuration Management

### 1. Configuration Structure

**Main Configuration File:** `storage/app/xray.json`

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
    },
    {
      "tag": "proxy",
      "protocol": "vmess",
      "port": 10001,
      "settings": {
        "clients": [
          {
            "id": "uuid-here",
            "alterId": 0
          }
        ]
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
  },
  "policy": {
    "levels": {
      "0": {
        "statsUserUplink": true,
        "statsUserDownlink": true
      }
    },
    "system": {
      "statsInboundUplink": true,
      "statsInboundDownlink": true
    }
  },
  "routing": {
    "domainStrategy": "IPIfNonMatch",
    "rules": []
  }
}
```

### 2. Configuration Components

**Logging Configuration:**
```go
type Log struct {
    LogLevel string `json:"logLevel" validate:"required,oneof=debug info warn error"`
    Access   string `json:"access" validate:"required"`
    Error    string `json:"error" validate:"required"`
}
```

**Inbound Configuration:**
```go
type Inbound struct {
    Tag      string            `json:"tag" validate:"required"`
    Protocol string            `json:"protocol" validate:"required"`
    Listen   string            `json:"listen,omitempty"`
    Port     int               `json:"port" validate:"required,min=1,max=65535"`
    Settings *InboundSettings  `json:"settings,omitempty"`
}
```

**Outbound Configuration:**
```go
type Outbound struct {
    Tag      string             `json:"tag" validate:"required"`
    Protocol string             `json:"protocol" validate:"required"`
    Settings *OutboundSettings  `json:"settings,omitempty"`
}
```

### 3. Dynamic Configuration Updates

**Configuration Comparison:**
```go
func (c *Config) Equals(other *Config) bool {
    thisBytes, _ := json.Marshal(c)
    otherBytes, _ := json.Marshal(other)
    return string(thisBytes) == string(otherBytes)
}
```

**Hot Reload Process:**
```go
func (c *Coordinator) Sync() error {
    remoteConfig, err := c.fetchConfig(c.d.Data.Manager)
    if err != nil {
        return err
    }

    if !c.xray.Config().Equals(remoteConfig) {
        c.xray.SetConfig(remoteConfig)
        go c.xray.Restart()  // Non-blocking restart
    }

    return nil
}
```

## API Communication

### 1. gRPC Connection

**Connection Setup:**
```go
func (x *Xray) connect() error {
    inbound := x.config.FindInbound("api")
    if inbound == nil {
        return errors.New("no api inbound found")
    }
    
    address := "127.0.0.1:" + strconv.Itoa(inbound.Port)
    
    for {
        x.connection, err = grpc.NewClient(address, 
            grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err == nil {
            return nil
        }
        time.Sleep(time.Second)
    }
}
```

**API Services:**
- **StatsService**: Traffic statistics
- **HandlerService**: Dynamic handler management
- **LoggerService**: Log level management

### 2. Statistics Collection

**Query Statistics:**
```go
func (x *Xray) QueryStats() ([]*stats.Stat, error) {
    client := stats.NewStatsServiceClient(x.connection)
    qs, err := client.QueryStats(context.Background(), 
        &stats.QueryStatsRequest{Reset_: true})
    if err != nil {
        return nil, err
    }
    return qs.GetStat(), nil
}
```

**Statistics Format:**
```json
{
  "stats": [
    {
      "name": "inbound>>>proxy>>>traffic>>>uplink",
      "value": 1024576
    },
    {
      "name": "inbound>>>proxy>>>traffic>>>downlink",
      "value": 2048192
    },
    {
      "name": "user>>>user1>>>traffic>>>uplink",
      "value": 512288
    }
  ]
}
```

## Protocol Support

### 1. VMess Protocol

**Configuration Example:**
```json
{
  "tag": "vmess-inbound",
  "protocol": "vmess",
  "port": 10001,
  "settings": {
    "clients": [
      {
        "id": "uuid-v4-here",
        "alterId": 0,
        "level": 0
      }
    ]
  },
  "streamSettings": {
    "network": "tcp",
    "security": "none"
  }
}
```

**Features:**
- UUID-based authentication
- Dynamic port allocation
- Multiple encryption methods
- WebSocket transport support

### 2. VLESS Protocol

**Configuration Example:**
```json
{
  "tag": "vless-inbound",
  "protocol": "vless",
  "port": 10002,
  "settings": {
    "clients": [
      {
        "id": "uuid-v4-here",
        "level": 0
      }
    ],
    "decryption": "none"
  },
  "streamSettings": {
    "network": "tcp",
    "security": "tls",
    "tlsSettings": {
      "certificates": [...]
    }
  }
}
```

**Features:**
- Lightweight protocol
- No encryption overhead
- TLS termination support
- XTLS support for better performance

### 3. Trojan Protocol

**Configuration Example:**
```json
{
  "tag": "trojan-inbound",
  "protocol": "trojan",
  "port": 10003,
  "settings": {
    "clients": [
      {
        "password": "password-here",
        "level": 0
      }
    ]
  },
  "streamSettings": {
    "network": "tcp",
    "security": "tls"
  }
}
```

## Routing and Traffic Management

### 1. Routing Rules

**Domain-based Routing:**
```json
{
  "routing": {
    "domainStrategy": "IPIfNonMatch",
    "rules": [
      {
        "type": "field",
        "domain": ["google.com", "youtube.com"],
        "outboundTag": "proxy"
      },
      {
        "type": "field",
        "domain": ["geosite:cn"],
        "outboundTag": "direct"
      }
    ]
  }
}
```

**IP-based Routing:**
```json
{
  "rules": [
    {
      "type": "field",
      "ip": ["geoip:cn"],
      "outboundTag": "direct"
    },
    {
      "type": "field",
      "ip": ["geoip:private"],
      "outboundTag": "block"
    }
  ]
}
```

### 2. Load Balancing

**Round Robin:**
```json
{
  "balancers": [
    {
      "tag": "proxy-balancer",
      "selector": ["proxy-1", "proxy-2", "proxy-3"],
      "strategy": {
        "type": "roundrobin"
      }
    }
  ]
}
```

**Least Load:**
```json
{
  "strategy": {
    "type": "leastload",
    "settings": {
      "healthCheck": {
        "interval": "1m",
        "timeout": "10s"
      }
    }
  }
}
```

## Monitoring and Logs

### 1. Access Logs

**Access Log Format:**
```
2024-01-15 10:30:45 [Info] [Transport] [tcp:127.0.0.1:10001] accepted connection
2024-01-15 10:30:45 [Info] [Proxy][VMess][In] received request to tcp:google.com:443
2024-01-15 10:30:45 [Info] [Router] picked outbound: proxy
```

**Log Configuration:**
```json
{
  "log": {
    "logLevel": "info",
    "access": "./storage/logs/xray-access.log",
    "error": "./storage/logs/xray-error.log"
  }
}
```

### 2. Error Handling

**Common Error Patterns:**
```bash
# Configuration errors
grep "config" storage/logs/xray-error.log

# Connection errors  
grep -i "connection\|timeout" storage/logs/xray-error.log

# Protocol errors
grep -i "vmess\|vless\|trojan" storage/logs/xray-error.log
```

### 3. Performance Monitoring

**Resource Usage:**
```bash
# Memory usage
ps aux | grep xray

# Network connections
netstat -tulpn | grep xray

# File descriptors
lsof -p $(pgrep xray)
```

## Security Considerations

### 1. Process Security

- **Isolation**: Xray runs as separate process
- **Resource Limits**: Configurable memory and CPU limits
- **File Permissions**: Strict permissions on configuration files
- **Network Binding**: API only binds to localhost

### 2. Configuration Security

- **Validation**: All configuration is validated before application
- **Secrets Management**: UUIDs and passwords generated securely
- **Access Control**: API requires authentication
- **Audit Logging**: All configuration changes logged

### 3. Network Security

- **TLS Support**: Full TLS/XTLS support for all protocols
- **Certificate Management**: Automated certificate handling
- **Port Management**: Dynamic port allocation prevents conflicts
- **Firewall Integration**: Proper firewall rule management

## Troubleshooting

### Common Issues

**1. Process Startup Failures:**
```bash
# Check binary permissions
ls -la third_party/xray-linux-64/xray

# Test binary manually
./third_party/xray-linux-64/xray -test -c storage/app/xray.json

# Check configuration syntax
./third_party/xray-linux-64/xray -test -c storage/app/xray.json
```

**2. API Connection Issues:**
```bash
# Check API port in configuration
grep -A 5 '"tag": "api"' storage/app/xray.json

# Test API connectivity
curl -v http://127.0.0.1:3411/

# Check process is running
ps aux | grep xray
```

**3. Configuration Validation Errors:**
```bash
# Validate JSON syntax
cat storage/app/xray.json | jq .

# Check for port conflicts
netstat -tulpn | grep :10001

# Review validation errors in logs
grep -i validation storage/logs/app-err.log
```

### Debug Commands

```bash
# Monitor Xray logs in real-time
tail -f storage/logs/xray-*.log

# Check Xray process status
systemctl status arch-node-1 | grep -A 5 xray

# Test configuration manually
./third_party/xray-linux-64/xray run -c storage/app/xray.json

# Monitor API calls
strace -p $(pgrep xray) -e trace=network

# Check memory usage
cat /proc/$(pgrep xray)/status | grep -i mem
```

### Performance Tuning

**1. System Optimization:**
```bash
# Increase file descriptor limits
ulimit -n 65536

# Optimize network settings
echo 'net.core.rmem_max = 134217728' >> /etc/sysctl.conf
echo 'net.core.wmem_max = 134217728' >> /etc/sysctl.conf
sysctl -p
```

**2. Xray Configuration:**
```json
{
  "policy": {
    "levels": {
      "0": {
        "connIdle": 300,
        "downlinkOnly": 5,
        "handshake": 4,
        "uplinkOnly": 2
      }
    }
  }
}
```

This comprehensive integration with Xray-core provides a robust, high-performance proxy solution that can be centrally managed through the Arch-Manager ecosystem.
