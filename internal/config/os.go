package config

// Import the OS-specific stuff
import oscfg "github.com/gromgit/litebrew/internal/config/os"

const DEFAULT_PREFIX = oscfg.DEFAULT_PREFIX
const JSON_URL = oscfg.JSON_URL
const JSON_PATH = DEFAULT_PREFIX + "/var/homebrew/bottles.json"
const CELLAR = DEFAULT_PREFIX + "/Cellar"
const PINDIR = DEFAULT_PREFIX + "/var/homebrew/pinned"
const LINKDIR = DEFAULT_PREFIX + "/var/homebrew/linked"

var OS_FIELD = oscfg.GetOS()
