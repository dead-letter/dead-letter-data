package pg

import (
	"testing"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	var testUser *data.User

	t.Run("TestCreate", func(t *testing.T) {
		var err error

		testUser, err = models.User.Create("johndoe@email.com", "super_secret_password")
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, 1, testUser.Version)
		assert.Equal(t, "johndoe@email.com", testUser.Email)
	})

	t.Run("TestRead", func(t *testing.T) {
		var readUser *data.User
		var err error

		readUser, err = models.User.Read(testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})
}
