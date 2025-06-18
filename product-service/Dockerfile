FROM golang:1.24 AS builder

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./product-service/main.go


FROM alpine:latest

COPY --from=builder /server /server

EXPOSE 50052 8081

ENTRYPOINT ["/server"]
