package cmd

import (
	"fmt"
	"github.com/gromgit/pour/internal/formula"
	"os"
	"path/filepath"
)

func List(allf formula.Formulas, args []string) error {
	if len(args) > 0 {
		for _, name := range args {
			f := allf[name]
			if f.Name != "" {
				err := filepath.Walk(f.InstallDir,
					func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						if info.IsDir() {
							// Skip .brew dir
							if info.Name() == ".brew" {
								return filepath.SkipDir
							}
						} else {
							fmt.Println(path)
						}
						return nil
					})
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		// List installed bottles
		allf.Filter(
			func(item formula.Formula) bool {
				return item.Installed()
			}).
			Ls()
	}
	return nil
}

func Outdated(allf formula.Formulas, args []string) error {
	allf.Filter(
		func(item formula.Formula) bool {
			return item.Status == formula.OUTDATED
		}).
		Ls()
	return nil
}
