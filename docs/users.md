# Users, profiles, and settings

This guide covers the authentication and account APIs currently implemented in
`internal/tables/user` and the matching database records. It is written from the
frontend point of view: what to send, what comes back, and when to call each
endpoint.

## How the records fit together

An account is represented by three one-to-one records. They all refer to the
same internal `user_id`.

| Table | Purpose | Examples |
| --- | --- | --- |
| `users` | Sign-in identity and account lifecycle | email, password hash, stable public ID, created/updated dates, account status |
| `user_profiles` | Details people see and edit | display name, role title, pronouns, bio, avatar URL, availability, status message |
| `user_preferences` | Personal defaults for this website | active workspace, notification defaults, time zone, week start, theme |

The frontend should use the opaque UUID `user_id` only when an API needs the
internal account key. `public_id` (for example `ORB-AB12CD34`) is the
shareable, exact-match ID used by the directory endpoint. Email is used to
register and log in. Passwords and password hashes are never returned.

`users.username`, `users.full_name`, and `users.avatar_url` are legacy account
columns retained for compatibility with existing auth responses. New profile
editing is handled by `user_profiles`. Use `profile.display_name` and
`profile.avatar_url` to render the current profile UI.

The migration `20261020120000_add_user_profiles_and_preferences.sql` fills in
profile and preferences rows for existing users. New registrations create all
three records together. Defaults for new preferences are: notifications on,
weekly digest off, timezone `UTC`, week starts Monday (`1`), no active
workspace, and no saved theme.

## API basics

In development, the API is usually served at `http://localhost:8080`. Prefix
the paths below with `/api`. Send JSON request bodies with
`Content-Type: application/json`.

All `/api/user/...` routes require authentication. The server accepts either:

- The `orbit_access_token` HttpOnly cookie set by login; or
- `Authorization: Bearer <access-token>`.

For browser `fetch`, include `credentials: "include"` so the browser sends and
stores the cookies. The access cookie is marked Secure by the backend, so a
plain HTTP local setup may need local HTTPS or a development cookie setting
before browsers will send it. The current CORS config allows
`http://localhost:3000` with credentials.

Errors use this JSON shape:

```json
{
  "error": true,
  "message": "Description of the error"
}
```

## Suggested frontend flow

1. Register once with an email and password.
2. Log in. The browser stores the access and refresh cookies; the response also
   includes basic account information and the public ID.
3. When rendering the account menu or profile screen, call
   `GET /api/user/profile` and `GET /api/user/settings`.
4. Save profile edits with `PATCH /api/user/profile`; save settings with
   `PATCH /api/user/settings`. Send only fields that changed.
5. For a protected request that returns `401`, call `POST /api/auth/refresh`
   with credentials, then retry the original request once.
6. To find a person by their shared ID, call
   `GET /api/user/directory/:publicID` while signed in.

## Authentication endpoints

### Register — `POST /api/auth/register`

No authentication is required.

Send:

```json
{
  "email": "person@example.com",
  "password": "example-password"
}
```

The current request validation requires an email and a password between 3 and
15 characters. The service normalizes email to lowercase and trims whitespace.
The current password length is a backend constraint; coordinate any change to
password policy with the frontend validation.

On success (`201 Created`):

```json
{
  "error": false,
  "message": "Registration successful"
}
```

The frontend should then send the user to login. Registration does not return a
token or log the user in automatically.

Common errors: `400` invalid input, `409` email already registered, `500`
unexpected server error.

### Login — `POST /api/auth/login`

No authentication is required.

Send:

```json
{
  "email": "person@example.com",
  "password": "example-password"
}
```

On success (`201 Created`), the response sets the `orbit_access_token` and
`orbit_refresh_token` cookies. The tokens are HttpOnly; frontend JavaScript
should not try to read them. The JSON body is:

```json
{
  "error": false,
  "message": "Login successful",
  "user": {
    "user_id": "00000000-0000-0000-0000-000000000000",
    "public_id": "ORB-AB12CD34",
    "email": "person@example.com",
    "username": null,
    "full_name": null,
    "created_at": "2026-10-20T12:00:00Z"
  }
}
```

`avatar_url` is included if present in the legacy account record. Use the
profile endpoint for the editable profile fields. Common errors: `400` user not
found or token generation failure, `409` invalid credentials, `500`
unexpected server error.

### Refresh access token — `POST /api/auth/refresh`

Send no JSON body. Include cookies (`credentials: "include"` in browser fetch).
On success (`201 Created`), the server replaces the access cookie and returns:

```json
{
  "error": false,
  "message": "Refreshed Access token successful"
}
```

If refresh fails with `401`, send the user to login. The refresh route currently
validates the refresh cookie with the same token validator used for access
tokens; token-type separation is not implemented yet.

## Profile endpoints (`user_profiles`)

### Read my profile — `GET /api/user/profile`

No body. Requires authentication. Success (`200 OK`):

```json
{
  "error": false,
  "profile": {
    "user_id": "00000000-0000-0000-0000-000000000000",
    "public_id": "ORB-AB12CD34",
    "display_name": "Alex Morgan",
    "role_title": "Product designer",
    "pronouns": "they/them",
    "bio": "Designing useful tools for teams.",
    "avatar_url": null,
    "availability": "available",
    "status_message": null,
    "updated_at": "2026-10-20T12:00:00Z"
  }
}
```

Nullable fields may be returned as `null`.

### Edit my profile — `PATCH /api/user/profile`

Send only the fields the user changed. Omitted fields remain as they were.

```json
{
  "display_name": "Alex Morgan",
  "role_title": "Product designer",
  "pronouns": "they/them",
  "bio": "Designing useful tools for teams.",
  "availability": "available",
  "status_message": "In focus time"
}
```

Fields:

| Field | Type | Rules |
| --- | --- | --- |
| `display_name` | string | Cannot be empty when supplied |
| `role_title` | string or null | Personal job title; not workspace permission |
| `pronouns` | string or null | User-entered display detail |
| `bio` | string or null | Up to 500 characters |
| `avatar_url` | string or null | Stored URL; upload handling is not provided by this endpoint |
| `availability` | string | `available`, `away`, or `do_not_disturb` |
| `status_message` | string or null | Up to 160 characters |

Success is `200 OK` with the saved `profile` in the same shape as the read
response. Invalid values return `400`; an unauthenticated request returns
`401`.

**Current clearing behavior:** profile updates use SQL `COALESCE`, so omitted
fields and JSON `null` both leave the existing value unchanged. An empty string
is stored as an empty string for optional text fields. If the UI needs a true
“clear this field” action, the API needs an explicit clear representation (for
example, a separate list of fields to clear).

### Find a user — `GET /api/user/directory/:publicID`

Requires authentication. Replace `:publicID` with the full public ID, such as
`ORB-AB12CD34`; lookup is an exact match. No body. Success returns the same
profile shape as `GET /api/user/profile`. An unknown or inactive account returns
`404`.

## Preference endpoints (`user_preferences`)

These are personal defaults and are not workspace-specific notification
settings. A `workspace_preferences` table would be a separate design if users
need different notification choices in each workspace.

### Read my settings — `GET /api/user/settings`

No body. Requires authentication. Success (`200 OK`):

```json
{
  "error": false,
  "preferences": {
    "user_id": "00000000-0000-0000-0000-000000000000",
    "active_workspace_id": null,
    "email_notifications": true,
    "task_updates": true,
    "mentions": true,
    "weekly_digest": false,
    "timezone": "America/Los_Angeles",
    "week_starts_on": 1,
    "theme": "system",
    "updated_at": "2026-10-20T12:00:00Z"
  }
}
```

`week_starts_on` uses `0` for Sunday and `1` for Monday. `timezone` must be an
IANA time zone name. Theme is `light`, `dark`, or `system`, or `null` if no
preference is saved.

### Save settings — `PATCH /api/user/settings`

Send only fields that changed:

```json
{
  "email_notifications": true,
  "task_updates": false,
  "mentions": true,
  "weekly_digest": false,
  "timezone": "America/Los_Angeles",
  "week_starts_on": 1,
  "theme": "system",
  "active_workspace_id": "00000000-0000-0000-0000-000000000000"
}
```

Every field is optional. `active_workspace_id` must identify a workspace where
the signed-in user has active membership. Send `"active_workspace_id": ""` to
clear it. For other fields, JSON `null` behaves like omission (keeps the
existing value); there is no explicit clear operation for `theme` yet.

Success is `200 OK` with the full saved `preferences` object. Invalid time
zones, themes, week start values, or non-member workspaces return `400`.

## Frontend implementation notes

- Keep the profile form and settings form separate. Save each to its matching
  endpoint and merge the returned object into local UI state.
- Treat `user_id`, `public_id`, and `updated_at` as server-owned fields; do not
  send them in update bodies.
- Persist the theme in `user_preferences` for cross-device behavior. The
  browser can apply the returned theme locally.
- Use the returned `active_workspace_id` as the user's selected/default
  workspace. The backend checks membership when it is changed; workspace
  endpoints still perform their own authorization.
- For cross-origin browser requests, use credentials and ensure the frontend
  origin matches the server's CORS configuration.

