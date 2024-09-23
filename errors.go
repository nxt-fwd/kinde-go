package kinde

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type NotFoundError struct {
	Kind string
	ID   string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("failed to find %s with id %s", err.Kind, err.ID)
}

type RequestError struct {
	Method     string
	Path       string
	StatusCode int
	Err        error
}

func (err RequestError) Error() string {
	return fmt.Sprintf("failed to execute %s %s: %s", err.Method, err.Path, err.Err)
}

type KindeErrors []KindeError

func (errs KindeErrors) Error() string {
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		messages = append(messages, err.Error())
	}

	return strings.Join(messages, ", ")
}

func (errs KindeErrors) Has(code string) bool {
	for _, err := range errs {
		if err.Code == code {
			return true
		}
	}

	return false
}

func (errs *KindeErrors) UnmarshalJSON(data []byte) error {
	rawErrs := gjson.GetBytes(data, "errors")

	// its possible the kinde api may return an array or a single object
	// we need to handle both

	if rawErrs.IsArray() {
		var target struct {
			Errors []KindeError `json:"errors"`
		}

		if err := json.Unmarshal(data, &target); err != nil {
			return fmt.Errorf("failed to parse error list: %w", err)
		}

		*errs = target.Errors
		return nil
	}

	var target struct {
		Errors KindeError `json:"errors"`
	}

	if err := json.Unmarshal(data, &target); err != nil {
		return fmt.Errorf("failed to parse error: %w", err)
	}

	*errs = KindeErrors{target.Errors}
	return nil
}

type KindeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e KindeError) Error() string {
	builder := strings.Builder{}
	builder.WriteString(e.Code)
	builder.WriteString(": ")
	if e.Message != "" {
		builder.WriteString(e.Message)
	} else {
		builder.WriteString("N/A")
	}

	return builder.String()
}
