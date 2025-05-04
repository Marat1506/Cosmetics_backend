package order

import (
	"encoding/json"
	"net/http"
	"server/handlers"
	"server/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	ordersPath = "/api/orders"
	orderPath  = "/api/orders/:id"
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
	router.POST(ordersPath, h.CreateOrder)
	router.GET(ordersPath+"/user/:userId", h.GetUserOrders)
	router.GET(orderPath, h.GetOrder)
	router.PUT(orderPath+"/status", h.UpdateStatus)
	router.DELETE(orderPath, h.CancelOrder)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var dto CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	orderID, err := h.service.CreateOrder(r.Context(), dto)
	if err != nil {
		h.logger.Error("Failed to create order", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": orderID})
}

func (h *handler) GetUserOrders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userId")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	orders, err := h.service.GetUserOrders(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user orders", err)
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *handler) GetOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orderID := ps.ByName("id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderID)
	if err != nil {
		h.logger.Error("Failed to get order", err)
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *handler) UpdateStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orderID := ps.ByName("id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	var dto UpdateStatusDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateStatus(r.Context(), orderID, dto.Status); err != nil {
		h.logger.Error("Failed to update status", err)
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) CancelOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orderID := ps.ByName("id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.CancelOrder(r.Context(), orderID); err != nil {
		h.logger.Error("Failed to cancel order", err)
		http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
