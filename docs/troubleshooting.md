# Troubleshooting Guide

## Overview

This guide provides solutions to common issues, debugging techniques, and troubleshooting procedures for Arch-Node deployments.

## Common Issues and Solutions

### 1. Service Startup Issues

#### Issue: Service fails to start

**Symptoms:**
```bash
systemctl status arch-node-1
● arch-node-1.service - Arch-Node Instance 1
   Loaded: loaded (/etc/systemd/system/arch-node-1.service; enabled; vendor preset: enabled)
   Active: failed (Result: exit-code) since Mon 2024-01-15 10:30:00 UTC; 5s ago
```

**Diagnostic Steps:**
```bash
# Check detailed logs
journalctl -u arch-node-1 --no-pager

# Check file permissions
ls -la /opt/arch-node-1/arch-node
ls -la /opt/arch-node-1/configs/

# Test manual startup
cd /opt/arch-node-1
./arch-node start
```

**Common Solutions:**

1. **Missing execute permissions:**
   ```bash
   chmod +x /opt/arch-node-1/arch-node
   chmod +x /opt/arch-node-1/scripts/*.sh
   ```

2. **Missing configuration files:**
   ```bash
   ls -la configs/main.defaults.json
   # If missing, restore from repository
   ```

3. **Port conflicts:**
   ```bash
   # Check what's using the port
   netstat -tulpn | grep :15888
   
   # Reset database to generate new port
   rm storage/database/app.json
   systemctl restart arch-node-1
   ```

#### Issue: Binary not found

**Symptoms:**
```
xray: binary not found at path: third_party/xray-linux-64/xray
```

**Solutions:**
```bash
# Setup Xray binary
make setup-xray

# Or manually download
./scripts/setup-xray.sh

# Verify binary exists and is executable
ls -la third_party/xray-linux-64/xray
chmod +x third_party/xray-linux-64/xray
```

### 2. Configuration Issues

#### Issue: Invalid JSON configuration

**Symptoms:**
```
Config validation error: invalid character '}' looking for beginning of object key string
```

**Diagnostic Steps:**
```bash
# Validate JSON syntax
cat configs/main.defaults.json | jq .
cat configs/main.json | jq .

# Check for common JSON errors
python -m json.tool configs/main.json
```

**Solutions:**
```bash
# Reset to default configuration
rm configs/main.json
systemctl restart arch-node-1

# Fix JSON syntax manually or restore from backup
cp configs/main.json.backup configs/main.json
```

#### Issue: Configuration validation errors

**Symptoms:**
```
Validation error: Level must be one of [debug info warn error]
```

**Solutions:**
```bash
# Check allowed values in documentation
# Fix invalid configuration values
vim configs/main.json

# Example fix:
{
  "logger": {
    "level": "info"  // Was "verbose" (invalid)
  }
}
```

### 3. Network and Connectivity Issues

#### Issue: API endpoints not accessible

**Symptoms:**
```bash
curl: (7) Failed to connect to localhost port 15888: Connection refused
```

**Diagnostic Steps:**
```bash
# Check if service is running
systemctl status arch-node-1

# Check port from database
cat storage/database/app.json | jq '.settings.http_port'

# Check if port is bound
netstat -tulpn | grep :15888
ss -tulpn | grep :15888
```

**Solutions:**
```bash
# Restart service
systemctl restart arch-node-1

# Check firewall
ufw status
iptables -L

# Test with correct port
PORT=$(cat storage/database/app.json | jq -r '.settings.http_port')
curl http://localhost:$PORT/
```

#### Issue: Manager connection failures

**Symptoms:**
```
coordinator: cannot sync: connection timeout
```

**Diagnostic Steps:**
```bash
# Check manager configuration
cat storage/database/app.json | jq '.manager'

# Test manager connectivity manually
curl -I https://manager.example.com/configs

# Check DNS resolution
nslookup manager.example.com
```

**Solutions:**
```bash
# Test network connectivity
ping manager.example.com

# Check SSL certificates
curl -v https://manager.example.com

# Update manager configuration
TOKEN=$(cat storage/database/app.json | jq -r '.settings.http_token')
PORT=$(cat storage/database/app.json | jq -r '.settings.http_port')

curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"url":"https://correct-manager.com","token":"correct-token"}' \
  "http://localhost:$PORT/v1/manager"
```

### 4. Xray Integration Issues

#### Issue: Xray process crashes

**Symptoms:**
```
xray: cannot execute the binary: signal: killed
```

**Diagnostic Steps:**
```bash
# Check Xray logs
tail -f storage/logs/xray-error.log

# Test Xray configuration manually
./third_party/xray-linux-64/xray -test -c storage/app/xray.json

# Check process status
ps aux | grep xray
```

**Solutions:**
```bash
# Validate Xray configuration
./third_party/xray-linux-64/xray run -c storage/app/xray.json

# Check for port conflicts
netstat -tulpn | grep :10001

# Restart with fresh configuration
rm storage/app/xray.json
systemctl restart arch-node-1
```

#### Issue: Xray API connection failures

**Symptoms:**
```
xray: trying to connect to api: connection refused
```

**Diagnostic Steps:**
```bash
# Check API inbound configuration
grep -A 10 '"tag": "api"' storage/app/xray.json

# Check if API port is available
PORT=$(grep -A 5 '"tag": "api"' storage/app/xray.json | grep '"port"' | cut -d: -f2 | tr -d ' ,')
netstat -tulpn | grep :$PORT
```

**Solutions:**
```bash
# Ensure API inbound is properly configured
# Check that Xray process has started successfully
# Wait for process to fully initialize before API connection
```

### 5. Database Issues

#### Issue: Database corruption

**Symptoms:**
```
json: cannot unmarshal object into Go struct field
```

**Diagnostic Steps:**
```bash
# Check database file integrity
cat storage/database/app.json | jq .

# Check file permissions
ls -la storage/database/app.json
```

**Solutions:**
```bash
# Restore from backup if available
cp storage/database/app.json.backup storage/database/app.json

# Or reset to defaults (will generate new random values)
rm storage/database/app.json
systemctl restart arch-node-1
```

#### Issue: Permission denied errors

**Symptoms:**
```
open storage/database/app.json: permission denied
```

**Solutions:**
```bash
# Fix file permissions
chmod 644 storage/database/app.json
chmod 755 storage/database/

# Fix directory ownership
chown -R $(whoami):$(whoami) storage/
```

### 6. Performance Issues

#### Issue: High memory usage

**Diagnostic Steps:**
```bash
# Check memory usage
ps aux | grep arch-node
ps aux | grep xray

# Monitor memory over time
top -p $(pgrep arch-node)
```

**Solutions:**
```bash
# Reduce log levels
vim configs/main.json
{
  "logger": {
    "level": "warn"
  },
  "xray": {
    "log_level": "error"
  }
}

# Clean old logs
make local-clean

# Restart service
systemctl restart arch-node-1
```

#### Issue: High CPU usage

**Diagnostic Steps:**
```bash
# Monitor CPU usage
htop
iostat 1

# Check for busy loops in logs
grep -i "error\|timeout" storage/logs/*.log
```

**Solutions:**
```bash
# Check for configuration issues causing restart loops
journalctl -u arch-node-1 | grep -i restart

# Optimize Xray configuration
# Reduce polling frequency if applicable
```

## Debugging Tools and Techniques

### 1. Log Analysis

#### Application Logs
```bash
# Monitor all logs
tail -f storage/logs/*.log

# Filter by component
tail -f storage/logs/app-std.log | grep coordinator
tail -f storage/logs/app-std.log | grep xray
tail -f storage/logs/app-std.log | grep database

# Filter by log level
tail -f storage/logs/app-std.log | grep -E "(ERROR|WARN)"

# Search for specific errors
grep -r "connection refused" storage/logs/
```

#### System Logs
```bash
# Systemd journal
journalctl -u arch-node-1 -f
journalctl -u arch-node-1 --since "1 hour ago"
journalctl -u arch-node-1 --no-pager | grep -i error

# Syslog
tail -f /var/log/syslog | grep arch-node
```

### 2. Network Debugging

#### Port and Connection Analysis
```bash
# Check listening ports
netstat -tulpn | grep arch-node
ss -tulpn | grep arch-node

# Check established connections
netstat -tupln | grep ESTABLISHED
ss -tuln | grep ESTABLISHED

# Monitor network traffic
tcpdump -i eth0 port 10001
```

#### API Testing
```bash
# Test health endpoint
curl -v http://localhost:15888/

# Test with authentication
TOKEN=$(cat storage/database/app.json | jq -r '.settings.http_token')
PORT=$(cat storage/database/app.json | jq -r '.settings.http_port')

curl -v -H "Authorization: Bearer $TOKEN" \
     "http://localhost:$PORT/v1/stats"

# Test manager connectivity
MANAGER_URL=$(cat storage/database/app.json | jq -r '.manager.url')
MANAGER_TOKEN=$(cat storage/database/app.json | jq -r '.manager.token')

curl -v -H "Authorization: Bearer $MANAGER_TOKEN" "$MANAGER_URL"
```

### 3. Process Monitoring

#### Process Information
```bash
# Find process IDs
pgrep -f arch-node
pgrep -f xray

# Detailed process information
ps aux | grep -E "(arch-node|xray)"

# Process tree
pstree -p $(pgrep arch-node)

# Resource usage
top -p $(pgrep -d, -f "arch-node|xray")
```

#### File Descriptor Usage
```bash
# Check open files
lsof -p $(pgrep arch-node)
lsof -p $(pgrep xray)

# Count file descriptors
ls /proc/$(pgrep arch-node)/fd | wc -l
```

### 4. Configuration Debugging

#### Configuration Validation
```bash
# Validate JSON files
find configs/ -name "*.json" -exec jq . {} \;
find storage/ -name "*.json" -exec jq . {} \;

# Test configuration loading
./arch-node start --dry-run

# Compare configurations
diff <(jq -S . configs/main.defaults.json) <(jq -S . configs/main.json)
```

#### Environment Analysis
```bash
# Check environment variables
env | grep -E "(LOG_LEVEL|HTTP_PORT|TAG)"

# Check working directory
pwd
ls -la

# Check file permissions
find . -name "*.json" -exec ls -la {} \;
```

## Recovery Procedures

### 1. Service Recovery

#### Restart Service
```bash
# Basic restart
systemctl restart arch-node-1

# Stop, reset, and start
systemctl stop arch-node-1
systemctl reset-failed arch-node-1
systemctl start arch-node-1

# Check status
systemctl status arch-node-1
```

#### Reload Configuration
```bash
# After configuration changes
systemctl reload arch-node-1

# Or restart if reload not supported
systemctl restart arch-node-1
```

### 2. Configuration Recovery

#### Reset to Defaults
```bash
# Backup current configuration
cp configs/main.json configs/main.json.backup

# Remove custom configuration
rm configs/main.json

# Reset database
cp storage/database/app.json storage/database/app.json.backup
rm storage/database/app.json

# Restart service
systemctl restart arch-node-1
```

#### Restore from Backup
```bash
# Restore configuration
cp configs/main.json.backup configs/main.json

# Restore database
cp storage/database/app.json.backup storage/database/app.json

# Restart service
systemctl restart arch-node-1
```

### 3. Complete Reinstall

#### Clean Reinstall
```bash
# Stop service
systemctl stop arch-node-1
systemctl disable arch-node-1

# Backup important data
tar -czf arch-node-backup-$(date +%Y%m%d).tar.gz storage/

# Remove installation
rm -rf /opt/arch-node-1

# Fresh install
git clone https://github.com/ebadidev/arch-node.git /opt/arch-node-1
cd /opt/arch-node-1
make setup

# Restore data if needed
tar -xzf arch-node-backup-$(date +%Y%m%d).tar.gz
```

## Emergency Procedures

### 1. Service Hanging

```bash
# Check if process is responsive
kill -0 $(pgrep arch-node)

# Send SIGTERM for graceful shutdown
kill -TERM $(pgrep arch-node)

# Wait and force kill if necessary
sleep 10
kill -KILL $(pgrep arch-node)

# Restart service
systemctl start arch-node-1
```

### 2. Disk Space Issues

```bash
# Check disk usage
df -h

# Clean log files
rm storage/logs/*.log

# Rotate logs immediately
logrotate -f /etc/logrotate.d/arch-node

# Clean old backups
find storage/ -name "*.backup" -mtime +7 -delete
```

### 3. Port Conflicts

```bash
# Find what's using the port
PORT=$(cat storage/database/app.json | jq -r '.settings.http_port')
lsof -i :$PORT

# Kill conflicting process
kill $(lsof -t -i :$PORT)

# Or regenerate with new port
rm storage/database/app.json
systemctl restart arch-node-1
```

## Monitoring and Alerting

### 1. Health Checks

#### Basic Health Check Script
```bash
#!/bin/bash
# health-check.sh

SERVICE="arch-node-1"
URL="http://localhost:$(cat storage/database/app.json | jq -r '.settings.http_port')/"

# Check service status
if ! systemctl is-active --quiet $SERVICE; then
    echo "CRITICAL: $SERVICE is not running"
    exit 2
fi

# Check HTTP endpoint
if ! curl -s --max-time 5 "$URL" > /dev/null; then
    echo "CRITICAL: HTTP endpoint not responding"
    exit 2
fi

echo "OK: Service is healthy"
exit 0
```

#### Continuous Monitoring
```bash
# Add to crontab for regular checks
*/5 * * * * /opt/arch-node-1/health-check.sh || systemctl restart arch-node-1
```

### 2. Performance Monitoring

#### Resource Usage Script
```bash
#!/bin/bash
# monitor-resources.sh

ARCH_PID=$(pgrep arch-node)
XRAY_PID=$(pgrep xray)

if [ -n "$ARCH_PID" ]; then
    echo "Arch-Node CPU: $(ps -p $ARCH_PID -o %cpu --no-headers)%"
    echo "Arch-Node Memory: $(ps -p $ARCH_PID -o %mem --no-headers)%"
fi

if [ -n "$XRAY_PID" ]; then
    echo "Xray CPU: $(ps -p $XRAY_PID -o %cpu --no-headers)%"
    echo "Xray Memory: $(ps -p $XRAY_PID -o %mem --no-headers)%"
fi

echo "Disk Usage: $(df -h . | tail -1 | awk '{print $5}')"
```

This troubleshooting guide provides comprehensive solutions for most issues you'll encounter with Arch-Node deployments. Always start with the most basic checks (service status, logs) before moving to more complex debugging procedures.
