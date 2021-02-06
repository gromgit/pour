// +build darwin,amd64

package config

import (
	"os/exec"
	s "strings"
)

const OS_NAME = "macOS"
const JSON_URL = "https://formulae.brew.sh/api/formula.json"
const DEFAULT_PREFIX = "/usr/local"

func GetDeps() []string {
	return []string{}
}

func GetOS() (os string) {
	verCmd := exec.Command("defaults", "read", "loginwindow", "SystemVersionStampAsString")
	if verOut, err := verCmd.Output(); err != nil {
		panic(err)
	} else {
		verStr := string(verOut)
		switch {
		case s.HasPrefix(verStr, "10.15"):
			os = "Catalina"
		case s.HasPrefix(verStr, "10.14"):
			os = "Mojave"
		case s.HasPrefix(verStr, "10.13"):
			os = "HighSierra"
		case s.HasPrefix(verStr, "10.12"):
			os = "Sierra"
		case s.HasPrefix(verStr, "10.11"):
			os = "ElCapitan"
		case s.HasPrefix(verStr, "10.10"):
			os = "Yosemite"
		case s.HasPrefix(verStr, "10.9"):
			os = "Mavericks"
		default:
			panic("Unrecognized macOS version")
		}
	}
	return
}
