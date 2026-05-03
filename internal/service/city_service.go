package service

import (
	"context"
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
)

type cityRepository interface {
	Create(ctx context.Context, input domain.CreateCityInput) (domain.City, error)
	Add2User(ctx context.Context, userID int64, input domain.AddCityInput) error
	ListOfUser(ctx context.Context, userID int64, filter domain.ListCitiesFilter) ([]domain.City, error)
	GetByName(ctx context.Context, name string) (domain.City, error)
	DeleteFromUser(ctx context.Context, userID, cityID int64) error
}

type CityService struct {
	repo cityRepository
}

func NewCityService(repo cityRepository) *CityService {
	return &CityService{repo: repo}
}

////// methods
////// methods
////// methods

func (s *CityService) Create(ctx context.Context, input domain.CreateCityInput) (domain.City, error) {
	if err := input.NormalizeAndValidate(); err != nil {
		return domain.City{}, err
	}
	return s.repo.Create(ctx, input)
}

func (s *CityService) Add2User(ctx context.Context, userID int64, input domain.AddCityInput) error {
	if err := input.NormalizeAndValidate(); err != nil {
		return err
	}
	return s.repo.Add2User(ctx, userID, input)
}

func (s *CityService) ListOfUser(ctx context.Context, userID int64, filter domain.ListCitiesFilter) ([]domain.City, error) {
	filter.Normalize()
	return s.repo.ListOfUser(ctx, userID, filter)
}

func (s *CityService) GetByName(ctx context.Context, name string) (domain.City, error) {
	// if err := input.NormalizeAndValidate(); err != nil {
	// 	return domain.City{}, err
	// }
	return s.repo.GetByName(ctx, strings.TrimSpace(strings.ToLower(name)))
}

func (s *CityService) DeleteFromUser(ctx context.Context, userID, cityID int64) error {
	if userID <= 0 {
		return domain.ErrInvalidUserID
	}
	if cityID <= 0 {
		return domain.ErrInvalidCityID
	}
	return s.repo.DeleteFromUser(ctx, userID, cityID)
}
