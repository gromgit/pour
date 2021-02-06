package cmd

import (
	"fmt"
	"github.com/gromgit/litebrew/internal/formula"
)

func Install(formulas *formula.Formulas, args []string) {
	for _, name := range args {
		f := (*formulas)[name]
		if f.Name != "" {
			// First install all dependencies
			Install(formulas, f.Dependencies)
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
}
