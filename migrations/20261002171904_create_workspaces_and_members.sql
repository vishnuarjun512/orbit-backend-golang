
-- +goose Up


-- =========================================
-- 1. WORKSPACES
-- =========================================

CREATE TABLE workspaces (
    workspace_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(100) NOT NULL,

    slug VARCHAR(100) NOT NULL UNIQUE,

    description TEXT,

    created_by UUID NOT NULL
        REFERENCES users(user_id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX workspaces_created_by_idx
ON workspaces(created_by);


-- =========================================
-- 2. WORKSPACE MEMBERS
-- =========================================

CREATE TABLE workspace_members (
    workspace_id UUID NOT NULL
        REFERENCES workspaces(workspace_id) ON DELETE CASCADE,

    user_id UUID NOT NULL
        REFERENCES users(user_id) ON DELETE CASCADE,

    role VARCHAR(20) NOT NULL DEFAULT 'member',

    status VARCHAR(20) NOT NULL DEFAULT 'active',

    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (workspace_id, user_id),

    CONSTRAINT workspace_members_role_check
        CHECK (role IN ('owner', 'admin', 'member', 'guest')),

    CONSTRAINT workspace_members_status_check
        CHECK (status IN ('active', 'suspended'))
);


CREATE INDEX workspace_members_user_id_idx
ON workspace_members(user_id);


-- =========================================
-- 3. WORKSPACE INVITATIONS
-- =========================================

CREATE TABLE workspace_invitations (
    invitation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    workspace_id UUID NOT NULL
        REFERENCES workspaces(workspace_id) ON DELETE CASCADE,

    email VARCHAR(255) NOT NULL,

    role VARCHAR(20) NOT NULL DEFAULT 'member',

    token_hash TEXT NOT NULL UNIQUE,

    invited_by UUID
        REFERENCES users(user_id) ON DELETE SET NULL,

    accepted_by UUID
        REFERENCES users(user_id) ON DELETE SET NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    expires_at TIMESTAMPTZ NOT NULL,

    accepted_at TIMESTAMPTZ,

    revoked_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT workspace_invitations_role_check
        CHECK (role IN ('admin', 'member', 'guest')),

    CONSTRAINT workspace_invitations_status_check
        CHECK (status IN ('pending', 'accepted', 'revoked', 'expired'))
);


CREATE INDEX workspace_invitations_workspace_id_idx
ON workspace_invitations(workspace_id);


CREATE INDEX workspace_invitations_email_idx
ON workspace_invitations(LOWER(email));


-- +goose Down

DROP TABLE IF EXISTS workspace_invitations;
DROP TABLE IF EXISTS workspace_members;
DROP TABLE IF EXISTS workspaces;

