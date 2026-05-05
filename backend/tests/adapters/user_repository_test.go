package adapters_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"community-forum/backend/internal/adapters/db"
	"community-forum/backend/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockGorm(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	require.NoError(t, err)
	return gormDB, mock
}

func TestGORMUserRepository_Create(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","email","password","avatar","bio","stacks","role") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), &domain.User{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hashed",
		Avatar:   "avatar.jpg",
		Bio:      "Hello!",
		Stacks:   []string{"Go", "React"},
		Role:     domain.RoleUser,
	})
	require.NoError(t, err)
}

func TestGORMUserRepository_Create_DBError(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnError(errors.New("constraint violation"))
	mock.ExpectRollback()

	err := repo.Create(context.Background(), &domain.User{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hash",
		Stacks:   []string{},
	})
	assert.Error(t, err)
}

func selectFromUsersWhere(col string) string {
	return `SELECT * FROM "users" WHERE ` + col + ` = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`
}

func TestGORMUserRepository_GetByID(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(selectFromUsersWhere(`"users"."id"`))).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "avatar", "bio", "stacks", "role"}).
			AddRow(1, "johndoe", "john@example.com", "hashed", "avatar.jpg", "Hello!", `["Go","React"]`, "user"))

	user, err := repo.GetByID(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, []string{"Go", "React"}, user.Stacks)
}

func TestGORMUserRepository_GetByID_NotFound(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(selectFromUsersWhere(`"users"."id"`))).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.GetByID(context.Background(), 999)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestGORMUserRepository_GetByUsername(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(selectFromUsersWhere("username"))).
		WithArgs("janedoe").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).
			AddRow(2, "janedoe", "jane@example.com"))

	user, err := repo.GetByUsername(context.Background(), "janedoe")
	require.NoError(t, err)
	assert.Equal(t, "janedoe", user.Username)
	assert.Equal(t, "jane@example.com", user.Email)
}

func TestGORMUserRepository_GetByUsername_NotFound(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(selectFromUsersWhere("username"))).
		WithArgs("nonexistent").
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.GetByUsername(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestGORMUserRepository_GetByEmail(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(selectFromUsersWhere("email"))).
		WithArgs("john@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).
			AddRow(1, "johndoe", "john@example.com"))

	user, err := repo.GetByEmail(context.Background(), "john@example.com")
	require.NoError(t, err)
	assert.Equal(t, "johndoe", user.Username)
}

func TestGORMUserRepository_GetByEmail_NotFound(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(selectFromUsersWhere("email"))).
		WithArgs("none@example.com").
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.GetByEmail(context.Background(), "none@example.com")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestGORMUserRepository_Update(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), &domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hash",
		Bio:      "Updated bio",
		Avatar:   "new.jpg",
		Stacks:   []string{"Go", "React"},
		Role:     domain.RoleUser,
	})
	require.NoError(t, err)
}

func TestGORMUserRepository_Update_NotFound(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), &domain.User{
		ID:       999,
		Username: "ghost",
		Email:    "ghost@example.com",
		Password: "hash",
		Stacks:   []string{},
	})
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestGORMUserRepository_StacksDeserialization(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(selectFromUsersWhere(`"users"."id"`))).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "stacks", "role"}).
			AddRow(1, "fullstack", "full@example.com", `["Go","React","TypeScript"]`, "user"))

	user, err := repo.GetByID(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, []string{"Go", "React", "TypeScript"}, user.Stacks)
}

func TestGORMUserRepository_EmptyStacks(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(selectFromUsersWhere(`"users"."id"`))).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "stacks", "role"}).
			AddRow(1, "user", "user@example.com", "", "user"))

	user, err := repo.GetByID(context.Background(), 1)
	require.NoError(t, err)
	assert.Empty(t, user.Stacks)
}

func TestGORMUserRepository_InvalidStacksJSON(t *testing.T) {
	gormDB, mock := newMockGorm(t)
	repo := db.NewGORMUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(selectFromUsersWhere(`"users"."id"`))).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "stacks", "role"}).
			AddRow(1, "user", "user@example.com", "not-json", "user"))

	_, err := repo.GetByID(context.Background(), 1)
	assert.Error(t, err)
}
