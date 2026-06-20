# SwalpaURL Configuration

## Build & Run Commands
- Run application: `go run main.go`
- Run local tests: `go test ./...`
- Spin up infrastructure: `docker-compose up -d`

## System Context
- Backend: Go (Golang) 
- Storage: PostgreSQL (Persistent storage) & Redis (Caching/Key Generation Service) both deployed using Docker

## Coding Standards
- Handle all errors explicitly; do not drop errors using blank identifiers (`_`).
- Maintain clean hexagonal architecture separation between handlers and repositories.

## Deployment(to do)
- Currently everything runs in docker locally. 
- There is an nginx configuration file, need to get ssl done for this tho
- deployment stack is docker, compose and env variables are stored in a .env file 
- constraint is to maintain network isolation between backend services. 


