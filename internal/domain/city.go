package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrCityNotFound          = errors.New("city not found")
	ErrInvalidCityID         = errors.New("invalid city id")
	ErrInvalidCityInput      = errors.New("invalid city input")
	ErrCityAlreadyAdded2User = errors.New("city already added to user")
)

type City struct {
	CityID    int64     `db:"city_id" json:"city_id,omitempty"`
	City      string    `db:"city" json:"city"`
	Lat       float64   `db:"lat" json:"lat"`
	Lon       float64   `db:"lon" json:"lon"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CityItem struct {
	City string  `db:"city" json:"city"`
	Lat  float64 `db:"lat" json:"lat"`
	Lon  float64 `db:"lon" json:"lon"`
}

type AddCityInput struct {
	City string `json:"city"`
}

type CreateCityInput struct {
	City string  `json:"city"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type ListCitiesFilter struct {
	Offset         int  `json:"offset"`
	IncludeDeleted bool `json:"include_deleted"`
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

func (in *CreateCityInput) NormalizeAndValidate() error {
	in.City = strings.TrimSpace(strings.ToLower(in.City))

	if in.City == "" {
		return ErrInvalidCityInput
	}

	return nil
}

func (in *AddCityInput) NormalizeAndValidate() error {
	in.City = strings.TrimSpace(strings.ToLower(in.City))

	if in.City == "" {
		return ErrInvalidCityInput
	}

	return nil
}

func (f *ListCitiesFilter) Normalize() {
	if f.Offset < 0 {
		f.Offset = 0
	}
}
