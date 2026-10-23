-- +goose Up

CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,

    display_name TEXT,
    role_title TEXT,
    pronouns TEXT,

    bio TEXT,
    avatar_url TEXT,
    availability TEXT NOT NULL DEFAULT 'available',
    status_message TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    CHECK (bio IS NULL OR CHAR_LENGTH(bio) <= 500),
    CHECK (availability IN ('available', 'away', 'do_not_disturb')),
    CHECK (status_message IS NULL OR CHAR_LENGTH(status_message) <= 160)
);


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

-- +goose Down

DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS user_profiles;
