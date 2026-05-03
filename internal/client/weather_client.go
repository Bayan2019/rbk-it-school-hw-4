package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
)

type WeatherClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewWeatherClient(httpClient *http.Client) *WeatherClient {
	return &WeatherClient{
		httpClient: httpClient,
		baseURL:    "https://api.open-meteo.com/v1/forecast",
	}
}

type openMeteoResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
		Weathercode int     `json:"weathercode"`
		Time        string  `json:"time"`
	} `json:"current_weather"`
}

////// methods
////// methods
////// methods

func (c *WeatherClient) GetCurrentWeather(ctx context.Context, lat, lon float64) (domain.ProviderWeatherResponse, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return domain.ProviderWeatherResponse{}, fmt.Errorf("parse base url: %w", err)
	}

	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%.4f", lat))
	q.Set("longitude", fmt.Sprintf("%.4f", lon))
	q.Set("current_weather", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return domain.ProviderWeatherResponse{}, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.ProviderWeatherResponse{}, fmt.Errorf("call external api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.ProviderWeatherResponse{}, fmt.Errorf("external api returned status: %d", resp.StatusCode)
	}

	var result openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return domain.ProviderWeatherResponse{}, fmt.Errorf("decode external api response: %w", err)
	}

	return domain.ProviderWeatherResponse{
		Temperature: result.CurrentWeather.Temperature,
		Description: mapWeatherCode(result.CurrentWeather.Weathercode),
	}, nil
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

func mapWeatherCode(code int) string {
	switch code {
	case 0:
		return "sunny"
	case 1, 2, 3:
		return "cloudy"
	case 45, 48:
		return "fog"
	case 51, 53, 55:
		return "drizzle"
	case 61, 63, 65:
		return "rain"
	case 71, 73, 75:
		return "snow"
	case 95:
		return "storm"
	default:
		return "unknown"
	}
}
