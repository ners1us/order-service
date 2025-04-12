FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/order-service
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:latest
EXPOSE 8080
EXPOSE 3000
WORKDIR /root
COPY --from=build /app/cmd/order-service/main .
CMD ["./main"]
