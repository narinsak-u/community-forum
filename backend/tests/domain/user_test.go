package domain_test

import (
	"testing"

	"community-forum/backend/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestUserRoles(t *testing.T) {
	assert.Equal(t, "user", domain.RoleUser)
	assert.Equal(t, "admin", domain.RoleAdmin)
}

func TestErrNotFound(t *testing.T) {
	err := domain.ErrNotFound
	assert.EqualError(t, err, "user not found")
}

func TestUserStructDefaults(t *testing.T) {
	u := domain.User{}
	assert.Equal(t, uint(0), u.ID)
	assert.Equal(t, "", u.Username)
	assert.Equal(t, "", u.Email)
	assert.Equal(t, "", u.Password)
	assert.Equal(t, "", u.Avatar)
	assert.Equal(t, "", u.Bio)
	assert.Nil(t, u.Stacks)
	assert.Equal(t, "", u.Role)
	assert.True(t, u.CreatedAt.IsZero())
	assert.True(t, u.UpdatedAt.IsZero())
}
