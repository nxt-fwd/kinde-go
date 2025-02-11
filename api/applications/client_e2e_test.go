//go:build e2e
// +build e2e

package applications_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nxt-fwd/kinde-go/api/applications"
	"github.com/nxt-fwd/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EList(t *testing.T) {
	client := applications.New(testutil.DefaultE2EClient(t))
	res, err := client.List(context.TODO(), applications.ListParams{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2ECreateGetUpdateDelete(t *testing.T) {
	client := applications.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	res, err := client.Create(context.TODO(), applications.CreateParams{Name: tempID, Type: applications.TypeRegular})
	assert.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.ID)

	id := res.ID

	t.Logf("created test application: %s\n", id)

	res, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, res)

	t.Logf("got test application: %+v\n", res)

	updateParams := applications.UpdateParams{
		Name:         tempID + "-updated",
		LoginURI:     "https://example.com",
		HomepageURI:  "https://example.com",
		LogoutURIs:   []string{"https://example.com"},
		RedirectURIs: []string{"https://example.com"},
	}

	err = client.Update(context.TODO(), id, updateParams)
	assert.NoError(t, err)

	// Verify the updated parameters
	updated, err := client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, updated)

	// Assert that the parameters were updated correctly
	assert.Equal(t, updateParams.Name, updated.Name)
	assert.Equal(t, updateParams.LoginURI, updated.LoginURI)
	assert.Equal(t, updateParams.HomepageURI, updated.HomepageURI)
	// Note: LogoutURIs and RedirectURIs are set via the update request
	// but are not returned in the GET response from the Kinde API

	t.Logf("updated test application: %+v\n", updated)

	err = client.Delete(context.TODO(), id)
	assert.NoError(t, err)

	t.Logf("deleted test application: %+v\n", updated)
}

func TestE2EApplicationTypes(t *testing.T) {
	client := applications.New(testutil.DefaultE2EClient(t))
	baseID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	testCases := []struct {
		name string
		type_ applications.Type
		validateSettings func(t *testing.T, app *applications.Application)
	}{
		{
			name: "Regular Web Application",
			type_: applications.TypeRegular,
			validateSettings: func(t *testing.T, app *applications.Application) {
				assert.Equal(t, applications.TypeRegular, app.Type)
				assert.NotEmpty(t, app.ClientSecret, "Regular apps should have a client secret")
			},
		},
		{
			name: "Single Page Application",
			type_: applications.TypeSinglePageApplication,
			validateSettings: func(t *testing.T, app *applications.Application) {
				assert.Equal(t, applications.TypeSinglePageApplication, app.Type)
				assert.Empty(t, app.ClientSecret, "SPA should not have a client secret")
			},
		},
		{
			name: "Machine to Machine Application",
			type_: applications.TypeMachineToMachine,
			validateSettings: func(t *testing.T, app *applications.Application) {
				assert.Equal(t, applications.TypeMachineToMachine, app.Type)
				assert.NotEmpty(t, app.ClientSecret, "M2M apps should have a client secret")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			appName := fmt.Sprintf("%s-%s", baseID, tc.type_)
			
			// Create application
			app, err := client.Create(context.TODO(), applications.CreateParams{
				Name: appName,
				Type: tc.type_,
			})
			assert.NoError(t, err)
			require.NotNil(t, app)
			require.NotEmpty(t, app.ID)
			t.Logf("created %s application: %s", tc.type_, app.ID)

			// Get and validate application
			app, err = client.Get(context.TODO(), app.ID)
			assert.NoError(t, err)
			require.NotNil(t, app)

			// Validate type-specific settings
			tc.validateSettings(t, app)

			// Test type-specific updates
			updateParams := applications.UpdateParams{
				Name: appName + "-updated",
				LoginURI: "https://example.com/login",
				HomepageURI: "https://example.com",
				LogoutURIs: []string{"https://example.com/logout"},
				RedirectURIs: []string{"https://example.com/callback"},
			}

			err = client.Update(context.TODO(), app.ID, updateParams)
			assert.NoError(t, err)
			t.Log("updated application settings")

			// Verify updates
			updated, err := client.Get(context.TODO(), app.ID)
			assert.NoError(t, err)
			require.NotNil(t, updated)
			assert.Equal(t, updateParams.Name, updated.Name)
			assert.Equal(t, updateParams.LoginURI, updated.LoginURI)
			assert.Equal(t, updateParams.HomepageURI, updated.HomepageURI)
			t.Log("verified application updates")

			// Clean up
			err = client.Delete(context.TODO(), app.ID)
			assert.NoError(t, err)
			t.Logf("deleted test application: %s", app.ID)
		})
	}
}
