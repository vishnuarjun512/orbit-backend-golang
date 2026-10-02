
-- +goose Up


-- =========================================
-- 1. ADD TASK WORKSPACE CONSTRAINT
-- =========================================

ALTER TABLE tasks
ADD CONSTRAINT tasks_workspace_task_unique
UNIQUE (workspace_id, task_id);


-- =========================================
-- 2. TEAMS
-- =========================================

CREATE TABLE teams (
    team_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    workspace_id UUID NOT NULL
        REFERENCES workspaces(workspace_id) ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,

    description TEXT,

    created_by UUID NOT NULL
        REFERENCES users(user_id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT teams_workspace_team_unique
        UNIQUE (workspace_id, team_id)
);


CREATE UNIQUE INDEX teams_workspace_name_unique_idx
ON teams(workspace_id, LOWER(name));


-- =========================================
-- 3. TEAM MEMBERS
-- =========================================

CREATE TABLE team_members (
    workspace_id UUID NOT NULL,

    team_id UUID NOT NULL,

    user_id UUID NOT NULL,

    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (team_id, user_id),

    CONSTRAINT team_members_team_fk
        FOREIGN KEY (workspace_id, team_id)
        REFERENCES teams(workspace_id, team_id)
        ON DELETE CASCADE,

    CONSTRAINT team_members_workspace_member_fk
        FOREIGN KEY (workspace_id, user_id)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE CASCADE
);


CREATE INDEX team_members_user_id_idx
ON team_members(user_id);


-- =========================================
-- 4. TAGS
-- =========================================

CREATE TABLE tags (
    tag_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    workspace_id UUID NOT NULL
        REFERENCES workspaces(workspace_id) ON DELETE CASCADE,

    name VARCHAR(50) NOT NULL,

    color VARCHAR(7) NOT NULL DEFAULT '#64748b',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT tags_color_check
        CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),

    CONSTRAINT tags_workspace_tag_unique
        UNIQUE (workspace_id, tag_id)
);


CREATE UNIQUE INDEX tags_workspace_name_unique_idx
ON tags(workspace_id, LOWER(name));


-- =========================================
-- 5. TASK TAGS
-- =========================================

CREATE TABLE task_tags (
    workspace_id UUID NOT NULL,

    task_id UUID NOT NULL,

    tag_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (task_id, tag_id),

    CONSTRAINT task_tags_task_fk
        FOREIGN KEY (workspace_id, task_id)
        REFERENCES tasks(workspace_id, task_id)
        ON DELETE CASCADE,

    CONSTRAINT task_tags_tag_fk
        FOREIGN KEY (workspace_id, tag_id)
        REFERENCES tags(workspace_id, tag_id)
        ON DELETE CASCADE
);


CREATE INDEX task_tags_tag_id_idx
ON task_tags(tag_id);


-- =========================================
-- 6. PROJECT TAGS
-- =========================================

CREATE TABLE project_tags (
    workspace_id UUID NOT NULL,

    project_id UUID NOT NULL,

    tag_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (project_id, tag_id),

    CONSTRAINT project_tags_project_fk
        FOREIGN KEY (workspace_id, project_id)
        REFERENCES projects(workspace_id, project_id)
        ON DELETE CASCADE,

    CONSTRAINT project_tags_tag_fk
        FOREIGN KEY (workspace_id, tag_id)
        REFERENCES tags(workspace_id, tag_id)
        ON DELETE CASCADE
);


CREATE INDEX project_tags_tag_id_idx
ON project_tags(tag_id);


-- +goose Down

DROP TABLE IF EXISTS project_tags;
DROP TABLE IF EXISTS task_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;

ALTER TABLE tasks
DROP CONSTRAINT IF EXISTS tasks_workspace_task_unique;

