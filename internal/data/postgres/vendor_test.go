package postgres

import (
	"context"
	"testing"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestVendorService(t *testing.T) {
	t.Parallel()
	pg := testDB(t)
	defer pg.Close()

	ctx := context.Background()
	testEmail := "test@email.com"
	validPassword := []byte("super_secret_password")

	var testUser *data.User
	var testVendor *data.Vendor
	var err error

	assert.Equal(t, "sarah", "sarah")

	t.Run("TestUserCreate", func(t *testing.T) {
		testUser, err = pg.Users.Create(ctx, testEmail, validPassword)
		assert.NoError(t, err)
		assert.NotNil(t, testUser)
		assert.Equal(t, int32(1), testUser.Version)
		assert.Equal(t, testEmail, testUser.Email)
	})

	t.Run("TestCreate", func(t *testing.T) {
		testVendor, err = pg.Vendors.Create(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, testVendor)
		assert.Equal(t, int32(1), testVendor.Version)
		assert.Equal(t, testUser.ID, testVendor.ID)
	})

	t.Run("TestRead", func(t *testing.T) {
		var readVendor *data.Vendor
		readVendor, err = pg.Vendors.Read(ctx, testVendor.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readVendor)
		assert.Equal(t, testVendor, readVendor)
	})

	t.Run("TestUpdate", func(t *testing.T) {
		currentVersion := testVendor.Version
		err = pg.Vendors.Update(ctx, testVendor)
		assert.NoError(t, err)
		assert.NotNil(t, testVendor)
		assert.Equal(t, currentVersion+1, testVendor.Version)

		var readVendor *data.Vendor
		readVendor, err = pg.Vendors.Read(ctx, testVendor.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readVendor)
		assert.Equal(t, testVendor, readVendor)
	})

	t.Run("TestUserDelete", func(t *testing.T) {
		err = pg.Users.Delete(ctx, testVendor.ID)
		assert.NoError(t, err)

		_, err = pg.Vendors.Read(ctx, testVendor.ID)
		assert.ErrorIs(t, err, data.ErrRecordNotFound)
	})
}
