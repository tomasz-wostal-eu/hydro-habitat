package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tomasz-wostal-eu/hydro-habitat/internal/model"
	"github.com/tomasz-wostal-eu/hydro-habitat/internal/service"
	"github.com/tomasz-wostal-eu/hydro-habitat/pkg/httputil"
)

// UserHandler handles HTTP requests for user resources
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes registers the user routes on the given router
func (h *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id:[0-9]+}", h.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/users", h.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id:[0-9]+}", h.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id:[0-9]+}", h.DeleteUser).Methods(http.MethodDelete)
}

// GetAllUsers handles GET requests to fetch all users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		httputil.SendError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	httputil.SendSuccess(w, http.StatusOK, "Users retrieved successfully", users)
}

// GetUserByID handles GET requests to fetch a specific user
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		httputil.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		httputil.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.SendSuccess(w, http.StatusOK, "User retrieved successfully", user)
}

// CreateUser handles POST requests to create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userData model.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		httputil.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	id, err := h.userService.CreateUser(r.Context(), userData)
	if err != nil {
		httputil.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get the complete user data to return
	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		httputil.SendError(w, http.StatusInternalServerError, "User created but failed to retrieve")
		return
	}

	httputil.SendSuccess(w, http.StatusCreated, "User created successfully", user)
}

// UpdateUser handles PUT requests to update an existing user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		httputil.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var userData model.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		httputil.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.userService.UpdateUser(r.Context(), id, userData); err != nil {
		httputil.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get the updated user data to return
	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		httputil.SendError(w, http.StatusInternalServerError, "User updated but failed to retrieve")
		return
	}

	httputil.SendSuccess(w, http.StatusOK, "User updated successfully", user)
}

// DeleteUser handles DELETE requests to remove a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		httputil.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(r.Context(), id); err != nil {
		httputil.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.SendSuccess(w, http.StatusOK, "User deleted successfully", nil)
}
