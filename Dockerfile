# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder

RUN addgroup -g 1000 -S appgroup && \
  adduser -u 1000 -S appuser -G appgroup

RUN mkdir app

# Install certificates
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

RUN chown -R appuser:appgroup /app

FROM scratch 

# Copy certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy user permissions from builder
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/main .

USER 1000

CMD ["./main"]