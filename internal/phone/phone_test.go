package phone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNumber(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantLocal     string
		wantCountryID string
		wantErr       bool
	}{
		{
			name:          "Armenian number",
			input:         "+37455251234",
			wantLocal:     "55251234",
			wantCountryID: "am",
			wantErr:       false,
		},
		{
			name:          "Australian number",
			input:         "+61412345678",
			wantLocal:     "412345678",
			wantCountryID: "au",
			wantErr:       false,
		},
		{
			name:          "US number",
			input:         "+12025550123",
			wantLocal:     "2025550123",
			wantCountryID: "us",
			wantErr:       false,
		},
		{
			name:          "Invalid number",
			input:         "+1234",
			wantLocal:     "",
			wantCountryID: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLocal, gotCountryID, err := ParseNumber(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantLocal, gotLocal)
			assert.Equal(t, tt.wantCountryID, gotCountryID)
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name        string
		localNumber string
		countryID   string
		want        string
		wantErr     bool
	}{
		{
			name:        "Armenian number",
			localNumber: "55251234",
			countryID:   "am",
			want:        "+37455251234",
			wantErr:     false,
		},
		{
			name:        "Australian number",
			localNumber: "412345678",
			countryID:   "au",
			want:        "+61412345678",
			wantErr:     false,
		},
		{
			name:        "US number",
			localNumber: "2025550123",
			countryID:   "us",
			want:        "+12025550123",
			wantErr:     false,
		},
		{
			name:        "Invalid number",
			localNumber: "123",
			countryID:   "us",
			want:        "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatNumber(tt.localNumber, tt.countryID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
