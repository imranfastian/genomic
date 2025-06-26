# Stage 1: Build Go app
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy all Go code
COPY . .

# Download dependencies
RUN go mod init genomic && go mod tidy

# Build binary
RUN go build -o genomic

# Stage 2: Run lightweight container
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/genomic .
EXPOSE 8080
ENTRYPOINT ["./genomic"]
