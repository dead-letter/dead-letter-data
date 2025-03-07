package pg

import (
	"testing"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestRiderService(t *testing.T) {
	t.Parallel()
	models, pool := testModels(t)
	defer pool.Close()

	testEmail := "test@email.com"
	validPassword := "super_secret_password"
	//nonExistantID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000")

	var testUser *data.User
	var testRider *data.Rider
	var err error

	assert.Equal(t, "sarah", "sarah")

	t.Run("TestUserCreate", func(t *testing.T) {
		testUser, err = models.User.Create(testEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, 1, testUser.Version)
		assert.Equal(t, testEmail, testUser.Email)
	})

	t.Run("TestCreate", func(t *testing.T) {
		testRider, err = models.Rider.Create(testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, testRider)
		assert.Equal(t, 1, testRider.Version)
		assert.Equal(t, testUser.ID, testRider.ID)
	})

	t.Run("TestRead", func(t *testing.T) {
		var readRider *data.Rider
		readRider, err = models.Rider.Read(testRider.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readRider)
		assert.Equal(t, testRider, readRider)
	})

	t.Run("TestUpdate", func(t *testing.T) {
		currentVersion := testRider.Version
		err = models.Rider.Update(testRider)
		assert.NoError(t, err)
		assert.NotNil(t, testRider)
		assert.Equal(t, currentVersion+1, testRider.Version)

		var readRider *data.Rider
		readRider, err = models.Rider.Read(testRider.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readRider)
		assert.Equal(t, testRider, readRider)
	})

	t.Run("TestUserDelete", func(t *testing.T) {
		err = models.User.Delete(testRider.ID)
		assert.NoError(t, err)

		_, err = models.Rider.Read(testRider.ID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})
}
