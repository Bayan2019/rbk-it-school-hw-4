package domain

import (
	"strings"
	"time"
)

type CityWeatherInput struct {
	City        string `db:"city" json:"city"`
	Temperature float64
	Description string
}

type CityWeatherResponse struct {
	City        string  `db:"city" json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

type ProviderWeatherResponse struct {
	Temperature float64
	Description string
}

type WeatherHistoryFilter struct {
	City   string `json:"city"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type WeatherResult struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Temperature    float64 `json:"temperature"`
	WindSpeed      float64 `json:"wind_speed"`
	WeatherCode    int     `json:"weather_code"`
	Time           string  `json:"time"`
	Description    string  `json:"description"`
	Recommendation string  `json:"recommendation"`
}

type WeatherHistoryResponse struct {
	UserID      int       `db:"user_id" json:"user_id"`
	City        string    `db:"city" json:"city,omitempty"`
	RequestedAt time.Time `db:"requested_at" json:"requested_at"`
	Temperature int       `db:"temperature" json:"temperature"`
	Description string    `db:"description" json:"description"`
}

type WeatherHistoryCityResponse struct {
	RequestedAt time.Time `db:"requested_at" json:"requested_at"`
	Temperature int       `db:"temperature" json:"temperature"`
	Description string    `db:"description" json:"description"`
}

type WeatherHistoryOfUserOfCityResponse struct {
	UserID  int                          `json:"user_id"`
	City    string                       `json:"city"`
	History []WeatherHistoryCityResponse `json:"history"`
}

type WeatherHistoryOfUserResponse struct {
	UserID  int                      `json:"user_id"`
	History []WeatherHistoryResponse `json:"history"`
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

func (in *CityWeatherInput) NormalizeAndValidate() error {
	in.City = strings.TrimSpace(strings.ToLower(in.City))

	if in.City == "" {
		return ErrInvalidCityInput
	}

	return nil
}

func (f *WeatherHistoryFilter) Normalize() {
	if f.Limit <= 0 {
		f.Limit = 0
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	f.City = strings.TrimSpace(
		strings.ToLower(f.City))
}
