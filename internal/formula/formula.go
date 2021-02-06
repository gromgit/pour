package formula

import (
	"encoding/json"
	"fmt"
	"github.com/gromgit/pour/internal/config"
	"github.com/gromgit/pour/internal/console"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func (allf *Formulas) Load(json_path string) error {
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
	*allf = make(Formulas)
	for _, i := range result {
		if !i.BottleDisabled && i.Installable() && i.Versions.Bottle {
			i.Status = MISSING
			// Check if installed
			optLink := filepath.Join(config.OPTDIR, i.Name)
			if stat, err := os.Stat(optLink); err == nil {
				// Some version is installed, but which one?
				instPath, err := filepath.EvalSymlinks(optLink)
				if err == nil {
					instVer := filepath.Base(instPath)
					if instVer == i.GetVersion() {
						// Latest-n-greatest
						i.Status = INSTALLED
						i.InstallDir = instPath
						i.InstallTime = stat.ModTime().Format("2006-01-02 at 15:04:05")
						i.Pinned = isPinned(i.Name, stat)
						i.Leaf = isLeaf(i.Name, stat)
						instcount++
					} else {
						i.Status = OUTDATED
						i.InstallDir = instPath
						i.InstallTime = stat.ModTime().Format("2006-01-02 at 15:04:05")
						i.Pinned = isPinned(i.Name, stat)
						i.Leaf = isLeaf(i.Name, stat)
						instcount++
					}
				}
			}
			// Look up proper URL / SHA256
			baseVal := reflect.ValueOf(i.Bottle.Stable.Files).FieldByName(config.OS_FIELD)
			i.Bottle.Stable.URL = baseVal.FieldByName("URL").String()
			i.Bottle.Stable.Sha256 = baseVal.FieldByName("Sha256").String()
			// Record it officially
			(*allf)[i.Name] = i
		}
	}
	// Now run through all bottles to populate Users list
	for b, _ := range *allf {
		for _, dep := range (*allf)[b].Dependencies {
			if (*allf)[dep].Name != "" {
				(*allf)[dep] = (*allf)[dep].AddUser(b)
			}
		}
	}
	fmt.Fprintf(os.Stderr, "FORMULAS: Total = %d, Bottled = %d, Installed = %d\n", len(result), len(*allf), instcount)
	return nil
}

func isPinned(name string, stat os.FileInfo) (result bool) {
	if _, err := os.Stat(filepath.Join(config.PINDIR, name)); err == nil {
		// Homebrew pinning: link in PREFIX/var/homebrew/pinned/
		result = true
	}
	return
}

func isLeaf(name string, stat os.FileInfo) (result bool) {
	if _, err := os.Stat(filepath.Join(config.LEAFDIR, name)); err == nil {
		// Homebrew pinning: link in PREFIX/var/homebrew/leaves/
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

func (allf Formulas) Filter(fn func(item Formula) bool) Formulas {
	result := make(Formulas)
	for _, i := range allf {
		if fn(i) {
			result[i.Name] = i
		}
	}
	return result
}

func (allf Formulas) MkStrList() (list console.FancyStrSlice) {
	for _, i := range allf {
		if i.Status == INSTALLED {
			list = append(list, console.FancyString{i.Out(), console.Bold})
		} else {
			list = append(list, console.FancyString{i.Out(), ""})
		}
	}
	sort.Sort(console.FancyStrSlice(list))
	return
}

func (allf Formulas) Ls() {
	fmt.Print(allf.MkStrList().Columnate())
}

func (formula Formula) GetCellar() string {
	result := formula.Bottle.Stable.Cellar
	if formula.Installable() {
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

func (formula Formula) Installable() bool {
	return strings.HasPrefix(formula.Bottle.Stable.Cellar, ":any") ||
		formula.Bottle.Stable.Cellar == config.CELLAR
}


// Formula reports (mostly for "info" cmd)
func (allf Formulas) GetDepStr(depList []string) (result string) {
	deps := make(Formulas)
	for _, d := range depList {
		f := allf[d]
		if f.Name != "" {
			deps[d] = allf[d]
		}
	}
	if len(deps) > 0 {
		result = deps.MkStrList().List()
	}
	return
}

func (f Formula) GetDepReport(allf Formulas) (results map[string]string) {
	results = make(map[string]string)
	for k, v := range map[string][]string{
		"Required":    f.Dependencies,
		"Recommended": f.RecommendedDependencies,
		"Optional":    f.OptionalDependencies,
	} {
		if len(v) > 0 {
			results[k] = allf.GetDepStr(v)
		}
	}
	return
}

func (f Formula) getKegReason() (result string) {
	if strings.Contains(f.Name, "@") {
		result = "this is an alternate version of another formula"
	} else {
		result = fmt.Sprintf("%s provides an older %s", config.OS_NAME, f.Name)
	}
	return
}

func (f Formula) GetCaveatReport() (results map[string]string) {
	results = make(map[string]string)
	if f.Caveats != "" {
		results["Specific"] = f.Caveats
	}
	if f.KegOnly {
		results["KegReason"] = f.getKegReason()
		// Let's see what needs to be highlighted
		baseDir := filepath.Join(f.Bottle.Stable.Prefix, "opt", f.Name)
		if _, err := os.Stat(filepath.Join(baseDir, "bin")); err == nil {
			results["Path"] = "true"
		}
		if _, err := os.Stat(filepath.Join(baseDir, "lib")); err == nil {
			results["Dev"] = "true"
		}
		if _, err := os.Stat(filepath.Join(baseDir, "lib/pkgconfig")); err == nil {
			results["Pkgconfig"] = "true"
		}
	}
	if len(results) > 0 {
		// Fill in the fixed stuff
		results["Name"] = f.Name
		results["OS"] = config.OS_NAME
		results["Prefix"] = f.Bottle.Stable.Prefix
	}
	return
}

func (f Formula) AddUser(u string) Formula {
	f.Users = append(f.Users, u)
	return f
}
