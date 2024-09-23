package enum

import (
	"fmt"
	"strings"
)

type Enum[T ~string] interface {
	Options() []T
	Valid(string) error
}

type InvalidEnumMemberError struct {
	Options []string
	Value   string
}

func (err InvalidEnumMemberError) Error() string {
	options := strings.Join(err.Options, ", ")
	return fmt.Sprintf("invalid enum member %s must be one of [%s]", err.Value, options)
}

func Valid[T ~string](options []T, value string) error {
	for _, v := range options {
		if string(v) == value {
			return nil
		}
	}

	optionsAsString := make([]string, 0, len(options))
	for _, v := range options {
		optionsAsString = append(optionsAsString, string(v))
	}

	return InvalidEnumMemberError{
		Options: optionsAsString,
		Value:   value,
	}
}
