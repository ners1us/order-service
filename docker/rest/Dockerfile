FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/rest-app
RUN CGO_ENABLED=0 go build -o rest-app .

FROM alpine:latest
EXPOSE 8080
WORKDIR /root
COPY --from=build /app/cmd/rest-app/rest-app .
CMD ["./rest-app"]