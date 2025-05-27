FROM golang:1.22 AS builder

# Set work directory
WORKDIR /app

# Copy and install the dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy content
COPY . .

# Create log directory
RUN mkdir -p logs

# Build the app
RUN go build -o app main.go

FROM debian:bookworm-slim

# Set work directory
WORKDIR /app

# Copy build file
COPY --from=builder /app/app .
COPY --from=builder /app/.env .

# Port
EXPOSE 8000

# Run build
CMD ["./app"]
