package formula

import (
	"encoding/json"
	"fmt"
	"github.com/gromgit/litebrew/internal/config"
	"github.com/gromgit/litebrew/internal/console"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Formula struct {
	Name              string        `json:"name"`
	FullName          string        `json:"full_name"`
	Oldname           interface{}   `json:"oldname"`
	Aliases           []interface{} `json:"aliases"`
	VersionedFormulae []interface{} `json:"versioned_formulae"`
	Desc              string        `json:"desc"`
	Homepage          string        `json:"homepage"`
	Versions          struct {
		Stable string      `json:"stable"`
		Devel  interface{} `json:"devel"`
		Head   interface{} `json:"head"`
		Bottle bool        `json:"bottle"`
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
	KegOnly                 bool          `json:"keg_only"`
	BottleDisabled          bool          `json:"bottle_disabled"`
	Options                 []interface{} `json:"options"`
	BuildDependencies       []interface{} `json:"build_dependencies"`
	Dependencies            []interface{} `json:"dependencies"`
	RecommendedDependencies []interface{} `json:"recommended_dependencies"`
	OptionalDependencies    []interface{} `json:"optional_dependencies"`
	UsesFromMacos           []interface{} `json:"uses_from_macos"`
	Requirements            []interface{} `json:"requirements"`
	ConflictsWith           []interface{} `json:"conflicts_with"`
	Caveats                 interface{}   `json:"caveats"`
	Installed               bool
	/*
		LinkedKeg               interface{}   `json:"linked_keg"`
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

type Formulas []Formula

func (formulas *Formulas) Load(json_path string) error {
	// Parse the formulas JSON
	var result Formulas
	if f, err := os.Open(json_path); err != nil {
		return err
	} else if fbytes, err := ioutil.ReadAll(f); err != nil {
		return err
	} else {
		json.Unmarshal(fbytes, &result)
	}
	// Post-process the results
	var instcount = 0
	for _, i := range result {
		if i.Versions.Bottle {
			var fpath = filepath.Join(i.GetCellar(), i.Name, i.GetVersion())
			if _, err := os.Stat(fpath); err == nil {
				i.Installed = true
				instcount++
			}
			*formulas = append(*formulas, i)
		}
	}
	fmt.Fprintf(os.Stderr, "FORMULAS: Total = %d, Bottled = %d, Installed = %d\n", len(result), len(*formulas), instcount)
	return nil
}

func (formulas Formulas) Filter(fn func(item Formula) bool) Formulas {
	result := make(Formulas, 0)
	for _, i := range formulas {
		if fn(i) {
			result = append(result, i)
		}
	}
	return result
}

func (formulas Formulas) Ls() {
	installed_fmt := console.Set(console.BOLD_ON) + "%s ✔︎\n" + console.Reset()
	for _, i := range formulas {
		if i.Installed {
			fmt.Printf(installed_fmt, i.Name)
		} else {
			fmt.Println(i.Name)
		}
	}
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
