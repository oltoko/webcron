FROM golang:1 AS builder
WORKDIR /src
COPY . .
RUN go build -o webcron

FROM alpine:3
COPY --from=builder /src/webcron /app/webcron
ENTRYPOINT ["/app/webcron"]
