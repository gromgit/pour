package main

import (
	"fmt"
	"github.com/gromgit/pour/internal/cmd"
	cfg "github.com/gromgit/pour/internal/config"
	"github.com/gromgit/pour/internal/formula"
	"os"
)

var allf formula.Formulas

func help(args []string) {
	fmt.Println(`Available subcommands:
  info [<formula>...]
  install <formula>...
  list, ls [<formula>...]
  outdated
  pin [<formula>...]
  search [--desc] [<text> | /<regex>/]
  shellenv
  uninstall, remove, rm <formula>...
  unpin [<formula>...]
  update, up
  upgrade [<formula>...]`)
}

func fatal(args ...interface{}) {
	cfg.Log("FATAL ERROR:", args)
	os.Exit(1)
}

func baseChecks() {
	// Check for basic writability in DEFAULT_PREFIX
	if _, err := os.Stat(cfg.DEFAULT_PREFIX); os.IsNotExist(err) {
		fatal("Pour base dir", cfg.DEFAULT_PREFIX, "doesn't exist")
	} else {
		failedDirs := []string{}
		for _, d := range cfg.SYSDIRS {
			if err := os.MkdirAll(d, 0775); err != nil {
				failedDirs = append(failedDirs, d)
			}
		}
		if len(failedDirs) > 0 {
			fatal("Can't create the following dirs:", failedDirs)
		}
	}
}

func doMeta(json_path string, args []string) (rtn int, quit bool) {
	rtn = 0
	quit = true
	// Check for subcommands
	if len(os.Args) < 2 {
		help(os.Args)
		rtn = 1
	} else {
		switch os.Args[1] {
		case "help", "-h", "--help":
			help(os.Args[2:])
		case "shellenv":
			if err := cmd.Shellenv(os.Args[2:]); err != nil {
				fatal(err)
			}
		case "update", "up":
			if err := cmd.Update(json_path); err != nil {
				fatal("Unable to update local JSON:", err)
			}
		default:
			// Not a metacommand, need to continue
			quit = false
		}
	}
	return
}

func main() {
	baseChecks()
	json_path := cfg.JSON_PATH
	if rtn, quit := doMeta(json_path, os.Args); quit {
		os.Exit(rtn)
	}

	if _, err := os.Stat(json_path); os.IsNotExist(err) {
		if err := cmd.Update(json_path); err != nil {
			fatal("Unable to update local JSON:", err)
		}
	}
	allf.Load(json_path)

	var err error
	switch os.Args[1] {
	case "search":
		err = cmd.Search(allf, os.Args[2:])
	case "info":
		err = cmd.Info(allf, os.Args[2:])
	case "install":
		err = cmd.Install(&allf, os.Args[2:])
	case "pin":
		err = cmd.Pin(allf, os.Args[2:])
	case "unpin":
		err = cmd.Unpin(allf, os.Args[2:])
	case "upgrade":
		err = cmd.Upgrade(allf, os.Args[2:])
	case "uninstall", "remove", "rm":
		err = cmd.Uninstall(allf, os.Args[2:])
	case "list", "ls":
		err = cmd.List(allf, os.Args[2:])
	case "outdated":
		err = cmd.Outdated(allf, os.Args[2:])
	default:
		err = fmt.Errorf("Unknown subcommand '%s'", os.Args[1])
	}
	if err != nil {
		fatal(err)
	}
}
