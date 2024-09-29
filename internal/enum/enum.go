package enum

import (
	"fmt"
	"strings"
)

type Enum[T ~string] interface {
	Options() []T
	Valid() error
}

type InvalidEnumMemberError struct {
	Options []string
	Value   string
}

func (err InvalidEnumMemberError) Error() string {
	options := strings.Join(err.Options, ", ")
	return fmt.Sprintf("invalid enum member %s must be one of [%s]", err.Value, options)
}

func Valid[T ~string](options []T, value T) error {
	for _, option := range options {
		if option == value {
			return nil
		}
	}

	optionsAsString := make([]string, 0, len(options))
	for _, v := range options {
		optionsAsString = append(optionsAsString, string(v))
	}

	return InvalidEnumMemberError{
		Options: optionsAsString,
		Value:   string(value),
	}
}
