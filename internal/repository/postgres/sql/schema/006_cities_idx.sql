-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS ux_cities_city
ON cities (LOWER(city));

-- +goose Down
DROP INDEX IF EXISTS ux_cities_city;