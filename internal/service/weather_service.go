package service

import (
	"context"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
)

type weatherRepository interface {
	CreateHistory(ctx context.Context, userID int64, cityWeather domain.CityWeatherInput) (domain.WeatherHistoryResponse, error)
	WeatherHistoryOfUser(ctx context.Context, userID int64, filter domain.WeatherHistoryFilter) ([]domain.WeatherHistoryResponse, error)
}

type WeatherService struct {
	repo weatherRepository
}

func NewWeatherService(repo weatherRepository) *WeatherService {
	return &WeatherService{repo: repo}
}

////// methods
////// methods
////// methods

func (s *WeatherService) CreateHistory(ctx context.Context, userID int64, cityWeather domain.CityWeatherInput) (domain.WeatherHistoryResponse, error) {
	if err := cityWeather.NormalizeAndValidate(); err != nil {
		return domain.WeatherHistoryResponse{}, err
	}
	return s.repo.CreateHistory(ctx, userID, cityWeather)
}

func (s *WeatherService) WeatherHistoryOfUser(ctx context.Context, userID int64, filter domain.WeatherHistoryFilter) ([]domain.WeatherHistoryResponse, error) {

	filter.Normalize()
	return s.repo.WeatherHistoryOfUser(ctx, userID, filter)
}
