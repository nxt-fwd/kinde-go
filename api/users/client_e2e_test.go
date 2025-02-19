//go:build e2e
// +build e2e

package users_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nxt-fwd/kinde-go/api/identities"
	"github.com/nxt-fwd/kinde-go/api/users"
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
	identitiesClient := identities.New(testutil.DefaultE2EClient(t))
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

	// Add primary email identity
	primaryIdentity, err := client.AddIdentity(context.TODO(), user.ID, users.AddIdentityParams{
		Type:  users.IdentityTypeEmail,
		Value: email,
	})
	assert.NoError(t, err)
	require.NotNil(t, primaryIdentity)
	require.NotEmpty(t, primaryIdentity.ID)

	// Add a secondary email identity
	secondaryEmail := fmt.Sprintf("%s-secondary@example.com", tempID)
	secondaryIdentity, err := client.AddIdentity(context.TODO(), user.ID, users.AddIdentityParams{
		Type:  users.IdentityTypeEmail,
		Value: secondaryEmail,
	})
	assert.NoError(t, err)
	require.NotNil(t, secondaryIdentity)
	require.NotEmpty(t, secondaryIdentity.ID)

	// Get initial identities
	initialIdentities, err := client.GetIdentities(context.TODO(), user.ID)
	assert.NoError(t, err)
	require.NotNil(t, initialIdentities)
	require.Len(t, initialIdentities, 2, "expected 2 initial identities")

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

	// Add a phone identity
	phoneNumber := "+61412345678" // Australian number
	phoneIdentity, err := client.AddPhoneIdentity(context.TODO(), user.ID, phoneNumber)
	assert.NoError(t, err)
	require.NotNil(t, phoneIdentity)
	require.NotEmpty(t, phoneIdentity.ID)

	t.Logf("added phone identity: %+v\n", phoneIdentity)

	// Get and verify all identities
	allIdentities, err := client.GetIdentities(context.TODO(), user.ID)
	assert.NoError(t, err)
	require.NotNil(t, allIdentities)
	require.Len(t, allIdentities, 4, "expected 4 identities total")

	// Make the secondary email primary
	secondaryEmailIdentity := allIdentities[1] // The second email we added
	updatedIdentity, err := identitiesClient.Update(context.TODO(), secondaryEmailIdentity.ID, true)
	assert.NoError(t, err)
	require.NotNil(t, updatedIdentity)
	t.Logf("made secondary email primary: %+v\n", updatedIdentity)

	// Delete the username identity
	err = identitiesClient.Delete(context.TODO(), usernameIdentity.ID)
	assert.NoError(t, err)
	t.Logf("deleted username identity: %s\n", usernameIdentity.ID)

	// Verify remaining identities
	remainingIdentities, err := client.GetIdentities(context.TODO(), user.ID)
	assert.NoError(t, err)
	require.NotNil(t, remainingIdentities)
	require.Len(t, remainingIdentities, 3, "expected 3 remaining identities")

	// Clean up - delete the test user
	err = client.Delete(context.TODO(), user.ID)
	assert.NoError(t, err)
	t.Logf("deleted test user: %s\n", user.ID)
}

func TestE2EUserManagement(t *testing.T) {
	client := users.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())
	email := fmt.Sprintf("%s@example.com", tempID)

	// Create test user
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

	id := user.ID
	t.Logf("created test user: %s\n", id)

	// Test user suspension
	isSuspended := true
	updateParams := users.UpdateParams{
		IsSuspended: &isSuspended,
	}
	user, err = client.Update(context.TODO(), id, updateParams)
	assert.NoError(t, err)
	require.NotNil(t, user)
	assert.True(t, user.IsSuspended)
	t.Log("suspended user")

	// Test user reactivation
	isSuspended = false
	updateParams = users.UpdateParams{
		IsSuspended: &isSuspended,
	}
	user, err = client.Update(context.TODO(), id, updateParams)
	assert.NoError(t, err)
	require.NotNil(t, user)
	assert.False(t, user.IsSuspended)
	t.Log("reactivated user")

	// Test profile updates
	updateParams = users.UpdateParams{
		GivenName:  "Updated",
		FamilyName: "Name",
		ProvidedID: tempID + "-custom",
	}
	user, err = client.Update(context.TODO(), id, updateParams)
	assert.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Updated", user.FirstName)
	assert.Equal(t, "Name", user.LastName)
	t.Log("updated user profile")

	// Clean up
	err = client.Delete(context.TODO(), id)
	assert.NoError(t, err)
	t.Log("deleted test user")
}

func TestE2EPhoneIdentity(t *testing.T) {
	client := users.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())
	email := fmt.Sprintf("%s@example.com", tempID)

	// Create test user
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

	// Test adding phone identities from different countries
	testPhones := []struct {
		name        string
		phoneNumber string
	}{
		{
			name:        "Armenian number",
			phoneNumber: "+37455251234",
		},
		{
			name:        "Australian number",
			phoneNumber: "+61412345678",
		},
		{
			name:        "US number",
			phoneNumber: "+12025550123",
		},
	}

	for _, tt := range testPhones {
		t.Run(tt.name, func(t *testing.T) {
			phoneIdentity, err := client.AddPhoneIdentity(context.TODO(), user.ID, tt.phoneNumber)
			assert.NoError(t, err)
			require.NotNil(t, phoneIdentity)
			require.NotEmpty(t, phoneIdentity.ID)
			t.Logf("added phone identity: %+v\n", phoneIdentity)

			// Verify the identity was added
			identities, err := client.GetIdentities(context.TODO(), user.ID)
			assert.NoError(t, err)
			require.NotNil(t, identities)

			var found bool
			for _, identity := range identities {
				if identity.Type == string(users.IdentityTypePhone) && identity.Name == tt.phoneNumber {
					found = true
					break
				}
			}
			assert.True(t, found, "phone identity not found")
		})
	}

	// Clean up
	err = client.Delete(context.TODO(), user.ID)
	assert.NoError(t, err)
	t.Log("deleted test user")
}

func TestE2EProfileReset(t *testing.T) {
	client := users.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())
	email := fmt.Sprintf("%s@example.com", tempID)

	// Create test user with full profile
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

	id := user.ID
	t.Logf("created test user: %s\n", id)

	// Verify initial state
	user, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "User", user.LastName)
	t.Logf("initial user state: %+v\n", user)

	// Attempt to reset first and last name to null
	updateParams := users.UpdateParams{
		GivenName:  "",  // Try to reset to empty
		FamilyName: "",  // Try to reset to empty
	}
	user, err = client.Update(context.TODO(), id, updateParams)
	assert.NoError(t, err)
	require.NotNil(t, user)
	t.Logf("attempted to update user with empty names: %+v\n", user)

	// Get user to verify the state
	user, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, user)
	t.Logf("user after empty name update: %+v\n", user)

	// Try updating with omitted fields
	updateParamsOmit := users.UpdateParams{
		// First and last name fields omitted
		ProvidedID: tempID + "-updated",
	}
	user, err = client.Update(context.TODO(), id, updateParamsOmit)
	assert.NoError(t, err)
	require.NotNil(t, user)
	t.Logf("attempted to update user with omitted names: %+v\n", user)

	// Get user again to verify the final state
	user, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, user)
	t.Logf("user after omitted name update: %+v\n", user)

	// Clean up
	err = client.Delete(context.TODO(), id)
	assert.NoError(t, err)
}
