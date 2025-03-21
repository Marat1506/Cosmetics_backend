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
	getUsers = "/users"
	getUserById = "/users/:uuid"
	createUser = "/createuser"
	login = "/login"
	
)

type handler struct {
	logger *logging.Logger
	service *Service
} 

func NewHandler(logger *logging.Logger, service *Service) handlers.Handler {
	return &handler {
		logger: logger,
		service: service,
		
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, getUsers, h.GetAllUsers)
	router.GET(getUserById, h.GetUserById)
	router.POST(createUser, h.CreateUser)
	router.POST(login, h.Login)
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
		http.Error(w,"failed to get all users", http.StatusInternalServerError)
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