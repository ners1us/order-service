FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/rest-app
RUN CGO_ENABLED=0 go build -o rest-app .
WORKDIR /app/cmd/grpc-app
RUN CGO_ENABLED=0 go build -o grpc-app .

FROM alpine:latest
EXPOSE 8080
EXPOSE 3000
WORKDIR /root
COPY --from=build /app/cmd/rest-app/rest-app .
COPY --from=build /app/cmd/grpc-app/grpc-app .
CMD ["sh", "-c", "./rest-app & ./grpc-app"]