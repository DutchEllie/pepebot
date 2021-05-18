# syntax=docker/dockerfile:1
FROM golang:1.16.4 AS builder
WORKDIR /go/src/quenten.nl/pepebot/
  
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./discord

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/quenten.nl/pepebot/app .
CMD ["./app"]  