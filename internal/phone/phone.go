package phone

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// ParseNumber takes a full international format phone number and returns
// the local number and country ID for use with Kinde API
func ParseNumber(fullNumber string) (localNumber, countryID string, err error) {
	// Parse the phone number
	num, err := phonenumbers.Parse(fullNumber, "")
	if err != nil {
		return "", "", fmt.Errorf("invalid phone number format: %w", err)
	}

	// Validate the number
	if !phonenumbers.IsValidNumber(num) {
		return "", "", fmt.Errorf("invalid phone number")
	}

	// Get the country code
	region := phonenumbers.GetRegionCodeForNumber(num)
	countryID = strings.ToLower(region)

	// Get the national number (without country code)
	localNumber = fmt.Sprint(num.GetNationalNumber())

	return localNumber, countryID, nil
}

// FormatNumber takes a local number and country ID and returns
// the full international format phone number
func FormatNumber(localNumber, countryID string) (string, error) {
	// Parse the phone number with the country code
	num, err := phonenumbers.Parse(localNumber, strings.ToUpper(countryID))
	if err != nil {
		return "", fmt.Errorf("invalid phone number format: %w", err)
	}

	// Validate the number
	if !phonenumbers.IsValidNumber(num) {
		return "", fmt.Errorf("invalid phone number")
	}

	// Format the number in international format
	formatted := phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
	
	// Remove any spaces and hyphens, ensure it starts with +
	formatted = strings.ReplaceAll(formatted, " ", "")
	formatted = strings.ReplaceAll(formatted, "-", "")
	if !strings.HasPrefix(formatted, "+") {
		formatted = "+" + formatted
	}

	return formatted, nil
} 