package product

import (
	"encoding/json"
	"net/http"
	"server/handlers"
	"server/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	createProduct = "/api/createProduct"
	getProducts   = "/api/products"
)

type handler struct {
	logger  *logging.Logger
	service *Service
}

func NewHandler(logger *logging.Logger, service *Service) handlers.Handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.POST(createProduct, h.CreateProduct)
	router.HandlerFunc(http.MethodGet, getProducts, h.GetAllProducts)
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	productID, err := h.service.Create(r.Context(), product)
	if err != nil {
		h.logger.Error("Failed to create product", err)
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": productID})
}

func (h *handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("Failed to get products", err)
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
