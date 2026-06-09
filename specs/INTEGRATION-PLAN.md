# API Integration Plan

## Executive Summary

The frontend has a polished UI built with **hardcoded mock data**, while the backend provides a full REST API. Only the Login page is integrated. This plan details what needs to be done to connect all frontend pages to the backend.

---

## Current State

### Frontend Pages and Their Data Sources

| Page | Current Data Source | Backend Endpoint Needed |
|------|-----------------|-------------------|
| `Login.tsx` | **WORKING** - Uses `useSignin`, `useSignup`, `useMe` hooks | All integrated |
| `Index.tsx` | **HARDCODED** - Mock arrays `trending[]`, `threads[]` | GET `/threads`, `/threads/featured`, `/threads/trending` |
| `ThreadDetail.tsx` | **HARDCODED** - Mock `replies[]` | GET `/threads/:slug` |
| `CreateEntry.tsx` | **NO SUBMIT** - UI only | POST `/threads` |
| `Profile.tsx` | **HARDCODED** - Mock user data | GET `/users/:username` |
| `Settings.tsx` | **NO SUBMIT** - UI only | PATCH `/users/:username` |

---

## Backend API Endpoints

### Auth (`/api/v1/auth/*`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/signup` | Register new user |
| POST | `/auth/signin` | Login user |
| POST | `/auth/signout` | Logout user |
| GET | `/auth/me` | Get current user profile |

### Threads (`/api/v1/threads/*`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/threads` | Create new thread |
| GET | `/threads` | List threads (paginated) |
| GET | `/threads/featured` | Get featured thread |
| GET | `/threads/trending` | Get top 3 trending threads |
| GET | `/threads/:slug` | Get single thread by slug |
| PATCH | `/threads/:slug` | Update thread (author only) |
| DELETE | `/threads/:slug` | Delete thread (author only) |

### Comments (`/api/v1/comments/*`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/threads/:slug/comments` | Create comment on thread |
| DELETE | `/comments/:id` | Delete comment |

### Votes (`/api/v1/*`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/threads/:slug/vote` | Vote on thread |
| POST | `/comments/:id/vote` | Vote on comment |

### Users (`/api/v1/users/*`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users/:username` | Get user profile |
| PATCH | `/users/:username` | Update own profile |
| GET | `/users/:username/threads` | Get user's threads |

### Tags (`/api/v1/tags`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tags` | List all tags |
| POST | `/tags` | Create tag (auth only) |

---

## Issues to Fix

### 1. `api.ts` Missing Methods

Current `src/lib/api.ts` only has `get` and `post`. Need to add:

```typescript
// src/lib/api.ts - Add these methods
patch: <T>(path: string, data?: unknown) => request<T>(path, { method: "PATCH", ... })
put: <T>(path: string, data?: unknown) => request<T>(path, { method: "PUT", ... })
delete: <T>(path: string) => request<T>(path, { method: "DELETE" })
```

### 2. Type Mismatch - User.stacks

| Location | Backend | Frontend (Current) | Frontend (Should) |
|---------|---------|-------------------|-------------------|
| `auth-store.ts` | `string[]` | `string` | `string[]` |

Backend sends stacks as JSON array, frontend should match.

### 3. Missing Frontend Hooks

Need to create these hooks in `src/hooks/`:

- `use-threads.ts` - fetch threads list, featured, trending
- `use-thread.ts` - fetch single thread by slug, CRUD operations
- `use-comments.ts` - create/delete comments
- `use-votes.ts` - voting on threads/comments

### 4. Hardcoded Data in Pages

These pages need API integration:

- **Index.tsx**: Replace hardcoded `threads[]`, `trending[]` with API calls
- **ThreadDetail.tsx**: Replace hardcoded `replies[]` with API call to `/threads/:slug`
- **CreateEntry.tsx**: Add form submit calling `POST /threads`
- **Profile.tsx**: Replace hardcoded user data with `/users/:username`
- **Settings.tsx**: Add form submit calling `PATCH /users/:username`

---

## Implementation Phases

### Phase 1: API Layer Fixes

**Tasks:**
1. Add `patch`, `put`, `delete` methods to `src/lib/api.ts`
2. Fix `User.stacks` type in `src/stores/auth-store.ts` to `string[]`

### Phase 2: Create Hooks (with Mock Fallback)

**Pattern for each hook:**
```typescript
// If API returns empty/fails, fall back to mock data for display
function useThreads() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["threads"],
    queryFn: () => api.get< threadsResponse>("/threads"),
  });

  // Fall back to mock data if API is empty or fails
  const threads = data?.threads?.length ? data.threads : MOCK_THREADS;
  const isEmpty = !data?.threads?.length;

  return { threads, isLoading, isEmpty };
}
```

**Tasks:**
1. Create `src/hooks/use-threads.ts`
   - `useThreads()` - list with pagination, sorting (falls back to mock if empty)
   - `useFeaturedThread()` - featured thread (falls back to mock if empty)
   - `useTrendingThreads()` - top 3 trending (falls back to mock if empty)
2. Create `src/hooks/use-thread.ts`
   - `useThread(slug)` - single thread (falls back to mock if not found)
   - `useCreateThread()` - create mutation (API only - no mock)
   - `useUpdateThread()` - update mutation
   - `useDeleteThread()` - delete mutation
3. Create `src/hooks/use-comments.ts`
   - `useCreateComment(slug)` - create comment
   - `useDeleteComment()` - delete comment
4. Create `src/hooks/use-votes.ts`
   - `useVoteThread(slug)` - vote on thread
   - `useVoteComment(id)` - vote on comment

### Phase 3: Page Integration

**Tasks:**
1. **Index.tsx**
   - Replace `trending` array with `GET /threads/trending`
   - Replace `threads` array with `GET /threads?page=1&pageSize=10`
2. **ThreadDetail.tsx**
   - Fetch thread by slug from `GET /threads/:slug`
   - Display comments from response
3. **CreateEntry.tsx**
   - Wire form to `POST /threads`
   - Add title, content, tags fields
   - Handle validation (title: 5-255 chars, content: 10-50000 chars)
4. **Profile.tsx**
   - Fetch user profile from `GET /users/:username`
   - Fetch user's threads from `GET /users/:username/threads`
5. **Settings.tsx**
   - Wire form to `PATCH /users/:username`
   - Fields: avatar, bio, stacks

---

## Response Format Examples

### GET /threads (list)
```json
{
  "threads": [
    {
      "id": 1,
      "title": "Thread Title",
      "slug": "thread-title",
      "content": "...",
      "status": "published",
      "view_count": 100,
      "upvotes": 10,
      "downvotes": 2,
      "replies_count": 5,
      "author": { "id": 1, "username": "user", "avatar": "..." },
      "tags": [{ "id": 1, "name": "tag", "color": "#6366f1" }],
      "created_at": "2026-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 5,
    "total": 100,
    "totalPages": 20
  }
}
```

### GET /threads/:slug (single)
```json
{
  "id": 1,
  "title": "Thread Title",
  "slug": "thread-title",
  "content": "...",
  "status": "published",
  "view_count": 101,
  "upvotes": 10,
  "downvotes": 2,
  "replies_count": 5,
  "author": { "id": 1, "username": "user", "avatar": "..." },
  "tags": [...],
  "comments": [
    {
      "id": 1,
      "content": "Comment text",
      "upvotes": 5,
      "downvotes": 0,
      "author": { "id": 2, "username": "commenter", "avatar": "..." },
      "replies": [...],
      "created_at": "2026-01-01T00:00:00Z"
    }
  ],
  "created_at": "2026-01-01T00:00:00Z"
}
```

### User Response
```json
{
  "id": 1,
  "username": "user",
  "email": "user@example.com",
  "avatar": "https://...",
  "bio": "User bio",
  "role": "user",
  "stacks": ["React", "Go", "PostgreSQL"],
  "created_at": "2026-01-01T00:00:00Z"
}
```

---

## Priority Order

1. **Highest**: Login already works - auth is integrated
2. **High**: Index.tsx and ThreadDetail.tsx - main content pages
3. **Medium**: CreateEntry.tsx - core functionality
4. **Medium**: Profile.tsx - user visibility
5. **Low**: Settings.tsx - user settings

---

## Notes

### Mock Fallback Pattern

For each data-fetching hook, implement this pattern:

```typescript
function useThreads() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["threads"],
    queryFn: () => api.get<ThreadsResponse>("/threads"),
  });

  // If API returns empty or fails, use mock data for display
  const threads = data?.threads?.length ? data.threads : MOCK_THREADS;
  const isEmpty = !data?.threads?.length;

  return { threads, isLoading, isEmpty, error };
}
```

**Rules:**
- **READS (GET)**: Fall back to mock if API returns empty or errors
- **WRITES (POST/PATCH/DELETE)**: Always use real API, no mock fallback
- Keep mock data in a separate file: `src/lib/mock-data.ts`

### Common Issues

- All authenticated routes require session cookie (handled by `credentials: "include"`)
- Thread listing supports pagination: `?page=1&pageSize=5`
- Thread listing supports sorting: `?sort=latest`, `?sort=oldest`, `?sort=votes`
- Votes: value `1` for upvote, `-1` for downvote, `0` to remove vote