# Orbit Backend API Reference

This document describes the API routes currently registered by the
backend. Unless noted otherwise, the server listens on port `8080` in
development.

## Conventions

- JSON request bodies use `Content-Type: application/json`.
- Error responses use this shape:

  ```json
  {
    "error": true,
    "message": "Description of the error"
  }
  ```

- Workspace routes require a valid JWT. Send it either as the
  `orbit_access_token` cookie or as an
  `Authorization: Bearer <JWT>` header.
- The JWT is issued at login and expires after 15 minutes. Login sets
  the cookie as `HttpOnly` and `Secure`. The login JSON response does
  not currently include the JWT; browser clients can use the cookie.
- Workspace reads and mutations are scoped to the authenticated user
  who created the workspace. A different user receives `404 Not Found`.

## Health

### `GET /health`

Authentication: not required. No request body, path parameters, or
query parameters.

Success response: `200 OK`

```json
{
  "status": "Success",
  "message": "API is running"
}
```

## Authentication

### `POST /api/auth/register`

Authentication: not required.

Request body:

| Field | Type | Required | Constraints |
|---|---|---:|---|
| `email` | string | Yes | Valid email address |
| `password` | string | Yes | 3–15 characters |

Example:

```json
{
  "email": "person@example.com",
  "password": "example-password"
}
```

Success response: `201 Created`

```json
{
  "error": false,
  "message": "Registration successful"
}
```

Other responses:

| Status | Meaning | Response message |
|---|---|---|
| `400 Bad Request` | Invalid request fields | `Invalid registration details` |
| `409 Conflict` | Email is already registered | `Email is already registered` |
| `500 Internal Server Error` | Unexpected server error | `Something went wrong` |

### `POST /api/auth/login`

Authentication: not required.

Request body:

| Field | Type | Required | Constraints |
|---|---|---:|---|
| `email` | string | Yes | Valid email address |
| `password` | string | Yes | 3–15 characters |

Example:

```json
{
  "email": "person@example.com",
  "password": "example-password"
}
```

Success response: `201 Created`. The response also sets the
`orbit_access_token` cookie.

```json
{
  "error": false,
  "message": "Login successful",
  "user": {
    "user_id": "00000000-0000-0000-0000-000000000000",
    "email": "person@example.com",
    "username": null,
    "full_name": null,
    "created_at": "2026-10-07T00:00:00Z"
  }
}
```

`avatar_url` is included when present. Passwords and password hashes are
never returned.

Other responses:

| Status | Meaning | Response message |
|---|---|---|
| `400 Bad Request` | User was not found or token generation failed | `User not found` or `Token generation failed` |
| `409 Conflict` | Credentials are invalid | `Invalid Credentials` |
| `500 Internal Server Error` | Unexpected server error | `Something went wrong` |

## Workspaces

All workspace endpoints require JWT authentication.

### `POST /api/workspace`

Create a workspace. There are no path parameters or query parameters.
The authenticated user is recorded as its creator and added as an active
workspace owner.

Request body:

| Field | Type | Required | Constraints |
|---|---|---:|---|
| `name` | string | Yes | 3–100 characters |
| `description` | string | No | Optional description |

The slug is generated from `name`, lowercased, with runs of
non-letter/digit characters normalized to hyphens.

Example:

```json
{
  "name": "Product Team",
  "description": "Workspace for product planning"
}
```

Success response: `201 Created`

```json
{
  "error": false,
  "workspace": {
    "workspace_id": "00000000-0000-0000-0000-000000000000",
    "name": "Product Team",
    "slug": "product-team",
    "description": "Workspace for product planning",
    "created_by": "00000000-0000-0000-0000-000000000000",
    "created_at": "2026-10-07T00:00:00Z",
    "updated_at": "2026-10-07T00:00:00Z"
  }
}
```

Other responses include `400 Bad Request` for invalid fields,
`401 Unauthorized` for a missing/invalid JWT, `409 Conflict` if the
globally unique slug is taken, and `500 Internal Server Error` for an
unexpected server error.

### `GET /api/workspace`

List workspaces created by the authenticated user. No request body or
path parameters.

Query parameters:

| Parameter | Type | Required | Default / constraints |
|---|---|---:|---|
| `page` | integer | No | 1-based; default `1` |
| `limit` | integer | No | 1–100; default `20` |

Example: `GET /api/workspace?page=1&limit=20`

Success response: `200 OK`

```json
{
  "error": false,
  "workspaces": [
    {
      "workspace_id": "00000000-0000-0000-0000-000000000000",
      "name": "Product Team",
      "slug": "product-team",
      "description": "Workspace for product planning",
      "created_by": "00000000-0000-0000-0000-000000000000",
      "created_at": "2026-10-07T00:00:00Z",
      "updated_at": "2026-10-07T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20
  }
}
```

An empty result returns `"workspaces": []`. Invalid pagination values
return `400 Bad Request`; missing or invalid JWTs return
`401 Unauthorized`.

### `GET /api/workspace/:workspaceID`

Fetch one workspace. `workspaceID` is a UUID path parameter. No request
body or query parameters.

Success response: `200 OK`

```json
{
  "error": false,
  "workspace": {
    "workspace_id": "00000000-0000-0000-0000-000000000000",
    "name": "Product Team",
    "slug": "product-team",
    "description": "Workspace for product planning",
    "created_by": "00000000-0000-0000-0000-000000000000",
    "created_at": "2026-10-07T00:00:00Z",
    "updated_at": "2026-10-07T00:00:00Z"
  }
}
```

Other responses: `400 Bad Request` for an invalid UUID,
`401 Unauthorized` for a missing/invalid JWT, `404 Not Found` if the
workspace does not exist or was created by another user, or
`500 Internal Server Error` for an unexpected server error.

### `PATCH /api/workspace/:workspaceID`

Update a workspace. `workspaceID` is a UUID path parameter. No query
parameters. Include at least one field in the JSON body:

| Field | Type | Required | Constraints / behavior |
|---|---|---:|---|
| `name` | string | No | If provided, 3–100 characters |
| `slug` | string | No | If provided, up to 100 characters; normalized like the create slug |
| `description` | string | No | If provided as `""`, clears the description |

Example:

```json
{
  "name": "Product Design",
  "description": ""
}
```

Success response: `200 OK`, with the updated workspace using the same
shape as the get-workspace response.

Other responses: `400 Bad Request` for an invalid UUID or request,
`401 Unauthorized` for a missing/invalid JWT, `404 Not Found` if the
workspace does not exist or was created by another user, `409 Conflict`
if the slug is already taken, or `500 Internal Server Error` for an
unexpected server error.

### `DELETE /api/workspace/:workspaceID`

Delete a workspace. `workspaceID` is a UUID path parameter. No request
body or query parameters.

Success response: `204 No Content` with an empty response body.

Other responses: `400 Bad Request` for an invalid UUID,
`401 Unauthorized` for a missing/invalid JWT, `404 Not Found` if the
workspace does not exist or was created by another user, or
`500 Internal Server Error` for an unexpected server error.
