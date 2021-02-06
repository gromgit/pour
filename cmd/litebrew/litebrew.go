package main

import (
	"fmt"
	"github.com/gromgit/litebrew/internal/cmd"
	"github.com/gromgit/litebrew/internal/config"
	"github.com/gromgit/litebrew/internal/formula"
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
	fmt.Fprintln(os.Stderr, "FATAL ERROR:", args)
	os.Exit(1)
}

func baseChecks() {
	// Check for basic writability in DEFAULT_PREFIX
	if _, err := os.Stat(config.DEFAULT_PREFIX); os.IsNotExist(err) {
		fatal("Litebrew base dir", config.DEFAULT_PREFIX, "doesn't exist")
	} else if err = os.MkdirAll(config.VAR_PATH, 0775); err != nil {
		fatal("Can't create", config.VAR_PATH, ":", err)
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
			cmd.Shellenv(os.Args[2:])
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
	json_path := config.JSON_PATH
	if rtn, quit := doMeta(json_path, os.Args); quit {
		os.Exit(rtn)
	}

	if _, err := os.Stat(json_path); os.IsNotExist(err) {
		if err := cmd.Update(json_path); err != nil {
			fatal("Unable to update local JSON:", err)
		}
	}
	allf.Load(json_path)

	switch os.Args[1] {
	case "search":
		cmd.Search(allf, os.Args[2:])
	case "info":
		cmd.Info(allf, os.Args[2:])
	case "install":
		cmd.Install(&allf, os.Args[2:])
	case "pin":
		cmd.Pin(allf, os.Args[2:])
	case "unpin":
		cmd.Unpin(allf, os.Args[2:])
	case "upgrade":
		cmd.Upgrade(allf, os.Args[2:])
	case "uninstall", "remove", "rm":
		cmd.Uninstall(allf, os.Args[2:])
	case "list", "ls":
		cmd.List(allf, os.Args[2:])
	case "outdated":
		cmd.Outdated(allf, os.Args[2:])
	default:
		fmt.Printf("Unknown subcommand '%s'\n", os.Args[1])
		os.Exit(1)
	}
}
