-- +goose Up
CREATE TABLE users_cities(
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    city_id BIGINT NOT NULL REFERENCES cities(city_id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(user_id, city_id)
);

-- +goose Down
DROP TABLE users_cities;