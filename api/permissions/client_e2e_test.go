//go:build e2e
// +build e2e

package permissions_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nxt-fwd/kinde-go/api/permissions"
	"github.com/nxt-fwd/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EList(t *testing.T) {
	client := permissions.New(testutil.DefaultE2EClient(t))
	res, err := client.List(context.TODO(), permissions.ListParams{})
	assert.NoError(t, err)
	// The permissions array might be nil if there are no permissions
	if res == nil {
		res = []permissions.Permission{}
	}
	t.Logf("found %d permissions", len(res))
}

func TestE2ECreateUpdateDelete(t *testing.T) {
	client := permissions.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	res, err := client.Create(context.TODO(), permissions.CreateParams{Name: tempID, Key: tempID})
	assert.NoError(t, err)
	require.NotNil(t, res)

	t.Logf("created test permission: %s\n", res.ID)

	err = client.Update(context.TODO(), res.ID, permissions.UpdateParams{Name: tempID + "2", Key: tempID + "2"})
	assert.NoError(t, err)

	t.Logf("updated test permission: %s\n", res.ID)

	err = client.Delete(context.TODO(), res.ID)
	assert.NoError(t, err)

	t.Logf("deleted test permission: %s\n", res.ID)
}

func TestE2EDescriptionReset(t *testing.T) {
	client := permissions.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Create permission with a description
	perm, err := client.Create(context.TODO(), permissions.CreateParams{
		Name:        tempID,
		Key:         tempID,
		Description: "Initial description",
	})
	assert.NoError(t, err)
	require.NotNil(t, perm)
	require.NotEmpty(t, perm.ID)
	t.Logf("created test permission with description: %s\n", perm.ID)

	// Get permission to verify initial state
	perm, err = client.Search(context.TODO(), permissions.SearchParams{
		Name: tempID,
		Key:  tempID,
	})
	assert.NoError(t, err)
	require.NotNil(t, perm)
	t.Logf("initial permission state: %+v\n", perm)

	// Attempt to update with empty description
	err = client.Update(context.TODO(), perm.ID, permissions.UpdateParams{
		Name:        tempID + "-updated",
		Key:         tempID + "-updated",
		Description: "", // Try to reset to empty
	})
	assert.NoError(t, err)
	t.Log("attempted to update permission with empty description")

	// Get permission to verify the state
	perm, err = client.Search(context.TODO(), permissions.SearchParams{
		Name: tempID + "-updated",
		Key:  tempID + "-updated",
	})
	assert.NoError(t, err)
	require.NotNil(t, perm)
	t.Logf("permission after update attempt: %+v\n", perm)

	// Try updating with null/omitted description
	err = client.Update(context.TODO(), perm.ID, permissions.UpdateParams{
		Name: tempID + "-updated-again",
		Key:  tempID + "-updated-again",
		// Description field omitted
	})
	assert.NoError(t, err)
	t.Log("attempted to update permission with omitted description")

	// Get permission again to verify the final state
	perm, err = client.Search(context.TODO(), permissions.SearchParams{
		Name: tempID + "-updated-again",
		Key:  tempID + "-updated-again",
	})
	assert.NoError(t, err)
	require.NotNil(t, perm)
	t.Logf("permission after second update attempt: %+v\n", perm)

	// Clean up
	err = client.Delete(context.TODO(), perm.ID)
	assert.NoError(t, err)
}
