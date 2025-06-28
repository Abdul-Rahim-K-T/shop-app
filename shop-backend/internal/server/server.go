package server

import (
	"net/http"

	"shop-backend/config"
	"shop-backend/pkg/database"

	"shop-backend/internal/handler"
	"shop-backend/internal/repository"
	"shop-backend/internal/routes"
	"shop-backend/internal/service"

	"github.com/gorilla/mux"
)

func NewServer(cfg *config.Config) *http.Server {
	db := database.ConnectDB()

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	kitRepo := repository.NewKitRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	productService := service.NewProductService(productRepo)
	kitService := service.NewKitService(kitRepo)
	orderService := service.NewOrderService(orderRepo)

	userHandler := handler.NewUserHandler(authService, productService, orderService)
	adminHandler := handler.NewAdminHandler(productService, kitService)

	router := mux.NewRouter()

	// Register routes
	routes.RegisterUserRoutes(router, userHandler)
	routes.RegisterAdminRoutes(router, adminHandler)

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
}
