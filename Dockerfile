# Build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /src/genote-watcher
COPY . .
ENV CGO_ENABLED=0
RUN go get -d -v ./...
RUN go build -o /bin/genote-watcher -v

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/genote-watcher /bin/genote-watcher/app
ENTRYPOINT [ "tail", "-f", "/dev/null" ]
LABEL Name=genotewatcher