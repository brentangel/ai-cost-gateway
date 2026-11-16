FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ai-gateway main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ai-gateway .
EXPOSE 8888
CMD ["./ai-gateway"]
