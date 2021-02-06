package formula

import (
	"encoding/json"
	"fmt"
	"github.com/gromgit/litebrew/internal/config"
	"github.com/gromgit/litebrew/internal/console"
	"io/ioutil"
	"os"
	"path/filepath"
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
		// Devel  interface{} `json:"devel"`
		// Head   interface{} `json:"head"`
		Bottle bool `json:"bottle"`
	} `json:"versions"`
	Revision      int `json:"revision"`
	VersionScheme int `json:"version_scheme"`
	Bottle        struct {
		Stable struct {
			Rebuild int    `json:"rebuild"`
			Cellar  string `json:"cellar"`
			Prefix  string `json:"prefix"`
			RootURL string `json:"root_url"`
			Files   struct {
				Catalina struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"catalina"`
				Mojave struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"mojave"`
				HighSierra struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"high_sierra"`
				Sierra struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"sierra"`
				ElCapitan struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"el_capitan"`
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
	Installed               bool
	/*
		LinkedKeg               string   `json:"linked_keg"`
		Pinned                  bool          `json:"pinned"`
		Outdated                bool          `json:"outdated"`
		Bottle                  struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
		Bottle struct {
		} `json:"bottle,omitempty"`
	*/
}

type Formulas map[string]Formula

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
		if i.Versions.Bottle {
			var fpath = filepath.Join(i.GetCellar(), i.Name, i.GetVersion())
			if _, err := os.Stat(fpath); err == nil {
				i.Installed = true
				instcount++
			}
			(*formulas)[i.Name] = i
		}
	}
	fmt.Fprintf(os.Stderr, "FORMULAS: Total = %d, Bottled = %d, Installed = %d\n", len(result), len(*formulas), instcount)
	return nil
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
	bold := console.Set(console.BOLD_ON)
	for _, i := range formulas {
		if i.Installed {
			flist = append(flist, console.FancyString{i.Name + " ✓", bold})
		} else {
			flist = append(flist, console.FancyString{i.Name, ""})
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
