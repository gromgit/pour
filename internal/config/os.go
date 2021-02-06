package config

// Import the OS-specific stuff
import (
	oscfg "github.com/gromgit/litebrew/internal/config/os"
	"log"
	"os"
)

const OS_NAME = oscfg.OS_NAME
const DEFAULT_PREFIX = oscfg.DEFAULT_PREFIX
const JSON_URL = oscfg.JSON_URL
const CELLAR = DEFAULT_PREFIX + "/Cellar"
const VAR_PATH = DEFAULT_PREFIX + "/var/litebrew"
const JSON_PATH = VAR_PATH + "/bottles.json"
const PINDIR = VAR_PATH + "/pinned"
const LINKDIR = VAR_PATH + "/linked"
const BOTTLEDIR = VAR_PATH + "/bottles"
const TEMPDIR = VAR_PATH + "/tmp"

var SYSDIRS = []string{CELLAR, PINDIR, LINKDIR, BOTTLEDIR, TEMPDIR}
var OS_FIELD = oscfg.GetOS()

// Private logger instance
var logger = log.New(os.Stderr, "litebrew", log.LstdFlags)

func Log(v ...interface{}) {
	logger.Println(v)
}
