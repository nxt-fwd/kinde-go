//go:build e2e
// +build e2e

package roles_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nxt-fwd/kinde-go/api/roles"
	"github.com/nxt-fwd/kinde-go/api/permissions"
	"github.com/nxt-fwd/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EList(t *testing.T) {
	client := roles.New(testutil.DefaultE2EClient(t))
	rolesList, err := client.List(context.TODO())
	assert.NoError(t, err)
	// The roles array might be nil if there are no roles
	if rolesList == nil {
		rolesList = []roles.Role{}
	}
	t.Logf("found %d roles", len(rolesList))
}

func TestE2ECreateGetUpdateDelete(t *testing.T) {
	client := roles.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Create test permissions first
	permClient := permissions.New(testutil.DefaultE2EClient(t))
	
	// Create read permission
	readPerm, err := permClient.Create(context.TODO(), permissions.CreateParams{
		Name: "Test Read Permission",
		Key:  "read:users",
		Description: "Test permission for reading users",
	})
	assert.NoError(t, err)
	require.NotNil(t, readPerm)
	t.Logf("created read permission with ID: %s\n", readPerm.ID)

	// Create write permission
	writePerm, err := permClient.Create(context.TODO(), permissions.CreateParams{
		Name: "Test Write Permission",
		Key:  "write:users",
		Description: "Test permission for writing users",
	})
	assert.NoError(t, err)
	require.NotNil(t, writePerm)
	t.Logf("created write permission with ID: %s\n", writePerm.ID)

	// Create role with read permission
	role, err := client.Create(context.TODO(), roles.CreateParams{
		Name:        tempID,
		Key:         tempID,
		Description: "Initial role description",
		Permissions: []string{readPerm.ID},
	})
	assert.NoError(t, err)
	require.NotNil(t, role)
	require.NotEmpty(t, role.ID)

	id := role.ID
	t.Logf("created test role: %s\n", id)

	// Get role
	role, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, role)
	t.Logf("got test role: %+v\n", role)

	// Update role with both permissions
	updateParams := roles.UpdateParams{
		Name:        tempID + "-updated",
		Description: "Updated role description",
		Permissions: []string{readPerm.ID, writePerm.ID},
	}
	role, err = client.Update(context.TODO(), id, updateParams)
	assert.NoError(t, err)
	require.NotNil(t, role)
	t.Logf("updated test role: %+v\n", role)

	// Delete role
	err = client.Delete(context.TODO(), id)
	assert.NoError(t, err)
	t.Logf("deleted test role: %s\n", id)

	// Verify role deletion
	_, err = client.Get(context.TODO(), id)
	assert.Error(t, err)

	// Clean up - delete the test permissions
	err = permClient.Delete(context.TODO(), readPerm.ID)
	assert.NoError(t, err)
	t.Logf("deleted read permission: %s\n", readPerm.ID)

	err = permClient.Delete(context.TODO(), writePerm.ID)
	assert.NoError(t, err)
	t.Logf("deleted write permission: %s\n", writePerm.ID)
}

func TestE2EPermissions(t *testing.T) {
	client := roles.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Create test permissions first
	permClient := permissions.New(testutil.DefaultE2EClient(t))
	
	// Create read permission
	readPermKey := "users:read"
	readPerm, err := permClient.Create(context.TODO(), permissions.CreateParams{
		Name: "Test Read Permission",
		Key:  readPermKey,
		Description: "Test permission for reading users",
	})
	assert.NoError(t, err)
	require.NotNil(t, readPerm)
	t.Logf("created read permission with ID: %s\n", readPerm.ID)

	// Create write permission
	writePermKey := "users:write"
	writePerm, err := permClient.Create(context.TODO(), permissions.CreateParams{
		Name: "Test Write Permission",
		Key:  writePermKey,
		Description: "Test permission for writing users",
	})
	assert.NoError(t, err)
	require.NotNil(t, writePerm)
	t.Logf("created write permission with ID: %s\n", writePerm.ID)

	// Create a test role
	role, err := client.Create(context.TODO(), roles.CreateParams{
		Name:        tempID,
		Key:         tempID,
		Description: "Test role",
	})
	assert.NoError(t, err)
	require.NotNil(t, role)
	require.NotEmpty(t, role.ID)

	id := role.ID
	t.Logf("created test role: %s\n", id)

	// Add both permissions to the role
	updateResponse, err := client.UpdatePermissions(context.TODO(), id, roles.UpdatePermissionsParams{
		Permissions: []roles.UpdatePermissionItem{
			{ID: readPerm.ID},
			{ID: writePerm.ID},
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, updateResponse)
	t.Logf("added permissions: %v (response: added=%v, removed=%v)\n", 
		[]string{readPerm.ID, writePerm.ID}, 
		updateResponse.PermissionsAdded, 
		updateResponse.PermissionsRemoved,
	)

	// Get the role to verify permissions
	role, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, role)
	t.Logf("role permissions after adding: %v\n", role.Permissions)

	// Remove one permission using the DELETE endpoint
	err = client.RemovePermission(context.TODO(), id, writePerm.ID)
	assert.NoError(t, err)
	t.Logf("removed permission with ID: %s\n", writePerm.ID)

	// Get the role again to verify the permission was removed
	role, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, role)
	t.Logf("role permissions after removal: %v\n", role.Permissions)

	// Clean up - delete the test permissions
	err = permClient.Delete(context.TODO(), readPerm.ID)
	assert.NoError(t, err)
	err = permClient.Delete(context.TODO(), writePerm.ID)
	assert.NoError(t, err)

	// Clean up - delete the test role
	err = client.Delete(context.TODO(), id)
	assert.NoError(t, err)
}

func TestE2EBulkPermissionChanges(t *testing.T) {
	client := roles.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Create test permissions first
	permClient := permissions.New(testutil.DefaultE2EClient(t))
	
	// Create first permission
	perm1Key := "users:read"
	perm1, err := permClient.Create(context.TODO(), permissions.CreateParams{
		Name: "Test Read Permission",
		Key:  perm1Key,
		Description: "Test permission for reading users",
	})
	assert.NoError(t, err)
	require.NotNil(t, perm1)
	t.Logf("created first permission with ID: %s\n", perm1.ID)

	// Create second permission
	perm2Key := "users:write"
	perm2, err := permClient.Create(context.TODO(), permissions.CreateParams{
		Name: "Test Write Permission",
		Key:  perm2Key,
		Description: "Test permission for writing users",
	})
	assert.NoError(t, err)
	require.NotNil(t, perm2)
	t.Logf("created second permission with ID: %s\n", perm2.ID)

	// Create third permission
	perm3Key := "users:delete"
	perm3, err := permClient.Create(context.TODO(), permissions.CreateParams{
		Name: "Test Delete Permission",
		Key:  perm3Key,
		Description: "Test permission for deleting users",
	})
	assert.NoError(t, err)
	require.NotNil(t, perm3)
	t.Logf("created third permission with ID: %s\n", perm3.ID)

	// Create a test role with initial permission (perm1)
	role, err := client.Create(context.TODO(), roles.CreateParams{
		Name:        tempID,
		Key:         tempID,
		Description: "Test role",
	})
	assert.NoError(t, err)
	require.NotNil(t, role)
	require.NotEmpty(t, role.ID)

	id := role.ID
	t.Logf("created test role: %s\n", id)

	// Add first permission
	updateResponse, err := client.UpdatePermissions(context.TODO(), id, roles.UpdatePermissionsParams{
		Permissions: []roles.UpdatePermissionItem{
			{ID: perm1.ID},
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, updateResponse)
	t.Logf("added first permission: %v (response: added=%v, removed=%v)\n", 
		perm1.ID, 
		updateResponse.PermissionsAdded, 
		updateResponse.PermissionsRemoved,
	)

	// In a single request: remove perm1 and add perm2 and perm3
	updateResponse, err = client.UpdatePermissions(context.TODO(), id, roles.UpdatePermissionsParams{
		Permissions: []roles.UpdatePermissionItem{
			{ID: perm1.ID, Operation: "delete"},
			{ID: perm2.ID},
			{ID: perm3.ID},
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, updateResponse)
	t.Logf("bulk update response: added=%v, removed=%v\n", 
		updateResponse.PermissionsAdded, 
		updateResponse.PermissionsRemoved,
	)

	// Get the role to verify permissions
	role, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, role)
	t.Logf("role permissions after bulk update: %v\n", role.Permissions)

	// Clean up - delete the test permissions
	err = permClient.Delete(context.TODO(), perm1.ID)
	assert.NoError(t, err)
	err = permClient.Delete(context.TODO(), perm2.ID)
	assert.NoError(t, err)
	err = permClient.Delete(context.TODO(), perm3.ID)
	assert.NoError(t, err)

	// Clean up - delete the test role
	err = client.Delete(context.TODO(), id)
	assert.NoError(t, err)
}

func TestE2EBulkPermissionOperations(t *testing.T) {
	client := roles.New(testutil.DefaultE2EClient(t))
	permClient := permissions.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Create test permissions
	testPerms := []struct {
		name string
		key  string
	}{
		{name: "Read Users", key: "users:read"},
		{name: "Write Users", key: "users:write"},
		{name: "Delete Users", key: "users:delete"},
		{name: "Read Orgs", key: "orgs:read"},
		{name: "Write Orgs", key: "orgs:write"},
	}

	var createdPerms []string
	for _, p := range testPerms {
		perm, err := permClient.Create(context.TODO(), permissions.CreateParams{
			Name: fmt.Sprintf("%s-%s", tempID, p.name),
			Key:  fmt.Sprintf("%s-%s", tempID, p.key),
		})
		assert.NoError(t, err)
		require.NotNil(t, perm)
		createdPerms = append(createdPerms, perm.ID)
		t.Logf("created permission: %s (%s)", perm.Name, perm.ID)
	}

	// Create a test role
	role, err := client.Create(context.TODO(), roles.CreateParams{
		Name:        tempID,
		Key:         tempID,
		Description: "Test role for bulk permission operations",
		Permissions: []string{createdPerms[0]},
	})
	assert.NoError(t, err)
	require.NotNil(t, role)
	t.Logf("created role: %s", role.ID)

	// Test bulk permission addition
	addResponse, err := client.UpdatePermissions(context.TODO(), role.ID, roles.UpdatePermissionsParams{
		Permissions: []roles.UpdatePermissionItem{
			{ID: createdPerms[0]},
			{ID: createdPerms[1]},
			{ID: createdPerms[2]},
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, addResponse)
	assert.Len(t, addResponse.PermissionsAdded, 3, "should have added 3 permissions")
	t.Log("added initial permissions")

	// Verify permissions were added
	role, err = client.Get(context.TODO(), role.ID)
	assert.NoError(t, err)
	require.NotNil(t, role)
	assert.Len(t, role.Permissions, 3, "role should have 3 permissions")
	t.Log("verified initial permissions")

	// Test mixed operation - add some and remove some in one call
	mixedResponse, err := client.UpdatePermissions(context.TODO(), role.ID, roles.UpdatePermissionsParams{
		Permissions: []roles.UpdatePermissionItem{
			{ID: createdPerms[0], Operation: "delete"}, // Remove first permission
			{ID: createdPerms[3]}, // Add new permission
			{ID: createdPerms[4]}, // Add new permission
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, mixedResponse)
	assert.Len(t, mixedResponse.PermissionsAdded, 2, "should have added 2 permissions")
	assert.Len(t, mixedResponse.PermissionsRemoved, 1, "should have removed 1 permission")
	t.Log("performed mixed permission update")

	// Verify final permission state
	role, err = client.Get(context.TODO(), role.ID)
	assert.NoError(t, err)
	require.NotNil(t, role)
	assert.Len(t, role.Permissions, 4, "role should have 4 permissions")
	t.Log("verified final permissions")

	// Clean up - delete the role
	err = client.Delete(context.TODO(), role.ID)
	assert.NoError(t, err)
	t.Log("deleted role")

	// Clean up - delete all test permissions
	for _, permID := range createdPerms {
		err = permClient.Delete(context.TODO(), permID)
		assert.NoError(t, err)
	}
	t.Log("deleted test permissions")
}

func TestE2EPermissionRemoval(t *testing.T) {
	client := roles.New(testutil.DefaultE2EClient(t))
	permClient := permissions.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Create test permissions
	perm1, err := permClient.Create(context.TODO(), permissions.CreateParams{
		Name: fmt.Sprintf("%s-perm1", tempID),
		Key:  fmt.Sprintf("%s-perm1", tempID),
	})
	assert.NoError(t, err)
	require.NotNil(t, perm1)
	t.Logf("created first permission: %s", perm1.ID)

	perm2, err := permClient.Create(context.TODO(), permissions.CreateParams{
		Name: fmt.Sprintf("%s-perm2", tempID),
		Key:  fmt.Sprintf("%s-perm2", tempID),
	})
	assert.NoError(t, err)
	require.NotNil(t, perm2)
	t.Logf("created second permission: %s", perm2.ID)

	// Create a role
	role, err := client.Create(context.TODO(), roles.CreateParams{
		Name:        fmt.Sprintf("test-%d", time.Now().UnixMilli()),
		Key:         fmt.Sprintf("test-%d", time.Now().UnixMilli()),
		Description: "Test role for permission removal",
		Permissions: []string{perm1.ID, perm2.ID},
	})
	assert.NoError(t, err)
	require.NotNil(t, role)
	t.Logf("created role: %s", role.ID)

	// Test Case 1: Add and remove permissions one by one
	t.Log("Test Case 1: Individual permission operations")

	// Add both permissions
	updateResponse, err := client.UpdatePermissions(context.TODO(), role.ID, roles.UpdatePermissionsParams{
		Permissions: []roles.UpdatePermissionItem{
			{ID: perm1.ID},
			{ID: perm2.ID},
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, updateResponse)
	assert.Len(t, updateResponse.PermissionsAdded, 2, "should have added 2 permissions")
	t.Log("added initial permissions")

	// Verify initial permissions
	role, err = client.Get(context.TODO(), role.ID)
	assert.NoError(t, err)
	require.NotNil(t, role)
	assert.Len(t, role.Permissions, 2, "role should have 2 permissions initially")
	t.Log("verified initial permissions")

	// Remove first permission using the delete endpoint
	err = client.RemovePermission(context.TODO(), role.ID, perm1.ID)
	assert.NoError(t, err)
	t.Log("removed first permission")

	// Verify only one permission remains
	role, err = client.Get(context.TODO(), role.ID)
	assert.NoError(t, err)
	require.NotNil(t, role)
	assert.Len(t, role.Permissions, 1, "role should have 1 permission after removal")
	assert.Equal(t, perm2.ID, role.Permissions[0], "remaining permission should be perm2")
	t.Log("verified one permission remains")

	// Remove the second permission using the delete endpoint
	err = client.RemovePermission(context.TODO(), role.ID, perm2.ID)
	assert.NoError(t, err)
	t.Log("removed second permission")

	// Verify no permissions remain
	role, err = client.Get(context.TODO(), role.ID)
	assert.NoError(t, err)
	require.NotNil(t, role)
	assert.Len(t, role.Permissions, 0, "role should have no permissions after removing all")
	t.Log("verified no permissions remain")

	// Test Case 2: Bulk permission removal
	t.Log("Test Case 2: Bulk permission removal")

	// Add permissions again for bulk removal test
	updateResponse, err = client.UpdatePermissions(context.TODO(), role.ID, roles.UpdatePermissionsParams{
		Permissions: []roles.UpdatePermissionItem{
			{ID: perm1.ID},
			{ID: perm2.ID},
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, updateResponse)

	// Verify permissions were added
	role, err = client.Get(context.TODO(), role.ID)
	assert.NoError(t, err)
	require.NotNil(t, role)
	assert.Len(t, role.Permissions, 2, "role should have 2 permissions")

	// Remove all permissions using bulk update
	updateResponse, err = client.UpdatePermissions(context.TODO(), role.ID, roles.UpdatePermissionsParams{
		Permissions: []roles.UpdatePermissionItem{
			{ID: perm1.ID, Operation: "delete"},
			{ID: perm2.ID, Operation: "delete"},
		},
	})
	assert.NoError(t, err)
	require.NotNil(t, updateResponse)
	assert.Len(t, updateResponse.PermissionsRemoved, 2, "should have removed 2 permissions")

	// Verify all permissions were removed
	role, err = client.Get(context.TODO(), role.ID)
	assert.NoError(t, err)
	require.NotNil(t, role)
	assert.Len(t, role.Permissions, 0, "role should have no permissions after bulk removal")

	// Clean up
	err = client.Delete(context.TODO(), role.ID)
	assert.NoError(t, err)
	t.Log("deleted role")

	err = permClient.Delete(context.TODO(), perm1.ID)
	assert.NoError(t, err)
	err = permClient.Delete(context.TODO(), perm2.ID)
	assert.NoError(t, err)
	t.Log("deleted test permissions")
}

func TestE2EDescriptionBehavior(t *testing.T) {
	ctx := context.TODO()
	client := roles.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Test Case 1: Create role with initial description
	role, err := client.Create(ctx, roles.CreateParams{
		Name:        tempID,
		Key:         tempID,
		Description: "Initial description",
	})
	assert.NoError(t, err)
	assert.NotNil(t, role)
	t.Log("Created role with initial description")

	// Get the role to verify description
	retrievedRole, err := client.Get(ctx, role.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Initial description", retrievedRole.Description)

	// Test Case 2: Update role with new description
	updatedRole, err := client.Update(ctx, role.ID, roles.UpdateParams{
		Name:        tempID + "-updated",
		Key:         tempID,
		Description: "Updated description",
	})
	assert.NoError(t, err)
	assert.NotNil(t, updatedRole)
	t.Log("Updated role with new description")

	// Get the role again to verify updated description
	retrievedRole, err = client.Get(ctx, role.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated description", retrievedRole.Description)

	// Cleanup
	err = client.Delete(ctx, role.ID)
	assert.NoError(t, err)
	t.Log("Cleaned up test role")
}

func TestE2EDescriptionValidation(t *testing.T) {
	client := roles.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Test Case 1: Create without description
	t.Log("Test Case 1: Create without description")
	_, err := client.Create(context.TODO(), roles.CreateParams{
		Name: tempID,
		Key:  tempID,
		// Description intentionally omitted
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "description is required")
	t.Log("Verified error when creating role without description")

	// Test Case 2: Create with empty description
	t.Log("Test Case 2: Create with empty description")
	_, err = client.Create(context.TODO(), roles.CreateParams{
		Name:        tempID,
		Key:         tempID,
		Description: "", // Empty description
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "description is required")
	t.Log("Verified error when creating role with empty description")

	// Create a valid role for update tests
	role, err := client.Create(context.TODO(), roles.CreateParams{
		Name:        tempID,
		Key:         tempID,
		Description: "Initial description",
	})
	assert.NoError(t, err)
	require.NotNil(t, role)

	// Test Case 3: Update without description
	t.Log("Test Case 3: Update without description")
	_, err = client.Update(context.TODO(), role.ID, roles.UpdateParams{
		Name: tempID + "-updated",
		Key:  tempID,
		// Description intentionally omitted
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "description is required")
	t.Log("Verified error when updating role without description")

	// Test Case 4: Update with empty description
	t.Log("Test Case 4: Update with empty description")
	_, err = client.Update(context.TODO(), role.ID, roles.UpdateParams{
		Name:        tempID + "-updated",
		Key:         tempID,
		Description: "", // Empty description
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "description is required")
	t.Log("Verified error when updating role with empty description")

	// Clean up
	err = client.Delete(context.TODO(), role.ID)
	assert.NoError(t, err)
	t.Log("Cleaned up test role")
} 
