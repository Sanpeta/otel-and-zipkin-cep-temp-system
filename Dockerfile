#build stage
FROM golang:1.22.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/

#run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env ../

EXPOSE 8080
CMD ["/app/main"]
