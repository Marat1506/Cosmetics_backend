package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/handlers"
	"server/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	getUsers    = "/api/users"
	getUserById = "/api/user/:uuid"
	createUser  = "/api/createuser"
	login       = "/api/login"
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
	router.HandlerFunc(http.MethodGet, getUsers, h.GetAllUsers)
	router.GET(getUserById, h.GetUserById)
	router.POST(createUser, h.CreateUser)
	router.POST(login, h.Login)

	router.POST("/api/user/:userID/favorites/add", h.AddToFavorites)
	router.POST("/api/user/:userID/favorites/remove", h.RemoveFromFavorites)
	router.POST("/api/user/:userID/cart/add", h.AddToCart)
	router.POST("/api/user/:userID/cart/remove", h.RemoveFromCart)
	router.POST("/api/user/:userID/cart/update", h.UpdateCart)
	router.GET("/api/user/:userID/favorites", h.GetFavorites)

}
func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("this is list of users"))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var dto CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := h.service.CreateUser(r.Context(), dto)
	if err != nil {
		h.logger.Error("Failed to create user", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": userID})
}

func (h *handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())

	if err != nil {
		h.logger.Error("failed to get all users")
		http.Error(w, "failed to get all users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *handler) GetUserById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("id = ", r.Body)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var dto LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("ошибка входе", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), dto)
	if err != nil {
		h.logger.Error("ошибка при входу", err)
		http.Error(w, "Ошибка:", http.StatusInternalServerError)
		return
	}

	fmt.Println("user login = ", user)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]User{"user": user})

}

// Добавляем новые методы в структуру handler

func (h *handler) AddToFavorites(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		ProductID string `json:"productId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.AddToFavorites(r.Context(), userID, request.ProductID); err != nil {
		h.logger.Error("Failed to add to favorites", err)
		http.Error(w, "Failed to add to favorites", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h *handler) GetFavorites(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	favorites, err := h.service.GetFavorites(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get favorites", err)
		http.Error(w, "Failed to get favorites", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]string{"favorites": favorites})
}

func (h *handler) RemoveFromFavorites(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		ProductID string `json:"productId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveFromFavorites(r.Context(), userID, request.ProductID); err != nil {
		h.logger.Error("Failed to remove from favorites", err)
		http.Error(w, "Failed to remove from favorites", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) AddToCart(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		ProductID string `json:"productId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.AddToCart(r.Context(), userID, request.ProductID); err != nil {
		h.logger.Error("Failed to add to cart", err)
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) RemoveFromCart(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		ProductID string `json:"productId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveFromCart(r.Context(), userID, request.ProductID); err != nil {
		h.logger.Error("Failed to remove from cart", err)
		http.Error(w, "Failed to remove from cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) UpdateCart(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		Cart []string `json:"cart"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateCart(r.Context(), userID, request.Cart); err != nil {
		h.logger.Error("Failed to update cart", err)
		http.Error(w, "Failed to update cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
