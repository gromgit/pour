package bottle

import (
	"bufio"
	cfg "github.com/gromgit/pour/internal/config"
	"github.com/gromgit/pour/internal/formula"
	"github.com/gromgit/pour/internal/net"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Unpack(tarPath, destPath string) error {
	if f, err := os.Open(tarPath); err != nil {
		return err
	} else {
		r := bufio.NewReader(f)
		if err := Untar(r, destPath); err != nil {
			return err
		}
	}
	return nil
}

func getFilelist(list *[]string) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			switch info.Name() {
			case ".brew", ".bottle", "libexec":
				log.Println("Skipping", path)
				return filepath.SkipDir
			}
		} else if filepath.Dir(path) != "." {
			// Only files at least one level deep are linkable
			*list = append(*list, path)
		}
		return nil
	}
}

func getLinkables(pourRoot string) (list []string, err error) {
	oldwd, err := os.Getwd()
	if err != nil {
		return
	}
	defer os.Chdir(oldwd)
	err = os.Chdir(pourRoot)
	if err != nil {
		return
	}
	err = filepath.Walk(".", getFilelist(&list))
	return
}

// Unlink("<name>/<version>")
func Unlink(pkgSubdir string) error {
	if list, err := getLinkables(filepath.Join(cfg.CELLAR, pkgSubdir)); err != nil {
		return err
	} else {
		cfg.Log("Unlink paths:", list)
		for _, p := range list {
			pf := filepath.Join(cfg.CELLAR, p)
			if err := os.Remove(pf); err != nil {
				cfg.Log("ERROR on unlink", pf, err)
			}
		}
	}
	os.Remove(filepath.Join(cfg.LINKDIR, filepath.Dir(pkgSubdir)))
	return nil
}

// Link("<name>/<version>")
func Link(pkgSubdir string) error {
	pkgDir := filepath.Join(cfg.CELLAR, pkgSubdir)
	if list, err := getLinkables(pkgDir); err != nil {
		return err
	} else {
		cfg.Log("Link paths:", list)
		for _, p := range list {
			dest := filepath.Join(cfg.PREFIX, p)
			src := filepath.Join(pkgDir, p)
			rel, err := filepath.Rel(filepath.Dir(dest), src)
			if err != nil {
				return err
			}
			os.Remove(dest)
			if err := os.MkdirAll(filepath.Dir(dest), 0775); err != nil {
				return err
			}
			if err := os.Symlink(rel, dest); err != nil {
				return err
			}
		}
		// Add link to LINKDIR
		os.Mkdir(cfg.LINKDIR, 0775)
		linkRel, err := filepath.Rel(cfg.LINKDIR, pkgDir)
		if err != nil {
			return err
		}
		link := filepath.Join(cfg.LINKDIR, filepath.Dir(pkgSubdir))
		if err := os.Symlink(linkRel, link); err != nil {
			return err
		}
		cfg.Log("Link the link:", linkRel, link)
	}
	return nil
}

func Install(f formula.Formula) error {
	log.SetOutput(os.Stdout)
	url := f.Bottle.Stable.URL
	tarName := filepath.Base(url)
	tarPath := filepath.Join(cfg.BOTTLEDIR, tarName)
	if _, err := os.Stat(tarPath); err != nil {
		// Download it first
		cfg.Log("Downloading", url)
		if err := net.DownloadFile(tarPath, url); err != nil {
			return err
		}
	}
	// Unpack into temp dir
	if tempDir, err := ioutil.TempDir(cfg.TEMPDIR, "inst"); err != nil {
		return err
	} else {
		defer os.RemoveAll(tempDir)
		if err := Unpack(tarPath, tempDir); err != nil {
			return err
		}
		cfg.Log("Unpacked to", tempDir)
		// Make sure we have the right dir
		tempPkgdir := filepath.Join(tempDir, f.Name, f.GetVersion())
		if _, err := os.Stat(tempPkgdir); err != nil {
			return err
		}
		// Move hierarchy over to Cellar
		pkgSubdir := filepath.Join(f.Name, f.GetVersion())
		finalPkgdir := filepath.Join(cfg.CELLAR, pkgSubdir)
		if err := os.MkdirAll(filepath.Dir(finalPkgdir), 0775); err != nil {
			return err
		}
		if err := os.Rename(tempPkgdir, finalPkgdir); err != nil {
			return err
		}
		// Unlink/remove old version if present
		if f.InstallDir != "" {
			oldPkgSubdir, err := filepath.Rel(cfg.CELLAR, f.InstallDir)
			if err != nil {
				return err
			}
			if err := Unlink(oldPkgSubdir); err != nil {
				return err
			} else if err := os.RemoveAll(f.InstallDir); err != nil {
				return err
			} else {
				f.InstallDir = ""
				f.Status = formula.MISSING
			}
		}
		// Link new version in
		if err := Link(pkgSubdir); err != nil {
			return err
		}
	}
	return nil
}
