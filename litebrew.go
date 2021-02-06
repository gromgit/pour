package main

import (
	"fmt"
	"github.com/gromgit/litebrew/internal/cmd"
	"github.com/gromgit/litebrew/internal/config"
	"github.com/gromgit/litebrew/internal/formula"
	"os"
)

var formulas formula.Formulas

func help(args []string) {
	fmt.Println(`LITEBREW SUBCOMMANDS:
  search [--desc] [<text> | /<regex>/]
  info [<formula>...]
  install <formula>...
  update, up
  outdated [-q|--quiet] [-v|--verbose] [formula]
  upgrade [<formula>...]
  uninstall, remove, rm <formula>...
  list, ls [<formula>...]
  shellenv`)
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
				panic(err)
			}
		default:
			// Not a metacommand, need to continue
			quit = false
		}
	}
	return
}

func main() {
	json_path := config.JSON_PATH
	if rtn, quit := doMeta(json_path, os.Args); quit {
		os.Exit(rtn)
	}

	if _, err := os.Stat(json_path); os.IsNotExist(err) {
		if err = cmd.Update(json_path); err != nil {
			panic(err)
		}
	}
	formulas.Load(json_path)

	switch os.Args[1] {
	case "search":
		cmd.Search(formulas, os.Args[2:])
	case "info":
		cmd.Info(formulas, os.Args[2:])
	case "install":
		cmd.Install(formulas, os.Args[2:])
	case "upgrade":
		cmd.Upgrade(formulas, os.Args[2:])
	case "uninstall", "remove", "rm":
		cmd.Uninstall(formulas, os.Args[2:])
	case "list", "ls":
		cmd.List(formulas, os.Args[2:])
	default:
		fmt.Printf("Unknown subcommand '%s'\n", os.Args[1])
		os.Exit(1)
	}
}
