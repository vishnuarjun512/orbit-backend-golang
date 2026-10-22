-- +goose Up

ALTER TABLE users ADD COLUMN public_id TEXT;
UPDATE users
SET public_id = 'ORB-' || UPPER(SUBSTRING(REPLACE(user_id::text, '-', '') FROM 1 FOR 8))
WHERE public_id IS NULL;
ALTER TABLE users ALTER COLUMN public_id SET NOT NULL;
CREATE UNIQUE INDEX users_public_id_unique_idx ON users (public_id);

CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    display_name TEXT NOT NULL,
    role_title TEXT,
    pronouns TEXT,
    bio TEXT CHECK (bio IS NULL OR CHAR_LENGTH(bio) <= 500),
    avatar_url TEXT,
    availability TEXT NOT NULL DEFAULT 'available'
        CHECK (availability IN ('available', 'away', 'do_not_disturb')),
    status_message TEXT CHECK (status_message IS NULL OR CHAR_LENGTH(status_message) <= 160),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO user_profiles (user_id, display_name, avatar_url)
SELECT user_id, COALESCE(NULLIF(BTRIM(full_name), ''), SPLIT_PART(email, '@', 1)), avatar_url
FROM users;

CREATE TABLE user_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    active_workspace_id UUID REFERENCES workspaces(workspace_id) ON DELETE SET NULL,
    email_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    task_updates BOOLEAN NOT NULL DEFAULT TRUE,
    mentions BOOLEAN NOT NULL DEFAULT TRUE,
    weekly_digest BOOLEAN NOT NULL DEFAULT FALSE,
    timezone TEXT NOT NULL DEFAULT 'UTC',
    week_starts_on SMALLINT NOT NULL DEFAULT 1 CHECK (week_starts_on IN (0, 1)),
    theme TEXT CHECK (theme IS NULL OR theme IN ('light', 'dark', 'system')),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO user_preferences (user_id) SELECT user_id FROM users;

-- +goose Down

DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS user_profiles;
DROP INDEX IF EXISTS users_public_id_unique_idx;
ALTER TABLE users DROP COLUMN IF EXISTS public_id;
