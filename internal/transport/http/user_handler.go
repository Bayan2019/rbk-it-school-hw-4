package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
)

type userService interface {
	Create(ctx context.Context, input domain.RegisterUserInput) (domain.User, error)
	List(ctx context.Context, filter domain.ListUsersFilter) ([]domain.User, error)
	GetByID(ctx context.Context, id int64, includeDeleted bool) (domain.User, error)
	GetByEmail(ctx context.Context, email string, includeDeleted bool) (domain.User, error)
	Update(ctx context.Context, id int64, input domain.UpdateUserInput) (domain.User, error)
	Delete(ctx context.Context, id int64) error
}

type UserHandler struct {
	service    userService
	JwtManager *auth.JWTManager
}

/// json
/// json
/// json
/// json
/// json

type usersResponse struct {
	Data []domain.User `json:"data"`
}

type userResponse struct {
	Data domain.User `json:"data"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

/// json
/// json
/// json
/// json
/// json

func NewUserHandler(service userService, jwtManager *auth.JWTManager) *UserHandler {
	return &UserHandler{
		service:    service,
		JwtManager: jwtManager,
	}
}

////// methods
////// methods
////// methods

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// 1. Аутентификация
	// - регистрация пользователя
	var input domain.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid json body", Error: err.Error()})
		return
	}

	user, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// w.Header().Set("Location", "/api/v1/users/"+strconv.FormatInt(user.ID, 10))
	WriteJSON(w, http.StatusCreated, userResponse{Data: user})
}

// 1. Аутентификация
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json body", err)
		return
	}

	user, err := h.service.GetByEmail(r.Context(), req.Email, false)
	if err == domain.ErrUserNotFound {
		WriteError(w, http.StatusUnauthorized, "email not found", err)
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal error", err)
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		WriteError(w, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	token, err := h.JwtManager.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "token generation failed", err)
		return
	}

	// 1. Аутентификация
	// - возвращает access_token (JWT)
	WriteJSON(w, http.StatusOK, loginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user, err := UserFromContext(r.Context())
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "this is protected profile endpoint",
		"user":    user,
	})
}

func (h *UserHandler) AdminReports(w http.ResponseWriter, r *http.Request) {
	user, err := UserFromContext(r.Context())
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "this is admin-only endpoint",
		"user":    user,
		"reports": []string{
			"daily-auth-report",
			"security-login-report",
		},
	})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := domain.ListUsersFilter{
		Limit:          parseIntQuery(r, "limit", 20),
		Offset:         parseIntQuery(r, "offset", 0),
		Query:          r.URL.Query().Get("q"),
		IncludeDeleted: parseBoolQuery(r, "include_deleted", false),
	}

	users, err := h.service.List(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, usersResponse{Data: users})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	user, err := h.service.GetByID(r.Context(), id, parseBoolQuery(r, "include_deleted", false))
	if err != nil {
		h.handleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, userResponse{Data: user})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	var input domain.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json body"})
		return
	}

	user, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, userResponse{Data: user})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions

func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidUserID), errors.Is(err, domain.ErrInvalidUserInput):
		WriteJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUserNotFound):
		WriteJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrEmailAlreadyTaken):
		WriteJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	default:
		WriteJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		// writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
	}
}
