package server

import (
	"context"
	"github.com/ners1us/order-service/internal/loggers"
	"github.com/ners1us/order-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"

	"github.com/ners1us/order-service/internal/repositories"
	"github.com/ners1us/order-service/pkg/generated/proto"
)

type pvzGrpcServer struct {
	server         *grpc.Server
	pvzGrpcService *service.PVZGrpcService
	listener       net.Listener
	pvzRepo        repositories.PVZRepository
}

func NewServer(
	pvzRepo repositories.PVZRepository,
	port string,
) (BackendServer, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(loggers.GrpcLogger))

	return &pvzGrpcServer{
		server:   grpcServer,
		pvzRepo:  pvzRepo,
		listener: lis,
	}, nil
}

func (pgs *pvzGrpcServer) ConfigureRoutes() {
	reflection.Register(pgs.server)
	pgs.pvzGrpcService = service.NewPVZGrpcService(pgs.pvzRepo)
	proto.RegisterPVZServiceServer(pgs.server, pgs.pvzGrpcService)
}

func (pgs *pvzGrpcServer) Start() error {
	log.Println("starting gRPC server...")
	return pgs.server.Serve(pgs.listener)
}

func (pgs *pvzGrpcServer) Stop(ctx context.Context) {
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
