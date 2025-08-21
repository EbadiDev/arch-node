# Arch-Node

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE.md)
[![Docker](https://img.shields.io/badge/docker-supported-blue.svg)](https://hub.docker.com)
[![Status](https://img.shields.io/badge/status-development-orange.svg)](https://github.com/miladrahimi/p-node)

**Arch-Node** is a high-performance, lightweight multi-core proxy node built with Go that integrates multiple VPN cores for advanced traffic management and tunneling capabilities. It serves as a distributed node component that works seamlessly with the Arch-Manager ecosystem.

> ⚠️ **Development Status**: This project is currently in active development and is **not production ready**. Use at your own risk in production environments.

## ✨ Features

- **🚀 High Performance**: Built with Go for optimal speed and resource efficiency
- **⚡ Multi-Core Support**: Supports multiple VPN cores (currently Xray-Core, more coming soon)
- **🔄 Xray Integration**: Full integration with Xray-core for advanced proxy protocols
- **🔮 Extensible**: Modular architecture designed for easy integration of additional proxy cores
- **📊 Real-time Monitoring**: Built-in HTTP API for status monitoring and statistics
- **🔧 Easy Management**: Simple configuration and deployment via Arch-Manager
- **🐳 Docker Support**: Ready-to-use Docker containers and compose files
- **📈 Auto-scaling**: Support for multiple instances on a single server
- **🔒 Secure**: Built-in authentication and secure communication protocols
- **📝 Comprehensive Logging**: Structured logging with configurable levels

## 🏗️ Architecture

Arch-Node consists of several key components:

- **HTTP Server**: RESTful API for management and monitoring
- **Xray Core**: Current traffic proxy and tunneling engine
- **Core Manager**: Pluggable architecture for additional proxy cores (planned)
- **Database Manager**: Configuration and state persistence
- **Coordinator**: Synchronization with Arch-Manager
- **Worker Pool**: Concurrent task processing with multi-core optimization

## 🚀 Quick Start

### Prerequisites

- **Operating System**: Debian 10+ or Ubuntu 18.04+
- **Architecture**: amd64 (x86_64)
- **Memory**: 1 GB RAM minimum (2 GB recommended)
- **CPU**: 1 core minimum (2+ cores recommended)
- **Network**: Stable internet connection

### Installation

1. **System Dependencies**
   ```bash
   apt-get update && apt-get install -y \
     make wget jq curl vim git openssl cron
   ```

2. **BBR TCP Optimization** (Optional but recommended)
   ```bash
   echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
   echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf
   sysctl -p
   ```

3. **Install Arch-Node**
   ```bash
   # Find available directory name
   for ((i=1;;i++)); do [ ! -d "arch-node-${i}" ] && break; done
   
   # Clone and setup
   git clone https://github.com/ebadidev/arch-node.git "arch-node-${i}"
   cd "arch-node-${i}"
   make setup
   ```

4. **Get Node Information**
   ```bash
   make info
   ```

### Docker Deployment

1. **Using Docker Compose**
   ```bash
   # Pull the latest image
   docker compose pull
   
   # Start the service
   docker compose up -d
   ```

2. **Custom Configuration**
   ```bash
   # Edit configuration
   vim configs/main.defaults.json
   
   # Restart with new config
   docker compose restart
   ```

## ⚙️ Configuration

### Main Configuration (`configs/main.defaults.json`)

```json
{
  "logger": {
    "level": "info",        // debug, info, warn, error
    "format": "2006-01-02 15:04:05.000"
  },
  "xray": {
    "log_level": "info"     // Xray logging level
  }
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TAG` | Docker image tag | `latest` |
| `HTTP_PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Application log level | `info` |

## 🔧 Management

### Service Operations

```bash
# Check service status
systemctl status arch-node-1

# Start/stop service
systemctl start arch-node-1
systemctl stop arch-node-1

# Restart service
systemctl restart arch-node-1
```

### Updates

```bash
# Automatic update (recommended)
make update

# Manual update
git pull
make setup
systemctl restart arch-node-1
```

### Logs and Monitoring

```bash
# View real-time logs
journalctl -f -u arch-node-1

# Check application logs
tail -f ./storage/logs/*.log

# View service status
make info
```

## 🛠️ Development

### Local Development

```bash
# Setup development environment
make local-setup

# Run locally
make local-run

# Clean logs
make local-clean

# Fresh start (clear all data)
make local-fresh
```

### Building

```bash
# Build for Linux amd64
make build

# The binary will be created as 'arch-node'
```

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Health check |
| `/api/v1/stats` | GET | Node statistics |
| `/api/v1/configs` | GET/POST | Configuration management |
| `/api/v1/manager` | GET/POST | Manager operations |

## 📊 Monitoring

Arch-Node provides comprehensive monitoring through:

- **HTTP API**: Real-time statistics and health endpoints
- **Structured Logs**: JSON formatted logs with correlation IDs
- **Metrics**: Performance and usage metrics
- **Health Checks**: Automated service health monitoring

## 🔒 Security

- **Authentication**: Token-based authentication for API access
- **Encryption**: TLS encryption for all communications
- **Isolation**: Process isolation and secure defaults
- **Audit Logs**: Comprehensive audit logging

## 🤝 Integration

### Arch-Manager Setup

1. Get node information:
   ```bash
   make info
   ```

2. Register with Arch-Manager:
   ```bash
   make set-manager URL="https://your-manager.com" TOKEN="your-token"
   ```

### Multiple Instances

Run multiple instances on the same server:

```bash
# Instance 1
git clone https://github.com/ebadidev/arch-node.git arch-node-1
cd arch-node-1 && make setup

# Instance 2
git clone https://github.com/ebadidev/arch-node.git arch-node-2  
cd arch-node-2 && make setup

# Each instance will have its own systemd service
```

## 📋 Troubleshooting

### Common Issues

1. **Service won't start**
   ```bash
   # Check logs
   journalctl -u arch-node-1 --no-pager
   
   # Verify configuration
   ./arch-node start --dry-run
   ```

2. **High memory usage**
   ```bash
   # Monitor resource usage
   systemctl status arch-node-1
   
   # Adjust log levels in config
   vim configs/main.defaults.json
   ```

3. **Connection issues**
   ```bash
   # Check network connectivity
   curl -I http://localhost:8080/
   
   # Verify firewall rules
   ufw status
   ```

## 📚 Related Projects

- **[Arch-Manager](https://github.com/ebadidev/arch-manager)** - Central management platform
- **[Xray-core](https://github.com/XTLS/Xray-core)** - Underlying proxy engine

## 📄 License

This project is licensed under the terms specified in the [LICENSE](LICENSE.md) file.

---

**Made with ❤️ by the Arch Net team**
