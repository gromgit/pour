package cmd

import (
	"errors"
	"fmt"
	cfg "github.com/gromgit/pour/internal/config"
	"github.com/gromgit/pour/internal/net"
	"os"
	"path/filepath"
)

func Update(path string) error {
	dirpath := filepath.Dir(path)
	if err := os.MkdirAll(dirpath, 0755); err != nil {
		return errors.New("cannot create directory " + dirpath)
	}
	fmt.Printf("Updating %s\n", path)
	if err := net.DownloadFile(path, cfg.JSON_URL); err != nil {
		return errors.New("cannot download " + cfg.JSON_URL + ": " + err.Error())
	}
	return nil
}
