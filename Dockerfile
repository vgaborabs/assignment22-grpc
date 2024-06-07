FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o users cmd/main.go

FROM alpine:latest
COPY --from=builder /app/users /users
CMD ["./users"]
