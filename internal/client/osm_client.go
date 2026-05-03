package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/config"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
	"golang.org/x/time/rate"
)

type OsmClient struct {
	baseURL   string
	userAgent string
	client    *http.Client
	limiter   *rate.Limiter
}

func NewOsmClient(cfg config.ApiConfig) *OsmClient {
	return &OsmClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL:   "https://nominatim.openstreetmap.org",
		limiter:   cfg.Limiter,
		userAgent: cfg.UserAgent,
	}
}

////// methods
////// methods
////// methods

func (client *OsmClient) GetInfoOfCity(ctx context.Context, city string) (domain.Place, error) {
	var place domain.Place

	url := fmt.Sprintf("%s/search?city=%s&format=json", client.baseURL, city)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// fmt.Println("Error creating request:", err)
		return place, fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("User-Agent", client.userAgent)

	res, err := client.client.Do(req)
	if err != nil {
		return place, err
	}
	defer res.Body.Close()

	var places []domain.Place

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&places)
	if err != nil {
		return place, err
	}

	if len(places) < 1 {
		return place, fmt.Errorf("city %s not found", city)
	}

	place = places[0]

	return place, nil
}
