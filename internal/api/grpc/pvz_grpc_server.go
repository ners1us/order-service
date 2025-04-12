package grpc

import (
	"context"
	"github.com/ners1us/order-service/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"

	"github.com/ners1us/order-service/internal/repositories"
	"github.com/ners1us/order-service/pkg/generated/proto"
)

type PVZGrpcServer struct {
	server    *grpc.Server
	pvzServer *services.PVZGrpcService
	listener  net.Listener
}

func NewServer(pvzRepo repositories.PVZRepository, port string) (*PVZGrpcServer, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pvzServer := services.NewPVZGrpcService(pvzRepo)
	proto.RegisterPVZServiceServer(grpcServer, pvzServer)

	return &PVZGrpcServer{
		server:    grpcServer,
		pvzServer: pvzServer,
		listener:  lis,
	}, nil
}

func (pgs *PVZGrpcServer) Start() error {
	log.Println("starting gRPC server...")
	return pgs.server.Serve(pgs.listener)
}

func (pgs *PVZGrpcServer) Stop(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)
		pgs.server.GracefulStop()
	}()

	select {
	case <-doneCh:
		log.Println("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		defer pgs.server.Stop()
		log.Println("gRPC server forced to stop")
	}
}
