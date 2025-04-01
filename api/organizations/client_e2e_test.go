//go:build e2e
// +build e2e

package organizations_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nxt-fwd/kinde-go/api/organizations"
	"github.com/nxt-fwd/kinde-go/api/roles"
	"github.com/nxt-fwd/kinde-go/api/users"
	"github.com/nxt-fwd/kinde-go/internal/testutil"
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

func TestE2EUserRoleManagement(t *testing.T) {
	// Setup clients
	orgClient := organizations.New(testutil.DefaultE2EClient(t))
	userClient := users.New(testutil.DefaultE2EClient(t))
	roleClient := roles.New(testutil.DefaultE2EClient(t))

	// Create test organization
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())
	org, err := orgClient.Create(context.TODO(), organizations.CreateParams{
		Name: tempID,
	})
	assert.NoError(t, err)
	require.NotNil(t, org)
	require.NotEmpty(t, org.Code)

	orgCode := org.Code
	t.Logf("created test organization: %s\n", orgCode)

	// Create test user
	user, err := userClient.Create(context.TODO(), users.CreateParams{
		Profile: users.Profile{
			GivenName:  "Test",
			FamilyName: "User",
			Email:      fmt.Sprintf("%s@example.com", tempID),
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, user)
	require.NotEmpty(t, user.ID)

	userID := user.ID
	t.Logf("created test user: %s\n", userID)

	// Create test role
	role, err := roleClient.Create(context.TODO(), roles.CreateParams{
		Name: tempID,
		Key:  tempID,
	})
	assert.NoError(t, err)
	require.NotNil(t, role)
	require.NotEmpty(t, role.ID)

	roleID := role.ID
	t.Logf("created test role: %s\n", roleID)

	// Test error case: Try to add role before adding user to organization
	err = orgClient.AddUserRole(context.TODO(), orgCode, userID, roleID)
	assert.Error(t, err, "should error when adding role to user not in organization")
	t.Log("verified error when adding role to user not in organization")

	// Add user to organization
	err = orgClient.AddUsers(context.TODO(), orgCode, organizations.AddUsersParams{
		Users: []organizations.AddUser{
			{
				ID: userID,
			},
		},
	})
	assert.NoError(t, err)
	t.Log("added user to organization")

	// Test error case: Try to add non-existent role
	err = orgClient.AddUserRole(context.TODO(), orgCode, userID, "non-existent-role-id")
	assert.Error(t, err, "should error when adding non-existent role")
	t.Log("verified error when adding non-existent role")

	// Test error case: Try to get roles for non-existent user
	_, err = orgClient.GetUserRoles(context.TODO(), orgCode, "non-existent-user-id")
	assert.Error(t, err, "should error when getting roles for non-existent user")
	t.Log("verified error when getting roles for non-existent user")

	// Test successful role addition
	err = orgClient.AddUserRole(context.TODO(), orgCode, userID, roleID)
	assert.NoError(t, err)
	t.Log("added role to user")

	// Test adding the same role again (should error or be idempotent)
	err = orgClient.AddUserRole(context.TODO(), orgCode, userID, roleID)
	if err != nil {
		t.Log("adding same role twice returns error (expected)")
	} else {
		t.Log("adding same role twice is idempotent (acceptable)")
	}

	// Get and verify user roles
	roles, err := orgClient.GetUserRoles(context.TODO(), orgCode, userID)
	assert.NoError(t, err)
	require.NotNil(t, roles)
	assert.Len(t, roles, 1)
	assert.Equal(t, roleID, roles[0].ID)
	t.Log("verified user roles")

	// Test error case: Try to remove non-existent role
	err = orgClient.RemoveUserRole(context.TODO(), orgCode, userID, "non-existent-role-id")
	assert.Error(t, err, "should error when removing non-existent role")
	t.Log("verified error when removing non-existent role")

	// Remove role successfully
	err = orgClient.RemoveUserRole(context.TODO(), orgCode, userID, roleID)
	assert.NoError(t, err)
	t.Log("removed role from user")

	// Verify role was removed
	roles, err = orgClient.GetUserRoles(context.TODO(), orgCode, userID)
	assert.NoError(t, err)
	require.NotNil(t, roles)
	assert.Len(t, roles, 0)
	t.Log("verified role was removed")

	// Test removing the same role again (should error)
	err = orgClient.RemoveUserRole(context.TODO(), orgCode, userID, roleID)
	assert.Error(t, err, "should error when removing already removed role")
	t.Log("verified error when removing already removed role")

	// Cleanup
	err = userClient.Delete(context.TODO(), userID)
	assert.NoError(t, err)
	t.Log("deleted test user")

	err = roleClient.Delete(context.TODO(), roleID)
	assert.NoError(t, err)
	t.Log("deleted test role")

	err = orgClient.Delete(context.TODO(), orgCode)
	assert.NoError(t, err)
	t.Log("deleted test organization")
}

func TestE2EOrganizationSettings(t *testing.T) {
	client := organizations.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Create test organization
	org, err := client.Create(context.TODO(), organizations.CreateParams{
		Name: tempID,
	})
	assert.NoError(t, err)
	require.NotNil(t, org)
	require.NotEmpty(t, org.Code)

	code := org.Code
	t.Logf("created test organization: %s\n", code)

	// Test basic update
	updateParams := organizations.UpdateParams{
		Name: tempID + "-updated",
	}
	org, err = client.Update(context.TODO(), code, updateParams)
	assert.NoError(t, err)
	require.NotNil(t, org)
	t.Log("updated organization name")

	// Get organization to verify update
	org, err = client.Get(context.TODO(), code)
	assert.NoError(t, err)
	require.NotNil(t, org)
	assert.Equal(t, tempID+"-updated", org.Name)
	t.Log("verified organization update")

	// Clean up
	err = client.Delete(context.TODO(), code)
	assert.NoError(t, err)
	t.Log("deleted test organization")
}
