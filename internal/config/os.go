package config

// Import the OS-specific stuff
import oscfg "github.com/gromgit/litebrew/internal/config/os"

const DEFAULT_PREFIX = oscfg.DEFAULT_PREFIX
const JSON_URL = oscfg.JSON_URL
const CELLAR = DEFAULT_PREFIX + "/Cellar"
const VAR_PATH = DEFAULT_PREFIX + "/var/litebrew"
const JSON_PATH = VAR_PATH + "/bottles.json"
const PINDIR = VAR_PATH + "/pinned"
const LINKDIR = VAR_PATH + "/linked"

var OS_FIELD = oscfg.GetOS()
