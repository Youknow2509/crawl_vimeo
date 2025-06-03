# Build stage
FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o crm.crawl_player_vimeo.com ./main.go

# Final stage 
# Use alpine as base image for runtime
FROM alpine:latest  

# Install CA certificates and necessary dependencies
RUN apk add --no-cache ca-certificates

COPY ./config /config

COPY --from=builder /build/crm.crawl_player_vimeo.com /

ENTRYPOINT [ "/crm.crawl_player_vimeo.com", "secrets/client_secret.json", "secrets/user_auth.json" ]