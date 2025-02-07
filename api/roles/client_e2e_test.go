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
		Description: "Test role",
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
		Description: "Updated test role",
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
