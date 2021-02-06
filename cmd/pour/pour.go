package main

import (
	"fmt"
	"github.com/gromgit/pour/internal/cmd"
	cfg "github.com/gromgit/pour/internal/config"
	"github.com/gromgit/pour/internal/formula"
	"github.com/gromgit/pour/internal/log"
	"os"
)

var allf formula.Formulas

func help(args []string) {
	fmt.Println(`Available subcommands:
  info [<formula>...]
  install <formula>...
  link <formula>...
  list, ls [<formula>...]
  outdated
  leaves
  deps [--installed] [--common] <formula>...
  uses [--installed] [--recursive] <formula>
  pin [<formula>...]
  search [--desc] [<text> | /<regex>/]
  shellenv
  uninstall, remove, rm <formula>...
  unlink <formula>...
  unpin [<formula>...]
  update, up
  upgrade [<formula>...]`)
}

func fatal(args ...interface{}) {
	fmt.Fprintf(os.Stderr, "FATAL ERROR: %+v", args)
	os.Exit(1)
}

func baseChecks() {
	// Check for basic writability in PREFIX
	if _, err := os.Stat(cfg.PREFIX); os.IsNotExist(err) {
		fatal("Pour base dir", cfg.PREFIX, "doesn't exist")
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
	if len(args) < 1 {
		help(args)
		rtn = 1
	} else {
		switch args[0] {
		case "help", "-h", "--help":
			help(args[1:])
		case "shellenv":
			if err := cmd.Shellenv(args[1:]); err != nil {
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
	// First check global options
	args := os.Args[1:]
GlobalOptions:
	for len(args) > 0 {
		switch args[0] {
		case "--debug":
			if f, err := os.Create(args[1]); err != nil {
				fatal("Can't create debug file", err)
			} else {
				defer f.Close()
				log.File(f)
				args = args[1:]
			}
		default:
			break GlobalOptions
		}
		args = args[1:]
	}
	baseChecks()
	json_path := cfg.JSON_PATH
	if rtn, quit := doMeta(json_path, args); quit {
		os.Exit(rtn)
	}

	if _, err := os.Stat(json_path); os.IsNotExist(err) {
		if err := cmd.Update(json_path); err != nil {
			fatal("Unable to update local JSON:", err)
		}
	}
	allf.Load(json_path)
	if err := cmd.Install(&allf, cfg.OS_DEPS); err != nil {
		fatal("Unable to install OS prerequisites:", err)
	}

	var err error
	switch args[0] {
	case "search":
		err = cmd.Search(allf, args[1:])
	case "info":
		err = cmd.Info(allf, args[1:])
	case "deps":
		err = cmd.Deps(allf, args[1:])
	case "uses":
		err = cmd.Uses(allf, args[1:])
	case "install":
		err = cmd.Install(&allf, args[1:])
	case "link":
		err = cmd.Link(&allf, args[1:])
	case "pin":
		err = cmd.Pin(allf, args[1:])
	case "unpin":
		err = cmd.Unpin(allf, args[1:])
	case "upgrade":
		err = cmd.Upgrade(allf, args[1:])
	case "uninstall", "remove", "rm":
		err = cmd.Uninstall(allf, args[1:])
	case "unlink":
		err = cmd.Unlink(&allf, args[1:])
	case "list", "ls":
		err = cmd.List(allf, args[1:])
	case "outdated":
		err = cmd.Outdated(allf, args[1:])
	case "leaves":
		err = cmd.Leaves(allf, args[1:])
	default:
		err = fmt.Errorf("Unknown subcommand '%s'", args[0])
	}
	if err != nil {
		fatal(err)
	}
}
