package cmd

import (
	"fmt"
	"github.com/gromgit/pour/internal/formula"
	"os"
	"regexp"
	"strings"
)

// pour search [TEXT|/REGEX/]
func StringMatcher(m string) func(s string) bool {
	return func(s string) bool {
		return strings.Index(s, m) >= 0
	}
}

func RegexMatcher(m string) func(s string) bool {
	return func(s string) bool {
		match, err := regexp.MatchString(m, s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "regexp error: %s\n", err)
			return false
		} else {
			return match
		}
	}
}

func NameGetter(item formula.Formula) string {
	return item.Name
}

func DescGetter(item formula.Formula) string {
	return item.Desc
}

func Search(allf formula.Formulas, args []string) error {
	var matcher func(s string) bool
	var getter func(item formula.Formula) string
	fmt.Printf("Doing search %v\n", args)
	getter = NameGetter
	// Handle options
SearchOptions:
	for len(args) > 0 {
		switch {
		case strings.HasPrefix(args[0], "--desc"):
			getter = DescGetter
		case strings.HasPrefix(args[0], "--inst"):
			// Filter out only the installed stuff
			allf = allf.Filter(func(f formula.Formula) bool {
				return f.Status == formula.INSTALLED
			})
		default:
			break SearchOptions
		}
		args = args[1:]
	}
	if len(args) == 0 {
		// Return all bottles
		allf.Ls()
	} else {
		spec := args[0]
		if spec[0] == '/' && spec[len(spec)-1] == '/' {
			// Regex search
			matcher = RegexMatcher(spec[1 : len(spec)-1])
		} else {
			// String search
			matcher = StringMatcher(spec)
		}
		allf.Filter(
			func(item formula.Formula) bool {
				return matcher(getter(item))
			}).
			Ls()
	}
	return nil
}
