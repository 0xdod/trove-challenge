DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPZ,
    updated_at TIMESTAMPZ
);
