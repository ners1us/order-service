package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ners1us/order-service/internal/config"
	"github.com/ners1us/order-service/internal/database"
	"github.com/ners1us/order-service/internal/handlers"
	"github.com/ners1us/order-service/internal/middleware"
	"github.com/ners1us/order-service/internal/repositories"
	"github.com/ners1us/order-service/internal/services"
	"log"
)

func main() {
	cfg := config.InitConfig()

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

	userHandler := handlers.NewUserHandler(userService)
	pvzHandler := handlers.NewPVZHandler(pvzService)
	receptionHandler := handlers.NewReceptionHandler(receptionService)
	productHandler := handlers.NewProductHandler(productService)

	r := gin.Default()

	r.POST("/dummyLogin", userHandler.DummyLogin)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	protectedApi := r.Group("/", middleware.AuthMiddleware(jwtService))
	protectedApi.POST("/pvz", pvzHandler.CreatePVZ)
	protectedApi.GET("/pvz", pvzHandler.GetPVZList)
	protectedApi.POST("/pvz/:pvzId/close_last_reception", receptionHandler.CloseLastReception)
	protectedApi.POST("/pvz/:pvzId/delete_last_product", productHandler.DeleteLastProduct)
	protectedApi.POST("/receptions", receptionHandler.CreateReception)
	protectedApi.POST("/products", productHandler.AddProduct)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("failed running merch store service: ", err)
	}
}
