# Milestone TODO: Backend API v1
## Tasks
- [ ] **task-1**: Update GORM models — add User.stacks, User.role, Thread.status, Thread.viewCount, Session table. Add Vote score computation helpers. Update main. Go to auto- migrate all models.
- [ ] **task-2**: Build session auth — integrate Fiber sessions middleware with PostgreSQL session store. Implement auth handlers: signup (bcrypt), signin, signout, me. Create auth middleware for protected routes.
- [ ] **task-3**: Build thread CRUD handlers — create, list (paginated), featured, trending, detail (by slug + view count), update, soft-delete. Auto-generate slugs from title.
- [ ] **task-4**: Build comment and vote handlers — comment CRUD on threads, nested replies, vote toggle for threads and comments with upsert logic.
- [ ] **task-5**: Build user and tag handlers — user profile CRUD, threads-by-user listing, tag CRUD. Refactor main.go to delegate to handler service layer.