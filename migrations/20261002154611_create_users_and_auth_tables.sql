
-- +goose Up

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create users table
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(255) NOT NULL,
    username VARCHAR(50) NOT NULL,
    full_name VARCHAR(100) NOT NULL,

    password_hash TEXT NOT NULL,

    avatar_url TEXT,

    email_verified_at TIMESTAMPTZ,

    status VARCHAR(20) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_status_check
        CHECK (status IN ('active', 'suspended', 'deleted'))
);

-- Prevent duplicate emails regardless of letter case
CREATE UNIQUE INDEX users_email_unique_idx
ON users (LOWER(email));

-- Prevent duplicate usernames regardless of letter case
CREATE UNIQUE INDEX users_username_unique_idx
ON users (LOWER(username));


-- +goose Down

DROP TABLE IF EXISTS users;
