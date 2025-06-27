package handler

import (
	"encoding/json"
	"net/http"

	"shop-backend/internal/model"
	"shop-backend/internal/service"
)

type UserHandler struct {
	AuthService    *service.AuthService
	ProductService *service.ProductService
	OrderService   *service.OrderService
}

func NewUserHandler(auth *service.AuthService, product *service.ProductService, order *service.OrderService) *UserHandler {
	return &UserHandler{
		AuthService:    auth,
		ProductService: product,
		OrderService:   order,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	if err := h.AuthService.Register(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	token, err := h.AuthService.Login(r.Context(), creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string `json:"email"`
		Otp   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := h.AuthService.VerifyOtp(r.Context(), data.Email, data.Otp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "OTP verified successfully"})
}

type resendOtpRequest struct {
	Email string `json:"email"`
}

func (h *UserHandler) ResendOtp(w http.ResponseWriter, r *http.Request) {
	var req resendOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.AuthService.ResendOtp(ctx, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OTP resent successfully"))
}

//
// func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
//Implemen JWT blacklist logic if needed
// w.Write([]byte("Logged out"))
// }
//
// func (h *UserHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
// products, err := h.ProductService.ListProducts(r.Context())
// if err != nil {
// http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
// return
// }
// json.NewEncoder(w).Encode(products)
// }
//
// func (h *UserHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
//You would extract user from context, then process cart logic
// w.Write([]byte("Item added to cart"))
// }
//
// func (h *UserHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
// w.Write([]byte("Item removed from cart"))
// }
//
// func (h *UserHandler) Checkout(w http.ResponseWriter, r *http.Request) {
// w.Write([]byte("Order placed"))
// }
//
// func (h *UserHandler) OrderHistory(w http.ResponseWriter, r *http.Request) {
// orders := []string{"Order 1", "Order 2"} // placeholder
// json.NewEncoder(w).Encode(orders)
// }
//
