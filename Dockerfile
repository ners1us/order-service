FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/main-app
RUN CGO_ENABLED=0 go build -o main-app .
WORKDIR /app/cmd/grpc-app
RUN CGO_ENABLED=0 go build -o grpc-app .

FROM alpine:latest
EXPOSE 8080
EXPOSE 3000
WORKDIR /root
COPY --from=build /app/cmd/main-app/main-app .
COPY --from=build /app/cmd/grpc-app/grpc-app .
CMD ["sh", "-c", "./order-app & ./grpc-app"]