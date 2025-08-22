# Deployment Guide

## Overview

Arch-Node supports multiple deployment methods to suit different environments and requirements. This guide covers Docker deployment, systemd service installation, and multi-instance setups.

## System Requirements

### Minimum Requirements

- **Operating System**: Debian 10+ or Ubuntu 18.04+
- **Architecture**: amd64 (x86_64)
- **Memory**: 1 GB RAM
- **CPU**: 1 core
- **Disk Space**: 500 MB
- **Network**: Stable internet connection

### Recommended Requirements

- **Memory**: 2 GB RAM
- **CPU**: 2+ cores
- **Disk Space**: 2 GB (for logs and storage)
- **Network**: Low-latency connection

## Deployment Methods

### 1. Docker Deployment (Recommended)

#### Quick Start

```bash
# Clone the repository
git clone https://github.com/ebadidev/arch-node.git
cd arch-node

# Start with Docker Compose
docker compose up -d

# Check status
docker compose ps
docker compose logs -f
```

#### Docker Compose Configuration

**File:** `docker-compose.yml`

```yaml
services:
  app:
    image: ghcr.io/ebadidev/arch-node:${TAG:-latest}
    restart: always
    network_mode: host
    volumes:
      - ./configs/:/app/configs/
      - ./storage/:/app/storage/
```

#### Environment Variables

```bash
# Set specific version
export TAG=v25.8.21
docker compose up -d

# Override configuration
export LOG_LEVEL=debug
docker compose up -d
```

#### Docker Commands

```bash
# Pull latest image
docker compose pull

# Start services
docker compose up -d

# Stop services
docker compose down

# View logs
docker compose logs -f

# Restart services
docker compose restart

# Update to latest
docker compose pull && docker compose up -d
```

### 2. Systemd Service Installation

#### Automated Setup

```bash
# Clone repository
git clone https://github.com/ebadidev/arch-node.git arch-node-1
cd arch-node-1

# Run setup script
make setup

# Check service status
systemctl status arch-node-1
```

#### Manual Installation

**1. Install Dependencies:**

```bash
apt-get update && apt-get install -y \
  make wget jq curl vim git openssl cron
```

**2. Setup Application:**

```bash
# Create directory
mkdir -p /opt/arch-node-1
cd /opt/arch-node-1

# Clone repository
git clone https://github.com/ebadidev/arch-node.git .

# Make scripts executable
chmod +x scripts/*.sh

# Run setup
./scripts/setup.sh
```

**3. Systemd Service:**

The setup script creates a systemd service file:

```ini
[Unit]
Description=Arch-Node Instance 1
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/arch-node-1
ExecStart=/opt/arch-node-1/arch-node start
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

**4. Enable and Start:**

```bash
systemctl enable arch-node-1
systemctl start arch-node-1
systemctl status arch-node-1
```

### 3. Manual Installation

#### Build from Source

```bash
# Install Go (1.24+)
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Clone and build
git clone https://github.com/ebadidev/arch-node.git
cd arch-node
make build

# The binary will be created as 'arch-node'
```

#### Manual Execution

```bash
# Setup Xray binary
./scripts/setup-xray.sh

# Create storage directories
mkdir -p storage/{database,app,logs}

# Run the application
./arch-node start
```

## Multi-Instance Deployment

### Running Multiple Instances

Each instance gets its own directory and systemd service:

```bash
# Instance 1
git clone https://github.com/ebadidev/arch-node.git arch-node-1
cd arch-node-1 && make setup

# Instance 2
git clone https://github.com/ebadidev/arch-node.git arch-node-2  
cd arch-node-2 && make setup

# Instance 3
git clone https://github.com/ebadidev/arch-node.git arch-node-3
cd arch-node-3 && make setup
```

### Managing Multiple Instances

```bash
# Check all instances
systemctl status arch-node-*

# Start/stop specific instance
systemctl start arch-node-2
systemctl stop arch-node-2

# Restart all instances
systemctl restart arch-node-*

# View logs for specific instance
journalctl -f -u arch-node-1
```

### Instance Information

```bash
# Get information for each instance
cd arch-node-1 && make info
cd arch-node-2 && make info
cd arch-node-3 && make info
```

Each instance will have:
- Unique HTTP port (randomly assigned)
- Unique authentication token
- Separate configuration and logs
- Independent systemd service

## Production Deployment

### 1. System Optimization

**BBR TCP Optimization:**

```bash
echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf
sysctl -p
```

**File Limits:**

```bash
# Add to /etc/security/limits.conf
* soft nofile 65536
* hard nofile 65536

# Add to /etc/sysctl.conf
fs.file-max = 65536
net.core.somaxconn = 65536
```

### 2. Security Hardening

**Firewall Configuration:**

```bash
# Allow SSH
ufw allow 22

# Allow specific ports for your proxy services
ufw allow 10001:10100/tcp

# Enable firewall
ufw enable
```

**User Security:**

```bash
# Create dedicated user (optional)
useradd -r -s /bin/false arch-node

# Set ownership
chown -R arch-node:arch-node /opt/arch-node-*
```

### 3. Monitoring Setup

**Log Rotation:**

```bash
# Add to /etc/logrotate.d/arch-node
/opt/arch-node-*/storage/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    copytruncate
}
```

**Resource Monitoring:**

```bash
# Monitor CPU and memory usage
systemctl status arch-node-*

# Monitor network connections
netstat -tulpn | grep arch-node

# Monitor log files
tail -f /opt/arch-node-*/storage/logs/*.log
```

### 4. Backup Strategy

**Configuration Backup:**

```bash
#!/bin/bash
# backup-arch-nodes.sh

BACKUP_DIR="/backup/arch-nodes"
DATE=$(date +%Y%m%d-%H%M%S)

mkdir -p "$BACKUP_DIR"

for instance in /opt/arch-node-*; do
    if [ -d "$instance" ]; then
        instance_name=$(basename "$instance")
        tar -czf "$BACKUP_DIR/${instance_name}-${DATE}.tar.gz" \
            -C "$instance" \
            configs/ storage/database/ storage/app/
    fi
done

# Keep only last 30 days
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +30 -delete
```

**Automated Backup:**

```bash
# Add to crontab
0 2 * * * /usr/local/bin/backup-arch-nodes.sh
```

## Load Balancing

### 1. Nginx Frontend

**Configuration:**

```nginx
upstream arch_nodes {
    server 127.0.0.1:10001;
    server 127.0.0.1:10002;
    server 127.0.0.1:10003;
}

server {
    listen 443 ssl;
    server_name proxy.example.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://arch_nodes;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. HAProxy Configuration

```haproxy
global
    daemon

defaults
    mode tcp
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms

frontend proxy_frontend
    bind *:443
    default_backend arch_nodes

backend arch_nodes
    balance roundrobin
    server node1 127.0.0.1:10001 check
    server node2 127.0.0.1:10002 check
    server node3 127.0.0.1:10003 check
```

## Troubleshooting Deployment

### Common Issues

**1. Port Conflicts:**

```bash
# Check port usage
netstat -tulpn | grep :10001

# Find free ports
for port in {10001..10100}; do
    if ! netstat -tuln | grep -q ":$port "; then
        echo "Port $port is free"
        break
    fi
done
```

**2. Permission Issues:**

```bash
# Check file permissions
ls -la /opt/arch-node-1/

# Fix permissions
chmod +x /opt/arch-node-1/arch-node
chmod +x /opt/arch-node-1/scripts/*.sh
chmod 644 /opt/arch-node-1/configs/*.json
```

**3. Service Startup Issues:**

```bash
# Check service status
systemctl status arch-node-1

# View detailed logs
journalctl -u arch-node-1 --no-pager

# Test manual startup
cd /opt/arch-node-1
./arch-node start
```

**4. Network Issues:**

```bash
# Test local connectivity
curl http://localhost:15888/

# Check firewall
ufw status
iptables -L

# Test external connectivity
curl -I http://external-ip:port/
```

### Debug Commands

```bash
# Check all arch-node services
systemctl list-units arch-node-*

# Monitor all logs
journalctl -f -u arch-node-*

# Check resource usage
top -p $(pgrep -f arch-node)

# Network connectivity test
ss -tulpn | grep arch-node
```

### Recovery Procedures

**1. Service Recovery:**

```bash
# Restart failed service
systemctl restart arch-node-1

# Reset failed state
systemctl reset-failed arch-node-1

# Disable and re-enable
systemctl disable arch-node-1
systemctl enable arch-node-1
systemctl start arch-node-1
```

**2. Configuration Recovery:**

```bash
# Reset to defaults
cd /opt/arch-node-1
rm -f storage/database/app.json
systemctl restart arch-node-1
```

**3. Complete Reinstall:**

```bash
# Stop service
systemctl stop arch-node-1
systemctl disable arch-node-1

# Remove files
rm -rf /opt/arch-node-1

# Reinstall
git clone https://github.com/ebadidev/arch-node.git /opt/arch-node-1
cd /opt/arch-node-1
make setup
```

## Performance Tuning

### 1. System Tuning

```bash
# Increase file descriptor limits
echo "arch-node soft nofile 65536" >> /etc/security/limits.conf
echo "arch-node hard nofile 65536" >> /etc/security/limits.conf

# Optimize network settings
echo "net.core.rmem_max = 134217728" >> /etc/sysctl.conf
echo "net.core.wmem_max = 134217728" >> /etc/sysctl.conf
echo "net.ipv4.tcp_rmem = 4096 87380 134217728" >> /etc/sysctl.conf
echo "net.ipv4.tcp_wmem = 4096 65536 134217728" >> /etc/sysctl.conf
sysctl -p
```

### 2. Application Tuning

**Configuration Optimization:**

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

**Resource Monitoring:**

```bash
# Monitor resource usage
watch -n 1 'ps aux | grep arch-node'

# Monitor network usage
iftop -i eth0

# Monitor disk I/O
iotop -o
```

This comprehensive deployment guide covers all aspects of getting Arch-Node running in production environments with proper security, monitoring, and performance optimization.
