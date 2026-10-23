
-- +goose Up

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create users table
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id TEXT,

    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,

    verified BOOLEAN NOT NULL DEFAULT 'FALSE',

    status VARCHAR(20) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_status_check
        CHECK (status IN ('active', 'suspended', 'deleted'))
);

-- Prevent duplicate emails regardless of letter case
CREATE UNIQUE INDEX users_email_unique_idx
ON users (LOWER(email));


-- Prevent duplicate public_ids regardless of letter case
CREATE UNIQUE INDEX users_public_id_unique_idx ON users (public_id);


-- +goose Down

DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS users_email_unique_idx;
DROP INDEX IF EXISTS users_public_id_unique_idx;