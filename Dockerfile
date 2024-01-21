# Stage 1: Build the application
FROM golang:1.21 AS builder
WORKDIR /app
# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o http-server ./main.go

# Stage 2: Run the application
FROM alpine:latest
WORKDIR /app

# Setting up timezone data
ENV TZ=Etc/UTC
RUN apk --no-cache add tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/http-server .

# Create a directory for resources and copy them
COPY --from=builder /app/resource /app/resource

# Copy other necessary files
COPY --from=builder /app/.env .

RUN ls

# Expose port 8080
EXPOSE 8080
# Run the binary
CMD ["./http-server"]
