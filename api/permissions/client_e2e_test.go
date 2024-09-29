//go:build e2e
// +build e2e

package permissions_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/axatol/kinde-go/api/permissions"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EList(t *testing.T) {
	client := permissions.New(testutil.DefaultE2EClient(t))
	res, err := client.List(context.TODO(), permissions.ListParams{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
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
