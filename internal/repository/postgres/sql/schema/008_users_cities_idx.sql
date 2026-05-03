-- +goose Up
CREATE INDEX IF NOT EXISTS idx_users_cities_user_city 
ON users_cities (user_id, city_id)
WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_users_cities_user_city;