package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
	httptransport "github.com/Bayan2019/rbk-it-school-hw-4/internal/transport/http"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// 1. Аутентификация
	r.Post("/auth/register", h.User.Register)
	r.Post("/auth/login", h.User.Login)

	// 4. Защита маршрутов
	r.Group(func(r chi.Router) {
		// Все операции должны работать через текущего пользователя из JWT.
		r.Use(httptransport.AuthMiddleware(h.User.JwtManager))
		// Убрать user_id из URL.
		r.Post("/cities", h.City.Add2User)
		r.Get("/cities", h.City.ListOfUser)
		r.Delete("/cities/{city_id}", h.City.DeleteFromUser)
		r.Get("/weather", h.Weather.GetWeatherOfUserCities)
		r.Get("/weather/history", h.Weather.GetWeatherHistoryOfUser)

		// 8. Новый endpoint
		r.Get("/users/me", h.User.Profile)

		// 5. Авторизация (Roles)
		r.Group(func(r chi.Router) {
			// Использовать middleware RequireRole("admin")
			r.Use(httptransport.RequireRole(domain.RolesAdmin))
			// Только admin может:
			r.Get("/users", h.User.List)
			r.Get("/users/{id}", h.User.GetByID)
			r.Delete("/users/{id}", h.User.Delete)
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Put("/{id}", h.User.Update)
		})
	})

	return r
}
