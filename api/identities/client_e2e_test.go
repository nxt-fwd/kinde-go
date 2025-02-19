package identities

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nxt-fwd/kinde-go/api/users"
	"github.com/nxt-fwd/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentitiesE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	ctx := context.Background()
	client := testutil.DefaultE2EClient(t)
	identitiesClient := New(client)
	usersClient := users.New(client)

	email := fmt.Sprintf("test-%d@example.com", time.Now().UnixMilli())

	// First create a user
	user, err := usersClient.Create(ctx, users.CreateParams{
		Profile: users.Profile{
			GivenName:  "John",
			FamilyName: "Doe",
			Email:      email,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, user)

	t.Cleanup(func() {
		err := usersClient.Delete(ctx, user.ID)
		assert.NoError(t, err)
	})

	// Add a primary email identity
	primaryIdentity, err := usersClient.AddIdentity(ctx, user.ID, users.AddIdentityParams{
		Type:  users.IdentityTypeEmail,
		Value: email,
	})
	require.NoError(t, err)
	require.NotNil(t, primaryIdentity)

	// Add a secondary email identity
	secondaryEmail := fmt.Sprintf("test-%d-secondary@example.com", time.Now().UnixMilli())
	secondaryIdentity, err := usersClient.AddIdentity(ctx, user.ID, users.AddIdentityParams{
		Type:  users.IdentityTypeEmail,
		Value: secondaryEmail,
	})
	require.NoError(t, err)
	require.NotNil(t, secondaryIdentity)

	// Get user's identities
	identities, err := usersClient.GetIdentities(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, identities, 2)

	// Get first identity details
	identity1, err := identitiesClient.Get(ctx, identities[0].ID)
	require.NoError(t, err)
	require.NotNil(t, identity1)
	assert.Equal(t, identities[0].ID, identity1.ID)

	// Make second identity primary
	identity2, err := identitiesClient.Update(ctx, identities[1].ID, true)
	require.NoError(t, err)
	require.NotNil(t, identity2)

	// Delete first identity
	err = identitiesClient.Delete(ctx, identities[0].ID)
	require.NoError(t, err)

	// Verify only one identity remains
	remainingIdentities, err := usersClient.GetIdentities(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, remainingIdentities, 1)
	assert.Equal(t, identities[1].ID, remainingIdentities[0].ID)
} 