# Contract: API Schemas
## API Version
`/api/v1`
## Base Response
All responses follow this shape. Success responses include `data` (or the resource directly). Error responses include `error`.

```json
{ "message": "OK" }
{ "error": "Error description" }
```
## Authentication
All endpoints under `/api/v1/auth/*` are public.
All other endpoints require an active session cookie (`session_id`).
Unauthenticated requests return `{ "error": "Unauthorized" }` with HTTP 401.
## Endpoints

### Auth

#### POST /api/v1/auth/signup
Register a new user.

**Request body:**
```json
{
  "username": "string (3-30 chars, alphanumeric + underscore)",
  "email": "string (valid email)",
  "password": "string (min 8 chars)"
}
```

**Response 201:**
```json
{
  "message": "User registered successfully"
}
```

**Errors:** 400 (validation failed), 409 (username or email already exists)

#### POST /api/v1/auth/signin
Authenticate and create a session.

**Request body:**
```json
{
  "login": "string (username or email)",
  "password": "string"
}
```

**Response 200:**
```json
{
  "message": "Login successful",
  "user": {
    "id": "number",
    "username": "string",
    "email": "string",
    "avatar": "string | null",
    "bio": "string | null",
    "stacks": "string[] | null",
    "role": "string",
    "created_at": "string (ISO 8601)"
  }
}
```

**Errors:** 400 (validation failed), 401 (invalid credentials), 500 (server error)

#### POST /api/v1/auth/signout
Destroy the current session.

**Response 200:**
```json
{
  "message": "Logged out successfully"
}
```

**Errors:** 401 (not authenticated), 500 (server error)

#### GET /api/v1/auth/me
Get the currently authenticated user.

**Response 200:**
```json
{
  "user": {
    "id": "number",
    "username": "string",
    "email": "string",
    "avatar": "string | null",
    "bio": "string | null",
    "stacks": "string[] | null",
    "role": "string",
    "created_at": "string (ISO 8601)"
  }
}
```

**Errors:** 401 (not authenticated)

---

### Threads

#### POST /api/v1/threads
Create a new thread. Requires auth.

**Request body:**
```json
{
  "title": "string (5-255 chars)",
  "content": "string (markdown text, 10-50000 chars)",
  "tags": "string[] (existing tag names, optional, max 5)",
  "status": "string (draft | published, default: draft)"
}
```

**Response 201:**
```json
{
  "message": "Thread created",
  "thread": {
    "id": "number",
    "title": "string",
    "slug": "string",
    "content": "string",
    "status": "string",
    "view_count": "number",
    "upvotes": "number",
    "downvotes": "number",
    "replies_count": "number",
    "author": {
      "id": "number",
      "username": "string",
      "avatar": "string | null"
    },
    "tags": [{ "id": "number", "name": "string", "color": "string" }],
    "created_at": "string (ISO 8601)"
  }
}
```

**Errors:** 400 (validation failed), 401 (not authenticated)

#### GET /api/v1/threads
List published threads with pagination.

**Query params:**
| Param | Type | Default | Description |
|------|------|---------|-------------|
| page | number | 1 | Page number (1-indexed) |
| pageSize | number | 5 | Items per page (max 50) |
| sort | string | latest | Sort order: `latest`, `oldest`, `votes` |

**Response 200:**
```json
{
  "threads": [{
    "id": "number",
    "title": "string",
    "slug": "string",
    "status": "string",
    "view_count": "number",
    "upvotes": "number",
    "downvotes": "number",
    "replies_count": "number",
    "author": {
      "id": "number",
      "username": "string",
      "avatar": "string | null"
    },
    "tags": [{ "id": "number", "name": "string", "color": "string" }],
    "created_at": "string (ISO 8601)"
  }],
  "pagination": {
    "page": "number",
    "pageSize": "number",
    "total": "number",
    "totalPages": "number"
  }
}
```

**Errors:** 400 (invalid params)

#### GET /api/v1/threads/featured
Get the featured thread (highest vote score within last 7 days).

**Response 200:**
```json
{
  "thread": {
    "id": "number",
    "title": "string",
    "slug": "string",
    "content": "string",
    "status": "string",
    "view_count": "number",
    "upvotes": "number",
    "downvotes": "number",
    "replies_count": "number",
    "author": {
      "id": "number",
      "username": "string",
      "avatar": "string | null"
    },
    "tags": [{ "id": "number", "name": "string", "color": "string" }],
    "created_at": "string (ISO 8601)"
  }
}
```

**Errors:** 404 (no featured thread found)

#### GET /api/v1/threads/trending
Get top 3 trending threads by combined score (`upvotes - downvotes + replies_count * 0.5`).

**Response 200:**
```json
{
  "threads": [{ /* same shape as GET /api/v1/threads item */ }]
}
```

#### GET /api/v1/threads/:slug
Get a single thread by slug. Increments view_count. Requires auth for voting status.

**Response 200:**
```json
{
  "thread": {
    "id": "number",
    "title": "string",
    "slug": "string",
    "content": "string (markdown)",
    "status": "string",
    "view_count": "number",
    "upvotes": "number",
    "downvotes": "number",
    "replies_count": "number",
    "author": {
      "id": "number",
      "username": "string",
      "avatar": "string | null",
      "bio": "string | null",
      "stacks": "string[] | null",
      "role": "string"
    },
    "tags": [{ "id": "number", "name": "string", "color": "string" }],
    "comments": [{
      "id": "number",
      "content": "string",
      "upvotes": "number",
      "downvotes": "number",
      "author": {
        "id": "number",
        "username": "string",
        "avatar": "string | null"
      },
      "replies": [{
        "id": "number",
        "content": "string",
        "upvotes": "number",
        "downvotes": "number",
        "author": {
          "id": "number",
          "username": "string",
          "avatar": "string | null"
        },
        "replies": []
      }],
      "created_at": "string (ISO 8601)"
    }],
    "created_at": "string (ISO 8601)",
    "updated_at": "string (ISO 8601)"
  }
}
```

**Errors:** 404 (thread not found)

#### PATCH /api/v1/threads/:slug
Update a thread. Requires auth + ownership.

**Request body:**
```json
{
  "title": "string (optional, 5-255 chars)",
  "content": "string (optional, 10-50000 chars)",
  "tags": "string[] (optional, max 5)",
  "status": "string (optional, draft | published | archived)"
}
```

**Response 200:**
```json
{
  "message": "Thread updated",
  "thread": { /* same shape as POST /api/v1/threads response */ }
}
```

**Errors:** 400, 401, 403 (not owner), 404

#### DELETE /api/v1/threads/:slug
Soft-delete a thread. Requires auth + ownership.

**Response 200:**
```json
{
  "message": "Thread deleted"
}
```

**Errors:** 401, 403, 404

---

### Comments

#### POST /api/v1/threads/:slug/comments
Add a comment (top-level or reply) to a thread. Requires auth.

**Request body:**
```json
{
  "content": "string (2-10000 chars)",
  "parentId": "number (optional, existing comment ID for nested reply)"
}
```

**Response 201:**
```json
{
  "message": "Comment posted",
  "comment": {
    "id": "number",
    "content": "string",
    "upvotes": "number",
    "downvotes": "number",
    "author": {
      "id": "number",
      "username": "string",
      "avatar": "string | null"
    },
    "replies": [],
    "created_at": "string (ISO 8601)"
  }
}
```

**Errors:** 400, 401, 404 (thread or parent comment not found)

#### DELETE /api/v1/comments/:id
Soft-delete a comment. Requires auth + ownership or admin.

**Response 200:**
```json
{
  "message": "Comment deleted"
}
```

**Errors:** 401, 403, 404

---

### Votes

#### POST /api/v1/threads/:slug/vote
Vote on a thread. Requires auth. One vote per user per thread. Re-voting updates the value.

**Request body:**
```json
{
  "value": "number (1 for upvote, -1 for downvote, 0 to remove vote)"
}
```

**Response 200:**
```json
{
  "message": "Vote recorded",
  "upvotes": "number",
  "downvotes": "number"
}
```

**Errors:** 400 (invalid value), 401, 404

#### POST /api/v1/comments/:id/vote
Vote on a comment. Requires auth. One vote per user per comment. Re-voting updates the value.

**Request body:**
```json
{
  "value": "number (1 for upvote, -1 for downvote, 0 to remove vote)"
}
```

**Response 200:**
```json
{
  "message": "Vote recorded",
  "upvotes": "number",
  "downvotes": "number"
}
```

**Errors:** 400, 401, 404

---

### Users

#### GET /api/v1/users/:username
Get a user profile.

**Response 200:**
```json
{
  "user": {
    "id": "number",
    "username": "string",
    "avatar": "string | null",
    "bio": "string | null",
    "stacks": "string[] | null",
    "role": "string",
    "created_at": "string (ISO 8601)"
  }
}
```

**Errors:** 404

#### PATCH /api/v1/users/:username
Update a user profile. Requires auth + ownership.

**Request body:**
```json
{
  "avatar": "string (URL, optional)",
  "bio": "string (optional, max 500 chars)",
  "stacks": "string[] (optional, max 10 items)"
}
```

**Response 200:**
```json
{
  "message": "Profile updated",
  "user": { /* same shape as GET /api/v1/users/:username response */ }
}
```

**Errors:** 400, 401, 403, 404

#### GET /api/v1/users/:username/threads
Get all threads by a user.

**Query params:** same as GET /api/v1/threads

**Response 200:**
```json
{
  "threads": [{ /* same shape as GET /api/v1/threads item */ }],
  "pagination": {
    "page": "number",
    "pageSize": "number",
    "total": "number",
    "totalPages": "number"
  }
}
```

---

### Tags

#### GET /api/v1/tags
List all tags.

**Response 200:**
```json
{
  "tags": [{
    "id": "number",
    "name": "string",
    "color": "string"
  }]
}
```

#### POST /api/v1/tags
Create a new tag. Requires auth + admin role.

**Request body:**
```json
{
  "name": "string (3-50 chars)",
  "color": "string (hex color, optional, e.g., #FF5733)"
}
```

**Response 201:**
```json
{
  "message": "Tag created",
  "tag": {
    "id": "number",
    "name": "string",
    "color": "string"
  }
}
```

**Errors:** 400, 401, 403 (not admin), 409 (tag name exists)