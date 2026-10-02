
-- +goose Up

-- =========================================
-- TASK COMMENTS
-- =========================================

CREATE TABLE task_comments (
    comment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL,
    task_id UUID NOT NULL,
    author_id UUID NOT NULL,
    content TEXT NOT NULL,
    edited_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT task_comments_task_fk
        FOREIGN KEY (workspace_id, task_id)
        REFERENCES tasks(workspace_id, task_id)
        ON DELETE CASCADE,

    CONSTRAINT task_comments_author_fk
        FOREIGN KEY (workspace_id, author_id)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE RESTRICT
);

CREATE INDEX task_comments_task_created_idx
ON task_comments(task_id, created_at DESC);

CREATE INDEX task_comments_author_idx
ON task_comments(author_id);


-- =========================================
-- TASK ATTACHMENTS
-- =========================================

CREATE TABLE task_attachments (
    attachment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL,
    task_id UUID NOT NULL,
    uploaded_by UUID NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_url TEXT NOT NULL,
    file_type VARCHAR(100),
    file_size BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT task_attachments_file_size_check
        CHECK (file_size >= 0),

    CONSTRAINT task_attachments_task_fk
        FOREIGN KEY (workspace_id, task_id)
        REFERENCES tasks(workspace_id, task_id)
        ON DELETE CASCADE,

    CONSTRAINT task_attachments_uploader_fk
        FOREIGN KEY (workspace_id, uploaded_by)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE RESTRICT
);

CREATE INDEX task_attachments_task_id_idx
ON task_attachments(task_id);


-- =========================================
-- ACTIVITY LOGS
-- =========================================

CREATE TABLE activity_logs (
    activity_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL
        REFERENCES workspaces(workspace_id) ON DELETE CASCADE,
    actor_id UUID,
    entity_type VARCHAR(30) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(100) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT activity_logs_entity_type_check
        CHECK (
            entity_type IN (
                'workspace',
                'project',
                'task',
                'comment'
            )
        ),

    CONSTRAINT activity_logs_actor_fk
        FOREIGN KEY (workspace_id, actor_id)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE SET NULL (actor_id)
);

CREATE INDEX activity_logs_workspace_created_idx
ON activity_logs(workspace_id, created_at DESC);

CREATE INDEX activity_logs_entity_idx
ON activity_logs(entity_type, entity_id);

CREATE INDEX activity_logs_actor_idx
ON activity_logs(actor_id);


-- =========================================
-- NOTIFICATIONS
-- =========================================

CREATE TABLE notifications (
    notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL,
    recipient_id UUID NOT NULL,
    actor_id UUID,
    task_id UUID,
    project_id UUID,
    notification_type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    body TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT notifications_recipient_fk
        FOREIGN KEY (workspace_id, recipient_id)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE CASCADE,

    CONSTRAINT notifications_actor_fk
        FOREIGN KEY (workspace_id, actor_id)
        REFERENCES workspace_members(workspace_id, user_id)
        ON DELETE SET NULL (actor_id),

    CONSTRAINT notifications_task_fk
        FOREIGN KEY (workspace_id, task_id)
        REFERENCES tasks(workspace_id, task_id)
        ON DELETE CASCADE,

    CONSTRAINT notifications_project_fk
        FOREIGN KEY (workspace_id, project_id)
        REFERENCES projects(workspace_id, project_id)
        ON DELETE CASCADE,

    CONSTRAINT notifications_target_check
        CHECK (num_nonnulls(task_id, project_id) <= 1)
);

CREATE INDEX notifications_recipient_created_idx
ON notifications(recipient_id, created_at DESC);

CREATE INDEX notifications_unread_idx
ON notifications(recipient_id, created_at DESC)
WHERE read_at IS NULL;

CREATE INDEX notifications_workspace_idx
ON notifications(workspace_id);


-- +goose Down

DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS activity_logs;
DROP TABLE IF EXISTS task_attachments;
DROP TABLE IF EXISTS task_comments;
