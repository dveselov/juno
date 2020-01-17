BEGIN;

CREATE TABLE IF NOT EXISTS authentication_user (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS authentication_provider_code (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    code VARCHAR(255) NUT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX authentication_provider_code_idx ON authentication_provider_code (
    provider, phone_number, code
);

COMMIT;
