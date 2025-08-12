# URL Shortener (Go)

A minimal URL shortening service written in Go.  
Built as a learning project to explore **Go**, **HTTP servers**, **clean architecture**, and scalable project structure.

## Features

- `/health` endpoint for quick status checks
- Project layout following Go conventions (`cmd/`, `internal/`, `pkg/`)
- Graceful shutdown on interrupt signals
- Basic logging middleware

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- `git`

### Run locally

```bash
git clone https://github.com/carelessly3/url-shortener.git
cd url-shortener
go mod tidy
go run ./cmd/server
```
