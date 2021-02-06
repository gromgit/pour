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
	InstallTime             string   `json:"-"`
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
			i.Status = MISSING
			// Check if installed
			targetParent := filepath.Join(i.GetCellar(), i.Name)
			targetDir := filepath.Join(targetParent, i.GetVersion())
			if stat, err := os.Stat(targetDir); err == nil {
				i.Status = INSTALLED
				i.InstallDir = targetDir
				i.InstallTime = stat.ModTime().Format("2006-01-02 at 15:04:05")
				i.Pinned = isPinned(i.Name, stat)
				instcount++
			} else if _, err := os.Stat(targetParent); err == nil {
				// Let's find the latest bottle installed here
				if bottles, err := filepath.Glob(filepath.Join(targetParent, "*")); err == nil && len(bottles) > 0 {
					// TODO: Find a better choice than the first one
					if stat, err := os.Stat(bottles[0]); err == nil {
						i.Status = OUTDATED
						i.InstallDir = bottles[0]
						i.InstallTime = stat.ModTime().Format("2006-01-02 at 15:04:05")
						i.Pinned = isPinned(i.Name, stat)
						instcount++
					}
				}
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
	if _, err := os.Stat(filepath.Join(config.PINDIR, name)); err == nil {
		// Homebrew pinning: link in PREFIX/var/homebrew/pinned/
		result = true
	}
	return
}

func (formula Formula) Pin() (e error) {
	if formula.Status == MISSING {
		fmt.Fprintf(os.Stderr, "Bottle '%s' not installed, cannot pin\n", formula.Name)
	} else if formula.Pinned {
		fmt.Fprintf(os.Stderr, "Bottle '%s' already pinned\n", formula.Name)
	} else {
		// Link the current version
		srcpath := formula.InstallDir
		destpath := filepath.Join(config.PINDIR, formula.Name)
		if err := os.Symlink(srcpath, destpath); err != nil {
			e = err
		}
	}
	return
}

func (formula Formula) Unpin() (e error) {
	if formula.Status == MISSING {
		fmt.Fprintf(os.Stderr, "Bottle '%s' not installed, cannot unpin\n", formula.Name)
	} else if !formula.Pinned {
		fmt.Fprintf(os.Stderr, "Bottle '%s' not pinned\n", formula.Name)
	} else {
		// Remove the existing link
		if err := os.Remove(filepath.Join(config.PINDIR, formula.Name)); err != nil {
			e = err
		}
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

func (formulas Formulas) MkStrList() (list console.FancyStrSlice) {
	for _, i := range formulas {
		if i.Status == INSTALLED {
			list = append(list, console.FancyString{i.Out(), console.Bold})
		} else {
			list = append(list, console.FancyString{i.Out(), ""})
		}
	}
	sort.Sort(console.FancyStrSlice(list))
	return
}

func (formulas Formulas) Ls() {
	fmt.Print(formulas.MkStrList().Columnate())
}

func (formula Formula) GetCellar() string {
	result := formula.Bottle.Stable.Cellar
	if strings.HasPrefix(result, ":any") {
		result = config.CELLAR
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
