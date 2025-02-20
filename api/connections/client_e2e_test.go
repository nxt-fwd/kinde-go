//go:build e2e
// +build e2e

package connections

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nxt-fwd/kinde-go/api/applications"
	"github.com/nxt-fwd/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EListConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	client := New(testutil.DefaultE2EClient(t))
	connections, err := client.List(context.TODO())
	assert.NoError(t, err)
	require.NotNil(t, connections)

	// Log the connections to see their structure
	for i, conn := range connections {
		t.Logf("Connection %d:", i+1)
		t.Logf("  ID: %s", conn.ID)
		t.Logf("  Name: %s", conn.Name)
		t.Logf("  DisplayName: %s", conn.DisplayName)
		t.Logf("  Strategy: %s", conn.Strategy)
		t.Logf("  EnabledApplications: %v", conn.EnabledApplications)
		t.Logf("  Options: %+v", conn.Options)
		t.Logf("---")
	}
}

func TestE2EGetConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	client := New(testutil.DefaultE2EClient(t))
	
	// First get all connections
	connections, err := client.List(context.TODO())
	assert.NoError(t, err)
	require.NotNil(t, connections)

	// Then get each connection individually to examine its structure
	for _, conn := range connections {
		t.Logf("Getting connection: %s (%s)", conn.Name, conn.ID)
		
		connection, err := client.Get(context.TODO(), conn.ID)
		assert.NoError(t, err)
		require.NotNil(t, connection)

		// Convert options to JSON for detailed inspection
		var optionsJSON []byte
		if connection.Options != nil {
			optionsJSON, err = json.MarshalIndent(connection.Options, "", "  ")
			assert.NoError(t, err)
		}

		t.Logf("Connection Details:")
		t.Logf("  ID: %s", connection.ID)
		t.Logf("  Name: %s", connection.Name)
		t.Logf("  DisplayName: %s", connection.DisplayName)
		t.Logf("  Strategy: %s", connection.Strategy)
		t.Logf("  EnabledApplications: %v", connection.EnabledApplications)
		if optionsJSON != nil {
			t.Logf("  Options: %s", string(optionsJSON))
		} else {
			t.Logf("  Options: nil")
		}
		t.Logf("---")
	}
}

func TestE2ECreateConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	client := New(testutil.DefaultE2EClient(t))
	
	// Create a unique name for the test connection
	timestamp := time.Now().UnixMilli()
	name := fmt.Sprintf("test-google-%d", timestamp)
	displayName := fmt.Sprintf("Test Google %d", timestamp)

	// Create a Google OAuth2 connection
	options := &SocialConnectionOptions{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	params := CreateParams{
		Name:        name,
		DisplayName: displayName,
		Strategy:    StrategyOAuth2Google,
		Options:     options,
	}

	// Create the connection - note that the create response only includes the ID
	connection, err := client.Create(context.TODO(), params)
	assert.NoError(t, err)
	require.NotNil(t, connection)
	require.NotEmpty(t, connection.ID)

	// Clean up the test connection after we're done
	defer func() {
		t.Logf("Cleaning up test connection: %s", connection.ID)
		err := client.Delete(context.TODO(), connection.ID)
		assert.NoError(t, err)
	}()

	t.Logf("Created connection with ID: %s", connection.ID)

	// Get the connection to verify all fields
	connection, err = client.Get(context.TODO(), connection.ID)
	assert.NoError(t, err)
	require.NotNil(t, connection)

	// Log the connection details
	t.Logf("Connection Details:")
	t.Logf("  ID: %s", connection.ID)
	t.Logf("  Name: %s", connection.Name)
	t.Logf("  DisplayName: %s", connection.DisplayName)
	t.Logf("  Strategy: %s", connection.Strategy)
	t.Logf("  EnabledApplications: %v", connection.EnabledApplications)

	// Convert options to JSON for detailed inspection
	var optionsJSON []byte
	if connection.Options != nil {
		optionsJSON, err = json.MarshalIndent(connection.Options, "", "  ")
		assert.NoError(t, err)
		t.Logf("  Options: %s", string(optionsJSON))
	} else {
		t.Logf("  Options: nil")
	}

	// Verify the connection was created with the correct values
	assert.Equal(t, name, connection.Name)
	assert.Equal(t, displayName, connection.DisplayName)
	assert.Equal(t, string(StrategyOAuth2Google), connection.Strategy)
}

func TestE2ECreateConnectionWithEmailTrust(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	client := New(testutil.DefaultE2EClient(t))
	
	// Create a unique name for the test connection
	timestamp := time.Now().UnixMilli()
	name := fmt.Sprintf("test-google-trust-%d", timestamp)
	displayName := fmt.Sprintf("Test Google Trust %d", timestamp)

	// Create a Google OAuth2 connection with email trust setting
	options := map[string]interface{}{
		"client_id":              "test-client-id",
		"client_secret":          "test-client-secret",
		"trust_email_verified":   true,
	}

	params := CreateParams{
		Name:        name,
		DisplayName: displayName,
		Strategy:    StrategyOAuth2Google,
		Options:     options,
	}

	// Create the connection
	connection, err := client.Create(context.TODO(), params)
	assert.NoError(t, err)
	require.NotNil(t, connection)
	require.NotEmpty(t, connection.ID)

	// Clean up the test connection after we're done
	defer func() {
		t.Logf("Cleaning up test connection: %s", connection.ID)
		err := client.Delete(context.TODO(), connection.ID)
		assert.NoError(t, err)
	}()

	t.Logf("Created connection with ID: %s", connection.ID)

	// Get the connection to verify all fields
	connection, err = client.Get(context.TODO(), connection.ID)
	assert.NoError(t, err)
	require.NotNil(t, connection)

	// Log the connection details
	t.Logf("Connection Details:")
	t.Logf("  ID: %s", connection.ID)
	t.Logf("  Name: %s", connection.Name)
	t.Logf("  DisplayName: %s", connection.DisplayName)
	t.Logf("  Strategy: %s", connection.Strategy)
	t.Logf("  EnabledApplications: %v", connection.EnabledApplications)

	// Convert options to JSON for detailed inspection
	var optionsJSON []byte
	if connection.Options != nil {
		optionsJSON, err = json.MarshalIndent(connection.Options, "", "  ")
		assert.NoError(t, err)
		t.Logf("  Options: %s", string(optionsJSON))
	} else {
		t.Logf("  Options: nil")
	}

	// Verify the connection was created with the correct values
	assert.Equal(t, name, connection.Name)
	assert.Equal(t, displayName, connection.DisplayName)
	assert.Equal(t, string(StrategyOAuth2Google), connection.Strategy)
}

// TestE2EDeleteConnection tests the delete connection functionality
func TestE2EDeleteConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	client := New(testutil.DefaultE2EClient(t))
	
	// Create a connection to delete
	timestamp := time.Now().UnixMilli()
	name := fmt.Sprintf("test-delete-%d", timestamp)
	displayName := fmt.Sprintf("Test Delete %d", timestamp)

	params := CreateParams{
		Name:        name,
		DisplayName: displayName,
		Strategy:    StrategyOAuth2Google,
		Options: &SocialConnectionOptions{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
	}

	// Create the connection
	connection, err := client.Create(context.TODO(), params)
	assert.NoError(t, err)
	require.NotNil(t, connection)
	require.NotEmpty(t, connection.ID)

	t.Logf("Created connection with ID: %s", connection.ID)

	// Delete the connection
	err = client.Delete(context.TODO(), connection.ID)
	assert.NoError(t, err)

	// Verify the connection is deleted by trying to get it
	_, err = client.Get(context.TODO(), connection.ID)
	assert.Error(t, err) // Should return an error since the connection no longer exists
}

// TestE2EConnectionWithApplications tests creating and updating connections with enabled applications
func TestE2EConnectionWithApplications(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	ctx := context.TODO()
	client := New(testutil.DefaultE2EClient(t))
	appsClient := applications.New(testutil.DefaultE2EClient(t))
	
	// First create some test applications
	timestamp := time.Now().UnixMilli()
	
	// Create first test application
	app1Name := fmt.Sprintf("test-app1-%d", timestamp)
	app1, err := appsClient.Create(ctx, applications.CreateParams{
		Name: app1Name,
		Type: applications.TypeRegular,
	})
	require.NoError(t, err)
	require.NotNil(t, app1)
	
	// Clean up first application after test
	defer func() {
		t.Logf("Cleaning up test application: %s", app1.ID)
		err := appsClient.Delete(ctx, app1.ID)
		assert.NoError(t, err)
	}()

	// Create second test application
	app2Name := fmt.Sprintf("test-app2-%d", timestamp)
	app2, err := appsClient.Create(ctx, applications.CreateParams{
		Name: app2Name,
		Type: applications.TypeRegular,
	})
	require.NoError(t, err)
	require.NotNil(t, app2)
	
	// Clean up second application after test
	defer func() {
		t.Logf("Cleaning up test application: %s", app2.ID)
		err := appsClient.Delete(ctx, app2.ID)
		assert.NoError(t, err)
	}()

	t.Logf("Created test applications:")
	t.Logf("  App1 ID: %s, Name: %s", app1.ID, app1.Name)
	t.Logf("  App2 ID: %s, Name: %s", app2.ID, app2.Name)

	// Now create a connection with these applications enabled
	name := fmt.Sprintf("test-apps-%d", timestamp)
	displayName := fmt.Sprintf("Test Apps %d", timestamp)

	// Use the newly created application IDs
	enabledApps := []string{
		app1.ID,
		app2.ID,
	}
	
	params := CreateParams{
		Name:                name,
		DisplayName:         displayName,
		Strategy:            StrategyOAuth2Google,
		EnabledApplications: enabledApps,
		Options: &SocialConnectionOptions{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
	}

	// Create the connection
	connection, err := client.Create(ctx, params)
	assert.NoError(t, err)
	require.NotNil(t, connection)
	require.NotEmpty(t, connection.ID)

	// Clean up the test connection after we're done
	defer func() {
		t.Logf("Cleaning up test connection: %s", connection.ID)
		err := client.Delete(ctx, connection.ID)
		assert.NoError(t, err)
	}()

	t.Logf("Created connection with ID: %s", connection.ID)

	// Get the connection to verify enabled applications
	connection, err = client.Get(ctx, connection.ID)
	assert.NoError(t, err)
	require.NotNil(t, connection)

	t.Logf("Connection Details after creation:")
	t.Logf("  ID: %s", connection.ID)
	t.Logf("  Name: %s", connection.Name)
	t.Logf("  DisplayName: %s", connection.DisplayName)
	t.Logf("  Strategy: %s", connection.Strategy)
	t.Logf("  EnabledApplications: %v", connection.EnabledApplications)

	// Verify connections through the applications endpoint
	// First application should have the connection
	app1Connections, err := appsClient.GetConnections(ctx, app1.ID)
	assert.NoError(t, err)
	require.NotNil(t, app1Connections)
	t.Logf("App1 Connections: %+v", app1Connections)
	
	var foundInApp1 bool
	for _, conn := range app1Connections {
		if conn.ID == connection.ID {
			foundInApp1 = true
			assert.Equal(t, name, conn.Name)
			assert.Equal(t, displayName, conn.DisplayName)
			assert.Equal(t, string(StrategyOAuth2Google), conn.Strategy)
			break
		}
	}
	assert.True(t, foundInApp1, "Connection should be found in app1's connections")

	// Second application should also have the connection
	app2Connections, err := appsClient.GetConnections(ctx, app2.ID)
	assert.NoError(t, err)
	require.NotNil(t, app2Connections)
	t.Logf("App2 Connections: %+v", app2Connections)
	
	var foundInApp2 bool
	for _, conn := range app2Connections {
		if conn.ID == connection.ID {
			foundInApp2 = true
			assert.Equal(t, name, conn.Name)
			assert.Equal(t, displayName, conn.DisplayName)
			assert.Equal(t, string(StrategyOAuth2Google), conn.Strategy)
			break
		}
	}
	assert.True(t, foundInApp2, "Connection should be found in app2's connections")

	// Create a third application for update test
	app3Name := fmt.Sprintf("test-app3-%d", timestamp)
	app3, err := appsClient.Create(ctx, applications.CreateParams{
		Name: app3Name,
		Type: applications.TypeRegular,
	})
	require.NoError(t, err)
	require.NotNil(t, app3)
	
	// Clean up third application after test
	defer func() {
		t.Logf("Cleaning up test application: %s", app3.ID)
		err := appsClient.Delete(ctx, app3.ID)
		assert.NoError(t, err)
	}()

	t.Logf("Created third test application:")
	t.Logf("  App3 ID: %s, Name: %s", app3.ID, app3.Name)

	// Update the enabled applications to only use the third app
	updatedApps := []string{app3.ID}
	updateParams := UpdateParams{
		EnabledApplications: updatedApps,
	}

	// Update the connection
	_, err = client.Update(ctx, connection.ID, updateParams)
	assert.NoError(t, err)

	// Get the connection again to verify the update
	connection, err = client.Get(ctx, connection.ID)
	assert.NoError(t, err)
	require.NotNil(t, connection)

	t.Logf("Connection Details after update:")
	t.Logf("  ID: %s", connection.ID)
	t.Logf("  Name: %s", connection.Name)
	t.Logf("  DisplayName: %s", connection.DisplayName)
	t.Logf("  Strategy: %s", connection.Strategy)

	// Verify connections through the applications endpoint after update
	// App1 and App2 should no longer have the connection
	app1Connections, err = appsClient.GetConnections(ctx, app1.ID)
	assert.NoError(t, err)
	foundInApp1 = false
	for _, conn := range app1Connections {
		if conn.ID == connection.ID {
			foundInApp1 = true
			break
		}
	}
	assert.False(t, foundInApp1, "Connection should not be found in app1's connections after update")

	app2Connections, err = appsClient.GetConnections(ctx, app2.ID)
	assert.NoError(t, err)
	foundInApp2 = false
	for _, conn := range app2Connections {
		if conn.ID == connection.ID {
			foundInApp2 = true
			break
		}
	}
	assert.False(t, foundInApp2, "Connection should not be found in app2's connections after update")

	// App3 should now have the connection
	app3Connections, err := appsClient.GetConnections(ctx, app3.ID)
	assert.NoError(t, err)
	require.NotNil(t, app3Connections)
	t.Logf("App3 Connections after update: %+v", app3Connections)
	
	var foundInApp3 bool
	for _, conn := range app3Connections {
		if conn.ID == connection.ID {
			foundInApp3 = true
			assert.Equal(t, name, conn.Name)
			assert.Equal(t, displayName, conn.DisplayName)
			assert.Equal(t, string(StrategyOAuth2Google), conn.Strategy)
			break
		}
	}
	assert.True(t, foundInApp3, "Connection should be found in app3's connections after update")
}

// TestE2EUpdateAndReplaceConnection tests updating and replacing a connection
func TestE2EUpdateAndReplaceConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	ctx := context.TODO()
	client := New(testutil.DefaultE2EClient(t))
	
	// Create a test connection
	timestamp := time.Now().UnixMilli()
	name := fmt.Sprintf("test-conn-%d", timestamp)
	displayName := fmt.Sprintf("Test Connection %d", timestamp)

	params := CreateParams{
		Name:        name,
		DisplayName: displayName,
		Strategy:    StrategyOAuth2Google,
		Options: &SocialConnectionOptions{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
	}

	// Create the connection
	connection, err := client.Create(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, connection)
	require.NotEmpty(t, connection.ID)

	// Clean up the test connection after we're done
	defer func() {
		t.Logf("Cleaning up test connection: %s", connection.ID)
		err := client.Delete(ctx, connection.ID)
		assert.NoError(t, err)
	}()

	t.Logf("Created connection with ID: %s", connection.ID)

	// Test Update - update only specific fields
	updatedName := fmt.Sprintf("test-conn-updated-%d", timestamp)
	updateParams := UpdateParams{
		Name: updatedName,
	}

	updatedConn, err := client.Update(ctx, connection.ID, updateParams)
	assert.NoError(t, err)
	require.NotNil(t, updatedConn)
	assert.Equal(t, updatedName, updatedConn.Name)
	assert.Equal(t, displayName, updatedConn.DisplayName) // Should remain unchanged
	t.Log("Successfully updated connection name")

	// Test Replace - replace entire connection configuration
	replaceParams := ReplaceParams{
		Name:        fmt.Sprintf("test-conn-replaced-%d", timestamp),
		DisplayName: fmt.Sprintf("Test Connection Replaced %d", timestamp),
		Options: &SocialConnectionOptions{
			ClientID:          "new-client-id",
			ClientSecret:      "new-client-secret",
			IsUseCustomDomain: true,
		},
	}

	replacedConn, err := client.Replace(ctx, connection.ID, replaceParams)
	assert.NoError(t, err)
	require.NotNil(t, replacedConn)
	assert.Equal(t, replaceParams.Name, replacedConn.Name)
	assert.Equal(t, replaceParams.DisplayName, replacedConn.DisplayName)
	t.Log("Successfully replaced connection configuration")

	// Verify the final state
	finalConn, err := client.Get(ctx, connection.ID)
	assert.NoError(t, err)
	require.NotNil(t, finalConn)
	assert.Equal(t, replaceParams.Name, finalConn.Name)
	assert.Equal(t, replaceParams.DisplayName, finalConn.DisplayName)
}

// TestE2EConnectionsOrder tests the order of connections returned by the List endpoint
func TestE2EConnectionsOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	ctx := context.TODO()
	client := New(testutil.DefaultE2EClient(t))

	// Create multiple test connections with different names
	timestamp := time.Now().UnixMilli()
	testConnections := []struct {
		name        string
		displayName string
	}{
		{
			name:        fmt.Sprintf("test-conn-a-%d", timestamp),
			displayName: fmt.Sprintf("Test Connection A %d", timestamp),
		},
		{
			name:        fmt.Sprintf("test-conn-b-%d", timestamp),
			displayName: fmt.Sprintf("Test Connection B %d", timestamp),
		},
		{
			name:        fmt.Sprintf("test-conn-c-%d", timestamp),
			displayName: fmt.Sprintf("Test Connection C %d", timestamp),
		},
	}

	var createdConnections []*Connection
	// Create the test connections
	for _, tc := range testConnections {
		params := CreateParams{
			Name:        tc.name,
			DisplayName: tc.displayName,
			Strategy:    StrategyOAuth2Google,
			Options: &SocialConnectionOptions{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
			},
		}

		conn, err := client.Create(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, conn)
		createdConnections = append(createdConnections, conn)

		// Clean up after test
		defer func(id string) {
			t.Logf("Cleaning up test connection: %s", id)
			err := client.Delete(ctx, id)
			assert.NoError(t, err)
		}(conn.ID)
	}

	// Get all connections and analyze their order
	connections, err := client.List(ctx)
	require.NoError(t, err)
	require.NotNil(t, connections)

	t.Log("Connections in order returned by API:")
	for i, conn := range connections {
		t.Logf("%d. ID: %s, Name: %s, DisplayName: %s", i+1, conn.ID, conn.Name, conn.DisplayName)
	}

	// Find our test connections in the list and check their relative positions
	var foundIndexes []int
	for _, created := range createdConnections {
		for i, conn := range connections {
			if conn.ID == created.ID {
				foundIndexes = append(foundIndexes, i)
				break
			}
		}
	}

	t.Logf("Test connections found at indexes: %v", foundIndexes)

	// Check if the connections maintain creation order
	for i := 1; i < len(foundIndexes); i++ {
		assert.Greater(t, foundIndexes[i], foundIndexes[i-1], 
			"Connections should maintain creation order (newer connections appear later)")
	}
} 