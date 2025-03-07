package pg

import (
	"testing"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/stretchr/testify/assert"
)

var testUser *data.User

func TestCreate(t *testing.T) {
	var err error
	testUser, err = s.Create("johndoe@email.com", "super_secret_password")
	assert.NoError(t, err)
	assert.NotNil(t, testUser)
	assert.Equal(t, 1, testUser.Version)
	assert.Equal(t, "johndoe@email.com", testUser.Email)
}
