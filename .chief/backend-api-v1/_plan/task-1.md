# Task Spec: task-1 — Update GORM Models
## Overview
Update the GORM models to match the schema defined in contracts 001 and 002. Add missing fields, add the Session model, add computed helpers, and update main.go to auto-migrate all models.
## Implementation Steps
### 1.1 — Update `internal/models/models.go`
1. **User** struct — add two fields before `Sessions`:
   - `Stacks string gorm:"type:jsonb" json:"-"`
   - `Role string gorm:"size:20;default:'user'" json:"role"`
2. **Thread** struct — add two fields:
   - `Status string gorm:"size:20;default:'draft'" json:"status"`
   - `ViewCount uint gorm:"default:0" json:"view_count"`
3. **Comment** struct — ensure field tags include length constraints:
   - `Content string gorm:"type:text;size:10000" json:"content"` (max chars enforced at API level)
4. **Tag** struct — add default color:
   - `Color string gorm:"size:20;default:'#6b366f1'" json:"color"`
5. **Session** struct — add as new struct at end of file:
   ```
   type Session struct {
       ID        string         gorm:"primaryKey;size:255" json:"id"`
       CreatedAt time.Time      json:"created_at"`
       UpdatedAt time.Time      json:"updated_at"`
       Data     string         gorm:"type:text" json:"data"`
       UserID   uint           gorm:"index" json:"user_id"`
       User     User           gorm:"foreignKey:UserID" json:"user"`
       ExpiresAt time.Time      gorm:"index" json:"expires_at"`
   }
   ```
6. Add `Session` relationship to `User` struct: `Sessions []Session gorm:"foreignKey:UserID" json:"sessions,omitempty"`
7. Add `Vote` struct — change `Value` default: `Value int8 gorm:"default:0" json:"value"` (no tag needed — zero is the default)
8. Ensure `Tag.Threads` relationship is removed (already handled by `thread_tags` many2many — do NOT add new field)

### 1.2 — Add Vote Computation Helpers
Add after the structs in `models.go`:
```go
func (t *Thread) Upvotes(db *gorm.DB) int64 {
    var count int64
    db.Model(&Vote{}).Where("thread_id = ? AND value = 1", t.ID).Count(&count)
    return count
}
func (t *Thread) Downvotes(db *gorm.DB) int64 {
    var count int64
    db.Model(&Vote{}).Where("thread_id = ? AND value = -1", t.ID).Count(&count)
    return count
}
func (t *Thread) RepliesCount(db *gorm.DB) int64 {
    var count int64
    db.Model(&Comment{}).Where("thread_id = ? AND parent_id IS NULL", t.ID).Count(&count)
    return count
}
func (c *Comment) Upvotes(db *gorm.DB) int64 {
    var count int64
    db.Model(&Vote{}).Where("comment_id = ? AND value = 1", c.ID).Count(&count)
    return count
}
func (c *Comment) Downvotes(db *gorm.DB) int64 {
    var count int64
    db.Model(&Vote{}).Where("comment_id = ? AND value = -1", c.ID).Count(&count)
    return count
}
```
Note: These are instance methods that accept `*gorm.DB` as a parameter to avoid circular imports. They compute counts from the database.

### 1.3 — Update `cmd/server/main.go`
1. Add `models.Session{}` to the `AutoMigrate` call.
2. Replace inline handlers with a clean import of handler packages (do not add logic yet — just import).
3. Keep the `ErrorHandler` and existing middleware (recover, logger, cors).
4. Keep the `/` and `/health` routes.
5. Move all `/api/v1` routes to handler packages (stubs allowed for now).

## Acceptance Criteria
1. `go build ./cmd/server` succeeds with zero errors.
2. `gorm:"-"` tag on sensitive fields prevents accidental exposure.
3. All 6 models are registered with AutoMigrate.
4. Soft delete (`DeletedAt`) is indexed on all models.
5. Vote helpers accept `*gorm.DB` and return `int64`.
6. Session model uses `string` as primary key (UUID token).
7. No business logic in main.go — only wiring and stubs.

## Edge Cases
- Session model: if `ID` primary key causes issues, use `uint` with a UUID stored as a regular field and add a unique index instead. Verify with PostgreSQL compatibility.
- Vote helper: `Value` field default of `0` means no vote. Ensure the Go struct default (zero value) matches the database default. `int8` zero value is `0` — this is correct.