# Contract: Database Models
## ORM
GORM v2 with PostgreSQL. Auto-migrate runs on server startup.
## Model Definitions
### User
| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| ID | uint | GORM primaryKey | Auto-increment |
| CreatedAt | time.Time | — | — |
| UpdatedAt | time.Time | — | — |
| DeletedAt | gorm.DeletedAt | GORM index | Soft delete |
| Username | string | uniqueIndex, size:50 | Required, 3-30 chars |
| Email | string | uniqueIndex, size:255 | Required, valid email |
| Password | string | size:255 | bcrypt-hashed, never exposed |
| Avatar | string | size:255, nullable | Avatar URL |
| Bio | string | type:text, nullable | Max 500 chars |
| Stacks | string | type:jsonb, nullable | JSON-encoded string[] |
| Role | string | size:20, default:'user' | Enum: user, admin |
| Sessions | []Session | foreignKey:UserID | One-to-many |
| Threads | []Thread | foreignKey:AuthorID | One-to-many |
| Comments | []Comment | foreignKey:AuthorID | One-to-many |
| Votes | []Vote | foreignKey:UserID | One-to-many |
### Thread
| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| ID | uint | GORM primaryKey | Auto-increment |
| CreatedAt | time.Time | — | — |
| UpdatedAt | time.Time | — | — |
| DeletedAt | gorm.DeletedAt | GORM index | Soft delete |
| Title | string | size:255 | Required, 5-255 chars |
| Slug | string | uniqueIndex, size:255 | Auto-generated from title |
| Content | string | type:text | Markdown, 10-50000 chars |
| Status | string | size:20, default:'draft' | Enum: draft, published, archived |
| ViewCount | uint | default:0 | Incremented on GET /threads/:slug |
| AuthorID | uint | GORM index | FK → User |
| Author | User | foreignKey:AuthorID | Belongs-to |
| Tags | []Tag | many2many:thread_tags | Many-to-many |
| Comments | []Comment | foreignKey:ThreadID | One-to-many (top-level only) |
| Votes | []Vote | foreignKey:ThreadID | One-to-many |
### Comment
| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| ID | uint | GORM primaryKey | Auto-increment |
| CreatedAt | time.Time | — | — |
| UpdatedAt | time.Time | — | — |
| DeletedAt | gorm.DeletedAt | GORM index | Soft delete |
| Content | string | type:text | 2-10000 chars |
| ThreadID | uint | GORM index | FK → Thread |
| Thread | Thread | foreignKey:ThreadID | Belongs-to |
| AuthorID | uint | GORM index | FK → User |
| Author | User | foreignKey:AuthorID | Belongs-to |
| ParentID | *uint | GORM index, nullable | Self-reference for replies |
| Parent | *Comment | foreignKey:ParentID | Belongs-to |
| Replies | []Comment | foreignKey:ParentID | One-to-many (nested) |
| Votes | []Vote | foreignKey:CommentID | One-to-many |
### Tag
| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| ID | uint | GORM primaryKey | Auto-increment |
| CreatedAt | time.Time | — | — |
| UpdatedAt | time.Time | — | — |
| DeletedAt | gorm.DeletedAt | GORM index | Soft delete |
| Name | string | uniqueIndex, size:50 | Required, 3-50 chars |
| Color | string | size:20, default:'#6366f1' | Hex color |
| Threads | []Thread | many2many:thread_tags | Many-to-many |
### Vote
| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| ID | uint | GORM primaryKey | Auto-increment |
| CreatedAt | time.Time | — | — |
| UpdatedAt | time.Time | — | — |
| DeletedAt | gorm.DeletedAt | GORM index | Soft delete |
| UserID | uint | GORM index | FK → User |
| User | User | foreignKey:UserID | Belongs-to |
| ThreadID | *uint | GORM index, nullable | FK → Thread |
| Thread | *Thread | foreignKey:ThreadID | Belongs-to |
| CommentID | *uint | GORM index, nullable | FK → Comment |
| Comment | *Comment | foreignKey:CommentID | Belongs-to |
| Value | int8 | default:0 | 1 (up), -1 (down), 0 (removed) |
| **Constraint** | — | — | ThreadID XOR CommentID must be set |
### Session
| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| ID | string | GORM primaryKey | Session token (UUID) |
| CreatedAt | time.Time | — | — |
| UpdatedAt | time.Time | — | — |
| Data | string | type:text | Session data (JSON) |
| UserID | uint | GORM index | FK → User |
| User | User | foreignKey:UserID | Belongs-to |
| ExpiresAt | time.Time | GORM index | Expiry time |
## Indexes
| Table | Index | Columns | Type |
|-------|-------|---------|------|
| users | uq_users_username | username | unique |
| users | uq_users_email | email | unique |
| threads | uq_threads_slug | slug | unique |
| tags | uq_tags_name | name | unique |
| sessions | ix_sessions_expires_at | expires_at | btree |
| votes | ix_votes_user_thread | user_id, thread_id | unique (partial where thread_id IS NOT NULL) |
| votes | ix_votes_user_comment | user_id, comment_id | unique (partial where comment_id IS NOT NULL) |
## Vote Score Computations
| Metric | Formula |
|--------|---------|
| upvotes | COUNT(Vote) WHERE value=1 AND (ThreadID = target OR CommentID = target) |
| downvotes | COUNT(Vote) WHERE value=-1 AND (ThreadID = target OR CommentID = target) |
| score | upvotes - downvotes |
| replies_count | COUNT(Comment) WHERE threadId = target AND parentId IS NULL |
| featured score | upvotes + (replies_count * 0.5) WHERE createdAt >= NOW() - 7 days |
## Slug Generation
Convert title to lowercase, replace spaces with hyphens, remove non-alphanumeric chars except hyphens, trim hyphens from ends, append short random suffix if duplicate.
## Cascading Rules
- Deleting a User: sessions (cascade), threads (cascade soft-delete), comments (cascade soft-delete), votes (cascade soft-delete)
- Deleting a Thread: comments (cascade soft-delete), votes (cascade soft-delete)
- Deleting a Comment: replies (cascade soft-delete), votes (cascade soft-delete)
- Deleting a Tag: unlinked from threads via many2many