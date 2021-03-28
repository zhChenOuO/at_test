
FROM golang:latest AS builder
WORKDIR /server
ENV GO111MODULE=on 

# COPY go.mod . 
# COPY go.sum .
# RUN go mod download
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build  -mod=vendor -o main

FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk add --no-cache curl tzdata && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /server
COPY --from=builder /server/main /server/main

RUN ls

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /server/main
USER appuser


