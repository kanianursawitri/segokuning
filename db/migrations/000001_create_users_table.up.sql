CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR NULL,
    phone VARCHAR NULL,
    "name" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_credential_type ON users (id);
CREATE INDEX idx_credential_type ON users (credential_type);
CREATE INDEX idx_credential_value ON users (credential_value);
