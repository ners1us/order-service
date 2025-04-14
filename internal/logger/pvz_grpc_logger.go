package logger

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	resp, err := handler(ctx, req)
	statusCode := codes.OK
	if err != nil {
		statusCode = status.Code(err)
	}

	duration := time.Since(start)
	log.Printf("[gRPC] %s | %s | %v", statusCode, info.FullMethod, duration)

	return resp, err
}
