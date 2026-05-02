# Task Spec: task-3 — Thread CRUD Handlers
## Overview
Implement all thread handlers: create, list (paginated), featured, trending, detail (by slug + view count increment), update, soft-delete. Slugs are auto-generated from title.
## Implementation Steps
### 3.1 — Slug Generation Utility
Create `internal/lib/slug.go`:
```go
func GenerateSlug(title string) string {
    slug := strings.ToLower(title)
    slug = regex.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
    slug = strings.Trim(slug, "-")
    return slug
}
```
For duplicate slugs, generate with a random suffix (4 chars alphanumeric). In the Thread create handler, attempt to save; on unique violation, regenerate with suffix.

### 3.2 — Create `internal/handlers/threads.go`
Implement the following handlers:

**CreateThreadHandler** (protected):
1. Parse body: `{ title, content, tags, status }`.
2. Validate: title 5-255 chars, content 10-50000 chars, status in [draft, published] (default: draft).
3. Generate slug from title. If slug collision, append `-{random}` (4 chars).
4. Look up tag names → fetch Tag IDs from DB.
5. Create Thread record.
6. Associate tags via `db.Exec("INSERT INTO thread_tags (thread_id, tag_id) VALUES (?, ?)", thread.ID, tagIDs...)`.
7. Return 201 with thread data (computed vote counts via helpers).

**ListThreadsHandler**:
1. Parse query: `page` (default 1), `pageSize` (default 5, max 50), `sort` (default latest, options: latest, oldest, votes).
2. Query: `WHERE status = 'published'`. Apply sorting.
3. Preload `Author`, `Tags`.
4. Count total with the same filters (for pagination).
5. Compute vote counts per thread via Vote helpers.
6. Return threads + pagination metadata.

**FeaturedThreadHandler**:
1. Query: Thread with `status = 'published'` and `created_at >= NOW() - 7 days`, ordered by score (upvotes + replies_count * 0.5) DESC, limit 1.
2. Preload `Author`, `Tags`.
3. Compute vote counts.
4. Return 200 with thread, or 404 if none found.

**TrendingThreadsHandler**:
1. Query: `status = 'published'`, ordered by score DESC, limit 3.
2. Preload `Author`, `Tags`.
3. Compute vote counts.
4. Return 200 with threads array.

**GetThreadHandler** (protected):
1. Lookup by slug: `db.Preload(...).First(&thread, "slug = ?", slug)`.
2. If not found → 404.
3. Increment `ViewCount`: `db.Model(&thread).Update("view_count", gorm.Expr("view_count + 1"))`.
4. Preload `Author`, `Tags`, `Comments.Replies.Author`, `Comments.Author`.
5. Filter comments to top-level only (`parent_id IS NULL`).
6. Compute vote counts.
7. Return 200 with full thread data.

**UpdateThreadHandler** (protected):
1. Lookup thread by slug. If not found → 404.
2. Verify ownership: `thread.AuthorID == GetUserID(c)`. If not → 403.
3. Parse body — only update non-zero/non-empty fields (partial update).
4. If title changed → regenerate slug.
5. If tags changed → update associations.
6. Save and return updated thread.

**DeleteThreadHandler** (protected):
1. Lookup thread by slug. If not found → 404.
2. Verify ownership. If not → 403.
3. Soft delete: `db.Delete(&thread)` (GORM uses DeletedAt).
4. Return 200.

## Acceptance Criteria
1. Create thread with valid data → 201, auto-generated slug from title.
2. Create thread with duplicate title → slug collision resolved with suffix.
3. List threads returns pagination metadata, sorted correctly.
4. Featured returns the thread with highest score in last 7 days (or 404).
5. Trending returns top 3 threads by score.
6. Get thread increments view count.
7. Update thread by non-owner → 403.
8. Delete thread by non-owner → 403.
9. Soft-deleted threads do not appear in listings.

## Edge Cases
- Slug uniqueness: use GORM's `Clauses(clause.OnConflict{...})` for PostgreSQL upsert behavior, or catch the unique constraint error and retry with a suffix.
- View count race condition: use `gorm.Expr("view_count + 1")` for atomic increment.
- Pagination: ensure `PageSize > 50` is capped at 50.
- Sort by votes: compute score at query time using a subquery or join — this is expensive. For v1, it is acceptable to sort by upvotes descending as a proxy for score, with a `// TODO` note about proper score sorting.
- Thread detail with comments: GORM `Preload("Comments.Replies.Author")` works for one level of nesting. For deeper nesting, implement a recursive loader or limit to 2 levels.