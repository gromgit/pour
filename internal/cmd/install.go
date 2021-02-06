package cmd

import (
	"fmt"
	"github.com/gromgit/litebrew/internal/formula"
)

func Install(allf *formula.Formulas, args []string) (err error) {
	for _, name := range args {
		f := (*allf)[name]
		if f.Name != "" {
			// First install all dependencies
			Install(allf, f.Dependencies)
			switch f.Status {
			case formula.INSTALLED:
				fmt.Println("DEBUG: Already installed", f.Name)
				continue
			case formula.OUTDATED:
				fmt.Println("DEBUG: Upgrade", f.Name)
			case formula.MISSING:
				fmt.Println("DEBUG: Install", f.Name)
			}
		}
	}
	return
}
