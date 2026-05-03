-- +goose Up
CREATE INDEX IF NOT EXISTS idx_weather_history_user_city 
ON weather_history (user_id, city);

-- +goose Down
DROP INDEX IF EXISTS idx_weather_history_user_city;