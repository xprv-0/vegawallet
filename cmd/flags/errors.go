package flags

import (
	"fmt"
)

type FlagError struct {
	message string
}

func (f FlagError) Error() string {
	return f.message
}

func FlagsMutuallyExclusiveError(n1, n2 string) error {
	return FlagError{
		message: fmt.Sprintf("--%s and --%s flags are mutually exclusive", n1, n2),
	}
}

func FlagMustBeSpecifiedError(name string) error {
	return FlagError{
		message: fmt.Sprintf("--%s flag must be specified", name),
	}
}

func ArgMustBeSpecifiedError(name string) error {
	return FlagError{
		message: fmt.Sprintf("%s argument must be specified", name),
	}
}

func OneOfFlagsMustBeSpecifiedError(n1, n2 string) error {
	return FlagError{
		message: fmt.Sprintf("--%s or --%s flags must be specified", n1, n2),
	}
}

func InvalidFlagFormatError(name string) error {
	return FlagError{
		message: fmt.Sprintf("--%s flag has not a valid format", name),
	}
}

func UnsupportedFlagValueError(name string, unsupported interface{}, supported []interface{}) error {
	return FlagError{
		message: fmt.Sprintf("--%s flag doesn't support value %s in %v", name, unsupported, supported),
	}
}

func ParentFlagMustBeSpecifiedError(name string, parent string) error {
	return FlagError{
		message: fmt.Sprintf("--%s flag requires --%s flag to be set", name, parent),
	}
}

func MustBase64EncodedError(name string) error {
	return FlagError{
		message: fmt.Sprintf("--%s flag value must be base64-encoded", name),
	}
}