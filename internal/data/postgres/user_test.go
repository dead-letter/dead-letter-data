package postgres

import (
	"context"
	"testing"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	t.Parallel()
	pg := testDB(t)
	defer pg.Close()

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
		testUser, err = pg.Users.Create(ctx, testEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, int32(1), testUser.Version)
		assert.Equal(t, testEmail, testUser.Email)
	})

	t.Run("TestCreateDuplicateEmail", func(t *testing.T) {
		_, err = pg.Users.Create(ctx, testEmail, validPassword)
		assert.ErrorIs(t, err, data.ErrDuplicateEmail)
	})

	t.Run("TestRead", func(t *testing.T) {
		var readUser *data.User
		readUser, err = pg.Users.Read(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestReadUnkown", func(t *testing.T) {
		_, err = pg.Users.Read(ctx, nonExistantID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})

	t.Run("TestExistsWithEmail", func(t *testing.T) {
		exists, err := pg.Users.ExistsWithEmail(ctx, testEmail)
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = pg.Users.ExistsWithEmail(ctx, nonExistantEmail)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("TestReadWithCredentials", func(t *testing.T) {
		var readUser *data.User
		readUser, err = pg.Users.ReadWithCredentials(ctx, testEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestReadWithIncorrectCredentials", func(t *testing.T) {
		_, err = pg.Users.ReadWithCredentials(ctx, testEmail, incorrectPassword)
		assert.ErrorIs(t, err, data.ErrInvalidCredentials)
	})

	t.Run("TestUpdate", func(t *testing.T) {
		testUser.Email = updatedEmail
		currentVersion := testUser.Version
		err = pg.Users.Update(ctx, testUser)
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, currentVersion+1, testUser.Version)

		var readUser *data.User
		readUser, err = pg.Users.Read(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser, readUser)
	})

	t.Run("TestUpdateInvalidVersion", func(t *testing.T) {
		testUser.Version -= 1
		err = pg.Users.Update(ctx, testUser)
		assert.ErrorIs(t, err, data.ErrEditConflict)
	})

	t.Run("TestUpdateDuplicateEmail", func(t *testing.T) {
		var newUser *data.User
		newUser, err = pg.Users.Create(ctx, newEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, newUser)
		assert.Equal(t, newEmail, newUser.Email)

		newUser.Email = testEmail
		err = pg.Users.Update(ctx, testUser)
		assert.ErrorIs(t, err, data.ErrEditConflict)
	})

	t.Run("TestDelete", func(t *testing.T) {
		err = pg.Users.Delete(ctx, testUser.ID)
		assert.NoError(t, err)
	})

	t.Run("TestDeleteUnkown", func(t *testing.T) {
		err = pg.Users.Delete(ctx, nonExistantID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})
}
