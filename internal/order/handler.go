package order

import (
	"encoding/json"
	"net/http"
	"server/handlers"
	"server/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	createOrder = "/api/createOrder"
	getOrders   = "/api/getOrders"
	changeOrder = "/api/changeOrder"
	deleteOrder = "/api/deleteOrder"
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
	router.POST(createOrder, h.CreateOrder)
	router.HandlerFunc(http.MethodGet, getOrders, h.GetOrders)
	router.PATCH(changeOrder, h.ChangeOrder)
	router.DELETE(deleteOrder, h.DeleteOrder)
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var dto CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	orderID, err := h.service.CreateOrder(r.Context(), dto)
	if err != nil {
		h.logger.Error("Failed to create user", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": orderID})
}

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetOrders(r.Context())

	if err != nil {
		h.logger.Error("Failed to get orders", err)
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)

}
func (h *handler) DeleteOrder(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var requestData struct {
		Id string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteOrder(r.Context(), requestData.Id); err != nil {
		h.logger.Error("Failed to delete order", err)
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
}

func (h *handler) ChangeOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var requestData struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.ChangeOrder(r.Context(), requestData.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
