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
		// TODO: Move it into place
		cfg.Log("Unpacked to", tempDir)
	}
	return nil
}
