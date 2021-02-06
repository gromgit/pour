// +build linux,amd64

package config

import (
	"github.com/gromgit/pour/internal/file"
	"github.com/gromgit/pour/internal/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const OS_NAME = "Linux"
const JSON_URL = "https://formulae.brew.sh/api/formula-linux.json"
const DEFAULT_PREFIX = "/home/linuxbrew/.linuxbrew"

var LOADER string

func GetDeps() []string {
	return []string{"patchelf"}
}

func GetOS() (os string) {
	os = "Linux64"
	return
}

func patchIt(path, prefix string) error {
	// Get rpath and interpreter
	cmd := exec.Command("patchelf", "--print-rpath", path)
	out, _ := cmd.Output()
	rpath := strings.TrimSpace(string(out))
	cmd = exec.Command("patchelf", "--print-interpreter", path)
	out, _ = cmd.Output()
	interp := strings.TrimSpace(string(out))
	patchArgs := []string{}
	rpathNew := strings.ReplaceAll(rpath, "@@HOMEBREW_PREFIX@@", prefix)
	if rpath != rpathNew {
		patchArgs = append(patchArgs, "--force-rpath", "--set-rpath", rpathNew)
	}
	interpNew := strings.ReplaceAll(interp, "@@HOMEBREW_PREFIX@@", prefix)
	if strings.HasSuffix(interpNew, "/lib/ld.so") {
		// Replace with real suffix
		interpNew = LOADER
	} else {
	}
	if interp != interpNew {
		patchArgs = append(patchArgs, "--set-interpreter", interpNew)
	}
	if len(patchArgs) > 0 {
		patchArgs = append(patchArgs, path)
		log.Log("Running patchelf", patchArgs)
		// First save existing perms
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}
		oldModes := stat.Mode()
		log.Logf("Chmod: %o -> %o\n", oldModes, oldModes|0222)
		// We want to be writable
		if err := os.Chmod(path, oldModes|0222); err != nil {
			return err
		}
		// Now patch the executable
		cmd = exec.Command("patchelf", patchArgs...)
		_, err = cmd.Output()
		if err != nil {
			return err
		}
		// And return file metadata to its original state
		if err := os.Chmod(path, oldModes); err != nil {
			return err
		}
		if err := os.Chtimes(path, stat.ModTime(), stat.ModTime()); err != nil {
			return err
		}
	}
	return nil
}

func doPatch(prefix string) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
			ft, err := file.GetTypeFromPath(path)
			if err != nil {
				return err
			}
			if ft == "application/x-executable" {
				if err = patchIt(path, prefix); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func PatchExec(base, prefix string) (err error) {
	if !strings.Contains(base, "/patchelf/") {
		err = filepath.Walk(base, doPatch(prefix))
	}
	return
}

func init() {
	// Find actual path of loader
	if loaders, err := filepath.Glob("/lib64/ld-linux*.so*"); err == nil {
		// Just pick the first match
		LOADER = loaders[0]
	}
}
