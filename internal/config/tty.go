package config

import (
	"github.com/mattn/go-isatty"
	"os"
	"strconv"
)

var Fancy = isatty.IsTerminal(os.Stdout.Fd())
var ScreenWidth int

func init() {
	var err error
	// Find the screen width
	ScreenWidth, err = strconv.Atoi(os.Getenv("COLUMNS"))
	if err != nil || ScreenWidth == 0 {
		ScreenWidth = 80 // A decent default
	}
}
