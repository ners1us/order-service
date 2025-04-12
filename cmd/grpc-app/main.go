package main

import (
	"context"
	"github.com/ners1us/order-service/internal/api/grpc"
	"github.com/ners1us/order-service/internal/config"
	"github.com/ners1us/order-service/internal/database"
	"github.com/ners1us/order-service/internal/repositories"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.NewConfig()

	db, err := database.NewDB(cfg.DbUrl)
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}
	defer db.Close()

	pvzRepo := repositories.NewPVZRepository(db)

	grpcServer, err := grpc.NewServer(pvzRepo, cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to initialize gRPC server: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	sig := <-sigCh
	log.Printf("gRPC server shutting down. received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.Stop(ctx)
}
