# Task Spec: task-3 — Thread CRUD Handlers
## Overview
Implement thread handlers: create, list (paginated), featured, trending, detail (by slug + view count increment), update, soft-delete. Slugs auto-generated from title.
## Implementation Steps
### 3.1 — Create `internal/lib/slug.go`
```go
import (
    "regexp"
    "strings"
)
var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)
func GenerateSlug(title string) string {
    slug := strings.ToLower(title)
    slug = nonAlphanumeric.ReplaceAllString(slug, "-")
    slug = strings.Trim(slug, "-")
    return slug
}
func GenerateUniqueSlug(title string, db *gorm.DB, model interface{}) (string, error) {
    base := GenerateSlug(title)
    slug := base
    for i := 1; ; i++ {
        var count int64
        db.Model(model).Where("slug = ?", slug).Count(&count)
        if count == 0 { return slug, nil }
        slug = fmt.Sprintf("%s-%d", base, i)
        if i > 10 { return fmt.Sprintf("%s-%s", base, randStr(4)), nil }
    }
}
func randStr(n int) string { ... } // n random alphanumeric chars
```
### 3.2 — Create `internal/handlers/threads.go`
**CreateThreadHandler** (protected): validate `{ title, content, tags, status }`, generate unique slug, look up tag IDs, create Thread with tags (via `db.Exec` or GORM association), return 201 with computed vote counts.
**ListThreadsHandler**: parse `page` (default 1), `pageSize` (default 5, max 50), `sort` (latest/oldest/votes). `WHERE status = 'published'`. Preload `Author`, `Tags`. Count total. Compute vote counts per thread. Return threads + pagination.
**FeaturedThreadHandler**: `status = 'published' AND created_at >= NOW() - 7 days`, sort by score DESC, limit 1. Preload `Author`, `Tags`. Compute vote counts. Return thread or 404.
**TrendingThreadsHandler**: `status = 'published'`, sort by score DESC, limit 3. Preload `Author`, `Tags`. Compute vote counts. Return threads array.
**GetThreadHandler** (protected): lookup by slug, increment `view_count` atomically via `gorm.Expr("view_count + 1")`, preload `Author`, `Tags`, `Comments` (top-level only), `Comments.Replies. Author`, `Comments.Author`. Compute vote counts. Return full thread.
**UpdateThreadHandler** (protected): lookup by slug, verify ownership (403), partial update (only non-nil fields), regenerate slug if title changed, update tag associations, save, return 200.
**DeleteThreadHandler** (protected): lookup by slug, verify ownership (403), `db.Delete(&thread)` (soft delete via GORM), return 200.
## Acceptance Criteria
1. Create thread with valid data → 201 with auto-generated slug.
2. Create thread with duplicate title → slug collision resolved with numeric suffix.
3. List threads → pagination metadata + correct sort.
4. Featured → highest score in last 7 days (or 404).
5. Trending → top 3 by score.
6. Get thread increments view count.
7. Update/delete by non-owner → 403.
8. Soft-deleted threads absent from listings.
9. `go build ./cmd/server` succeeds.
## Edge Cases
- Slug uniqueness: use `GenerateUniqueSlug` with retry loop.
- View count race: `gorm.Expr("view_count + 1")` is atomic in PostgreSQL.
- Sort by votes: for v1, sort by upvotes DESC as proxy. Add `// TODO: proper score sorting`.
- Comment nesting: `Preload("Comments. Replies. Author")` handles 1 level. For deeper, implement a recursive loader or limit to 2 levels with a `// TODO` note.