package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CityRepository struct {
	db *sqlx.DB
}

func NewCityRepository(db *sqlx.DB) *CityRepository {
	return &CityRepository{db: db}
}

////// methods
////// methods
////// methods

func (r *CityRepository) Create(ctx context.Context, input domain.CreateCityInput) (domain.City, error) {
	query := `
		INSERT INTO cities (city, lat, lon)
		VALUES (:city, :lat, :lon)
		RETURNING city_id, city, lat, lon,  created_at, updated_at
	`

	args := map[string]any{
		"city": input.City,
		"lat":  input.Lat,
		"lon":  input.Lon,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.City{}, domain.ErrEmailAlreadyTaken
		}
		return domain.City{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var city domain.City
		if err := rows.StructScan(&city); err != nil {
			return domain.City{}, err
		}
		return city, nil
	}

	return domain.City{}, errors.New("failed to create city")
}

func (r *CityRepository) Add2User(ctx context.Context, userID int64, input domain.AddCityInput) error {
	city, err := r.GetByName(ctx, input.City)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO users_cities (user_id, city_id)
		VALUES (:user_id, :city_id)
		RETURNING user_id, city_id
	`

	args := map[string]any{
		"city_id": city.CityID,
		"user_id": userID,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrCityAlreadyAdded2User
		}
		return err
	}
	defer rows.Close()

	if rows.Next() {
		// var city domain.City
		// if err := rows.StructScan(&city); err != nil {
		// 	return domain.City{}, err
		// }
		return nil
	}

	return errors.New("failed to add city")
}

func (r *CityRepository) ListOfUser(ctx context.Context, userID int64, filter domain.ListCitiesFilter) ([]domain.City, error) {
	builder := strings.Builder{}
	builder.WriteString(`
		SELECT c.city_id, c.city, c.lat, c.lon, c.created_at, c.updated_at
		FROM cities AS c
		JOIN users_cities AS uc
		ON c.city_id = uc.city_id
		WHERE uc.user_id = :user_id
	`)

	args := map[string]any{
		"user_id": userID,
		"offset":  filter.Offset,
	}

	if !filter.IncludeDeleted {
		builder.WriteString(" AND uc.deleted_at IS NULL")
	}

	builder.WriteString(" ORDER BY uc.added_at ASC OFFSET :offset")

	query, queryArgs, err := sqlx.Named(builder.String(), args)
	if err != nil {
		return nil, errors.New("sqlx.Named")
	}
	query = r.db.Rebind(query)

	var cities []domain.City
	if err := r.db.SelectContext(ctx, &cities, query, queryArgs...); err != nil {
		return nil, errors.New("r.db.SelectContext")
	}

	return cities, nil
}

func (r *CityRepository) GetByName(ctx context.Context, name string) (domain.City, error) {

	query := `
		SELECT city_id, city, lat, lon, created_at, updated_at
		FROM cities
		WHERE city = $1
	`

	var city domain.City
	if err := r.db.GetContext(ctx, &city, query, name); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return domain.City{}, domain.ErrCityNotFound
		}
		return domain.City{}, err
	}

	return city, nil
}

func (r *CityRepository) DeleteFromUser(ctx context.Context, userID, cityID int64) error {
	query := `
		UPDATE users_cities
		SET deleted_at = NOW()
		WHERE user_id = $1
			AND city_id = $2
			AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, userID, cityID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("result.RowsAffected()")
	}
	if rowsAffected == 0 {
		return domain.ErrCityNotFound
	}

	return nil
}
