# db.xyz - Postgres-as-a-Service

A complete Postgres-as-a-Service platform built on Proxmox LXC containers.

## Architecture

- **API Service** (`/api`) - Go backend with HTTP/JSON API
- **Web Console** (`/web`) - React + TypeScript SPA 
- **CLI Tool** (`/cli`) - Go-based `dbx` command-line interface
- **Migrations** (`/migrations`) - Database schema management
- **Deployments** (`/deployments`) - Infrastructure and deployment configs

## Getting Started

See the [specification document](./db_xyz_paa_s_spec_v_0.md) for full details.

## Development

```bash
# API server
cd api && go run cmd/server/main.go

# Web console  
cd web && npm run dev

# CLI tool
cd cli && go run cmd/dbx/main.go
```