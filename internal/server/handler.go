package server

import (
	"context"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"

	httptransport "github.com/Bayan2019/rbk-it-school-hw-4/internal/transport/http"
)

type userService interface {
	Create(ctx context.Context, input domain.RegisterUserInput) (domain.User, error)
	List(ctx context.Context, filter domain.ListUsersFilter) ([]domain.User, error)
	GetByID(ctx context.Context, id int64, includeDeleted bool) (domain.User, error)
	GetByEmail(ctx context.Context, email string, includeDeleted bool) (domain.User, error)
	Update(ctx context.Context, id int64, input domain.UpdateUserInput) (domain.User, error)
	Delete(ctx context.Context, id int64) error
}

type cityService interface {
	Create(ctx context.Context, input domain.CreateCityInput) (domain.City, error)
	Add2User(ctx context.Context, userID int64, filter domain.AddCityInput) error
	ListOfUser(ctx context.Context, userID int64, filter domain.ListCitiesFilter) ([]domain.City, error)
	GetByName(ctx context.Context, name string) (domain.City, error)
	DeleteFromUser(ctx context.Context, userID, cityID int64) error
}

type weatherService interface {
	CreateHistory(ctx context.Context, userID int64, cityWeather domain.CityWeatherInput) (domain.WeatherHistoryResponse, error)
	WeatherHistoryOfUser(ctx context.Context, userID int64, filter domain.WeatherHistoryFilter) ([]domain.WeatherHistoryResponse, error)
}

type osmProvider interface {
	GetInfoOfCity(ctx context.Context, city string) (domain.Place, error)
}

type weatherProvider interface {
	GetCurrentWeather(ctx context.Context, lat, lon float64) (domain.ProviderWeatherResponse, error)
}

type Handler struct {
	User    *httptransport.UserHandler
	City    *httptransport.CityHandler
	Weather *httptransport.WeatherHandler
	// jwtManager *auth.JWTManager
}

func NewHandler(userService userService, cityService cityService,
	weatherService weatherService,
	osmProvider osmProvider,
	weatherProvider weatherProvider,
	jwtManager *auth.JWTManager,
) *Handler {
	return &Handler{
		User: httptransport.NewUserHandler(userService, jwtManager),
		City: httptransport.NewCityHandler(cityService, osmProvider),
		Weather: &httptransport.WeatherHandler{
			CityService:     cityService,
			WeatherService:  weatherService,
			WeatherProvider: weatherProvider,
		},
		// jwtManager: jwtManager,
	}
}

/// json
/// json
/// json
/// json
/// json

/// json
/// json
/// json
/// json
/// json
