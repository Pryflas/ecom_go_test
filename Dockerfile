# Этап сборки
FROM golang:1.25 as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ecom_go_test .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/ecom_go_test .

EXPOSE 8080

CMD ["./ecom_go_test"]

