# API Reference

## Overview

Arch-Node provides a RESTful HTTP API for management and monitoring. The API uses JSON for data exchange and token-based authentication for security.

## Base URL

```
http://localhost:<port>/
```

The port is randomly generated on first startup and stored in `storage/database/app.json`.

## Authentication

All protected endpoints require Bearer token authentication:

```http
Authorization: Bearer <token>
```

The token is automatically generated and stored in `storage/database/app.json`.

### Getting Authentication Token

```bash
# Extract token from database
cat storage/database/app.json | jq -r '.settings.http_token'

# Get port and token together
cat storage/database/app.json | jq '.settings'
```

## Endpoints

### 1. Health Check

**GET /** - Health check endpoint

No authentication required.

**Request:**
```http
GET / HTTP/1.1
Host: localhost:15888
```

**Response:**
```http
HTTP/1.1 200 OK
Content-Type: text/plain

OK
```

**Usage:**
```bash
curl http://localhost:15888/
```

---

### 2. Statistics

**GET /v1/stats** - Get node statistics

Returns Xray traffic statistics and performance metrics.

**Request:**
```http
GET /v1/stats HTTP/1.1
Host: localhost:15888
Authorization: Bearer 9CwH8bSQDR1nNtcO
```

**Response:**
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
    },
    {
      "name": "user>>>user1>>>traffic>>>downlink",
      "value": 1024576
    }
  ]
}
```

**Usage:**
```bash
TOKEN=$(cat storage/database/app.json | jq -r '.settings.http_token')
PORT=$(cat storage/database/app.json | jq -r '.settings.http_port')

curl -H "Authorization: Bearer $TOKEN" \
     "http://localhost:$PORT/v1/stats"
```

---

### 3. Configuration Management

**POST /v1/configs** - Update Xray configuration

Updates the Xray configuration and restarts the Xray process.

**Request Headers:**
```http
Content-Type: application/json
Authorization: Bearer <token>
X-App-Name: Arch-Manager
```

**Request Body:**
```json
{
  "log": {
    "logLevel": "info",
    "access": "./storage/logs/xray-access.log",
    "error": "./storage/logs/xray-error.log"
  },
  "inbounds": [
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
      },
      "streamSettings": {
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
    "servers": ["8.8.8.8", "8.8.4.4"]
  },
  "routing": {
    "rules": []
  }
}
```

**Response:**
```json
{
  "message": "The configs stored successfully."
}
```

**Error Responses:**

*400 Bad Request - Invalid JSON:*
```json
{
  "message": "Cannot parse the request body."
}
```

*422 Unprocessable Entity - Validation Error:*
```json
{
  "message": "Validation error: Port is required"
}
```

*422 Unprocessable Entity - Port Conflict:*
```json
{
  "message": "The port 'proxy.10001' is already in use"
}
```

*400 Bad Request - Unknown Client:*
```json
{
  "message": "Unknown client."
}
```

**Usage:**
```bash
TOKEN=$(cat storage/database/app.json | jq -r '.settings.http_token')
PORT=$(cat storage/database/app.json | jq -r '.settings.http_port')

curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-App-Name: Arch-Manager" \
  -d @config.json \
  "http://localhost:$PORT/v1/configs"
```

---

### 4. Manager Configuration

**POST /v1/manager** - Configure manager connection

Sets or updates the connection information for Arch-Manager.

**Request:**
```json
{
  "url": "https://manager.example.com/v1/nodes/1",
  "token": "manager-auth-token"
}
```

**Response:**
```json
{
  "manager": {
    "url": "https://manager.example.com/v1/nodes/1",
    "token": "manager-auth-token"
  }
}
```

**Clear Manager Configuration:**
```json
{
  "url": "",
  "token": ""
}
```

**Error Responses:**

*400 Bad Request:*
```json
{
  "message": "Cannot parse the request body."
}
```

*422 Unprocessable Entity:*
```json
{
  "message": "Validation error: Url must be a valid URL"
}
```

**Usage:**
```bash
TOKEN=$(cat storage/database/app.json | jq -r '.settings.http_token')
PORT=$(cat storage/database/app.json | jq -r '.settings.http_port')

# Set manager
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"url":"https://manager.example.com/v1/nodes/1","token":"secret"}' \
  "http://localhost:$PORT/v1/manager"

# Clear manager
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"url":"","token":""}' \
  "http://localhost:$PORT/v1/manager"
```

## Request/Response Format

### Content Type

All API endpoints use JSON for data exchange:

```http
Content-Type: application/json
```

### Request Headers

**Required for protected endpoints:**
- `Authorization: Bearer <token>`
- `Content-Type: application/json` (for POST requests)

**Required for configuration endpoint:**
- `X-App-Name: Arch-Manager` (must match exactly)

### Response Format

**Success Response:**
```json
{
  "data": {},
  "message": "Success message"
}
```

**Error Response:**
```json
{
  "message": "Error description"
}
```

## Error Handling

### HTTP Status Codes

| Code | Description | Usage |
|------|-------------|-------|
| 200 | OK | Successful GET requests |
| 201 | Created | Successful POST requests |
| 400 | Bad Request | Invalid request format |
| 401 | Unauthorized | Missing or invalid token |
| 422 | Unprocessable Entity | Validation errors |
| 500 | Internal Server Error | Server-side errors |

### Error Response Examples

**401 Unauthorized:**
```json
{
  "message": "Unauthorized"
}
```

**400 Bad Request:**
```json
{
  "message": "Cannot parse the request body."
}
```

**422 Unprocessable Entity:**
```json
{
  "message": "Validation error: Port must be between 1 and 65536"
}
```

## Authentication Examples

### Get Token and Port

```bash
# Extract from database file
DB_FILE="storage/database/app.json"
TOKEN=$(jq -r '.settings.http_token' "$DB_FILE")
PORT=$(jq -r '.settings.http_port' "$DB_FILE")

echo "Token: $TOKEN"
echo "Port: $PORT"
```

### Test Authentication

```bash
# Test with valid token
curl -H "Authorization: Bearer $TOKEN" \
     "http://localhost:$PORT/v1/stats"

# Test with invalid token
curl -H "Authorization: Bearer invalid-token" \
     "http://localhost:$PORT/v1/stats"
```

## Configuration Validation

### Xray Configuration Schema

The configuration endpoint validates against the Xray configuration schema:

**Required Fields:**
- `log`: Logging configuration
- `inbounds`: Array of inbound configurations
- `outbounds`: Array of outbound configurations
- `dns`: DNS configuration

**Inbound Configuration:**
```json
{
  "tag": "unique-tag",
  "protocol": "vmess|vless|trojan|shadowsocks|...",
  "port": 1024,
  "settings": {},
  "streamSettings": {}
}
```

**Port Validation:**
- Ports must be between 1 and 65536
- Ports must not be in use by other services
- API port is automatically assigned
- Remote port conflicts are handled specially

### Manager Configuration Schema

```json
{
  "url": "https://manager.example.com/path",  // Required, valid URL
  "token": "auth-token"                       // Required, 1-128 characters
}
```

## Rate Limiting

Currently, no rate limiting is implemented, but considerations for production:

- Implement rate limiting middleware
- Set appropriate limits for different endpoints
- Consider IP-based or token-based limiting
- Monitor for abuse patterns

## Security Considerations

### API Security

1. **Token Security**
   - Tokens are randomly generated (16 characters)
   - Tokens should be rotated regularly
   - Never log tokens in plain text

2. **Network Security**
   - API binds to all interfaces (0.0.0.0)
   - Use firewall rules to restrict access
   - Consider VPN access for management

3. **Request Validation**
   - All inputs are validated
   - JSON parsing is secure
   - Configuration is validated before application

### Client Authentication

```bash
# Example of secure API usage
TOKEN=$(cat storage/database/app.json | jq -r '.settings.http_token')
PORT=$(cat storage/database/app.json | jq -r '.settings.http_port')

# Use HTTPS in production
curl -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "https://node.example.com:$PORT/v1/stats"
```

## SDK Examples

### Bash/Shell

```bash
#!/bin/bash

# Load configuration
DB_FILE="storage/database/app.json"
TOKEN=$(jq -r '.settings.http_token' "$DB_FILE")
PORT=$(jq -r '.settings.http_port' "$DB_FILE")
BASE_URL="http://localhost:$PORT"

# Function to make API calls
api_call() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    
    if [ -n "$data" ]; then
        curl -s -X "$method" \
             -H "Authorization: Bearer $TOKEN" \
             -H "Content-Type: application/json" \
             -d "$data" \
             "$BASE_URL$endpoint"
    else
        curl -s -X "$method" \
             -H "Authorization: Bearer $TOKEN" \
             "$BASE_URL$endpoint"
    fi
}

# Usage examples
api_call "GET" "/v1/stats"
api_call "POST" "/v1/manager" '{"url":"https://manager.com","token":"secret"}'
```

### Python

```python
import json
import requests

class ArchNodeAPI:
    def __init__(self, host="localhost", port=None, token=None):
        if port is None or token is None:
            # Load from database file
            with open("storage/database/app.json") as f:
                db = json.load(f)
                port = db["settings"]["http_port"]
                token = db["settings"]["http_token"]
        
        self.base_url = f"http://{host}:{port}"
        self.headers = {
            "Authorization": f"Bearer {token}",
            "Content-Type": "application/json"
        }
    
    def get_stats(self):
        response = requests.get(f"{self.base_url}/v1/stats", headers=self.headers)
        return response.json()
    
    def set_manager(self, url, token):
        data = {"url": url, "token": token}
        response = requests.post(f"{self.base_url}/v1/manager", 
                               headers=self.headers, json=data)
        return response.json()
    
    def update_config(self, config):
        headers = self.headers.copy()
        headers["X-App-Name"] = "Arch-Manager"
        response = requests.post(f"{self.base_url}/v1/configs", 
                               headers=headers, json=config)
        return response.json()

# Usage
api = ArchNodeAPI()
stats = api.get_stats()
print(json.dumps(stats, indent=2))
```

This API provides comprehensive management capabilities for Arch-Node instances while maintaining security and ease of use.
