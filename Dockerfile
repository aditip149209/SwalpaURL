# --- Stage 1: High-Performance Build Layer ---
FROM golang:1.25-alpine AS builder

# Create an isolated construction workspace inside the container
WORKDIR /workspace

# Copy your dependency sheets first to leverage Docker's layer caching mechanism
COPY go.mod go.sum ./
RUN go mod download

# Copy your entire internal project tree structure into the workspace
COPY . .

# Compile a statically linked, self-contained binary file
# Turning off CGO ensures complete compatibility with Alpine's minimal runtime footprint
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o swalpaurl ./cmd/server/main.go


# --- Stage 2: Ultra-Lightweight Runtime Layer ---
FROM alpine:3.20

# Set a clean working directory for your application engine root
WORKDIR /app

# Install root CA security certificates so your app can safely perform outgoing HTTPS calls
RUN apk --no-cache add ca-certificates

# Copy ONLY the final compiled binary executable out of the builder construction site
# Your templates and wordlists are already baked straight inside this binary file!
COPY --from=builder /workspace/swalpaurl .

# Inform Docker that this container application layer listens on network port 8080
EXPOSE 8080

# Execute the self-sufficient engine program on container boot
CMD ["./swalpaurl"]