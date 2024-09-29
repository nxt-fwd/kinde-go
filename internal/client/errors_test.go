package client_test

import (
	"fmt"
	"testing"

	"github.com/axatol/kinde-go"
	"github.com/stretchr/testify/assert"
)

func TestRequestError(t *testing.T) {
	err0 := &kinde.KindeErrors{
		{Code: "FOO"},
		{Code: "BAR"},
	}

	err1 := &kinde.RequestError{
		Method:     "GET",
		Path:       "/api/v1/apis",
		StatusCode: 404,
		Err:        err0,
	}

	err2 := fmt.Errorf("failed to execute request: %w", err1)

	assert.ErrorIs(t, err2, err1)
}
