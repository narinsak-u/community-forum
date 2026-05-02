# Task Spec: task-4 — Comment & Vote Handlers
## Overview
Implement comment CRUD (create, delete) on threads and vote toggling for threads and comments.
## Implementation Steps
### 4.1 — Create `internal/handlers/comments.go`
**CreateCommentHandler** (protected):
1. Parse body: `{ content, parentId }`.
2. Validate: content 2-10000 chars. If `parentId` set, verify parent exists and belongs to the same thread.
3. Enforce nesting depth: if parent has a parent, reject (max 1 level deep).
4. Create Comment record with `ThreadID` and `AuthorID`.
5. Return 201 with comment data.
**DeleteCommentHandler** (protected):
1. Lookup comment by ID. If not found → 404.
2. Verify ownership OR admin role. If not → 403.
3. Soft delete comment + all nested replies (`parent_id = comment.ID`).
4. Return 200.
### 4.2 — Create `internal/handlers/votes.go`
**VoteThreadHandler** (protected):
1. Parse body: `{ value }` where value ∈ {-1, 0, 1}. If invalid → 400.
2. Lookup thread by slug. If not found → 404.
3. Upsert: `db.Where("user_id = ? AND thread_id = ?", userID, thread.ID). First(&vote)`.
   - If value == 0: delete vote if exists.
   - If value != 0: find existing, update value, or create new.
4. Return 200 with `{ upvotes, downvotes }`.
**VoteCommentHandler** (protected):
1. Parse body: `{ value }`.
2. Lookup comment by ID. If not found → 404.
3. Same upsert logic.
4. Return 200 with `{ upvotes, downvotes }`.
### 4.3 — Partial Unique Indexes
Add via GORM AutoMigrate or raw SQL migration:
```sql
CREATE UNIQUE INDEX ix_votes_user_thread ON votes (user_id, thread_id) WHERE thread_id IS NOT NULL;
CREATE UNIQUE INDEX ix_votes_user_comment ON votes (user_id, comment_id) WHERE comment_id IS NOT NULL;
```
This ensures one vote per user per target at the DB level.
## Acceptance Criteria
1. Create top-level comment → 201.
2. Create reply (valid parentId) → 201, linked to parent.
3. Create reply (wrong thread parentId) → 400.
4. Create reply (3rd level) → 400 (nesting depth exceeded).
5. Vote upvote on thread → recorded, counts returned.
6. Re-vote upvote → updated.
7. Vote with value=0 → vote removed.
8. Delete own comment → 200 + replies deleted.
9. Delete other's comment → 403.
10. Vote with invalid value → 400.
## Edge Cases
- Vote upsert: `ON CONFLICT (user_id, thread_id) DO UPDATE SET value = ?`. Use GORM `Clauses(clause.OnConflict{...})`.
- Nested reply depth: enforce at handler level by checking `parent.ParentID != nil`.
- Race condition: PostgreSQL upsert is atomic — no special handling needed.
- Admin vote deletion: ownership OR admin role check in delete handler.