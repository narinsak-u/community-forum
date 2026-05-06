# Task Spec: task-1 — Update GORM Models
## Overview
Update the GORM models to match the schema defined in contracts 001 and 002. Add missing fields, add the Session model, add computed helpers, and update main.go to auto-migrate all models.
## Implementation Steps
### 1.1 — Update `internal/models/models.go`
1. **User** struct — add two fields: `Stacks string gorm:"type:jsonb" json:"-"` and `Role string gorm:"size:20;default:'user'" json:"role"`.
2. **Thread** struct — add: `Status string gorm:"size:20;default:'draft'" json:"status"` and `ViewCount uint gorm:"default:0" json:"view_count"`.
3. **Comment** struct — ensure: `Content string gorm:"type:text" json:"content"`.
4. **Tag** struct — update `Color`: `Color string gorm:"size:20;default:'#6b366f1'" json:"color"`.
5. **Session** struct — add as new struct:
```
type Session struct {
    ID        string         gorm:"primaryKey;size:255" json:"id"`
    CreatedAt time. Time      json:"created_at"`
    UpdatedAt time. Time      json:"updated_at"`
    Data     string         gorm:"type:text" json:"data"`
    UserID   uint           gorm:"index" json:"user_id"`
    User     User           gorm:"foreignKey:UserID" json:"user"`
    ExpiresAt time. Time      gorm:"index" json:"expires_at"`
}
```
6. Add `Sessions []Session gorm:"foreignKey:UserID" json:"sessions,omitempty"` to User struct.
7. **Vote** struct — change `Value` default: `Value int8 gorm:"default:0" json:"value"` (zero = no vote).
### 1.2 — Add Vote Computation Helpers
Add after structs in `models.go`:
```go
func (t *Thread) Upvotes(db *gorm.DB) int64 { ... }
func (t *Thread) Downvotes(db *gorm.DB) int64 { ... }
func (t *Thread) RepliesCount(db *gorm.DB) int64 { ... }
func (c *Comment) Upvotes(db *gorm.DB) int64 { ... }
func (c *Comment) Downvotes(db *gorm.DB) int64 { ... }
```
Instance methods accepting `*gorm.DB` to avoid circular imports. Return `int64`.
### 1.3 — Update `cmd/server/main.go`
1. Add `models.Session{}` to `AutoMigrate`.
2. Remove inline handlers — replace with import statements for handler packages (stubs allowed for now).
3. Keep `ErrorHandler`, middleware (recover, logger, cors), `/` and `/health` routes.
## Acceptance Criteria
1. `go build ./cmd/server` succeeds with zero errors.
2. All 6 models are registered with AutoMigrate.
3. Soft delete (`DeletedAt`) is indexed on all models.
4. Vote helpers accept `*gorm.DB` and return `int64`.
5. Session model uses `string` as primary key.
6. No business logic in main.go.
## Edge Cases
- Session `ID` as primary key (string UUID): verify PostgreSQL compatibility. If issues arise, use `uint` with `gorm:"uniqueIndex"` instead.
- Vote `int8` default: Go zero value is `0` — matches database default. Correct.