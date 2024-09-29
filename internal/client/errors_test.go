package client_test

import (
	"fmt"
	"testing"

	"github.com/axatol/kinde-go/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestRequestError(t *testing.T) {
	err0 := &client.KindeErrors{
		{Code: "FOO"},
		{Code: "BAR"},
	}

	err1 := &client.RequestError{
		Method:     "GET",
		Path:       "/api/v1/apis",
		StatusCode: 404,
		Err:        err0,
	}

	err2 := fmt.Errorf("failed to execute request: %w", err1)

	assert.ErrorIs(t, err2, err1)
}
