package routes

import (
	"shop-backend/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router, h *handler.UserHandler) {
	user := r.PathPrefix("/api/user").Subrouter()

	// Public routes
	user.HandleFunc("/register", h.Register).Methods("POST")
	user.HandleFunc("/login", h.Login).Methods("POST")
	user.HandleFunc("/verifyotp", h.VerifyOtp).Methods("POST")
	user.HandleFunc("/resendotp", h.ResendOtp).Methods("POST")
	// user.HandleFunc("/products", h.ListProducts).Methods("GET")

	//Potected routes (apply middleware to subrouter)
	// protected := user.NewRoute().Subrouter()
	// protected.Use(middleware.AuthMiddleware)

	// protected.HandleFunc("/logout", h.Logout).Methods("POST")
	//
	// protected.HandleFunc("/cart/add", h.AddToCart).Methods("POST")
	// protected.HandleFunc("/cart/remove", h.RemoveFromCart).Methods("DELETE")
	// protected.HandleFunc("/checkout", h.Checkout).Methods("POST")
	// protected.HandleFunc("/orders", h.OrderHistory).Methods("GET")
}
