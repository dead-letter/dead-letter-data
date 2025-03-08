package pg

import (
	"context"
	"testing"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	t.Parallel()
	pool := testPool(t)
	defer pool.Close()

	us := &UserService{pool}
	ctx := context.Background()

	testEmail := "test@email.com"
	updatedEmail := "updated@gmail.com"
	newEmail := "new@email.com"
	nonExistantEmail := "unknown@email.com"
	validPassword := []byte("super_secret_password")
	incorrectPassword := []byte("incorrect_password")
	nonExistantID := uuid.Nil

	var testUser *data.User
	var err error

	t.Run("TestCreate", func(t *testing.T) {
		testUser, err = us.Create(ctx, testEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, int32(1), testUser.Version)
		assert.Equal(t, testEmail, testUser.Email)
	})

	t.Run("TestCreateDuplicateEmail", func(t *testing.T) {
		_, err = us.Create(ctx, testEmail, validPassword)
		assert.ErrorIs(t, err, data.ErrDuplicateEmail)
	})

	t.Run("TestRead", func(t *testing.T) {
		var readUser *data.User
		readUser, err = us.Read(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestReadUnkown", func(t *testing.T) {
		_, err = us.Read(ctx, nonExistantID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})

	t.Run("TestReadWithEmail", func(t *testing.T) {
		var readUser *data.User
		readUser, err = us.ReadWithEmail(ctx, testEmail)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestReadWithNonExistantEmail", func(t *testing.T) {
		_, err := us.ReadWithEmail(ctx, nonExistantEmail)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})

	t.Run("TestReadWithCredentials", func(t *testing.T) {
		var readUser *data.User
		readUser, err = us.ReadWithCredentials(ctx, testEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestReadWithIncorrectCredentials", func(t *testing.T) {
		_, err = us.ReadWithCredentials(ctx, testEmail, incorrectPassword)
		assert.ErrorIs(t, err, data.ErrInvalidCredentials)
	})

	t.Run("TestUpdate", func(t *testing.T) {
		testUser.Email = updatedEmail
		currentVersion := testUser.Version
		err = us.Update(ctx, testUser)
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, currentVersion+1, testUser.Version)

		var readUser *data.User
		readUser, err = us.Read(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestUpdateInvalidVersion", func(t *testing.T) {
		testUser.Version -= 1
		err = us.Update(ctx, testUser)
		assert.ErrorIs(t, err, data.ErrEditConflict)
	})

	t.Run("TestUpdateDuplicateEmail", func(t *testing.T) {
		var newUser *data.User
		newUser, err = us.Create(ctx, newEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, newUser)
		assert.Equal(t, newEmail, newUser.Email)

		newUser.Email = testEmail
		err = us.Update(ctx, testUser)
		assert.ErrorIs(t, err, data.ErrEditConflict)
	})

	t.Run("TestDelete", func(t *testing.T) {
		err = us.Delete(ctx, testUser.ID)
		assert.NoError(t, err)
	})

	t.Run("TestDeleteUnkown", func(t *testing.T) {
		err = us.Delete(ctx, nonExistantID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})
}
