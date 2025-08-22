# Arch-Node Documentation

Welcome to the Arch-Node documentation. This project is a high-performance, lightweight multi-core proxy node built with Go that integrates with the Arch-Manager ecosystem for distributed proxy management.

## Table of Contents

- [Architecture Overview](./architecture.md)
- [Node-Manager Communication](./node-manager-communication.md)
- [Database System](./database.md)
- [Configuration Guide](./configuration.md)
- [API Reference](./api-reference.md)
- [Deployment Guide](./deployment.md)
- [Development Guide](./development.md)
- [Xray Integration](./xray-integration.md)
- [Troubleshooting](./troubleshooting.md)

## Quick Links

- **Architecture**: Understanding the core components and how they work together
- **Node-Manager Communication**: How nodes connect to and communicate with the Arch-Manager
- **Database**: Understanding the JSON-based persistence system
- **API**: REST API endpoints for management and monitoring
- **Deployment**: Docker and systemd deployment options

## Key Features

- **Multi-Core Proxy Support**: Currently supports Xray-core with extensible architecture
- **Distributed Management**: Seamless integration with Arch-Manager
- **High Performance**: Built with Go for optimal speed and resource efficiency
- **Real-time Monitoring**: Built-in HTTP API for status and statistics
- **Secure Communication**: Token-based authentication and TLS encryption
- **Docker Support**: Ready-to-use containers and compose files

## Getting Started

1. Review the [Architecture Overview](./architecture.md) to understand the system
2. Follow the [Deployment Guide](./deployment.md) for installation
3. Configure your node using the [Configuration Guide](./configuration.md)
4. Connect to an Arch-Manager using [Node-Manager Communication](./node-manager-communication.md)

## Project Structure

```
arch-node/
├── cmd/                    # CLI commands (Cobra)
├── configs/               # Configuration files
├── docs/                  # Documentation (this directory)
├── internal/              # Private application code
│   ├── app/              # Application orchestration
│   ├── config/           # Configuration management
│   ├── coordinator/      # Manager synchronization
│   ├── database/         # JSON file-based persistence
│   └── http/             # HTTP server and handlers
├── pkg/                   # Public/reusable packages
│   ├── http/             # HTTP client and middleware
│   ├── logger/           # Structured logging
│   ├── worker/           # Background workers
│   └── xray/             # Xray-core integration
├── scripts/              # Setup and management scripts
├── storage/              # Runtime data and logs
└── third_party/          # External binaries (Xray)
```

## Support

For issues, questions, or contributions, please refer to the project repository or contact the development team.
