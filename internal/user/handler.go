package user

import (
	"encoding/json"
	"net/http"
	"server/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	registerPath  = "/api/auth/register"
	loginPath     = "/api/auth/login"
	profilePath   = "/api/user/profile"
	favoritesPath = "/api/user/favorites"
	cartPath      = "/api/user/cart"
)

type Handler struct {
	logger  *logging.Logger
	service *Service
}

func NewHandler(logger *logging.Logger, service *Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) RegisterRouter(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, registerPath, h.RegisterUser)
	router.HandlerFunc(http.MethodPost, loginPath, h.LoginUser)
	router.HandlerFunc(http.MethodGet, profilePath, h.GetUserProfile)
	router.HandlerFunc(http.MethodPost, favoritesPath, h.AddToUserFavorites)
	router.HandlerFunc(http.MethodDelete, favoritesPath, h.RemoveFromUserFavorites)
	router.HandlerFunc(http.MethodPost, cartPath, h.AddToUserCart)
	router.HandlerFunc(http.MethodDelete, cartPath, h.RemoveFromUserCart)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var dto CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := h.service.Register(r.Context(), dto)
	if err != nil {
		h.logger.Error("Failed to register user", err)
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": userID})
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var dto LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), dto)
	if err != nil {
		h.logger.Error("Failed to login", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get profile", err)
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) AddToUserFavorites(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var dto UpdateFavoritesDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.AddToFavorites(r.Context(), userID, dto.ProductID); err != nil {
		h.logger.Error("Failed to add to favorites", err)
		http.Error(w, "Failed to add to favorites", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RemoveFromUserFavorites(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var dto UpdateFavoritesDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveFromFavorites(r.Context(), userID, dto.ProductID); err != nil {
		h.logger.Error("Failed to remove from favorites", err)
		http.Error(w, "Failed to remove from favorites", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddToUserCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var dto UpdateCartDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.AddToCart(r.Context(), userID, dto.ProductID); err != nil {
		h.logger.Error("Failed to add to cart", err)
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RemoveFromUserCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var dto UpdateCartDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to decode request", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveFromCart(r.Context(), userID, dto.ProductID); err != nil {
		h.logger.Error("Failed to remove from cart", err)
		http.Error(w, "Failed to remove from cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
