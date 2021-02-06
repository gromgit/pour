package cmd

import (
	"fmt"
	"github.com/gromgit/pour/internal/bottle"
	"github.com/gromgit/pour/internal/config"
	"github.com/gromgit/pour/internal/formula"
	"path/filepath"
)

func Uninstall(allf formula.Formulas, args []string) error {
	fmt.Println("This command has not been implemented yet.")
	return nil
}

func Unlink(allf *formula.Formulas, args []string) (err error) {
	for _, name := range args {
		f := (*allf)[name]
		if f.Installed() && f.InstallDir != "" {
			rel, err := filepath.Rel(config.CELLAR, f.InstallDir)
			if err != nil {
				return err
			}
			err = bottle.Unlink(rel)
			if err != nil {
				return err
			}
		}
	}
	return
}
