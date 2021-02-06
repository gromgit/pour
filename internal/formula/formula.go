package formula

import (
	"encoding/json"
	"fmt"
	"github.com/gromgit/litebrew/internal/config"
	"github.com/gromgit/litebrew/internal/console"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Formula struct {
	Name              string   `json:"name"`
	FullName          string   `json:"full_name"`
	Oldname           string   `json:"oldname"`
	Aliases           []string `json:"aliases"`
	VersionedFormulae []string `json:"versioned_formulae"`
	Desc              string   `json:"desc"`
	Homepage          string   `json:"homepage"`
	Versions          struct {
		Stable string `json:"stable"`
		Bottle bool   `json:"bottle"`
	} `json:"versions"`
	Revision      int `json:"revision"`
	VersionScheme int `json:"version_scheme"`
	Bottle        struct {
		Stable struct {
			Rebuild int    `json:"rebuild"`
			Cellar  string `json:"cellar"`
			Prefix  string `json:"prefix"`
			RootURL string `json:"root_url"`
			URL     string `json:"-"`
			Sha256  string `json:"-"`
			Files   struct {
				Catalina struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"catalina,omitempty"`
				Mojave struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"mojave,omitempty"`
				HighSierra struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"high_sierra,omitempty"`
				Sierra struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"sierra,omitempty"`
				ElCapitan struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"el_capitan,omitempty"`
				Yosemite struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"yosemite,omitempty"`
				Mavericks struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"mavericks,omitempty"`
				Linux64 struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"x86_64_linux,omitempty"`
			} `json:"files"`
		} `json:"stable"`
	} `json:"bottle,omitempty"`
	KegOnly                 bool     `json:"keg_only"`
	BottleDisabled          bool     `json:"bottle_disabled"`
	Options                 []string `json:"options"`
	BuildDependencies       []string `json:"build_dependencies"`
	Dependencies            []string `json:"dependencies"`
	RecommendedDependencies []string `json:"recommended_dependencies"`
	OptionalDependencies    []string `json:"optional_dependencies"`
	UsesFromMacos           []string `json:"uses_from_macos"`
	Requirements            []string `json:"requirements"`
	ConflictsWith           []string `json:"conflicts_with"`
	Caveats                 string   `json:"caveats"`
	Status                  int      `json:"-"`
	InstallDir              string   `json:"-"`
	Pinned                  bool     `json:"-"`
}

type Formulas map[string]Formula

var StatusMap = map[int]string{
	INSTALLED: " ✓",
	OUTDATED:  " ✗",
}

func (formulas *Formulas) Load(json_path string) error {
	// Parse the formulas JSON
	var result []Formula
	if f, err := os.Open(json_path); err != nil {
		return err
	} else if fbytes, err := ioutil.ReadAll(f); err != nil {
		return err
	} else {
		json.Unmarshal(fbytes, &result)
	}
	// Post-process the results
	var instcount = 0
	*formulas = make(Formulas)
	for _, i := range result {
		if !i.BottleDisabled && i.Versions.Bottle {
			// Check if installed
			i.InstallDir = filepath.Join(i.GetCellar(), i.Name, i.GetVersion())
			if stat, err := os.Stat(i.InstallDir); err == nil {
				i.Status = INSTALLED
				i.Pinned = isPinned(i.Name, stat)
				instcount++
			} else if stat, err := os.Stat(filepath.Dir(i.InstallDir)); err == nil {
				i.Status = OUTDATED
				i.Pinned = isPinned(i.Name, stat)
				instcount++
			} else {
				i.Status = MISSING
			}
			// Look up proper URL / SHA256
			baseVal := reflect.ValueOf(i.Bottle.Stable.Files).FieldByName(config.OS_FIELD)
			i.Bottle.Stable.URL = baseVal.FieldByName("URL").String()
			i.Bottle.Stable.Sha256 = baseVal.FieldByName("Sha256").String()
			// Record it officially
			(*formulas)[i.Name] = i
		}
	}
	fmt.Fprintf(os.Stderr, "FORMULAS: Total = %d, Bottled = %d, Installed = %d\n", len(result), len(*formulas), instcount)
	return nil
}

func isPinned(name string, stat os.FileInfo) (result bool) {
	if stat.Mode()&os.ModeSticky > 0 {
		// Litebrew pinning: sticky bit on install dir
		result = true
	} else if _, err := os.Stat(config.DEFAULT_PREFIX + "/var/homebrew/pinned/" + name); err == nil {
		// Homebrew pinning: link in PREFIX/var/homebrew/pinned/
		result = true
	}
	return
}

func (formula Formula) Out() (out string) {
	out = formula.Name
	if config.Fancy {
		out = out + StatusMap[formula.Status]
	}
	return
}

func (formulas Formulas) Filter(fn func(item Formula) bool) Formulas {
	result := make(Formulas)
	for _, i := range formulas {
		if fn(i) {
			result[i.Name] = i
		}
	}
	return result
}

func (formulas Formulas) Ls() {
	var flist []console.FancyString
	for _, i := range formulas {
		if i.Status == INSTALLED {
			flist = append(flist, console.FancyString{i.Out(), console.Bold})
		} else {
			flist = append(flist, console.FancyString{i.Out(), ""})
		}
	}
	sort.Sort(console.FancyStrSlice(flist))
	fmt.Print(console.Columnate(flist))
}

func (formula Formula) GetCellar() string {
	result := formula.Bottle.Stable.Cellar
	if strings.HasPrefix(result, ":any") {
		result = config.DEFAULT_CELLAR
	}
	return result
}

func (formula Formula) GetVersion() string {
	result := formula.Versions.Stable
	if formula.Revision > 0 {
		result = result + "_" + strconv.Itoa(formula.Revision)
	}
	return result
}

// Various formula-related enumerations
const (
	RUN = iota
	BUILD
	RECOMMENDED
	OPTIONAL
	INSTALLED
	OUTDATED
	MISSING
)
