package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"shop-backend/config"
	"shop-backend/internal/model"
	"shop-backend/internal/service"
	jwtutil "shop-backend/pkg/jwt"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminHandler struct {
	productService *service.ProductService
	kitService     *service.KitService
}

func NewAdminHandler(productService *service.ProductService, kitService *service.KitService) *AdminHandler {
	return &AdminHandler{
		productService: productService,
		kitService:     kitService,
	}
}

func (h *AdminHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the string ID to MongoDB ObjectID
	objId, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetByIDProduct(r.Context(), objId)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)

}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	cfg := config.LoadConfig()

	// Basic static admin check - could use JWT here
	if creds.Email == cfg.AdminEmail && creds.Password == cfg.AdminPass {
		// Generate JWT token
		token, err := jwtutil.GenerateToken(creds.Email, "admin", cfg.JWTSecret, 24*time.Hour) // optional helper
		if err != nil {
			http.Error(w, "Token generation failed", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,  // Prevent JS access (protects against XSS)
			Secure:   false, // Set to true if using HTTPS
			Path:     "/",
			SameSite: http.SameSiteLaxMode, // Or SameSiteStrictMode
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Login successful",
			"token":   token,
		})

		return
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func (h *AdminHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.ListProducts(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch products"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *AdminHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// 1.Parse multipart form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	// 2.Read form fields
	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")

	// 3.Parse numeric values
	price, _ := strconv.ParseFloat(priceStr, 64)
	stock, _ := strconv.Atoi(stockStr)

	// 4. Read the uploaded file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure uploads/ directory exists
	os.Mkdir("uploads", os.ModePerm)

	// 5. Save the file to the server
	filePath := "uploads/" + handler.Filename
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// 6. Create product model
	product := model.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImageURL:    filePath,
	}

	// 7. Call service layer
	if err := h.productService.CreateProduct(r.Context(), &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product created"))
}

func (h *AdminHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID from string to ObjectID
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Fetch existing product from DB
	existingProduct, err := h.productService.GetByIDProduct(r.Context(), objID)
	if err != nil {
		http.Error(w, "Poduct not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract form fields
	if name := r.FormValue("name"); name != "" {
		existingProduct.Name = name
	}
	if description := r.FormValue("description"); description != "" {
		existingProduct.Description = description
	}
	if priceStr := r.FormValue("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			existingProduct.Price = price
		}
	}
	if stockStr := r.FormValue("stock"); stockStr != "" {
		if stock, err := strconv.Atoi(stockStr); err == nil {
			existingProduct.Stock = stock
		}
	}

	// Optional image upload
	imagePath := ""
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		// Create uploads directory if not exists
		os.MkdirAll("uploads", os.ModePerm)

		imagePath = "uploads/" + handler.Filename
		dst, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		io.Copy(dst, file)
		existingProduct.ImageURL = imagePath
	}

	// Call service method to update product
	err = h.productService.UpdateProduct(r.Context(), existingProduct)
	if err != nil {
		http.Error(w, "Failed to update product: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("product Updated"))
}

func (h *AdminHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL parameter
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert string to ObjecID
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// call service method to delete product
	if err := h.productService.DeleteProduct(r.Context(), objID); err != nil {
		http.Error(w, "Failed to delete product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Product deleted"))
}

func (h *AdminHandler) CreateKit(w http.ResponseWriter, r *http.Request) {
	// Parse multipart from for image and fields
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Extract form values
	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	productIDsRaw := r.FormValue("product_ids") // this will be multiple string values
	productIDsStr := strings.Split(productIDsRaw, ",")

	// convert price to float64
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Failed to extract kit image", http.StatusBadRequest)
		return
	}

	// convert product_ids to ObjecIDs
	var productIDs []primitive.ObjectID
	for _, idStr := range productIDsStr {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue // skip empty str
		}
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			http.Error(w, "Invalid product ID:"+idStr, http.StatusBadRequest)
			return
		}
		productIDs = append(productIDs, id)
	}

	// Handle image file upload
	imagePath := ""
	file, handler, err := r.FormFile("image")
	if err != nil {
		defer file.Close()
		os.MkdirAll("uploads/kits", os.ModePerm)
		imagePath = "uploads/kits/" + handler.Filename

		dst, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		io.Copy(dst, file)
	}

	// Create Kit model
	kit := &model.Kit{
		Name:        name,
		Description: description,
		ProductIDs:  productIDs,
		Price:       price,
		ImageURL:    imagePath,
	}

	// Call service
	err = h.kitService.CreateKit(r.Context(), kit)
	if err != nil {
		http.Error(w, "Failed to create kit", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Kit created"))

}
