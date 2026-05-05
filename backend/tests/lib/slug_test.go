package lib_test

import (
	"errors"
	"testing"

	"community-forum/backend/internal/lib"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple title", "Hello World", "hello-world"},
		{"special characters", "Go & React!", "go-react"},
		{"multiple spaces", "a   b   c", "a-b-c"},
		{"trailing hyphens", "-hello-", "hello"},
		{"uppercase", "UPPERCASE Title", "uppercase-title"},
		{"already slug", "hello-world", "hello-world"},
		{"numbers", "Hello 123 World", "hello-123-world"},
		{"apostrophe", "What's new", "what-s-new"},
		{"double hyphens", "a--b", "a-b"},
		{"multiple hyphens", "a---b", "a-b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.GenerateSlug(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateSlug_Empty(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"   ", ""},
		{"---", ""},
		{"!@#$", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := lib.GenerateSlug(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateUniqueSlug_FirstTry(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "threads" WHERE slug = \$1`).
		WithArgs("hello-world").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	slug, err := lib.GenerateUniqueSlug("Hello World", gormDB, "threads", "slug")
	require.NoError(t, err)
	assert.Equal(t, "hello-world", slug)
}

func TestGenerateUniqueSlug_WithCollision(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "threads" WHERE slug = \$1`).
		WithArgs("hello-world").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`SELECT count\(\*\) FROM "threads" WHERE slug = \$1`).
		WithArgs("hello-world-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	slug, err := lib.GenerateUniqueSlug("Hello World", gormDB, "threads", "slug")
	require.NoError(t, err)
	assert.Equal(t, "hello-world-1", slug)
}

func TestGenerateUniqueSlug_FallbackAfter10Collisions(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	expectedSlugs := []string{
		"hello-world",
		"hello-world-1", "hello-world-2", "hello-world-3",
		"hello-world-4", "hello-world-5", "hello-world-6",
		"hello-world-7", "hello-world-8", "hello-world-9",
	}

	for _, slug := range expectedSlugs {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "threads" WHERE slug = \$1`).
			WithArgs(slug).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "threads" WHERE slug = \$1`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	slug, err := lib.GenerateUniqueSlug("Hello World", gormDB, "threads", "slug")
	require.NoError(t, err)
	assert.Contains(t, slug, "hello-world-")
}

func TestGenerateUniqueSlug_EmptyTitleFallback(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "threads" WHERE slug = \$1`).
		WithArgs("thread").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	slug, err := lib.GenerateUniqueSlug("!@#$", gormDB, "threads", "slug")
	require.NoError(t, err)
	assert.Equal(t, "thread", slug)
}

func TestGenerateUniqueSlug_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "threads" WHERE slug = \$1`).
		WithArgs("hello-world").
		WillReturnError(errors.New("db error"))

	_, err = lib.GenerateUniqueSlug("Hello World", gormDB, "threads", "slug")
	assert.ErrorContains(t, err, "failed to check slug uniqueness")
}
