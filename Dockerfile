# First stage: build the app
FROM golang:1.9.1 as builder

# Copy source files
WORKDIR /go/src/github.com/bobinette/tonight
COPY . .

# Install dependencies
RUN go get -u -v github.com/kardianos/govendor/...
RUN govendor sync -v

# Build binary
# Needed to be runnable in the alpine:latest container, see https://stackoverflow.com/a/36308464
ENV CGO_ENABLED 0
RUN ["go", "build", "-tags", "netgo", "-a", "-o", "./bin/tonight", "./cmd/main.go"]

# --------—--------—--------—--------—--------—--------—--------—--------—--------—
# Second stage: for a lighter container
FROM alpine:latest

# Update certificates
RUN apk --no-cache add ca-certificates

# Copy binary from builder into /auth
WORKDIR /tonight
COPY --from=builder /go/src/github.com/bobinette/tonight/bin/tonight .

COPY ./bleve/mapping.json ./bleve/mapping.json
COPY ./templates ./templates
COPY ./assets ./assets
COPY ./external ./external
COPY ./fonts ./fonts

# Start the server
CMD ["./tonight"]
