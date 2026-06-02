FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY web/ ./web/
COPY sql/ ./sql/
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# --- RELEASE STAGE ---
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server ./server
COPY --from=builder /app/web ./web
COPY --from=builder /app/sql ./sql
EXPOSE 8080
ENTRYPOINT ["/app/server"]
