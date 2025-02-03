//go:build e2e
// +build e2e

package organizations_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/axatol/kinde-go/api/organizations"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EList(t *testing.T) {
	client := organizations.New(testutil.DefaultE2EClient(t))
	res, err := client.List(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2ECreateGetUpdateDelete(t *testing.T) {
	client := organizations.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Create organization
	org, err := client.Create(context.TODO(), organizations.CreateParams{
		Name: tempID,
	})
	assert.NoError(t, err)
	require.NotNil(t, org)
	require.NotEmpty(t, org.Code)

	code := org.Code
	t.Logf("created test organization: %s\n", code)

	// Get organization
	org, err = client.Get(context.TODO(), code)
	assert.NoError(t, err)
	require.NotNil(t, org)
	t.Logf("got test organization: %+v\n", org)

	// Update organization
	updateParams := organizations.UpdateParams{
		Name: tempID + "-updated",
	}
	org, err = client.Update(context.TODO(), code, updateParams)
	assert.NoError(t, err)
	require.NotNil(t, org)
	t.Logf("updated test organization: %+v\n", org)

	// Delete organization
	err = client.Delete(context.TODO(), code)
	assert.NoError(t, err)
	t.Logf("deleted test organization: %+v\n", org)

	// Verify deletion
	_, err = client.Get(context.TODO(), code)
	assert.Error(t, err)
}

func TestE2EAddUsers(t *testing.T) {
	client := organizations.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// First create a test organization
	org, err := client.Create(context.TODO(), organizations.CreateParams{
		Name: tempID,
	})
	assert.NoError(t, err)
	require.NotNil(t, org)
	require.NotEmpty(t, org.Code)

	code := org.Code
	t.Logf("created test organization: %s\n", code)

	// Add users to the organization
	err = client.AddUsers(context.TODO(), code, organizations.AddUsersParams{
		Users: []organizations.AddUser{
			{
				ID:          "kp_a89a650856d144e8ad115ee58d377017", // Use an existing user ID
				Roles:       []string{"manager"},
				Permissions: []string{"admin"},
			},
		},
	})
	assert.NoError(t, err)
	t.Log("added users to organization")

	// Clean up - delete the test organization
	err = client.Delete(context.TODO(), code)
	assert.NoError(t, err)
} 