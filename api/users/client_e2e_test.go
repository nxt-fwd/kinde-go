//go:build e2e
// +build e2e

package users_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nxt-fwd/kinde-go/api/users"
	"github.com/nxt-fwd/kinde-go/internal/client"
	"github.com/nxt-fwd/kinde-go/internal/testutil"
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

func TestE2EIdentities(t *testing.T) {
	client := users.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())
	email := fmt.Sprintf("%s@example.com", tempID)

	// First create a test user
	user, err := client.Create(context.TODO(), users.CreateParams{
		Profile: users.Profile{
			GivenName:  "Test",
			FamilyName: "User",
			Email:      email,
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, user)
	require.NotEmpty(t, user.ID)

	t.Logf("created test user: %s\n", user.ID)

	// Add a secondary email identity
	secondaryEmail := fmt.Sprintf("%s-secondary@example.com", tempID)
	emailIdentity, err := client.AddIdentity(context.TODO(), user.ID, users.AddIdentityParams{
		Type:  users.IdentityTypeEmail,
		Value: secondaryEmail,
	})
	assert.NoError(t, err)
	require.NotNil(t, emailIdentity)
	require.NotEmpty(t, emailIdentity.ID)

	t.Logf("added secondary email identity: %+v\n", emailIdentity)

	// Add a username identity
	username := fmt.Sprintf("user_%s", tempID)
	usernameIdentity, err := client.AddIdentity(context.TODO(), user.ID, users.AddIdentityParams{
		Type:  users.IdentityTypeUsername,
		Value: username,
	})
	assert.NoError(t, err)
	require.NotNil(t, usernameIdentity)
	require.NotEmpty(t, usernameIdentity.ID)

	t.Logf("added username identity: %+v\n", usernameIdentity)

	// Get and verify all identities
	identities, err := client.GetIdentities(context.TODO(), user.ID)
	assert.NoError(t, err)
	require.NotNil(t, identities)

	// We should have 2 identities (secondary email and username)
	assert.Equal(t, 2, len(identities), "expected 2 identities")

	// Verify the identities
	var foundSecondaryEmail, foundUsername bool
	for _, identity := range identities {
		if identity.Type == string(users.IdentityTypeEmail) && identity.Name == secondaryEmail {
			foundSecondaryEmail = true
		} else if identity.Type == string(users.IdentityTypeUsername) && identity.Name == username {
			foundUsername = true
		}
	}

	assert.True(t, foundSecondaryEmail, "secondary email identity not found")
	assert.True(t, foundUsername, "username identity not found")
	t.Logf("verified identities: %+v\n", identities)

	// Clean up - delete the test user
	err = client.Delete(context.TODO(), user.ID)
	assert.NoError(t, err)
	t.Logf("deleted test user: %s\n", user.ID)
}
