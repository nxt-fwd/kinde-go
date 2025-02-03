//go:build e2e
// +build e2e

package users_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/axatol/kinde-go/api/users"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EList(t *testing.T) {
	client := users.New(testutil.DefaultE2EClient(t))
	res, err := client.List(context.TODO(), users.ListParams{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2ECreateGetUpdateDelete(t *testing.T) {
	client := users.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())
	email := fmt.Sprintf("%s@example.com", tempID)

	// Create user
	user, err := client.Create(context.TODO(), users.CreateParams{
		Profile: users.Profile{
			GivenName:  "Test",
			FamilyName: "User",
			Email:      email,
			ProvidedID: tempID,
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, user)
	require.NotEmpty(t, user.ID)

	id := user.ID
	t.Logf("created test user: %s\n", id)

	// Get user
	user, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, user)
	t.Logf("got test user: %+v\n", user)

	// Update user
	isSuspended := true
	updateParams := users.UpdateParams{
		GivenName:   "Updated",
		FamilyName:  "User",
		ProvidedID:  tempID + "-updated",
		IsSuspended: &isSuspended,
	}

	user, err = client.Update(context.TODO(), id, updateParams)
	assert.NoError(t, err)
	require.NotNil(t, user)

	// Verify the updated parameters
	assert.Equal(t, updateParams.GivenName, user.FirstName)
	assert.Equal(t, updateParams.FamilyName, user.LastName)
	assert.Equal(t, *updateParams.IsSuspended, user.IsSuspended)

	t.Logf("updated test user: %+v\n", user)

	// Delete user
	err = client.Delete(context.TODO(), id)
	assert.NoError(t, err)
	t.Logf("deleted test user: %+v\n", user)

	// Verify deletion
	_, err = client.Get(context.TODO(), id)
	assert.Error(t, err)
}
