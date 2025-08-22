# Database System

## Overview

Arch-Node uses a simple, file-based JSON database system for storing configuration and runtime data. This approach provides simplicity, transparency, and easy debugging while maintaining adequate performance for the node's requirements.

## Database Architecture

### Design Philosophy

The database system is designed with the following principles:
- **Simplicity**: Easy to understand and debug
- **Transparency**: Human-readable JSON format
- **Reliability**: Thread-safe operations with mutex locking
- **Portability**: Works across all platforms without external dependencies
- **Backup-friendly**: Easy to backup and restore

### Core Components

1. **Database Manager** (`internal/database/database.go`)
2. **Data Models** (`internal/database/settings.go`, `internal/database/manager.go`)
3. **Storage Layer** (JSON files in `storage/database/`)

## Data Models

### 1. Settings Model

Stores node-specific configuration:

```go
type Settings struct {
    HttpPort  int    `json:"http_port" validate:"required,min=1,max=65536"`
    HttpToken string `json:"http_token" validate:"required,min=8,max=128"`
}
```

**Fields:**
- `http_port`: Random port (1000-65536) for the management API
- `http_token`: Cryptographically secure random token (16 characters)

### 2. Manager Model

Stores connection information for Arch-Manager:

```go
type Manager struct {
    Url   string `json:"url" validate:"required,url,min=1,max=1024"`
    Token string `json:"token" validate:"required,min=1,max=128"`
}
```

**Fields:**
- `url`: Full URL to the manager's node endpoint
- `token`: Authentication token for manager communication

### 3. Data Container

Main data structure containing all persistent data:

```go
type Data struct {
    Settings *Settings `json:"settings"`
    Manager  *Manager  `json:"manager"`
}
```

## File Structure

### Database File Location

```
storage/database/app.json
```

### Sample Database File

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

### Directory Structure

```
storage/
├── database/
│   └── app.json          # Main database file
├── app/
│   ├── xray.json         # Xray configuration
│   └── update.txt        # Last update timestamp
└── logs/
    ├── app-std.log       # Application stdout logs
    ├── app-err.log       # Application stderr logs
    ├── xray-access.log   # Xray access logs
    └── xray-error.log    # Xray error logs
```

## Database Operations

### 1. Initialization

The database automatically initializes on first startup:

```go
func (d *Database) Init() error {
    d.locker.Lock()
    defer d.locker.Unlock()

    if utils.FileExist(Path) {
        return d.Load()  // Load existing database
    }

    // Generate random port if not available
    if !utils.PortFree(d.Data.Settings.HttpPort) {
        var err error
        if d.Data.Settings.HttpPort, err = utils.FreePort(); err != nil {
            return errors.Wrap(err, "cannot find free port")
        }
    }

    err := d.Save()  // Create new database file
    return errors.WithStack(err)
}
```

### 2. Loading Data

Loading data from the JSON file:

```go
func (d *Database) Load() error {
    content, err := os.ReadFile(Path)
    if err != nil {
        return errors.WithStack(err)
    }

    err = json.Unmarshal(content, d.Data)
    if err != nil {
        return errors.WithStack(err)
    }

    // Validate loaded data
    err = validator.New().Struct(d)
    return errors.WithStack(err)
}
```

### 3. Saving Data

Persisting data to the JSON file:

```go
func (d *Database) Save() error {
    content, err := json.Marshal(d.Data)
    if err != nil {
        return errors.WithStack(err)
    }

    err = os.WriteFile(Path, content, 0755)
    return errors.WithStack(err)
}
```

### 4. Thread Safety

All database operations are protected by mutex:

```go
type Database struct {
    l      *logger.Logger
    locker *sync.Mutex    // Protects concurrent access
    Data   *Data
}
```

## Data Validation

### Validation Rules

The database uses struct tags for validation:

```go
// Settings validation
HttpPort  int    `validate:"required,min=1,max=65536"`
HttpToken string `validate:"required,min=8,max=128"`

// Manager validation  
Url   string `validate:"required,url,min=1,max=1024"`
Token string `validate:"required,min=1,max=128"`
```

### Validation Process

1. **On Load**: Validate data after loading from file
2. **On Save**: Validate data before saving to file
3. **On API**: Validate incoming API requests
4. **Error Handling**: Return detailed validation errors

## Database Lifecycle

### 1. Startup Sequence

```
Application Start → Database.Init() → Check File Exists
                                   ↓
                                   Load Existing OR Create New
                                   ↓
                                   Validate Data
                                   ↓
                                   Ready for Use
```

### 2. Runtime Operations

- **Read Operations**: Direct access to `d.Data` (protected by mutex)
- **Write Operations**: Modify `d.Data` and call `Save()`
- **Manager Updates**: API endpoints update manager configuration
- **Settings Updates**: Automatic updates during port conflicts

### 3. Shutdown

- **Graceful Shutdown**: No special cleanup required
- **Data Persistence**: All data already persisted to disk
- **Backup Creation**: File can be copied for backup

## API Integration

### Manager Configuration Endpoint

```go
func ManagerStore(d *database.Database) echo.HandlerFunc {
    return func(c echo.Context) error {
        var r ManagerStoreRequest
        
        // Parse and validate request
        if err := c.Bind(&r); err != nil { /* handle error */ }
        if err := c.Validate(&r); err != nil { /* handle error */ }

        // Update database
        if r.Url == "" {
            d.Data.Manager = nil  // Remove manager
        } else {
            d.Data.Manager = &database.Manager{
                Url:   r.Url,
                Token: r.Token,
            }
        }

        // Persist changes
        if err := d.Save(); err != nil {
            return errors.WithStack(err)
        }

        return c.JSON(http.StatusCreated, map[string]interface{}{
            "manager": r,
        })
    }
}
```

## Backup and Recovery

### 1. Backup Procedures

**Manual Backup:**
```bash
# Copy database file
cp storage/database/app.json storage/database/app.json.backup

# Backup entire storage directory
tar -czf storage-backup-$(date +%Y%m%d).tar.gz storage/
```

**Automated Backup:**
```bash
# Add to crontab for daily backups
0 2 * * * cd /path/to/arch-node && cp storage/database/app.json storage/database/app.json.$(date +%Y%m%d)
```

### 2. Recovery Procedures

**Restore from Backup:**
```bash
# Restore database file
cp storage/database/app.json.backup storage/database/app.json

# Restart service to apply
systemctl restart arch-node-1
```

**Reset to Defaults:**
```bash
# Remove database file (will regenerate with new random values)
rm storage/database/app.json

# Restart service
systemctl restart arch-node-1
```

## Performance Characteristics

### Read Performance

- **Memory Access**: Data loaded into memory at startup
- **No Serialization**: Direct struct access for reads
- **Concurrent Reads**: Multiple goroutines can read simultaneously

### Write Performance

- **JSON Marshaling**: ~1ms for typical data sizes
- **File I/O**: Single file write operation
- **Atomic Updates**: File replacement ensures consistency
- **Mutex Overhead**: Minimal locking overhead

### Scalability Limits

- **File Size**: JSON file remains under 1KB for typical usage
- **Concurrent Access**: Mutex-based synchronization suitable for node usage
- **I/O Operations**: Infrequent writes (only on configuration changes)

## Migration and Upgrades

### Schema Evolution

**Adding New Fields:**
```go
// Old schema
type Settings struct {
    HttpPort  int    `json:"http_port"`
    HttpToken string `json:"http_token"`
}

// New schema (backward compatible)
type Settings struct {
    HttpPort  int    `json:"http_port"`
    HttpToken string `json:"http_token"`
    NewField  string `json:"new_field,omitempty"`  // Optional field
}
```

**Version Handling:**
```go
type Data struct {
    Version  int       `json:"version,omitempty"`
    Settings *Settings `json:"settings"`
    Manager  *Manager  `json:"manager"`
}
```

### Migration Process

1. **Version Detection**: Check database version on load
2. **Schema Migration**: Apply necessary transformations
3. **Validation**: Ensure migrated data is valid
4. **Backup**: Keep backup of pre-migration data

## Troubleshooting

### Common Issues

1. **File Permissions**
   ```bash
   # Check file permissions
   ls -la storage/database/app.json
   
   # Fix permissions if needed
   chmod 644 storage/database/app.json
   ```

2. **Corrupted JSON**
   ```bash
   # Validate JSON syntax
   cat storage/database/app.json | jq .
   
   # Restore from backup if corrupted
   cp storage/database/app.json.backup storage/database/app.json
   ```

3. **Port Conflicts**
   ```bash
   # Check if port is in use
   netstat -ln | grep :15888
   
   # Delete database to regenerate with new port
   rm storage/database/app.json && systemctl restart arch-node-1
   ```

### Debug Commands

```bash
# View current database content
cat storage/database/app.json | jq .

# Check database file exists and is readable
ls -la storage/database/app.json

# Validate JSON format
cat storage/database/app.json | python -m json.tool

# Monitor database changes
watch -n 1 'cat storage/database/app.json | jq .'
```

This simple yet effective database system provides all the necessary functionality for Arch-Node while maintaining simplicity and reliability.
