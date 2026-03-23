# Builder
# Uses $BUILDPLATFORM to run the compiler on the native CI runner architecture (fast)
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine3.23 AS builder
WORKDIR /app

# This ARG acts as a gate. If you pass a changing value here, 
# everything below this line is forced to run without cache.
ARG CACHEBUST=1

# Copy all of the application code
COPY . .

# Build the app
RUN go build -o serve main.go

# Run
FROM alpine:3.23

WORKDIR /app

# Copy templates directory
COPY templates/ ./templates

# Copy built binary
COPY --from=builder /app/serve .

EXPOSE 80

CMD ["/app/serve"]