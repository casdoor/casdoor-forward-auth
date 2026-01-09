# Casdoor Forward Auth

<p align="center">
  <a href="#badge">
    <img alt="semantic-release" src="https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg">
  </a>
  <a href="https://github.com/casdoor/casdoor-forward-auth/actions/workflows/ci.yml">
    <img alt="GitHub Workflow Status (branch)" src="https://img.shields.io/github/actions/workflow/status/casdoor/casdoor-forward-auth/ci.yml?branch=master">
  </a>
  <a href="https://github.com/casdoor/casdoor-forward-auth/releases/latest">
    <img alt="GitHub Release" src="https://img.shields.io/github/v/release/casdoor/casdoor-forward-auth.svg">
  </a>
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/casdoor/casdoor-forward-auth">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/casdoor/casdoor-forward-auth?style=flat-square">
  </a>
  <a href="https://github.com/casdoor/casdoor-forward-auth/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/casdoor/casdoor-forward-auth?style=flat-square" alt="license">
  </a>
  <a href="https://github.com/casdoor/casdoor-forward-auth/issues">
    <img alt="GitHub issues" src="https://img.shields.io/github/issues/casdoor/casdoor-forward-auth?style=flat-square">
  </a>
  <a href="#">
    <img alt="GitHub stars" src="https://img.shields.io/github/stars/casdoor/casdoor-forward-auth?style=flat-square">
  </a>
  <a href="https://github.com/casdoor/casdoor-forward-auth/network">
    <img alt="GitHub forks" src="https://img.shields.io/github/forks/casdoor/casdoor-forward-auth?style=flat-square">
  </a>
  <a href="https://discord.gg/5rPsrAzK7S">
    <img alt="Casdoor" src="https://img.shields.io/discord/1022748306096537660?style=flat-square&logo=discord&label=discord&color=5865F2">
  </a>
</p>

A standalone forward authentication service that integrates [Casdoor](https://casdoor.org/) authentication for protecting HTTP services behind reverse proxies like Traefik, Nginx, or Caddy. This service provides seamless SSO (Single Sign-On) capabilities without requiring changes to your backend services.

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Configuration](#-configuration)
- [API Endpoints](#-api-endpoints)
- [Deployment](#-deployment)
- [Development](#-development)
- [Testing](#-testing)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

- **Zero Backend Changes**: Add authentication to any HTTP service without modifying your application code
- **Seamless SSO**: Integrate with Casdoor for centralized authentication and user management
- **Reverse Proxy Compatible**: Works with Traefik, Nginx, Caddy, and other reverse proxies that support forward auth
- **Session Management**: Automatic session handling with secure cookies
- **OAuth 2.0 Flow**: Complete OAuth 2.0 implementation
- **Stateless Architecture**: In-memory state storage for scalability
- **Kubernetes Ready**: Easy deployment in containerized environments

## 🏗 Architecture

This service acts as a forward authentication handler that sits between your reverse proxy and backend services:

```
┌─────────┐      ┌──────────────┐      ┌───────────────┐      ┌─────────┐
│ Client  │─────▶│ Reverse Proxy│─────▶│ Forward Auth  │─────▶│ Casdoor │
│         │      │  (Traefik/   │      │   Service     │      │  Server │
│         │      │   Nginx)     │      └───────────────┘      └─────────┘
└─────────┘      └──────────────┘             │
                        │                     │
                        └─────────────────────┘
                         Authentication Flow
```

### How It Works

1. **Initial Request**: Client requests a protected resource
2. **Forward Auth Check**: Reverse proxy forwards request to this service
3. **Session Validation**: Service checks for valid authentication cookies
4. **Redirect to Login**: If no valid session, service returns redirect to Casdoor
5. **OAuth Flow**: User authenticates with Casdoor
6. **Callback Handling**: Service handles OAuth callback and sets session cookies
7. **Access Granted**: Subsequent requests pass authentication and proceed to backend

## 📦 Installation

### Prerequisites

- Go 1.23.0 or higher
- A running [Casdoor](https://casdoor.org/) instance
- A reverse proxy (Traefik, Nginx, Caddy, etc.)

### From Source

```bash
git clone https://github.com/casdoor/casdoor-forward-auth.git
cd casdoor-forward-auth
go build -o casdoor-forward-auth ./cmd/main.go
```

### Using Go Install

```bash
go install github.com/casdoor/casdoor-forward-auth/cmd@latest
```

## 🚀 Quick Start

### Step 1: Configure Casdoor

1. Access your Casdoor admin panel
2. Create a new application for forward authentication
3. Configure the redirect URL: `http://your-auth-service:9999/callback`
4. Note down the following details:
   - **Client ID**
   - **Client Secret**
   - **Organization Name**
   - **Application Name**

For detailed instructions, see [Casdoor Application Configuration](https://casdoor.org/docs/application/config/).

### Step 2: Create Configuration File

Create a configuration file `conf/config.json`:

```json
{
  "casdoorEndpoint": "http://localhost:8000",
  "casdoorClientId": "YOUR_CLIENT_ID",
  "casdoorClientSecret": "YOUR_CLIENT_SECRET",
  "casdoorOrganization": "YOUR_ORGANIZATION",
  "casdoorApplication": "YOUR_APPLICATION",
  "pluginEndpoint": "http://localhost:9999"
}
```

**Configuration Parameters:**

| Parameter | Required | Description |
|-----------|----------|-------------|
| `casdoorEndpoint` | Yes | URL of your Casdoor server |
| `casdoorClientId` | Yes | Client ID from Casdoor application |
| `casdoorClientSecret` | Yes | Client secret from Casdoor application |
| `casdoorOrganization` | Yes | Organization name in Casdoor |
| `casdoorApplication` | Yes | Application name in Casdoor |
| `pluginEndpoint` | Yes | URL where this service is accessible |

### Step 3: Run the Service

```bash
./casdoor-forward-auth -configFile=conf/config.json
```

The service will start on port 9999 by default.

### Step 4: Configure Your Reverse Proxy

#### Traefik Example

```yaml
http:
  routers:
    my-router:
      rule: "Host(`app.example.com`)"
      service: my-service
      middlewares:
        - forward-auth

  middlewares:
    forward-auth:
      forwardAuth:
        address: "http://localhost:9999/auth"
        trustForwardHeader: true

  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://localhost:8080"
```

#### Nginx Example

```nginx
location / {
    auth_request /auth;
    proxy_pass http://backend:8080;
}

location = /auth {
    internal;
    proxy_pass http://localhost:9999/auth;
    proxy_pass_request_body off;
    proxy_set_header Content-Length "";
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-Host $host;
    proxy_set_header X-Forwarded-URI $request_uri;
}
```

## ⚙️ Configuration

### Environment Variables

You can also use environment variables to override configuration:

- `CASDOOR_ENDPOINT`: Casdoor server URL
- `CASDOOR_CLIENT_ID`: Application client ID
- `CASDOOR_CLIENT_SECRET`: Application client secret
- `CASDOOR_ORGANIZATION`: Organization name
- `CASDOOR_APPLICATION`: Application name
- `PLUGIN_ENDPOINT`: Forward auth service URL

### Command-Line Options

```bash
casdoor-forward-auth [options]

Options:
  -configFile string
        path to the config file (default "conf/config.json")
```

## 🔌 API Endpoints

### `GET/POST /auth`

Forward authentication endpoint. Your reverse proxy should forward authentication requests here.

**Response:**
- `200 OK`: Authentication successful (with replacement headers/body)
- `307 Temporary Redirect`: Redirect to Casdoor login

### `GET /callback`

OAuth callback endpoint. Casdoor redirects here after successful authentication.

**Query Parameters:**
- `code`: OAuth authorization code
- `state`: State parameter for CSRF protection

### `GET/POST /test`

Test endpoint for development and debugging.

## 🚢 Deployment

### Docker

Create a `Dockerfile`:

```dockerfile
FROM golang:1.23.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o casdoor-forward-auth ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/casdoor-forward-auth .
COPY --from=builder /app/conf ./conf
EXPOSE 9999
CMD ["./casdoor-forward-auth"]
```

Build and run:

```bash
docker build -t casdoor-forward-auth .
docker run -p 9999:9999 -v $(pwd)/conf:/root/conf casdoor-forward-auth
```

### Kubernetes

Example deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: casdoor-forward-auth
spec:
  replicas: 2
  selector:
    matchLabels:
      app: casdoor-forward-auth
  template:
    metadata:
      labels:
        app: casdoor-forward-auth
    spec:
      containers:
      - name: casdoor-forward-auth
        image: casdoor-forward-auth:latest
        ports:
        - containerPort: 9999
        volumeMounts:
        - name: config
          mountPath: /root/conf
      volumes:
      - name: config
        configMap:
          name: casdoor-config
---
apiVersion: v1
kind: Service
metadata:
  name: casdoor-forward-auth
spec:
  selector:
    app: casdoor-forward-auth
  ports:
  - port: 9999
    targetPort: 9999
```

## 💻 Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/casdoor/casdoor-forward-auth.git
cd casdoor-forward-auth

# Install dependencies
go mod download

# Build
go build -o bin/casdoor-forward-auth ./cmd/main.go

# Run
./bin/casdoor-forward-auth -configFile=conf/config.json
```

### Project Structure

```
casdoor-forward-auth/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   │   ├── config.go
│   │   └── token.pem        # Default certificate
│   ├── handler/             # HTTP request handlers
│   │   ├── casdoor_handler.go
│   │   ├── init.go
│   │   └── util_handler.go
│   └── httpstate/           # Session state management
│       ├── state.go
│       └── state_memory_storage.go
├── conf/
│   └── config.json          # Configuration file
├── .github/
│   └── workflows/
│       └── ci.yml           # CI/CD pipeline
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## 🧪 Testing

### Run All Tests

```bash
go test -v ./...
```

### Run Tests with Coverage

```bash
go test -v -cover ./...
```

### Run Specific Package Tests

```bash
go test -v ./internal/httpstate/
go test -v ./internal/config/
go test -v ./internal/handler/
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes using [Conventional Commits](https://www.conventionalcommits.org/)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Message Format

We use [semantic-release](https://github.com/semantic-release/semantic-release) for automated versioning and releases. Please follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New features (triggers minor version bump)
- `fix:` - Bug fixes (triggers patch version bump)
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks
- `test:` - Test additions or modifications
- `refactor:` - Code refactoring

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Copyright 2026 The Casdoor Authors.

## 🔗 Related Projects

- [Casdoor](https://github.com/casdoor/casdoor) - Open-source Identity and Access Management (IAM) / Single-Sign-On (SSO) platform
- [Casdoor Go SDK](https://github.com/casdoor/casdoor-go-sdk) - Official Go SDK for Casdoor
- [Traefik Casdoor Auth Plugin](https://github.com/casdoor/traefik-casdoor-auth) - Traefik middleware plugin for Casdoor authentication

## 💬 Community

- [Discord](https://discord.gg/5rPsrAzK7S) - Join our Discord community
- [GitHub Discussions](https://github.com/casdoor/casdoor/discussions) - Ask questions and share ideas
- [GitHub Issues](https://github.com/casdoor/casdoor-forward-auth/issues) - Report bugs and request features