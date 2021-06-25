FROM golang:1.15-alpine AS build_base

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test ./... -v

# Build the Go app
RUN go build -o ./bin/webscraping_go cmd/main.go

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_base /tmp/go-app/bin/webscraping_go /app/webscraping_go
COPY --from=build_base /tmp/go-app/asset /asset

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/webscraping_go"]