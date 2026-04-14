// Package opts implements unicode-aware getopt(3)- and getopt_long(3)
// flag parsing.
//
// The opts package aims to provide as simple an API as possible.  If
// your usecase requires more advanced argument-parsing or a more robust
// API, this may not be the ideal package for you.
//
// While the opts package aims to closely follow the POSIX getopt(3) and
// GNU getopt_long(3) C functions, there are some notable differences.
// This package properly supports unicode flags, but also does not
// support a leading ‘:’ in the [Get] function’s option string; all
// user-facing I/O is delegrated to the caller.
package opts

import (
	"slices"
	"strings"
)

// ArgMode represents whether or not a long-option takes an argument.
type ArgMode int

// These tokens can be used to specify whether or not a long-option takes
// an argument.
const (
	None     ArgMode = iota // long opt takes no argument
	Required                // long opt takes an argument
	Optional                // long opt optionally takes an argument
)

// Flag represents a parsed command-line flag.  Key corresponds to the
// rune that was passed on the command-line, and Value corresponds to the
// flags argument if one was provided.  In the case of long-options Key
// will map to the corresponding short-code, even if a long-option was
// used.
type Flag struct {
	Key   rune   // the flag that was passed
	Value string // the flags argument
}

// LongOpt represents a long-option to attempt to parse.  All long
// options have a short-hand form represented by Short and a long-form
// represented by Long.  Arg is used to represent whether or not the
// long-option takes an argument.
//
// In the case that you want to parse a long-option which doesn’t have a
// short-hand form, you can set Short to a negative integer.
type LongOpt struct {
	Short rune
	Long  string
	Arg   ArgMode
}

// Get parses the command-line arguments in args according to optstr.
// Unlike POSIX-getopt(3), a leading ‘:’ in optstr is not supported and
// will be ignored and no I/O is ever performed.
//
// Get will look for the flags listed in optstr (i.e., it will look for
// ‘-a’, ‘-ß’, and ‘λ’ given optstr == "aßλ").  The optstr need not be
// sorted in any particular order.  If an option takes a required
// argument, it can be suffixed by a colon.  If an option takes an
// optional argument, it can be suffixed by two colons.  As an example,
// optstr == "a::ßλ:" will search for ‘-a’ with an optional argument,
// ‘-ß’ with no argument, and ‘-λ’ with a required argument.
//
// A successful parse returns the flags in the flags slice and a slice of
// the remaining non-option arguments in rest.  In the case of failure,
// err will be one of [BadOptionError] or [NoArgumentError].
func Get(args []string, optstr string) (flags []Flag, rest []string, err error) {
	if len(args) == 0 {
		return
	}

	optrs := []rune(optstr)
	var i int
	for i = 1; i < len(args); i++ {
		arg := args[i]
		if len(arg) == 0 || arg == "-" || arg[0] != '-' {
			break
		} else if arg == "--" {
			i++
			break
		}

		rs := []rune(arg[1:])
		for j, r := range rs {
			k := slices.Index(optrs, r)
			if k == -1 {
				return nil, nil, BadOptionError{r: r}
			}

			var s string
			switch am := colonsToArgMode(optrs[k+1:]); {
			case am != None && j < len(rs)-1:
				s = string(rs[j+1:])
			case am == Required:
				i++
				if i >= len(args) {
					return nil, nil, NoArgumentError{r: r}
				}
				s = args[i]
			default:
				flags = append(flags, Flag{Key: r})
				continue
			}

			flags = append(flags, Flag{r, s})
			break
		}
	}

	return flags, args[i:], nil
}

// GetLong parses the command-line arguments in args according to opts.
//
// This function is identical to [Get] except it parses according to a
// [LongOpt] slice instead of an opt-string, and it parses long-options.
// When parsing, GetLong will also accept incomplete long-options if they
// are unambiguous.  For example, given the following definition of opts:
//
//	opts := []LongOpt{
//		{Short: 'a', Long: "add", Arg: None},
//		{Short: 'd', Long: "delete", Arg: None},
//		{Short: 'D', Long: "defer", Arg: None},
//	}
//
// The options ‘--a’ and ‘--ad’ will parse as ‘--add’.  The option ‘--de’
// will not parse however as it is ambiguous.
func GetLong(args []string, opts []LongOpt) (flags []Flag, rest []string, err error) {
	if len(args) == 0 {
		return
	}

	var i int
	for i = 1; i < len(args); i++ {
		arg := args[i]
		if len(arg) == 0 || arg == "-" || arg[0] != '-' {
			break
		} else if arg == "--" {
			i++
			break
		}

		if strings.HasPrefix(arg, "--") {
			arg = arg[2:]

			n := arg
			j := strings.IndexByte(n, '=')
			if j != -1 {
				n = arg[:j]
			}

			var s string
			o, ok := optStruct(opts, n)

			switch {
			case !ok:
				return nil, nil, BadOptionError{s: n}
			case o.Arg != None && j != -1:
				s = arg[j+1:]
			case o.Arg == Required:
				i++
				if i >= len(args) {
					return nil, nil, NoArgumentError{s: n}
				}
				s = args[i]
			}

			flags = append(flags, Flag{Key: o.Short, Value: s})
		} else {
			rs := []rune(arg[1:])
			for j, r := range rs {
				var s string

				switch am, ok := getModeRune(opts, r); {
				case !ok:
					return nil, nil, BadOptionError{r: r}
				case am != None && j < len(rs)-1:
					s = string(rs[j+1:])
				case am == Required:
					i++
					if i >= len(args) {
						return nil, nil, NoArgumentError{r: r}
					}
					s = args[i]
				default:
					flags = append(flags, Flag{Key: r})
					continue
				}

				flags = append(flags, Flag{r, s})
				break
			}
		}
	}

	return flags, args[i:], nil
}

func getModeRune(os []LongOpt, r rune) (ArgMode, bool) {
	for _, o := range os {
		if o.Short == r {
			return o.Arg, true
		}
	}
	return 0, false
}

func optStruct(os []LongOpt, s string) (LongOpt, bool) {
	i := -1
	for j, o := range os {
		if strings.HasPrefix(o.Long, s) {
			if i != -1 {
				return LongOpt{}, false
			}
			i = j
		}
	}
	if i == -1 {
		return LongOpt{}, false
	}
	return os[i], true
}

func colonsToArgMode(rs []rune) ArgMode {
	if len(rs) >= 2 && rs[0] == ':' && rs[1] == ':' {
		return Optional
	}
	if len(rs) >= 1 && rs[0] == ':' {
		return Required
	}
	return None
}
