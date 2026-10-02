# Orbit Backend

Backend service for **Orbit**, a multi-tenant project and task
management platform inspired by tools such as Linear, ClickUp, Asana,
and Notion.

The backend is being built with **Go**, **Gin**, and **PostgreSQL**. It
uses `pgx/pgxpool` for database access and Goose for version-controlled
SQL migrations. Authentication is handled by the Go backend; Supabase is
currently used as the hosted PostgreSQL provider, not as the
authentication provider.

> **Project status:** Backend foundation and initial PostgreSQL schema
> are in place. The API currently includes a health-check endpoint and a
> PostgreSQL connection. Authentication and the remaining application
> APIs are under development.

------------------------------------------------------------------------

## Table of contents

-   [Technology stack](#technology-stack)
-   [Current features and progress](#current-features-and-progress)
-   [Project structure](#project-structure)
-   [Getting started](#getting-started)
-   [Environment variables](#environment-variables)
-   [Running the API](#running-the-api)
-   [Database and migrations](#database-and-migrations)
-   [Database schema](#database-schema)
    -   [Entity relationship overview](#entity-relationship-overview)
    -   [Users and authentication](#1-users-and-authentication)
    -   [Workspaces and membership](#2-workspaces-and-membership)
    -   [Projects and tasks](#3-projects-and-tasks)
    -   [Teams and tags](#4-teams-and-tags)
    -   [Collaboration and
        notifications](#5-collaboration-and-notifications)
-   [Multi-tenancy and data
    integrity](#multi-tenancy-and-data-integrity)
-   [Authentication design](#authentication-design)
-   [API endpoint](#api-endpoint)
-   [Development commands](#development-commands)
-   [Security notes](#security-notes)
-   [Roadmap](#roadmap)

------------------------------------------------------------------------

## Technology stack

  --------------------------------------------------------------------------------
  Area                    Technology                       Purpose
  ----------------------- -------------------------------- -----------------------
  Language                Go                               Backend application

  HTTP framework          Gin                              HTTP routing and
                                                           handlers

  Database                PostgreSQL                       Persistent application
                                                           data

  Database driver         pgx / pgxpool                    PostgreSQL connectivity
                                                           and connection pooling

  Database hosting        Supabase PostgreSQL              Hosted PostgreSQL
                                                           instance

  Configuration           godotenv                         Load local environment
                                                           variables

  Migrations              Goose                            Version-controlled SQL
                                                           schema migrations

  Development             Air                              Live-reload development
                                                           workflow

  Password hashing        Argon2id                         Password hashing for
                          (`golang.org/x/crypto/argon2`)   authentication
  --------------------------------------------------------------------------------

The project uses SQL migrations and direct `pgx` queries rather than an
ORM.

## Current features and progress

### Implemented

-   Go module and initial application structure.
-   Gin HTTP server.
-   Environment-based application configuration.
-   PostgreSQL connection using `pgxpool`.
-   Database connection health check during startup.
-   `/health` endpoint.
-   Goose migration workflow.
-   Initial relational database schema for users, authentication,
    workspaces, projects, tasks, teams, tags, comments, attachments,
    activity logs, and notifications.
-   Argon2id password-hashing utility (initial implementation).

### In progress / planned

-   User registration and login.
-   Access-token and refresh-token authentication.
-   Authentication and authorization middleware.
-   Workspace creation, membership, and invitations.
-   Project and task CRUD APIs.
-   Team and tag management.
-   Comments and attachments APIs.
-   Activity feed and notifications APIs.
-   Input validation, centralized error handling, and API tests.
-   Email verification and password reset flows.

## Project structure

The backend is organized by responsibility. The structure will grow as
the API is implemented.

``` text
orbit-backend-golang/
├── cmd/
│   └── api/
│       └── main.go             # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Environment-based configuration
│   ├── database/
│   │   └── postgres.go         # PostgreSQL connection pool
│   ├── security/
│   │   └── password.go         # Argon2id password hashing and verification
│   ├── handlers/               # HTTP handlers (to be implemented)
│   ├── services/               # Business logic (to be implemented)
│   ├── repositories/            # Database queries (to be implemented)
│   ├── models/                  # Application/domain models (to be implemented)
│   └── routes/                  # Route registration (to be implemented)
├── migrations/                  # Goose SQL migrations
├── .env                         # Local secrets (never commit)
├── .env.example                 # Environment variable template
├── .gitignore
├── go.mod
└── go.sum
```

The intended request flow is:

``` text
HTTP Request
    |
    v
Gin Route
    |
    v
Handler  -> parses request and returns HTTP response
    |
    v
Service  -> validates business rules and coordinates operations
    |
    v
Repository -> executes SQL using pgx
    |
    v
PostgreSQL
```

Handlers should not contain SQL queries or complex business logic.
Services should enforce application rules, and repositories should own
database access.

## Getting started

### Prerequisites

-   Go (use a version compatible with the dependencies in `go.mod`).
-   PostgreSQL database. The current development setup uses Supabase
    PostgreSQL.
-   Goose CLI.
-   Air (optional, for live reload).

Install the development tools if they are not already installed:

``` bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/air-verse/air@latest
```

Make sure Go's binary directory is included in your shell `PATH` so the
`goose` and `air` commands can be found.

### Clone the repository

``` bash
git clone <YOUR_GITHUB_REPOSITORY_URL>
cd orbit-backend-golang
```

### Install Go dependencies

``` bash
go mod download
```

### Configure environment variables

Create a local `.env` file in the project root. Use `.env.example` as a
template.

## Environment variables

The application configuration currently reads these values:

  ---------------------------------------------------------------------------------------------------------------
  Variable                Description             Example
  ----------------------- ----------------------- ---------------------------------------------------------------
  `APP_ENV`               Application environment `development`

  `APP_PORT`              HTTP server port        `8080`

  `DATABASE_URL`          PostgreSQL connection   `postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require`
                          URL                     
  ---------------------------------------------------------------------------------------------------------------

Example `.env.example`:

``` dotenv
APP_ENV=development
APP_PORT=8080
DATABASE_URL=postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require
```

Replace the database URL with your own PostgreSQL connection string. Do
not commit `.env`, real passwords, access tokens, API keys, or other
credentials.

## Running the API

From the project root:

``` bash
go run ./cmd/api
```

The default development server listens on port `8080`, unless overridden
by `APP_PORT`.

The health endpoint is:

``` http
GET /health
```

Expected response:

``` json
{
  "status": "success",
  "message": "API is running"
}
```

## Database and migrations

Database schema changes are managed using
[Goose](https://github.com/pressly/goose). Each migration is a numbered
SQL file in `migrations/`. Goose records applied migrations in its own
`goose_db_version` table.

### Current migration history

  -----------------------------------------------------------------------
  Migration                           Purpose
  ----------------------------------- -----------------------------------
  001                                 Users and authentication-related
                                      user fields

  002                                 Authentication sessions,
                                      email-verification tokens, and
                                      password-reset tokens

  003                                 Workspaces, workspace members, and
                                      workspace invitations

  004                                 Projects, project members, tasks,
                                      and task dependencies

  005                                 Teams, team members, tags, task
                                      tags, and project tags

  006                                 Task comments, task attachments,
                                      activity logs, and notifications
  -----------------------------------------------------------------------

The migration filenames may have descriptive names; the numbers identify
their execution order.

### Apply migrations

Run these commands from the project root:

``` bash
set -a
source .env
set +a

goose -dir ./migrations postgres "$DATABASE_URL" up
```

### Check migration status

``` bash
set -a
source .env
set +a

goose -dir ./migrations postgres "$DATABASE_URL" status
```

### Create a new migration

``` bash
goose -dir ./migrations create descriptive_migration_name sql
```

Add the schema change to the generated file, including both
`-- +goose Up` and `-- +goose Down` sections.

**Important:** Do not edit an already-applied migration to make schema
changes. Create a new migration instead. This keeps database history
reproducible across environments.

------------------------------------------------------------------------

# Database schema

The current schema contains **16 application tables**, plus Goose's
migration bookkeeping table (`goose_db_version`).

The database is designed around workspaces. A user can belong to
multiple workspaces, and each workspace can contain its own projects,
tasks, teams, and tags.

## Entity relationship overview

``` text
users
 ├── auth_sessions
 ├── email_verification_tokens
 ├── password_reset_tokens
 └── workspace_members
       |
       v
    workspaces
       ├── workspace_invitations
       ├── projects
       │    ├── project_members
       │    ├── project_tags ───────── tags
       │    └── tasks
       │         ├── task_dependencies
       │         ├── task_comments
       │         ├── task_attachments
       │         └── task_tags ─────── tags
       ├── teams
       │    └── team_members
       ├── activity_logs
       └── notifications
```

The diagram is a conceptual overview. The detailed relationships and
constraints are described below.

## 1. Users and authentication

### `users`

Stores the platform's user accounts.

  -----------------------------------------------------------------------
  Column                  Type                    Description
  ----------------------- ----------------------- -----------------------
  `user_id`               UUID                    Primary key; generated
                                                  automatically

  `email`                 VARCHAR(255)            User email address;
                                                  unique without case
                                                  sensitivity

  `username`              VARCHAR(50)             Public username; unique
                                                  without case
                                                  sensitivity

  `full_name`             VARCHAR(100)            User's display name

  `password_hash`         TEXT                    Argon2id password hash;
                                                  never store plaintext
                                                  passwords

  `avatar_url`            TEXT                    Optional
                                                  profile-picture URL

  `email_verified_at`     TIMESTAMPTZ             Email verification
                                                  timestamp; nullable

  `status`                VARCHAR(20)             `active`, `suspended`,
                                                  or `deleted`

  `created_at`            TIMESTAMPTZ             Account creation time

  `updated_at`            TIMESTAMPTZ             Last update time
  -----------------------------------------------------------------------

**Important constraints:** `user_id` is the primary key. Email and
username have unique indexes on `LOWER(...)`, so values differing only
by letter case are treated as duplicates.

### `auth_sessions`

Represents a user's authenticated session, typically one per device or
browser.

  -----------------------------------------------------------------------
  Column                  Type                    Description
  ----------------------- ----------------------- -----------------------
  `session_id`            UUID                    Primary key

  `user_id`               UUID                    References
                                                  `users.user_id`

  `refresh_token_hash`    TEXT                    Unique hash of the
                                                  refresh token; raw
                                                  token is not stored

  `device_name`           VARCHAR(100)            Optional device
                                                  description

  `user_agent`            TEXT                    Client user-agent
                                                  string

  `ip_address`            INET                    Client IP address

  `last_used_at`          TIMESTAMPTZ             Last refresh/use time

  `expires_at`            TIMESTAMPTZ             Session expiration time

  `revoked_at`            TIMESTAMPTZ             Set when the session is
                                                  revoked

  `created_at`            TIMESTAMPTZ             Session creation time
  -----------------------------------------------------------------------

Deleting a user cascades to their sessions. Indexes support looking up
sessions by user and expiration time.

### `email_verification_tokens`

Stores hashed, time-limited email-verification tokens.

  Column          Type          Description
  --------------- ------------- ----------------------------
  `token_id`      UUID          Primary key
  `user_id`       UUID          References `users.user_id`
  `token_hash`    TEXT          Unique token hash
  `expires_at`    TIMESTAMPTZ   Token expiration time
  `consumed_at`   TIMESTAMPTZ   Set after successful use
  `created_at`    TIMESTAMPTZ   Token creation time

### `password_reset_tokens`

Stores hashed, time-limited password-reset tokens.

  Column          Type          Description
  --------------- ------------- ----------------------------
  `token_id`      UUID          Primary key
  `user_id`       UUID          References `users.user_id`
  `token_hash`    TEXT          Unique token hash
  `expires_at`    TIMESTAMPTZ   Token expiration time
  `consumed_at`   TIMESTAMPTZ   Set after successful use
  `created_at`    TIMESTAMPTZ   Token creation time

Both token tables index `user_id`. Tokens should be random, single-use,
expire, and be stored as hashes rather than raw token values.

## 2. Workspaces and membership

### `workspaces`

A workspace is the main tenant boundary in Orbit. Projects, tasks,
teams, and tags belong to a workspace.

  Column           Type           Description
  ---------------- -------------- --------------------------------
  `workspace_id`   UUID           Primary key
  `name`           VARCHAR(100)   Workspace name
  `slug`           VARCHAR(100)   Globally unique workspace slug
  `description`    TEXT           Optional workspace description
  `created_by`     UUID           References the creating user
  `created_at`     TIMESTAMPTZ    Creation time
  `updated_at`     TIMESTAMPTZ    Last update time

`created_by` references `users.user_id` with `ON DELETE RESTRICT`. A
workspace cannot silently lose its creator through user deletion.

### `workspace_members`

Connects users to workspaces and assigns a role within each workspace.

  Column           Type          Description
  ---------------- ------------- ----------------------------------------
  `workspace_id`   UUID          References `workspaces.workspace_id`
  `user_id`        UUID          References `users.user_id`
  `role`           VARCHAR(20)   `owner`, `admin`, `member`, or `guest`
  `status`         VARCHAR(20)   `active` or `suspended`
  `joined_at`      TIMESTAMPTZ   Membership creation time

**Primary key:** (`workspace_id`, `user_id`). This prevents a user from
being added to the same workspace more than once.

A user's role is workspace-specific. The same user can be an owner in
one workspace and a member in another.

### `workspace_invitations`

Tracks invitations sent to people who may join a workspace.

  -----------------------------------------------------------------------
  Column                  Type                    Description
  ----------------------- ----------------------- -----------------------
  `invitation_id`         UUID                    Primary key

  `workspace_id`          UUID                    Workspace being joined

  `email`                 VARCHAR(255)            Invitee email address

  `role`                  VARCHAR(20)             Invited role: `admin`,
                                                  `member`, or `guest`

  `token_hash`            TEXT                    Unique hash of the
                                                  invitation token

  `invited_by`            UUID                    User who sent the
                                                  invitation; nullable

  `accepted_by`           UUID                    User who accepted;
                                                  nullable

  `status`                VARCHAR(20)             `pending`, `accepted`,
                                                  `revoked`, or `expired`

  `expires_at`            TIMESTAMPTZ             Invitation expiration

  `accepted_at`           TIMESTAMPTZ             Acceptance time

  `revoked_at`            TIMESTAMPTZ             Revocation time

  `created_at`            TIMESTAMPTZ             Creation time
  -----------------------------------------------------------------------

The invitation references its workspace and optionally the users who
sent and accepted it. The token is stored as a hash.

## 3. Projects and tasks

### `projects`

Stores projects within a workspace.

  -----------------------------------------------------------------------
  Column                  Type                    Description
  ----------------------- ----------------------- -----------------------
  `project_id`            UUID                    Primary key

  `workspace_id`          UUID                    Owning workspace

  `name`                  VARCHAR(150)            Project name

  `project_key`           VARCHAR(20)             Workspace-scoped
                                                  project identifier

  `description`           TEXT                    Optional project
                                                  description

  `status`                VARCHAR(30)             `active`, `on_hold`,
                                                  `completed`, or
                                                  `archived`

  `visibility`            VARCHAR(20)             `workspace` or
                                                  `private`

  `start_date`            DATE                    Optional start date

  `due_date`              DATE                    Optional due date

  `created_by`            UUID                    Creating user ID

  `created_at`            TIMESTAMPTZ             Creation time

  `updated_at`            TIMESTAMPTZ             Last update time
  -----------------------------------------------------------------------

**Important constraints:** - (`workspace_id`, `project_key`) is unique,
so project keys are unique within a workspace. - (`workspace_id`,
`project_id`) is also unique to support composite foreign keys from
workspace-owned records. - The due date cannot be earlier than the start
date when both are provided.

### `project_members`

Connects workspace members to projects.

  Column           Type          Description
  ---------------- ------------- -------------------------------------------
  `workspace_id`   UUID          Workspace context
  `project_id`     UUID          Project being joined
  `user_id`        UUID          Member being added
  `added_by`       UUID          Workspace member who added them; nullable
  `joined_at`      TIMESTAMPTZ   Membership time

**Primary key:** (`project_id`, `user_id`).

Composite foreign keys ensure that the project and its members belong to
the same workspace. A project member must already be a workspace member.

### `tasks`

Stores tasks and subtasks inside projects.

  -----------------------------------------------------------------------
  Column                  Type                    Description
  ----------------------- ----------------------- -----------------------
  `task_id`               UUID                    Primary key

  `workspace_id`          UUID                    Owning workspace

  `project_id`            UUID                    Parent project

  `parent_task_id`        UUID                    Optional parent task
                                                  for subtasks

  `title`                 VARCHAR(255)            Task title

  `description`           TEXT                    Optional task
                                                  description

  `status`                VARCHAR(30)             Defaults to `todo`;
                                                  intentionally flexible
                                                  for future custom
                                                  workflows

  `priority`              VARCHAR(20)             `low`, `medium`,
                                                  `high`, or `urgent`

  `assignee_id`           UUID                    Assigned workspace
                                                  member; nullable

  `created_by`            UUID                    Workspace member who
                                                  created the task

  `start_date`            TIMESTAMPTZ             Optional start time

  `due_date`              TIMESTAMPTZ             Optional due time

  `completed_at`          TIMESTAMPTZ             Completion time

  `position`              BIGINT                  Ordering value for
                                                  lists and boards

  `created_at`            TIMESTAMPTZ             Creation time

  `updated_at`            TIMESTAMPTZ             Last update time
  -----------------------------------------------------------------------

**Important constraints:** - A task must belong to a project in the same
workspace. - An assignee must be a member of the task's workspace. - The
creator must be a member of the task's workspace. - A parent task must
belong to the same project. - A task's due date cannot precede its start
date when both are present. - `priority` is constrained to the four
supported values. - (`workspace_id`, `task_id`) and (`project_id`,
`task_id`) are unique to support composite foreign keys.

Task `status` is not restricted to a fixed set in the database yet. This
leaves room for workspace- or project-specific workflow statuses to be
introduced later.

### `task_dependencies`

Represents dependencies between tasks. A task can depend on another task
in the same project.

  Column                 Type          Description
  ---------------------- ------------- -----------------------------------
  `workspace_id`         UUID          Workspace context
  `project_id`           UUID          Project context
  `task_id`              UUID          Task that has a dependency
  `depends_on_task_id`   UUID          Task that must be completed first
  `created_at`           TIMESTAMPTZ   Dependency creation time

**Primary key:** (`task_id`, `depends_on_task_id`).

The schema prevents a task from depending on itself and ensures both
tasks belong to the same project. Application logic should also prevent
longer dependency cycles.

## 4. Teams and tags

### `teams`

Stores named teams within a workspace.

  Column           Type           Description
  ---------------- -------------- ----------------------
  `team_id`        UUID           Primary key
  `workspace_id`   UUID           Owning workspace
  `name`           VARCHAR(100)   Team name
  `description`    TEXT           Optional description
  `created_by`     UUID           Creating user
  `created_at`     TIMESTAMPTZ    Creation time
  `updated_at`     TIMESTAMPTZ    Last update time

Team names are unique within a workspace without case sensitivity.
(`workspace_id`, `team_id`) is unique for composite references.

### `team_members`

Connects workspace members to teams.

  Column           Type          Description
  ---------------- ------------- -------------------
  `workspace_id`   UUID          Workspace context
  `team_id`        UUID          Team
  `user_id`        UUID          Team member
  `joined_at`      TIMESTAMPTZ   Membership time

**Primary key:** (`team_id`, `user_id`). Composite foreign keys ensure
the team and member belong to the same workspace.

### `tags`

Stores reusable, workspace-scoped tags.

  Column           Type          Description
  ---------------- ------------- ----------------------------------
  `tag_id`         UUID          Primary key
  `workspace_id`   UUID          Owning workspace
  `name`           VARCHAR(50)   Tag name
  `color`          VARCHAR(7)    Hex color, defaults to `#64748b`
  `created_at`     TIMESTAMPTZ   Creation time

Tag names are unique within a workspace without case sensitivity. The
color constraint requires a six-digit hexadecimal color in `#RRGGBB`
format.

### `task_tags`

Many-to-many relationship between tasks and tags.

  Column           Type          Description
  ---------------- ------------- ---------------------
  `workspace_id`   UUID          Workspace context
  `task_id`        UUID          Tagged task
  `tag_id`         UUID          Applied tag
  `created_at`     TIMESTAMPTZ   Tag assignment time

**Primary key:** (`task_id`, `tag_id`). Composite foreign keys ensure
the task and tag belong to the same workspace.

### `project_tags`

Many-to-many relationship between projects and tags.

  Column           Type          Description
  ---------------- ------------- ---------------------
  `workspace_id`   UUID          Workspace context
  `project_id`     UUID          Tagged project
  `tag_id`         UUID          Applied tag
  `created_at`     TIMESTAMPTZ   Tag assignment time

**Primary key:** (`project_id`, `tag_id`). Composite foreign keys ensure
the project and tag belong to the same workspace.

## 5. Collaboration and notifications

### `task_comments`

Stores comments and discussion on tasks.

  Column           Type          Description
  ---------------- ------------- ----------------------------------------
  `comment_id`     UUID          Primary key
  `workspace_id`   UUID          Workspace context
  `task_id`        UUID          Commented task
  `author_id`      UUID          Workspace member who wrote the comment
  `content`        TEXT          Comment content
  `edited_at`      TIMESTAMPTZ   Last edit time; nullable
  `deleted_at`     TIMESTAMPTZ   Soft-delete time; nullable
  `created_at`     TIMESTAMPTZ   Creation time
  `updated_at`     TIMESTAMPTZ   Last update time

Comments reference a task and an author in the same workspace. The
`deleted_at` column allows soft deletion so the application can preserve
a record instead of immediately removing it.

### `task_attachments`

Stores metadata and storage URLs for files attached to tasks.

  Column            Type           Description
  ----------------- -------------- ------------------------------------------
  `attachment_id`   UUID           Primary key
  `workspace_id`    UUID           Workspace context
  `task_id`         UUID           Attached task
  `uploaded_by`     UUID           Workspace member who uploaded the file
  `file_name`       VARCHAR(255)   Original or display filename
  `file_url`        TEXT           Object-storage URL or storage key
  `file_type`       VARCHAR(100)   MIME type; nullable
  `file_size`       BIGINT         File size in bytes; must be non-negative
  `created_at`      TIMESTAMPTZ    Upload record creation time

The database stores file metadata, not binary file contents. File
uploads will be handled by an object-storage service such as Amazon S3.

### `activity_logs`

Stores a historical record of actions performed in a workspace.

  -----------------------------------------------------------------------
  Column                  Type                    Description
  ----------------------- ----------------------- -----------------------
  `activity_id`           UUID                    Primary key

  `workspace_id`          UUID                    Workspace where the
                                                  action occurred

  `actor_id`              UUID                    Workspace member who
                                                  performed the action;
                                                  nullable

  `entity_type`           VARCHAR(30)             `workspace`, `project`,
                                                  `task`, or `comment`

  `entity_id`             UUID                    ID of the affected
                                                  entity

  `action`                VARCHAR(100)            Action name, such as
                                                  `task.status_changed`

  `metadata`              JSONB                   Additional action
                                                  details; defaults to an
                                                  empty JSON object

  `created_at`            TIMESTAMPTZ             Event time
  -----------------------------------------------------------------------

`metadata` can hold action-specific values. For example:

``` json
{
  "old_status": "todo",
  "new_status": "in_progress"
}
```

The `entity_id` is intentionally not a foreign key because it can refer
to different entity tables depending on `entity_type`. The application
must validate the target entity and workspace before writing an activity
record.

### `notifications`

Stores in-app notifications for workspace members.

  -----------------------------------------------------------------------
  Column                  Type                    Description
  ----------------------- ----------------------- -----------------------
  `notification_id`       UUID                    Primary key

  `workspace_id`          UUID                    Workspace context

  `recipient_id`          UUID                    Member receiving the
                                                  notification

  `actor_id`              UUID                    Member who triggered
                                                  it; nullable

  `task_id`               UUID                    Related task; nullable

  `project_id`            UUID                    Related project;
                                                  nullable

  `notification_type`     VARCHAR(50)             Machine-readable
                                                  notification category

  `title`                 VARCHAR(200)            Notification title

  `body`                  TEXT                    Optional notification
                                                  message

  `metadata`              JSONB                   Additional structured
                                                  information

  `read_at`               TIMESTAMPTZ             Read time; `NULL` means
                                                  unread

  `created_at`            TIMESTAMPTZ             Creation time
  -----------------------------------------------------------------------

A notification can reference a task or a project, but the database
constraint prevents it from referencing both at the same time. The
recipient must be a member of the specified workspace.

Indexes support listing notifications for a recipient and retrieving
unread notifications efficiently.

------------------------------------------------------------------------

## Multi-tenancy and data integrity

Orbit uses a **workspace-based multi-tenant model**.

-   A user can belong to multiple workspaces.
-   Workspace roles are stored in `workspace_members`, not globally on
    the user.
-   Projects belong to workspaces.
-   Tasks belong to projects and workspaces.
-   Teams and tags are scoped to workspaces.
-   Comments, attachments, and notifications carry workspace context.
-   Composite foreign keys are used where needed to prevent records from
    referencing resources in a different workspace.

### Important application-level rule

A `workspace_id` supplied by the frontend must never be treated as proof
of access.

For every protected request, the backend must:

1.  Authenticate the user.
2.  Determine the requested workspace.
3.  Verify that the authenticated user is an active member of that
    workspace.
4.  Check the user's role and resource-level permissions.
5.  Perform the database operation within that verified workspace
    context.

Database constraints provide an additional integrity layer, but they do
not replace authorization checks in the Go service layer.

Some relationships intentionally need application-level validation. For
example, the activity log's polymorphic `entity_id` is not backed by a
single foreign key, and the `created_by` field on teams references a
user rather than a workspace membership. The service layer must verify
those relationships.

## Authentication design

Authentication is intended to be handled by the Go API.

Planned flow:

-   Registration: validate input, hash the password with Argon2id, and
    create the user.
-   Login: verify the submitted password against the stored Argon2id
    hash.
-   Access token: issue a short-lived signed JWT.
-   Refresh token: issue a longer-lived random token and store only its
    hash in `auth_sessions`.
-   Logout: revoke the relevant session.
-   Email verification and password reset: use random, expiring,
    single-use tokens whose hashes are stored in their respective
    tables.

The current password utility is located at:

``` text
internal/security/password.go
```

It provides password hashing and verification helpers. Authentication
handlers, services, token issuance, and session management are still to
be implemented.

## API endpoint

  Method   Endpoint    Purpose                  Current status
  -------- ----------- ------------------------ ----------------
  `GET`    `/health`   Basic API health check   Implemented

Example:

``` bash
curl http://localhost:8080/health
```

Expected response:

``` json
{
  "status": "success",
  "message": "API is running"
}
```

## Development commands

Run the API:

``` bash
go run ./cmd/api
```

Run all Go package tests from the project root:

``` bash
go test ./...
```

Format Go files:

``` bash
gofmt -w .
```

Download dependencies:

``` bash
go mod download
```

Inspect module dependencies:

``` bash
go mod tidy
```

Create a migration:

``` bash
goose -dir ./migrations create descriptive_migration_name sql
```

Apply migrations:

``` bash
set -a
source .env
set +a
goose -dir ./migrations postgres "$DATABASE_URL" up
```

Check migration status:

``` bash
set -a
source .env
set +a
goose -dir ./migrations postgres "$DATABASE_URL" status
```

## Security notes

-   Never commit `.env` or real database credentials.
-   Keep `.env.example` limited to placeholder values.
-   Never store plaintext passwords. Use a password-hashing algorithm
    such as Argon2id.
-   Store hashes of refresh, verification, and reset tokens rather than
    their raw values.
-   Validate token expiry and single-use status when consuming
    verification or reset tokens.
-   Verify workspace membership and permissions on every protected
    operation.
-   Validate uploaded file type and size in the application; do not
    trust client-provided MIME types.
-   Use HTTPS in deployed environments.
-   Configure CORS with the specific frontend origins that should be
    allowed.
-   Use parameterized SQL queries with `pgx`; do not concatenate user
    input into SQL.
-   Avoid logging passwords, tokens, connection strings, or other
    secrets.

## Roadmap

-   [ ] Complete user repository and user model.
-   [ ] Implement registration with Argon2id.
-   [ ] Implement login and JWT access tokens.
-   [ ] Implement refresh-token rotation and logout.
-   [ ] Add authentication middleware.
-   [ ] Add request validation and centralized error responses.
-   [ ] Implement workspace creation and membership APIs.
-   [ ] Implement workspace invitations.
-   [ ] Implement project and project-member APIs.
-   [ ] Implement task CRUD, assignment, ordering, and dependencies.
-   [ ] Implement teams and tags.
-   [ ] Implement comments and attachment metadata APIs.
-   [ ] Implement activity feed and notifications.
-   [ ] Add pagination, filtering, search, and sorting.
-   [ ] Add rate limiting and security hardening.
-   [ ] Add unit and integration tests.
-   [ ] Add deployment configuration and observability.

------------------------------------------------------------------------

## License

Add a license before distributing this project publicly. Until then, no
license is specified.
