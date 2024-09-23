//go:build e2e
// +build e2e

package e2e

import (
	"encoding/json"
	"fmt"

	_ "embed"
)

type e2eValues struct {
	Domain       string
	Audience     string
	ClientID     string
	ClientSecret string
}

var (
	//go:embed e2e.json
	e2eRaw []byte

	Domain       string
	Audience     string
	ClientID     string
	ClientSecret string
)

func init() {
	var values e2eValues
	if err := json.Unmarshal(e2eRaw, &values); err != nil {
		panic(fmt.Errorf("failed to parse e2e credentials: %w", err))
	}

	Domain = values.Domain
	Audience = values.Audience
	ClientID = values.ClientID
	ClientSecret = values.ClientSecret
}
