package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ners1us/order-service/internal/middleware"
	"github.com/ners1us/order-service/internal/services"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	server           *http.Server
	engine           *gin.Engine
	userHandler      UserHandler
	pvzHandler       PVZHandler
	receptionHandler ReceptionHandler
	productHandler   ProductHandler
	jwtService       services.JWTService
}

func NewHTTPServer(
	port string,
	userHandler UserHandler,
	pvzHandler PVZHandler,
	receptionHandler ReceptionHandler,
	productHandler ProductHandler,
	jwtService services.JWTService,
) *HTTPServer {
	r := gin.Default()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &HTTPServer{
		server:           srv,
		engine:           r,
		userHandler:      userHandler,
		pvzHandler:       pvzHandler,
		receptionHandler: receptionHandler,
		productHandler:   productHandler,
		jwtService:       jwtService,
	}
}

func (hs *HTTPServer) ConfigureRoutes() {
	hs.engine.POST("/dummyLogin", hs.userHandler.DummyLogin)
	hs.engine.POST("/register", hs.userHandler.Register)
	hs.engine.POST("/login", hs.userHandler.Login)

	secured := hs.engine.Group("/", middleware.AuthMiddleware(hs.jwtService))
	secured.POST("/pvz", hs.pvzHandler.CreatePVZ)
	secured.GET("/pvz", hs.pvzHandler.GetPVZList)
	secured.POST("/pvz/:pvzId/close_last_reception", hs.receptionHandler.CloseLastReception)
	secured.POST("/pvz/:pvzId/delete_last_product", hs.productHandler.DeleteLastProduct)
	secured.POST("/receptions", hs.receptionHandler.CreateReception)
	secured.POST("/products", hs.productHandler.AddProduct)
}

func (hs *HTTPServer) Start() error {
	log.Println("starting HTTP server...")
	return hs.server.ListenAndServe()
}

func (hs *HTTPServer) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := hs.server.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
		return err
	}

	log.Println("HTTP server stopped gracefully")
	return nil
}
