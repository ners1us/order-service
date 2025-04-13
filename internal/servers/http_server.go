package servers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ners1us/order-service/internal/api/rest"
	"github.com/ners1us/order-service/internal/middleware"
	"github.com/ners1us/order-service/internal/services"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"log"
	"net/http"
	"time"
)

type httpServer struct {
	server           *http.Server
	engine           *gin.Engine
	userHandler      rest.UserHandler
	pvzHandler       rest.PVZHandler
	receptionHandler rest.ReceptionHandler
	productHandler   rest.ProductHandler
	jwtService       services.JWTService
}

func NewHTTPServer(
	port string,
	userHandler rest.UserHandler,
	pvzHandler rest.PVZHandler,
	receptionHandler rest.ReceptionHandler,
	productHandler rest.ProductHandler,
	jwtService services.JWTService,
) BackendServer {
	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &httpServer{
		server:           srv,
		engine:           r,
		userHandler:      userHandler,
		pvzHandler:       pvzHandler,
		receptionHandler: receptionHandler,
		productHandler:   productHandler,
		jwtService:       jwtService,
	}
}

func (hs *httpServer) ConfigureRoutes() {
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

func (hs *httpServer) Start() error {
	log.Println("starting HTTP server...")
	return hs.server.ListenAndServe()
}

func (hs *httpServer) Stop(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		err := hs.server.Shutdown(shutdownCtx)
		errCh <- err
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("HTTP server shutdown error: %v\n", err)
		} else {
			log.Println("HTTP server stopped gracefully")
		}
	case <-shutdownCtx.Done():
		log.Println("HTTP server forced to stop")
	}
}
