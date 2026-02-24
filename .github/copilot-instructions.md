# Project: WackyURL (Distributed System Simulation)
**Role:** Senior Backend Architect, Chaos Engineer & Go Expert.
**Goal:** Build a high-scale, fault-tolerant URL shortener that simulates real-world distributed system challenges, served via a hyper-fast server-side rendered UI.

---

# 1. Technology Stack & Standards (STRICT)

### Core Backend
* **Language:** Go (Golang) 1.22+.
* **Web Framework:** Gin (github.com/gin-gonic/gin) or Standard `net/http` with `Chi`.
* **Database:**
    * **Primary:** PostgreSQL (Persist mappings: `wacky_name` -> `long_url`).
    * **Cache/KGS:** Redis (Hot cache + "Pre-generated Keys" Set).
* **Config:** Viper or `os.Getenv` (Strict 12-Factor App compliance).

### Frontend (Server-Side Rendered)
* **Engine:** Go `html/template` (Standard Library).
* **Interactivity:** HTMX (for AJAX-like behavior without writing JS).
* **Styling:** Tailwind CSS (via CDN is acceptable for this simulation).
* **Rule:** NO React, NO Vue, NO Node.js build steps. The Go binary serves the HTML.

### Infrastructure
* **Containerization:** Docker (Multi-stage builds, `scratch` or `alpine` base images).
* **Orchestration:** Kubernetes (K3s for local, EKS for AWS).
* **Gateway:** Nginx (Reverse Proxy/Ingress).

### Testing & Simulation Tools
* **Integration:** Testcontainers-Go (Real Redis/Postgres instances).
* **Load Testing:** K6 (JavaScript-based load scripts).
* **Network Chaos:** Toxiproxy (Simulate latency/packet loss).

---

# 2. Functional Requirements (Features)

### A. The "Wacky" Name Service (KGS)
* **Constraint:** Do NOT use Base62 encoding/random hashes.
* **Logic:** Generate `Adjective-Noun` pairs (e.g., `Grumpy-Pink-Falcon`).
* **Uniqueness:** Must be collision-proof.
    * *Strategy:* Use a "Pre-generation Worker" that fills a Redis Set (`SADD available_names`).
    * *Assignment:* Use `SPOP` to atomically claim a name from Redis.
    * *Refill:* If Redis count < 1000, trigger async worker to generate more.

### B. URL Shortening & Redirection
* **Shorten (POST):** Accepts Long URL -> Returns Wacky URL.
* **Redirect (GET):** Accepts Wacky URL -> 302 Redirects to Long URL.
    * *Optimization:* Cache-Aside pattern. Check Redis first. If miss, DB -> Redis -> Redirect.
    * *Analytics:* Async increment of "click count" (do not block the redirect response).

### C. User Interface (Go + HTMX)
* **Landing Page:** Single input field.
* **UX Flow:**
    1.  User types URL.
    2.  HTMX sends `POST /shorten`.
    3.  Server returns HTML snippet (`<div id="result">...</div>`).
    4.  HTMX swaps the input form with the result snippet.
* **Constraint:** No full page reloads.

---

# 3. Non-Functional Requirements (System Design)

### A. Performance & Scalability
* **Latency Target:** * Redirect (Cache Hit): < 20ms.
    * Shorten (Write): < 100ms.
* **Throughput:** System must withstand 1000 Requests Per Second (RPS) in load tests.
* **Statelessness:** The Go binary must hold NO state. Session/Cache must be in Redis.

### B. Reliability & Chaos
* **Graceful Degradation:** If Redis is down, the system should fallback to PostgreSQL (slower, but working).
* **Geo-Simulation:**
    * Middleware must read `X-Simulated-Region` header.
    * If header is distinct from server region, inject `time.Sleep` (e.g., 200ms) to mimic speed of light latency.

### C. Observability
* **Structured Logging:** Use `slog` (Go 1.21+) or `Zap`. output JSON.
* **Tracing:** Every request must have a `trace_id` propagated to DB logs and error messages.
* **Health Checks:**
    * `/health/live`: Returns 200 OK (Server is running).
    * `/health/ready`: Returns 200 OK ONLY if Redis + Postgres are connected.

---

# 4. Coding Best Practices (Go Specific)

### Code Structure
* Follow "Standard Go Project Layout":
    * `/cmd`: Entry point.
    * `/internal`: Private app code (handlers, services).
    * `/pkg`: Public library code (wacky logic).
    * `/templates`: HTML files.

### Error Handling
* **NEVER** ignore errors.
* **NEVER** panic (except on startup).
* Wrap errors with context: `fmt.Errorf("failed to claim wacky name: %w", err)`.

### Testing Guidelines
* **Table-Driven Tests:** Use them for all Unit Tests (Wacky logic).
* **Integration Tests:**
    * Must use `testcontainers-go`.
    * **Race Condition Test:** Spawn 50 goroutines trying to shorten the *same* URL simultaneously. Assert only 1 DB write occurs.
* **Load Tests (K6):**
    * Create scripts for "Viral Traffic" (Hot Key) and "Global Traffic" (High Latency).

---

# 5. AWS Deployment Guidelines
* **Terraform/IaC:** Infrastructure must be defined as code.
* **Security:**
    * App runs in Private Subnet.
    * Load Balancer in Public Subnet.
    * Secrets (DB Passwords) injected via AWS Secrets Manager or Environment Variables.