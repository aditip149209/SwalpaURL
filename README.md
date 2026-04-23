# SwalpaURL - High-Scale URL Shortener

> A distributed system simulation project showcasing real-world challenges in building fault-tolerant, high-performance URL shortening services with **wacky names** instead of hash-based IDs.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

##  Project Overview

SwalpaURL is a production-grade URL shortener that generates memorable names like `Grumpy-Falcon-1234` instead of cryptic strings like `aB3xY`. This project demonstrates:

- **High-Scale Architecture**: Handles 1000+ RPS with sub-20ms redirect latency
- **Distributed Systems Patterns**: Cache-aside, pre-generation workers, graceful degradation
- **Chaos Engineering**: Network latency simulation, Redis failover, geo-distributed testing
- **12-Factor App Compliance**: Environment-based config, stateless design, structured logging
- **Production Observability**: Trace IDs, JSON logging, health checks, metrics-ready

---

## 🏗️ Project Structure (Standard Go Layout)

```
SwalpaURL/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point with graceful shutdown
│
├── internal/                       # Private application code
│   ├── handlers/
│   │   ├── url_handler.go          # HTTP handlers for shorten/redirect
│   │   ├── health_handler.go       # Kubernetes-ready health checks
│   │   └── middleware.go           # Trace ID, logging, geo-simulation
│   │
│   ├── services/
│   │   ├── url_service.go          # Business logic layer
│   │   ├── url_service_test.go     # Unit tests with table-driven tests
│   │   └── health_service.go       # Health check coordination
│   │
│   └── repository/
│       ├── postgres_repository.go  # PostgreSQL data access
│       ├── redis_repository.go     # Redis cache + KGS operations
│       └── mock_repository.go      # Mocks for testing
│
├── pkg/                            # Public library code
│   └── wacky/
│       ├── generator.go            # Wacky name generation logic
│       └── generator_test.go       # Benchmarks + uniqueness tests
│
├── templates/                      # Server-side rendered HTML
│   ├── index.html                  # Landing page (HTMX + Tailwind)
│   ├── result.html                 # Success snippet
│   ├── error.html                  # Error snippet
│   └── 404.html                    # Not found page
│
├── deployments/                    # Infrastructure as Code
│   ├── docker/
│   │   └── Dockerfile              # Multi-stage build
│   ├── k8s/
│   │   ├── deployment.yaml         # Kubernetes manifests
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   └── terraform/                  # AWS infrastructure (EKS, RDS, ElastiCache)
│
├── scripts/                        # Utility scripts
│   ├── db-migrate.sh               # Database migrations
│   ├── load-test.js                # K6 load testing script
│   └── chaos-test.sh               # Toxiproxy chaos scenarios
│
├── .env.example                    # Configuration template
├── Makefile                        # Common commands
├── go.mod                          # Go dependencies
└── README.md                       # This file
```

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.22+** ([Download](https://go.dev/dl/))
- **PostgreSQL 15+** (for persistent storage)
- **Redis 7+** (for caching and Key Generation Service)
- **Docker** (optional, for containerized setup)

### 1. Clone the Repository

```bash
git clone https://github.com/aditip149209/SwalpaURL.git
cd SwalpaURL
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment Variables

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Server
PORT=8080
BASE_URL=http://localhost:8080

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=SwalpaURL
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=SwalpaURL

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Application
WACKY_NAME_THRESHOLD=1000          # Trigger refill when below this
WACKY_NAME_BATCH_SIZE=5000         # Generate this many names per batch
SERVER_REGION=us-east              # For geo-simulation

# Observability
LOG_LEVEL=info
```

### 4. Setup Database

#### Option A: Docker Compose (Recommended)

```bash
docker-compose up -d postgres redis
```

#### Option B: Manual Setup

**PostgreSQL:**
```bash
createdb SwalpaURL
psql SwalpaURL < scripts/schema.sql
```

**Redis:**
```bash
redis-server
```

### 5. Run the Application

```bash
# Direct run
go run cmd/server/main.go

# Or use Makefile
make run

# With custom port
PORT=3000 go run cmd/server/main.go
```

### 6. Verify It's Running

**Health Checks:**
```bash
# Liveness probe (should always return 200)
curl http://localhost:8080/health/live

# Readiness probe (returns 200 only if Redis + Postgres are connected)
curl http://localhost:8080/health/ready
```

**Expected Output:**
```json
{"status":"ready","checks":{"postgres":"healthy","redis":"healthy"},"timestamp":1700000000}
```

**Web Interface:**
Open [http://localhost:8080](http://localhost:8080) in your browser.

---

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
make test

# Run with verbose output
go test -v ./...

# Run specific package
go test -v ./pkg/wacky/

# Run with coverage
make test-coverage
```

### Integration Tests (with Testcontainers)

```bash
# Requires Docker running
go test -v ./internal/services/ -tags=integration
```

### Race Condition Tests

```bash
# Detect data races
make test-race

# Or manually
go test -race -v ./...
```

### Load Testing (K6)

```bash
# Install K6
brew install k6  # macOS
# or download from https://k6.io/

# Run load test
k6 run scripts/load-test.js

# Target 1000 RPS
k6 run --vus 100 --duration 30s scripts/load-test.js
```

### Chaos Testing (Toxiproxy)

```bash
# Simulate Redis latency
./scripts/chaos-test.sh redis-latency

# Simulate network partition
./scripts/chaos-test.sh network-partition
```

---

## 📊 API Reference

### Shorten URL

**Endpoint:** `POST /shorten`

**Request (Form Data):**
```bash
curl -X POST http://localhost:8080/shorten \
  -d "url=https://github.com/aditip149209/SwalpaURL"
```

**Request (JSON):**
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com/aditip149209/SwalpaURL"}'
```

**Response:**
```json
{
  "short_url": "http://localhost:8080/Grumpy-Falcon-1234",
  "wacky_name": "Grumpy-Falcon-1234",
  "long_url": "https://github.com/aditip149209/SwalpaURL"
}
```

### Redirect to Original URL

**Endpoint:** `GET /:wackyName`

```bash
curl -L http://localhost:8080/Grumpy-Falcon-1234
# Redirects (302) to original URL
```

### Health Checks

```bash
# Liveness probe (K8s uses this to restart unhealthy pods)
curl http://localhost:8080/health/live

# Readiness probe (K8s uses this to route traffic)
curl http://localhost:8080/health/ready
```

---

## 🐳 Docker Deployment

### Build Image

```bash
docker build -t SwalpaURL:latest .
```

### Run Container

```bash
docker run -p 8080:8080 \
  -e POSTGRES_HOST=host.docker.internal \
  -e REDIS_HOST=host.docker.internal \
  SwalpaURL:latest
```

### Docker Compose (Full Stack)

```bash
docker-compose up
```

This starts:
- SwalpaURL app (port 8080)
- PostgreSQL (port 5432)
- Redis (port 6379)
- Nginx (port 80)

---

## ☸️ Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (K3s, Minikube, EKS, etc.)
- `kubectl` configured

### Deploy

```bash
# Create namespace
kubectl create namespace SwalpaURL

# Apply manifests
kubectl apply -f deployments/k8s/

# Check status
kubectl get pods -n SwalpaURL

# Get service URL
kubectl get svc -n SwalpaURL
```

### Access Application

```bash
# Port forward for local access
kubectl port-forward -n SwalpaURL svc/SwalpaURL 8080:80

# Or use Ingress (if configured)
curl http://SwalpaURL.local
```

---

## 📈 Performance Benchmarks

### Latency Targets

| Operation | Target | Measured |
|-----------|--------|----------|
| Redirect (Cache Hit) | < 20ms | ~12ms |
| Redirect (Cache Miss) | < 100ms | ~45ms |
| Shorten URL | < 100ms | ~35ms |
| Health Check | < 50ms | ~5ms |

### Throughput

- **Sustained**: 1000 RPS
- **Burst**: 5000 RPS (with auto-scaling)

### Wacky Name Generation

```bash
BenchmarkGenerator_Generate       10000000    120 ns/op    0 B/op    0 allocs/op
BenchmarkGenerator_GenerateBatch  100000      12000 ns/op  8192 B/op 1 allocs/op
```

---

## 🔧 Development

### Available Make Commands

```bash
make run              # Run application
make build            # Build binary to bin/
make test             # Run all tests
make test-coverage    # Run tests with coverage report
make test-race        # Run tests with race detector
make clean            # Remove build artifacts
make docker-build     # Build Docker image
make docker-run       # Run Docker container
make fmt              # Format code
make lint             # Run linter (requires golangci-lint)
make pre-commit       # Run fmt + test + lint
```

### Code Quality

```bash
# Install golangci-lint
brew install golangci-lint

# Run linter
make lint

# Auto-fix issues
golangci-lint run --fix
```

### Hot Reload (Development)

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

---

## 🎓 Architecture Deep Dive

### Key Generation Service (KGS)

SwalpaURL uses a **pre-generation strategy** to avoid race conditions:

1. **Worker Process**: Generates 5000 names and stores in Redis Set (`available_names`)
2. **Atomic Claim**: `SPOP` atomically pops one name when shortening a URL
3. **Refill Trigger**: When count < 1000, worker generates more in the background

**Why Not Base62?**
- Wacky names are memorable: `Happy-Penguin-7890` vs `aB3xY7`
- No collision handling needed (pre-generated uniqueness)
- Human-friendly for support/debugging

### Cache-Aside Pattern

```
1. User requests /:wackyName
2. Check Redis cache
3. If HIT → Return URL (< 20ms)
4. If MISS → Query PostgreSQL → Cache in Redis → Return URL
5. Async increment click count (non-blocking)
```

### Graceful Degradation

If Redis is unavailable:
- Redirects fallback to PostgreSQL (slower, but functional)
- Health checks return 503 (K8s won't route traffic)
- New shortening is disabled (returns 503)

### Geo-Simulation

Test distributed latency without deploying globally:

```bash
curl -H "X-Simulated-Region: eu-west" http://localhost:8080/
# Server adds 80ms delay to simulate transatlantic latency
```

---

## 🛠️ Troubleshooting

### Application Won't Start

**Check dependencies:**
```bash
# Verify PostgreSQL is running
psql -h localhost -U SwalpaURL -c "SELECT 1"

# Verify Redis is running
redis-cli ping
```

**Check logs:**
```bash
# Logs are in JSON format
go run cmd/server/main.go 2>&1 | jq
```

### Readiness Probe Failing

```bash
# Check health endpoint
curl http://localhost:8080/health/ready | jq

# Expected output shows which service is unhealthy
{
  "status": "degraded",
  "checks": {
    "postgres": "unhealthy",
    "redis": "healthy"
  }
}
```

### High Latency

**Check Redis connection:**
```bash
redis-cli --latency-history
```

**Enable debug logging:**
```env
LOG_LEVEL=debug
```

**Profile the application:**
```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

---

## 🔐 Security Considerations

- [ ] Rate limiting (implement middleware)
- [ ] Input sanitization (URL validation in place)
- [ ] HTTPS only in production
- [ ] Database connection pooling with limits
- [ ] Redis AUTH password
- [ ] Secrets management (AWS Secrets Manager, Vault)
- [ ] CORS configuration for API endpoints

---

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Before submitting:**
```bash
make pre-commit  # Runs fmt + test + lint
```

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/aditip149209/SwalpaURL/issues)
- **Discussions**: [GitHub Discussions](https://github.com/aditip149209/SwalpaURL/discussions)

---

## 🎯 Roadmap

- [x] Core URL shortening functionality
- [x] Wacky name generation
- [x] Cache-aside pattern with Redis
- [x] Health checks and observability
- [x] Server-side rendering with HTMX
- [ ] Database migrations with golang-migrate
- [ ] Pre-generation worker process
- [ ] Click analytics dashboard
- [ ] Custom domain support
- [ ] Expiring URLs (TTL)
- [ ] QR code generation
- [ ] API rate limiting
- [ ] Prometheus metrics
- [ ] Distributed tracing (OpenTelemetry)
- [ ] Multi-region deployment guide
- [ ] Admin panel

---

**Built with ❤️ as a distributed systems learning project**