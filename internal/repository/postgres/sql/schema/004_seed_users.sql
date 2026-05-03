-- +goose Up
INSERT INTO users (email, password_hash, first_name, last_name, is_active, role)
VALUES
    ('admin@example.com', '$2a$12$19E8eNXYZEZGelkKqNBIPuGIqptOh4lZsjQHi.Y2V2vZV4V8dFFJe',
        'Admin', 'Role', TRUE, 'admin'),
    ('user@example.com', '$2a$12$wb9gOgmR4OwdVfl2DmfsY.AIK53qR33nYdkPuMAWT8HsIXm7dEpBS', 
        'User', 'Role', TRUE, 'user'),
    ('ivan@example.com', '$2a$12$UJoiT8R0p2bLv3BHbDYhquWuRb33bmEFMJEz5p3bv3l0Ygvu9M2b6', 
        'Ivan', 'Ivanov', TRUE, 'user'),
    ('dana@example.com', '$2a$12$4BjlkDs6ZX.6n3eeNAp/y.b9ItwSLzVynAtHIqSU2a.svD7Szs9xa',
        'Dana', 'Sadykova', TRUE, 'user'),
    ('ayan@example.com', '$2a$12$LqHfxRBEcd7FU869.dVsn.2dmtZo/AukWuKnnopKrqT/jULxcDF22', 
        'Ayan', 'Kaliyev', FALSE, 'user')
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM users 
WHERE email IN ('admin@example.com', 'user@example.com', 'ivan@example.com', 'dana@example.com', 'ayan@example.com');