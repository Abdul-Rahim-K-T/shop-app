package routes

import (
	"shop-backend/internal/handler"
	"shop-backend/internal/middleware"

	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(r *mux.Router, h *handler.AdminHandler) {
	admin := r.PathPrefix("/api/admin").Subrouter()

	// Public admin login
	admin.HandleFunc("/login", h.Login).Methods("POST")

	// Protected admin routes
	protected := admin.PathPrefix("").Subrouter()
	protected.Use(middleware.AdminMiddleware)

	protected.HandleFunc("/products/{id}", h.GetProductByID).Methods("GET")
	protected.HandleFunc("/products", h.ListProducts).Methods("GET")
	protected.HandleFunc("/products", h.CreateProduct).Methods("POST")
	protected.HandleFunc("/products/{id}", h.UpdateProduct).Methods("PUT")
	protected.HandleFunc("/products/{id}", h.DeleteProduct).Methods("DELETE")
}
