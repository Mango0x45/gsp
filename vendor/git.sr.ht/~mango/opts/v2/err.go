package opts

import "fmt"

// A BadOptionError describes an option that the user attempted to pass
// which the developer did not register.
type BadOptionError struct {
	r rune
	s string
}

func (e BadOptionError) Error() string {
	if e.r != 0 {
		return fmt.Sprintf("unknown option ‘-%c’", e.r)
	}
	return fmt.Sprintf("unknown option ‘--%s’", e.s)
}

// A NoArgumentError describes an option that the user attempted to pass
// without an argument, which required an argument.
type NoArgumentError struct {
	r rune
	s string
}

func (e NoArgumentError) Error() string {
	if e.r != 0 {
		return fmt.Sprintf("expected argument for option ‘-%c’", e.r)
	}
	return fmt.Sprintf("expected argument for option ‘--%s’", e.s)
}
