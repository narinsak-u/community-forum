# Hexagonal Architecture Phase 1: User & Auth Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the User and Authentication modules to follow Hexagonal Architecture, decoupling business logic from Fiber and GORM.

**Architecture:** We will implement the Ports and Adapters pattern. Core domain entities will live in `internal/domain`, interfaces in `internal/ports`, business logic in `internal/usecase`, and infrastructure implementations (HTTP/DB) in `internal/adapters`.

**Tech Stack:** Go 1.24, Fiber (HTTP), GORM (PostgreSQL), bcrypt (hashing).

---

### Task 1: Project Structure and User Domain Entity

**Files:**
- Create: `backend/internal/domain/user.go`
- Modify: `backend/internal/models/models.go` (remove User helper methods if any, though none currently use DB in User struct itself)

- [ ] **Step 1: Create the User domain entity**
Create `backend/internal/domain/user.go` with a clean User struct and Role constants.

```go
package domain

import "time"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Email     string
	Password  string
	Avatar    string
	Bio       string
	Stacks    []string // Domain uses actual slice, not JSON string
	Role      string
}
```

- [ ] **Step 2: Verify compilation**
Run: `go build ./internal/domain/...`
Expected: Success

- [ ] **Step 3: Commit**
```bash
git add backend/internal/domain/user.go
git commit -m "feat: add User domain entity"
```

---

### Task 2: User and Auth Port Definitions

**Files:**
- Create: `backend/internal/ports/user.go`
- Create: `backend/internal/ports/auth.go`

- [ ] **Step 1: Define User Repository and Service interfaces**
Create `backend/internal/ports/user.go`.

```go
package ports

import (
	"context"
	"community-forum/backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type UserService interface {
	GetUserProfile(ctx context.Context, username string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uint, updates *domain.User) (*domain.User, error)
}
```

- [ ] **Step 2: Define Auth Service interface**
Create `backend/internal/ports/auth.go`.

```go
package ports

import (
	"context"
	"community-forum/backend/internal/domain"
)

type AuthService interface {
	Signup(ctx context.Context, username, email, password string) error
	Signin(ctx context.Context, login, password string) (*domain.User, error)
}
```

- [ ] **Step 3: Commit**
```bash
git add backend/internal/ports/*.go
git commit -m "feat: define User and Auth ports"
```

---

### Task 3: User Repository GORM Implementation

**Files:**
- Create: `backend/internal/adapters/db/user_repository.go`

- [ ] **Step 1: Implement UserRepository using GORM**
Create `backend/internal/adapters/db/user_repository.go`. This will map between `domain.User` and `models.User`.

```go
package db

import (
	"context"
	"encoding/json"
	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/models"
	"gorm.io/gorm"
)

type GORMUserRepository struct {
	db *gorm.DB
}

func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{db: db}
}

func (r *GORMUserRepository) Create(ctx context.Context, u *domain.User) error {
	m := toModel(u)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *GORMUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return fromModel(&m), nil
}

// ... Implement GetByID, GetByEmail, Update similarly ...

func toModel(u *domain.User) *models.User {
	stacks, _ := json.Marshal(u.Stacks)
	return &models.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Stacks:    string(stacks),
		Role:      u.Role,
	}
}

func fromModel(m *models.User) *domain.User {
	var stacks []string
	json.Unmarshal([]byte(m.Stacks), &stacks)
	return &domain.User{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		Avatar:    m.Avatar,
		Bio:       m.Bio,
		Stacks:    stacks,
		Role:      m.Role,
	}
}
```

- [ ] **Step 2: Commit**
```bash
git add backend/internal/adapters/db/user_repository.go
git commit -m "feat: implement GORM UserRepository"
```

---

### Task 4: Auth Use Case Implementation

**Files:**
- Create: `backend/internal/usecase/auth_service.go`

- [ ] **Step 1: Implement AuthService with validation and hashing logic**
Create `backend/internal/usecase/auth_service.go`. Move logic from `auth.go` handler here.

```go
package usecase

import (
	"context"
	"errors"
	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo ports.UserRepository
}

func NewAuthService(repo ports.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Signup(ctx context.Context, username, email, password string) error {
	// Validation logic from AuthHandler.SignupHandler goes here
	// Password hashing goes here
	// Call s.repo.Create
	return nil
}

func (s *AuthService) Signin(ctx context.Context, login, password string) (*domain.User, error) {
	// Auth logic from AuthHandler.SigninHandler goes here
	// Call s.repo.GetByUsername or s.repo.GetByEmail
	// bcrypt.CompareHashAndPassword
	return nil, nil
}
```

- [ ] **Step 2: Commit**
```bash
git add backend/internal/usecase/auth_service.go
git commit -m "feat: implement Auth use case"
```

---

### Task 5: Refactor Auth and User Handlers

**Files:**
- Modify: `backend/internal/handlers/auth.go`
- Modify: `backend/internal/handlers/user.go`
- Modify: `backend/cmd/server/main.go` (to wire up dependencies)

- [ ] **Step 1: Update AuthHandler to use ports.AuthService**
Modify `backend/internal/handlers/auth.go` to remove DB dependency and call service.

- [ ] **Step 2: Update UserHandler to use ports.UserService**
Modify `backend/internal/handlers/user.go`.

- [ ] **Step 3: Update main.go to wire everything up**
```go
userRepo := db.NewGORMUserRepository(dbConn)
authService := usecase.NewAuthService(userRepo)
authHandler := handlers.NewAuthHandler(authService)
```

- [ ] **Step 4: Verify application starts**
Run: `go run ./cmd/server/main.go`
Expected: Server starts without errors.

- [ ] **Step 5: Commit**
```bash
git add .
git commit -m "refactor: finish Phase 1 hexagonal migration for User/Auth"
```
