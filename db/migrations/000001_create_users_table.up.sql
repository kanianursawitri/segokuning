CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR NULL DEFAULT '',
    phone VARCHAR NULL DEFAULT '',
    "name" VARCHAR NOT NULL UNIQUE,
    "password" VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX id ON users (id);
CREATE INDEX phone ON users (phone);
CREATE INDEX email ON users (email);
