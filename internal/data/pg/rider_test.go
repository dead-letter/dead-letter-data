package pg

import (
	"context"
	"testing"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestRiderService(t *testing.T) {
	t.Parallel()
	pool := testPool(t)
	defer pool.Close()

	us := &UserService{pool}
	rs := &RiderService{pool}
	ctx := context.Background()

	testEmail := "test@email.com"
	validPassword := []byte("super_secret_password")

	var testUser *data.User
	var testRider *data.Rider
	var err error

	assert.Equal(t, "sarah", "sarah")

	t.Run("TestUserCreate", func(t *testing.T) {
		testUser, err = us.Create(ctx, testEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, int32(1), testUser.Version)
		assert.Equal(t, testEmail, testUser.Email)
	})

	t.Run("TestCreate", func(t *testing.T) {
		testRider, err = rs.Create(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, testRider)
		assert.Equal(t, int32(1), testRider.Version)
		assert.Equal(t, testUser.ID, testRider.ID)
	})

	t.Run("TestRead", func(t *testing.T) {
		var readRider *data.Rider
		readRider, err = rs.Read(ctx, testRider.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readRider)
		assert.Equal(t, testRider, readRider)
	})

	t.Run("TestUpdate", func(t *testing.T) {
		currentVersion := testRider.Version
		err = rs.Update(ctx, testRider)
		assert.NoError(t, err)
		assert.NotNil(t, testRider)
		assert.Equal(t, currentVersion+1, testRider.Version)

		var readRider *data.Rider
		readRider, err = rs.Read(ctx, testRider.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readRider)
		assert.Equal(t, testRider, readRider)
	})

	t.Run("TestUserDelete", func(t *testing.T) {
		err = us.Delete(ctx, testRider.ID)
		assert.NoError(t, err)

		_, err = rs.Read(ctx, testRider.ID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})
}
