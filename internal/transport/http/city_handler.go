package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type cityService interface {
	Create(ctx context.Context, input domain.CreateCityInput) (domain.City, error)
	Add2User(ctx context.Context, userID int64, input domain.AddCityInput) error
	ListOfUser(ctx context.Context, userID int64, filter domain.ListCitiesFilter) ([]domain.City, error)
	GetByName(ctx context.Context, name string) (domain.City, error)
	DeleteFromUser(ctx context.Context, userID, cityID int64) error
}

type osmProvider interface {
	GetInfoOfCity(ctx context.Context, city string) (domain.Place, error)
}

type CityHandler struct {
	service  cityService
	provider osmProvider
}

type citiesResponse struct {
	Data []domain.City `json:"data"`
}

type cityResponse struct {
	Data domain.City `json:"data"`
}

func NewCityHandler(service cityService, provider osmProvider) *CityHandler {
	return &CityHandler{
		service:  service,
		provider: provider,
	}
}

////// methods
////// methods
////// methods

func (h *CityHandler) Add2User(w http.ResponseWriter, r *http.Request) {

	user, err := UserFromContext(r.Context())
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	var input domain.AddCityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json body", Message: "couldn't parse input Body"})
		return
	}

	city, err := h.service.GetByName(r.Context(), input.City)
	if err != nil {
		if notFound(err) {
			place, err := h.provider.GetInfoOfCity(r.Context(), strings.TrimSpace(strings.ToLower(input.City)))
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error(), Message: "couldn't get info from osm"})
				return
			}

			lat, err := strconv.ParseFloat(place.Lat, 64)
			if err != nil {
				h.handleError(w, err)
				return
			}

			lon, err := strconv.ParseFloat(place.Lon, 64)
			if err != nil {
				h.handleError(w, err)
				return
			}

			city, err := h.service.Create(r.Context(), domain.CreateCityInput{
				City: input.City,
				Lat:  lat,
				Lon:  lon,
			})
			if err != nil {
				h.handleError(w, err)
				return
			}

			err = h.service.Add2User(r.Context(), user.ID, domain.AddCityInput{City: city.City})

			WriteJSON(w, http.StatusCreated, cityResponse{Data: city})
			return
		}
		h.handleError(w, err)
		return
	}

	err = h.service.Add2User(r.Context(), user.ID, input)
	if err != nil {
		h.handleError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, cityResponse{Data: city})
}

func (h *CityHandler) ListOfUser(w http.ResponseWriter, r *http.Request) {
	user, err := UserFromContext(r.Context())
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	filter := domain.ListCitiesFilter{
		Offset:         parseIntQuery(r, "offset", 0),
		IncludeDeleted: parseBoolQuery(r, "include_deleted", false),
	}

	cities, err := h.service.ListOfUser(r.Context(), user.ID, filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, citiesResponse{Data: cities})
}

func (h *CityHandler) DeleteFromUser(w http.ResponseWriter, r *http.Request) {
	userID, err := parseIDParam(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	cityID, err := strconv.ParseInt(chi.URLParam(r, "city_id"), 10, 64)
	if err != nil || cityID <= 0 {
		h.handleError(w, domain.ErrInvalidCityID)
	}

	if err := h.service.DeleteFromUser(r.Context(), userID, cityID); err != nil {
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

func (h *CityHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCityID), errors.Is(err, domain.ErrInvalidCityInput):
		WriteJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrCityNotFound):
		WriteJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	default:
		// writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		WriteJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
	}
}

func notFound(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "42P01"
	}
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}

func parseIDParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, domain.ErrInvalidUserID
	}
	return id, nil
}

func parseIntQuery(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseBoolQuery(r *http.Request, key string, fallback bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
