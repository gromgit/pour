package config

// Import the OS-specific stuff
import (
	oscfg "github.com/gromgit/pour/internal/config/os"
	"os"
)

const OS_NAME = oscfg.OS_NAME
const DEFAULT_PREFIX = oscfg.DEFAULT_PREFIX
const JSON_URL = oscfg.JSON_URL
const VAR_SUBPATH = "/var/pour"
const DEFAULT_VAR_PATH = DEFAULT_PREFIX + VAR_SUBPATH

var PREFIX = DEFAULT_PREFIX
var OS_FIELD = oscfg.GetOS()
var OS_DEPS = oscfg.GetDeps()
var PatchExec = oscfg.PatchExec

var CELLAR string
var OPTDIR string
var VAR_PATH string
var JSON_PATH string
var PINDIR string
var LINKDIR string
var LEAFDIR string
var BOTTLEDIR string
var TEMPDIR string
var SYSDIRS []string

func Cellar() string {
	return PREFIX + "/Cellar"
}

func OptDir() string {
	return PREFIX + "/opt"
}

func VarDir() string {
	return PREFIX + VAR_SUBPATH
}

func PinDir() string {
	return VarDir() + "/pinned"
}

func LinkDir() string {
	return VarDir() + "/linked"
}

func LeafDir() string {
	return VarDir() + "/leaves"
}

func BottleDir() string {
	return VarDir() + "/bottles"
}

func TempDir() string {
	return VarDir() + "/tmp"
}

func Json() string {
	return VarDir() + "/bottles.json"
}

func SysDirs() []string {
	return []string{CELLAR, PINDIR, LINKDIR, LEAFDIR, BOTTLEDIR, TEMPDIR}
}

func init() {
	// Check for POUR_PREFIX env var
	if prefix := os.Getenv("POUR_PREFIX"); prefix != "" {
		PREFIX = prefix
	}
	CELLAR = Cellar()
	OPTDIR = OptDir()
	VAR_PATH = VarDir()
	JSON_PATH = Json()
	PINDIR = PinDir()
	LINKDIR = LinkDir()
	LEAFDIR = LeafDir()
	BOTTLEDIR = BottleDir()
	TEMPDIR = TempDir()
	SYSDIRS = SysDirs()
}
