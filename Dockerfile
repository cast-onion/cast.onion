FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cast-onion ./cmd/server

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/cast-onion .
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations
EXPOSE 8443
ENTRYPOINT ["./cast-onion"]
