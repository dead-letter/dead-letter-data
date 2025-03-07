package pg

import (
	"testing"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

const ()

func TestUserService(t *testing.T) {
	t.Parallel()
	models, pool := testModels(t)
	defer pool.Close()

	testEmail := "test@email.com"
	updatedEmail := "updated@gmail.com"
	newEmail := "new@email.com"
	nonExistantEmail := "unknown@email.com"
	validPassword := "super_secret_password"
	incorrectPassword := "incorrect_password"
	nonExistantID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000")

	var testUser *data.User
	var err error

	t.Run("TestCreate", func(t *testing.T) {
		testUser, err = models.User.Create(testEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, 1, testUser.Version)
		assert.Equal(t, testEmail, testUser.Email)
	})

	t.Run("TestCreateDuplicateEmail", func(t *testing.T) {
		_, err = models.User.Create(testEmail, validPassword)
		assert.ErrorIs(t, err, data.ErrDuplicateEmail)
	})

	t.Run("TestRead", func(t *testing.T) {
		var readUser *data.User
		readUser, err = models.User.Read(testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestReadUnkown", func(t *testing.T) {
		_, err = models.User.Read(nonExistantID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})

	t.Run("TestReadWithEmail", func(t *testing.T) {
		var readUser *data.User
		readUser, err = models.User.ReadWithEmail(testEmail)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestReadWithNonExistantEmail", func(t *testing.T) {
		_, err := models.User.ReadWithEmail(nonExistantEmail)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})

	t.Run("TestReadWithCredentials", func(t *testing.T) {
		var readUser *data.User
		readUser, err = models.User.ReadWithCredentials(testEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestReadWithIncorrectCredentials", func(t *testing.T) {
		_, err = models.User.ReadWithCredentials(testEmail, incorrectPassword)
		assert.ErrorIs(t, err, data.ErrInvalidCredentials)
	})

	t.Run("TestUpdate", func(t *testing.T) {
		testUser.Email = updatedEmail
		currentVersion := testUser.Version
		err = models.User.Update(testUser)
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, currentVersion+1, testUser.Version)

		var readUser *data.User
		readUser, err = models.User.Read(testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestUpdateInvalidVersion", func(t *testing.T) {
		testUser.Version -= 1
		err = models.User.Update(testUser)
		assert.ErrorIs(t, err, data.ErrEditConflict)
	})

	t.Run("TestUpdateDuplicateEmail", func(t *testing.T) {
		var newUser *data.User
		newUser, err = models.User.Create(newEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, newUser)
		assert.Equal(t, 1, newUser.Version)
		assert.Equal(t, newEmail, newUser.Email)

		newUser.Email = testEmail
		err = models.User.Update(testUser)
		assert.ErrorIs(t, err, data.ErrEditConflict)
	})

	t.Run("TestDelete", func(t *testing.T) {
		err = models.User.Delete(testUser.ID)
		assert.NoError(t, err)
	})

	t.Run("TestDeleteUnkown", func(t *testing.T) {
		err = models.User.Delete(nonExistantID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})
}
