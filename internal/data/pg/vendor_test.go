package pg

import (
	"context"
	"testing"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestVendorService(t *testing.T) {
	t.Parallel()
	pool := testPool(t)
	defer pool.Close()

	us := &UserService{pool}
	vs := &VendorService{pool}
	ctx := context.Background()

	testEmail := "test@email.com"
	validPassword := "super_secret_password"
	//nonExistantID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000")

	var testUser *data.User
	var testVendor *data.Vendor
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
		testVendor, err = vs.Create(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, testVendor)
		assert.Equal(t, int32(1), testVendor.Version)
		assert.Equal(t, testUser.ID, testVendor.ID)
	})

	t.Run("TestRead", func(t *testing.T) {
		var readVendor *data.Vendor
		readVendor, err = vs.Read(ctx, testVendor.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readVendor)
		assert.Equal(t, testVendor, readVendor)
	})

	t.Run("TestUpdate", func(t *testing.T) {
		currentVersion := testVendor.Version
		err = vs.Update(ctx, testVendor)
		assert.NoError(t, err)
		assert.NotNil(t, testVendor)
		assert.Equal(t, currentVersion+1, testVendor.Version)

		var readVendor *data.Vendor
		readVendor, err = vs.Read(ctx, testVendor.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readVendor)
		assert.Equal(t, testVendor, readVendor)
	})

	t.Run("TestUserDelete", func(t *testing.T) {
		err = us.Delete(ctx, testVendor.ID)
		assert.NoError(t, err)

		_, err = vs.Read(ctx, testVendor.ID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})
}
