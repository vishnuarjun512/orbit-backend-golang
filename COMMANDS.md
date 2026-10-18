# Go Backend Command Reference — Orbit

A practical command cheat sheet for working on the Orbit Go/Gin backend. It covers Go, Air, Goose migrations, PostgreSQL, testing, formatting, dependency management, and common Windows/macOS terminal commands.

> **Run commands from your backend project root** — the directory containing `go.mod`, unless a section says otherwise.
>
> Examples use `cmd/api` and `migrations`, matching your project setup. Replace database URLs and example values with your own. Never commit real credentials.

## 1. Quick daily workflow

### Start the API normally

```bash
go run ./cmd/api
```

### Start with Air hot reload

```bash
air
```

Air watches your files, rebuilds the API, and restarts it after changes. Don't run both commands at the same time if both would listen on port `4000`.

### Run all tests

```bash
go test ./...
```

### Format Go code

```bash
go fmt ./...
```

### Check for common Go issues

```bash
go vet ./...
```

### Build the API

macOS/Linux:

```bash
go build -o ./tmp/main ./cmd/api
```

Windows PowerShell:

```powershell
go build -o .\tmp\main.exe .\cmd\api
```

---

## 2. Go commands

| Command                              | What it does                                                       |
| ------------------------------------ | ------------------------------------------------------------------ |
| `go version`                         | Shows the installed Go version.                                    |
| `go env`                             | Shows Go environment settings.                                     |
| `go env GOPATH`                      | Shows the Go workspace path.                                       |
| `go env GOMOD`                       | Shows the active module file path.                                 |
| `go run ./cmd/api`                   | Builds and runs the API entry point.                               |
| `go build ./cmd/api`                 | Builds the API executable in the current directory.                |
| `go build ./...`                     | Builds all packages in the module.                                 |
| `go test ./...`                      | Runs tests across all packages.                                    |
| `go test ./internal/tables/projects` | Runs tests for one package.                                        |
| `go test -v ./...`                   | Runs tests with verbose output.                                    |
| `go test -race ./...`                | Checks for data races where supported.                             |
| `go test -cover ./...`               | Reports package test coverage.                                     |
| `go test ./... -count=1`             | Runs tests without reusing cached test results.                    |
| `go fmt ./...`                       | Formats Go source files.                                           |
| `go vet ./...`                       | Runs static checks for suspicious code.                            |
| `go list ./...`                      | Lists packages in the module.                                      |
| `go clean`                           | Cleans Go build caches and temporary files (use only when needed). |

### Why `./cmd/api` instead of `cmd/api`?

The `./` prefix tells Go to use the local package directory. For this project, use:

```bash
go run ./cmd/api
```

### Dependency commands

```bash
# Add a dependency (example)
go get github.com/gin-gonic/gin

# Update a dependency to a selected version
go get github.com/gin-gonic/gin@latest

# Remove unused dependencies and update go.mod/go.sum
go mod tidy

# Download dependencies declared by the module
go mod download

# Verify downloaded dependencies against go.sum
go mod verify

# Show module dependency graph
go mod graph

# Show why a module is required
go mod why github.com/jackc/pgx/v5
```

`go mod tidy` can change `go.mod` and `go.sum`; review the diff before committing.

---

## 3. Air — hot reload

Air rebuilds and restarts your Go app when watched files change.

### Install Air

```bash
go install github.com/air-verse/air@latest
```

Check the installation:

```bash
air -v
```

If your terminal says `air` is not recognized, make sure Go's binary directory is in `PATH`. On macOS/Linux, it is commonly `$(go env GOPATH)/bin`. Restart the terminal after changing `PATH`.

### Initialize configuration

```bash
air init
```

This creates `.air.toml`. Review it and set the build output to match your operating system.

Example configuration for **Windows**:

```toml
[build]
  cmd = "go build -o ./tmp/main.exe ./cmd/api"
  bin = "./tmp/main.exe"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["tmp", "vendor", "node_modules", ".git"]
  delay = 1000
  stop_on_error = true
```

Example configuration for **macOS/Linux**:

```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/api"
  bin = "./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["tmp", "vendor", "node_modules", ".git"]
  delay = 1000
  stop_on_error = true
```

If your `.air.toml` has an `entrypoint` setting, make sure it doesn't point to an old or extensionless Windows executable. For a typical setup, remove the `entrypoint` override or leave it empty so Air uses `bin`.

Start Air with the default config:

```bash
air
```

Choose a specific config file:

```bash
air -c .air.toml
```

Stop it with `Ctrl+C`.

---

## 4. Goose — SQL migrations

Goose manages versioned SQL migrations. Creating a migration file does **not** run it; `up` applies it to the database.

### Install Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Check it:

```bash
goose -version
```

### Create migration files

```bash
# Create a migration to add the projects table
goose -dir migrations create create_projects_table sql

# Create a migration to add a column
goose -dir migrations create add_project_color_to_projects sql

# Create a migration to add workspace membership
goose -dir migrations create create_workspace_members_table sql
```

These commands create timestamped `.sql` files under `migrations/`. Edit the generated file to include `-- +goose Up` and `-- +goose Down` sections.

Example:

```sql
-- +goose Up
ALTER TABLE projects
ADD COLUMN project_color TEXT NOT NULL DEFAULT 'bg-rose-500';

-- +goose Down
ALTER TABLE projects
DROP COLUMN project_color;
```

**Caution:** A `Down` migration can permanently remove data. Check it before running rollback commands.

### Set the database URL

macOS/Linux (zsh/bash):

```bash
export DATABASE_URL='postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require'
```

Windows PowerShell:

```powershell
$env:DATABASE_URL = 'postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require'
```

Use the connection URL appropriate for your database. For local PostgreSQL, SSL options may differ from hosted PostgreSQL.

### Run migrations

macOS/Linux:

```bash
goose -dir migrations postgres "$DATABASE_URL" up
goose -dir migrations postgres "$DATABASE_URL" up-by-one
goose -dir migrations postgres "$DATABASE_URL" down
goose -dir migrations postgres "$DATABASE_URL" status
goose -dir migrations postgres "$DATABASE_URL" version
goose -dir migrations postgres "$DATABASE_URL" redo
goose -dir migrations postgres "$DATABASE_URL" down-to 0
```

Windows PowerShell:

```powershell
goose -dir migrations postgres "$env:DATABASE_URL" up
goose -dir migrations postgres "$env:DATABASE_URL" status
```

Validate migration files:

```bash
goose -dir migrations validate
```

**Important:** `down-to 0`, `redo`, and `down` can remove tables or data depending on the migration. Don't run them against production unless you understand the impact and have a backup.

### Migration conventions

- Use one migration for one logical schema change.
- Prefer new migrations over editing migrations that have already been applied to shared databases.
- Keep both `Up` and `Down` sections.
- Test migrations against a development database first.
- Use `IF EXISTS` / `IF NOT EXISTS` only when it fits the intended schema and doesn't hide unexpected drift.
- Never commit database passwords or real connection strings.

---

## 5. PostgreSQL and database checks

### Check whether PostgreSQL is reachable

If `psql` is installed:

macOS/Linux:

```bash
psql "$DATABASE_URL"
```

Windows PowerShell:

```powershell
psql "$env:DATABASE_URL"
```

Inside `psql`, useful commands include:

```sql
-- List tables
\dt

-- Describe a table (psql meta-command)
\d projects

-- List schemas
\dn

-- List roles
\du

-- Quit psql
\q
```

These `\` commands are `psql` commands, not SQL statements; run them inside the `psql` shell.

### Common SQL checks

```sql
-- Inspect project columns
SELECT column_name, data_type, is_nullable
FROM information_schema.columns
WHERE table_name = 'projects'
ORDER BY ordinal_position;

-- Check the current database
SELECT current_database();

-- Check the current user
SELECT current_user;
```

### Docker Compose (if PostgreSQL runs in Docker)

```bash
# Start services in the background
docker compose up -d

# Show running containers
docker compose ps

# Follow logs
docker compose logs -f

# Follow logs for one service (replace db with the service name)
docker compose logs -f db

# Stop services but preserve containers and volumes
docker compose stop

# Stop and remove containers/network (named volumes normally remain)
docker compose down

# WARNING: also removes declared named volumes and database data
docker compose down -v
```

The service name might be `postgres`, `database`, or something else in your `compose.yaml`. Use `docker compose ps` to see the actual names.

---

## 6. Run and test API endpoints

### Health check

If your API has a `/health` endpoint:

```bash
curl -i http://localhost:4000/health
```

### Test a protected route

```bash
curl -i http://localhost:4000/api/workspace
```

A `401 Unauthorized` response can be expected when the endpoint requires the `access_token` cookie.

### Test CORS preflight

```bash
curl -i -X OPTIONS \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET" \
  http://localhost:4000/api/workspace
```

`curl` doesn't enforce browser CORS rules; it only displays the server's response.

### Send JSON in a POST request

```bash
curl -i -X POST http://localhost:4000/api/projects \
  -H "Content-Type: application/json" \
  -d '{"name":"Project 1","workspace_id":"YOUR_WORKSPACE_UUID","project_color":"bg-rose-500","due_date":null}'
```

Adjust the route and JSON fields to match your actual router and request struct. If the endpoint is protected, send a valid cookie as well.

### Include a cookie in curl

```bash
curl -i http://localhost:4000/api/workspace \
  -H "Cookie: access_token=YOUR_ACCESS_TOKEN"
```

Don't paste real tokens into shared logs or commit them to files.

---

## 7. Windows-specific commands

Run these in PowerShell unless noted.

### Find what is using port 4000

```powershell
netstat -ano | findstr :4000
```

The final column is the PID. Identify it before stopping it:

```powershell
Get-Process -Id 12345
```

Stop only the process you confirmed should be terminated:

```powershell
Stop-Process -Id 12345
```

### Check environment variables

```powershell
$env:DATABASE_URL
$env:PORT
```

Set a variable for the current terminal session:

```powershell
$env:PORT = "4000"
```

These assignments are temporary for that terminal session. Your Go app may instead load values from `.env` using `godotenv`.

### Inspect files

```powershell
Get-Content .air.toml
Get-ChildItem
Get-ChildItem -Recurse -Filter *.go
```

### Stop the running development server

Press `Ctrl+C` in the terminal running Air or the Go server.

---

## 8. macOS/Linux-specific commands

### Find what is using port 4000

```bash
lsof -i :4000
```

Stop a process only after identifying it:

```bash
kill <PID>
```

If it won't stop, investigate before using `kill -9`.

### Check environment variables

```bash
echo "$DATABASE_URL"
echo "$PORT"
```

Set one for the current shell:

```bash
export PORT=4000
```

### Inspect files

```bash
cat .air.toml
find . -name '*.go' -not -path './vendor/*'
```

---

## 9. Troubleshooting quick reference

| Error or symptom                                                 | First thing to check                                                                                |
| ---------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| `go: command not found` / not recognized                         | Install Go and check `PATH`.                                                                        |
| `go run cmd/api` resolves to the Go toolchain                    | Use `go run ./cmd/api` from the module root.                                                        |
| `air` not recognized                                             | Check Go's bin directory is in `PATH`.                                                              |
| Air says it cannot execute a non-`.exe` file on Windows          | Ensure build output and `bin` both use `main.exe`; inspect `entrypoint` and `full_bin`.             |
| `bind: address already in use`                                   | Find the process listening on the API port; don't run two server instances.                         |
| Goose can't connect                                              | Check `DATABASE_URL`, credentials, network access, SSL settings, and database availability.         |
| Goose migration isn't applied                                    | Run `goose ... status`, then `up`; creating the file alone doesn't apply it.                        |
| `relation does not exist`                                        | Check the table name, schema, database URL, and migration status.                                   |
| `column does not exist`                                          | Check the migration and the actual database schema.                                                 |
| `number of field descriptions must equal number of destinations` | Match selected SQL columns to `Scan()` destinations; avoid `SELECT *`.                              |
| `cannot scan date into *string`                                  | Use a compatible date type such as `pgtype.Date`, or explicitly cast to text if you only need text. |
| Browser reports CORS but curl works                              | Inspect browser Network details and extensions; remember curl does not enforce CORS.                |
| `401 Unauthorized` on a protected route                          | Check whether the expected cookie is present, valid, and not expired.                               |

---

## 10. Suggested daily workflow

1. Open a terminal at the backend root.
2. Confirm you're on the intended Git branch: `git status`.
3. Start PostgreSQL or other dependencies if needed: `docker compose up -d`.
4. Check migration status and apply pending migrations if needed.
5. Run the backend with `air`.
6. Test the endpoint using the frontend, `curl`, or your API client.
7. Run `go test ./...` after meaningful changes.
8. Run `go fmt ./...` and `go vet ./...`.
9. Review `git diff` before committing.

## 11. Safety reminders

- Never commit `.env` or real database credentials.
- Keep `.env.example` populated with placeholder values, not real secrets.
- Don't print passwords, access tokens, or refresh tokens in logs.
- Be careful with `goose down`, `redo`, `down-to 0`, and `docker compose down -v`.
- Prefer migrations for schema changes instead of manually altering shared or production databases.
- A successful build doesn't prove the endpoint works; test the behavior too.

---

## Quick command list

```bash
# Run
go run ./cmd/api

# Hot reload
air

# Create a Goose SQL migration
goose -dir migrations create migration_name sql

# Apply migrations (macOS/Linux)
goose -dir migrations postgres "$DATABASE_URL" up

# Migration status (macOS/Linux)
goose -dir migrations postgres "$DATABASE_URL" status

# Tests, format, static checks
go test ./...
go fmt ./...
go vet ./...

# Dependencies
go mod tidy
go mod download

# API health check
curl -i http://localhost:4000/health
```

For PowerShell, use `$env:DATABASE_URL` in place of `$DATABASE_URL` in Goose commands.
