# syntax=docker/dockerfile:1
FROM golang:1.17 AS builder
WORKDIR /go/src/quenten.nl/pepeservice

COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/quenten.nl/pepeservice/app .
CMD ["./app"]