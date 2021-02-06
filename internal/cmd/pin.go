package cmd

import (
	"fmt"
	"github.com/gromgit/litebrew/internal/formula"
	"os"
)

func Pin(allf formula.Formulas, args []string) {
	if len(args) > 0 {
		// Pin some formulas
		for _, i := range args {
			if f := allf[i]; f.Name == "" {
				fmt.Fprintf(os.Stderr, "WARNING: Formula '%s' not found\n", i)
			} else {
				if err := f.Pin(); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR pinning '%s': %s\n", i, err.Error())
				}
			}

		}
	} else {
		// List all pinned bottles
		allf.Filter(
			func(item formula.Formula) bool {
				return item.Pinned
			}).
			Ls()
	}
}

func Unpin(allf formula.Formulas, args []string) {
	if len(args) > 0 {
		// Unpin some formulas
		for _, i := range args {
			if f := allf[i]; f.Name == "" {
				fmt.Fprintf(os.Stderr, "WARNING: Formula '%s' not found\n", i)
			} else {
				if err := f.Unpin(); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR unpinning '%s': %s\n", i, err.Error())
				}
			}

		}
	} else {
		// List all unpinned bottles
		allf.Filter(
			func(item formula.Formula) bool {
				return !item.Pinned
			}).
			Ls()
	}
}
