# Summary: Node-Manager Communication and Database Architecture

## How Nodes Connect to the Manager

### Connection Architecture

Arch-Node uses an HTTP-based client-server model to communicate with the Arch-Manager:

1. **HTTP Client Communication**: Nodes act as HTTP clients that periodically connect to the manager
2. **Token-Based Authentication**: Each node uses a manager-provided token for authentication
3. **Periodic Synchronization**: Automatic configuration sync every 30 seconds
4. **RESTful API**: Simple REST endpoints for configuration and status updates

### Connection Process

1. **Manual Registration**:
   ```bash
   make set-manager URL="https://manager.example.com/v1/nodes/1" TOKEN="manager-token"
   ```

2. **Configuration Storage**: Manager details stored in `storage/database/app.json`:
   ```json
   {
     "manager": {
       "url": "https://manager.example.com/v1/nodes/1",
       "token": "manager-auth-token"
     }
   }
   ```

3. **Automatic Synchronization**: The coordinator component handles sync:
   - Fetches configuration from manager every 30 seconds
   - Compares with local configuration
   - Applies changes and restarts Xray if needed

### Communication Flow

```
Node (Client) ←→ Manager (Server)
     │
     ├── GET /configs (fetch configuration)
     ├── POST /stats (send statistics)  
     └── Authentication via Bearer token
```

### Security Features

- **HTTPS Communication**: All communication encrypted
- **Token Authentication**: Bearer token for each request
- **Client Identification**: Custom headers (`X-App-Name: Arch-Node`)
- **Request Validation**: Configuration validated before application

## Database System Architecture

### Database Type: JSON File-Based Persistence

**Type**: File-based JSON database (NOT a traditional SQL/NoSQL database)

**Why JSON Files?**
- **Simplicity**: Easy to read, edit, and debug
- **Portability**: Works across all platforms without dependencies
- **Transparency**: Human-readable format
- **Backup-friendly**: Easy to backup and restore
- **No External Dependencies**: No need for external database servers

### Database Structure

#### Main Database File: `storage/database/app.json`

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

#### Data Models

1. **Settings**: Node-specific configuration
   - `http_port`: Random port (1000-65536) for management API
   - `http_token`: 16-character random authentication token

2. **Manager**: Connection information for Arch-Manager
   - `url`: Full URL to manager's node endpoint
   - `token`: Authentication token for manager communication

### How the Database Works

#### 1. **Initialization**
```go
func (d *Database) Init() error {
    if fileExists(Path) {
        return d.Load()  // Load existing data
    }
    
    // Generate new random values
    d.Data.Settings.HttpPort = randomPort()
    d.Data.Settings.HttpToken = randomToken()
    
    return d.Save()  // Create new database file
}
```

#### 2. **Thread-Safe Operations**
```go
type Database struct {
    locker *sync.Mutex  // Protects concurrent access
    Data   *Data
}

func (d *Database) Save() error {
    d.locker.Lock()
    defer d.locker.Unlock()
    
    content, _ := json.Marshal(d.Data)
    return os.WriteFile(Path, content, 0755)
}
```

#### 3. **Data Validation**
- Struct tags for validation rules
- JSON schema validation
- Type checking and constraints
- URL and token format validation

#### 4. **Persistence Strategy**
- **Read**: Load entire file into memory at startup
- **Write**: Marshal to JSON and write entire file
- **Atomic Updates**: File replacement ensures consistency
- **Backup**: Simple file copy for backups

### Storage Layout

```
storage/
├── database/
│   └── app.json           # Main database (settings + manager config)
├── app/  
│   ├── xray.json          # Xray configuration (managed separately)
│   └── update.txt         # Last update timestamp
└── logs/
    ├── app-std.log        # Application stdout logs
    ├── app-err.log        # Application stderr logs  
    ├── xray-access.log    # Xray access logs
    └── xray-error.log     # Xray error logs
```

### Performance Characteristics

- **Read Performance**: In-memory access after initial load
- **Write Performance**: ~1ms for typical file sizes (<1KB)
- **Concurrency**: Mutex-based thread safety
- **Scalability**: Suitable for single-node configuration data
- **Reliability**: Atomic file updates prevent corruption

### Database vs Traditional Systems

| Aspect | JSON File DB | Traditional DB |
|--------|--------------|----------------|
| **Setup** | Zero setup | Requires server setup |
| **Dependencies** | None | Database server required |
| **Debugging** | Human-readable | Requires tools |
| **Backup** | File copy | Database dumps |
| **Scaling** | Single node only | Multi-node support |
| **Performance** | Low latency | Higher latency |
| **Data Size** | < 1KB typical | Unlimited |

## Key Architecture Decisions

### Why This Approach?

1. **Simplicity**: No complex database setup or management
2. **Reliability**: File-based storage is very reliable
3. **Debugging**: Easy to inspect and modify data
4. **Portability**: Works everywhere Go works
5. **Resource Efficiency**: Minimal memory and CPU overhead
6. **Backup/Recovery**: Simple file operations

### Trade-offs

**Advantages**:
- Zero external dependencies
- Human-readable data format
- Easy backup and restore
- Fast read performance
- Simple debugging

**Limitations**:
- Not suitable for large datasets
- No complex queries
- Single-node only (no clustering)
- Write performance degrades with file size
- No built-in replication

### Comparison with Other Approaches

This design prioritizes simplicity and reliability over advanced database features, which is perfect for:
- Configuration storage
- Small datasets (< 1MB)
- Single-node applications
- Embedded systems
- Development and testing environments

For larger, more complex applications, you might consider:
- **SQLite**: For more complex queries while staying embedded
- **PostgreSQL/MySQL**: For multi-node deployments
- **Redis**: For high-performance caching
- **etcd**: For distributed configuration management

The JSON file approach perfectly fits Arch-Node's requirements: simple configuration storage with transparency, reliability, and ease of use.
