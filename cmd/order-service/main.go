package main

import (
	"context"
	"github.com/ners1us/order-service/internal/api/grpc"
	"github.com/ners1us/order-service/internal/api/rest"
	"github.com/ners1us/order-service/internal/config"
	"github.com/ners1us/order-service/internal/database"
	"github.com/ners1us/order-service/internal/repositories"
	"github.com/ners1us/order-service/internal/services"
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

	userRepo := repositories.NewUserRepository(db)
	pvzRepo := repositories.NewPVZRepository(db)
	receptionRepo := repositories.NewReceptionRepository(db)
	productRepo := repositories.NewProductRepository(db)

	jwtService := services.NewJWTService(cfg.JWTSecret)
	userService := services.NewUserService(userRepo, jwtService)
	pvzService := services.NewPVZService(pvzRepo, receptionRepo, productRepo)
	receptionService := services.NewReceptionService(receptionRepo, pvzRepo)
	productService := services.NewProductService(receptionRepo, productRepo)

	userHandler := rest.NewUserHandler(userService)
	pvzHandler := rest.NewPVZHandler(pvzService)
	receptionHandler := rest.NewReceptionHandler(receptionService)
	productHandler := rest.NewProductHandler(productService)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	grpcServer, err := grpc.NewServer(pvzRepo, cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to initialize gRPC server: %v", err)
	}

	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	httpServer := rest.NewHTTPServer(
		cfg.RestPort,
		userHandler,
		pvzHandler,
		receptionHandler,
		productHandler,
		jwtService,
	)
	httpServer.ConfigureRoutes()

	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	sig := <-sigCh
	log.Printf("server shutting down. received signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.Stop(ctx)

	if err := httpServer.Stop(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("servers stopped gracefully")
}
