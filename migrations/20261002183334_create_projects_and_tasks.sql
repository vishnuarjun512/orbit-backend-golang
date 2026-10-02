
-- +goose Up


-- =========================================
-- 1. PROJECTS
-- =========================================

CREATE TABLE projects (
    project_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    workspace_id UUID NOT NULL
        REFERENCES workspaces(workspace_id) ON DELETE CASCADE,

    name VARCHAR(150) NOT NULL,

    project_key VARCHAR(20) NOT NULL,

    description TEXT,

    status VARCHAR(30) NOT NULL DEFAULT 'active',

    visibility VARCHAR(20) NOT NULL DEFAULT 'workspace',

    start_date DATE,

    due_date DATE,

    created_by UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT projects_status_check
        CHECK (status IN ('active', 'on_hold', 'completed', 'archived')),

    CONSTRAINT projects_visibility_check
        CHECK (visibility IN ('workspace', 'private')),

    CONSTRAINT projects_dates_check
        CHECK (due_date IS NULL OR start_date IS NULL OR due_date >= start_date),

    CONSTRAINT projects_workspace_key_unique
        UNIQUE (workspace_id, project_key),

    -- Required for composite foreign keys from other tables.
    CONSTRAINT projects_workspace_project_unique
        UNIQUE (workspace_id, project_id)
);


CREATE INDEX projects_workspace_id_idx
ON projects(workspace_id);

CREATE INDEX projects_created_by_idx
ON projects(created_by);


-- =========================================
-- 2. PROJECT MEMBERS
-- =========================================

CREATE TABLE project_members (
    workspace_id UUID NOT NULL,

    project_id UUID NOT NULL,

    user_id UUID NOT NULL,

    added_by UUID,

    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (project_id, user_id),

    CONSTRAINT project_members_project_fk
        FOREIGN KEY (workspace_id, project_id)
        REFERENCES projects(workspace_id, project_id)
        ON DELETE CASCADE,

    CONSTRAINT project_members_workspace_member_fk
        FOREIGN KEY (workspace_id, user_id)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE CASCADE,

    CONSTRAINT project_members_added_by_fk
        FOREIGN KEY (workspace_id, added_by)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE SET NULL (added_by)
);


CREATE INDEX project_members_user_id_idx
ON project_members(user_id);


-- =========================================
-- 3. TASKS
-- =========================================

CREATE TABLE tasks (
    task_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    workspace_id UUID NOT NULL,

    project_id UUID NOT NULL,

    parent_task_id UUID,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    status VARCHAR(30) NOT NULL DEFAULT 'todo',

    priority VARCHAR(20) NOT NULL DEFAULT 'medium',

    assignee_id UUID,

    created_by UUID NOT NULL,

    start_date TIMESTAMPTZ,

    due_date TIMESTAMPTZ,

    completed_at TIMESTAMPTZ,

    position BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT tasks_priority_check
        CHECK (priority IN ('low', 'medium', 'high', 'urgent')),

    CONSTRAINT tasks_dates_check
        CHECK (due_date IS NULL OR start_date IS NULL OR due_date >= start_date),

    CONSTRAINT tasks_project_fk
        FOREIGN KEY (workspace_id, project_id)
        REFERENCES projects(workspace_id, project_id)
        ON DELETE CASCADE,

    CONSTRAINT tasks_assignee_fk
        FOREIGN KEY (workspace_id, assignee_id)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE SET NULL (assignee_id),

    CONSTRAINT tasks_creator_fk
        FOREIGN KEY (workspace_id, created_by)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE RESTRICT,

    CONSTRAINT tasks_parent_fk
        FOREIGN KEY (project_id, parent_task_id)
        REFERENCES tasks(project_id, task_id)
        ON DELETE SET NULL (parent_task_id),

    -- Required for task dependency foreign keys.
    CONSTRAINT tasks_project_task_unique
        UNIQUE (project_id, task_id)
);


CREATE INDEX tasks_workspace_id_idx
ON tasks(workspace_id);

CREATE INDEX tasks_project_id_idx
ON tasks(project_id);

CREATE INDEX tasks_assignee_id_idx
ON tasks(assignee_id);

CREATE INDEX tasks_status_idx
ON tasks(workspace_id, project_id, status);

CREATE INDEX tasks_due_date_idx
ON tasks(workspace_id, due_date);


-- =========================================
-- 4. TASK DEPENDENCIES
-- =========================================

CREATE TABLE task_dependencies (
    workspace_id UUID NOT NULL,

    project_id UUID NOT NULL,

    task_id UUID NOT NULL,

    depends_on_task_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (task_id, depends_on_task_id),

    CONSTRAINT task_dependencies_task_fk
        FOREIGN KEY (project_id, task_id)
        REFERENCES tasks(project_id, task_id)
        ON DELETE CASCADE,

    CONSTRAINT task_dependencies_dependency_fk
        FOREIGN KEY (project_id, depends_on_task_id)
        REFERENCES tasks(project_id, task_id)
        ON DELETE CASCADE,

    CONSTRAINT task_dependencies_workspace_fk
        FOREIGN KEY (workspace_id, project_id)
        REFERENCES projects(workspace_id, project_id)
        ON DELETE CASCADE,

    CONSTRAINT task_dependencies_no_self_reference
        CHECK (task_id <> depends_on_task_id)
);


CREATE INDEX task_dependencies_depends_on_idx
ON task_dependencies(depends_on_task_id);


-- +goose Down

DROP TABLE IF EXISTS task_dependencies;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
