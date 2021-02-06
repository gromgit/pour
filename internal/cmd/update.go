package cmd

import (
	"errors"
	"fmt"
	"github.com/gromgit/litebrew/internal/config"
	"github.com/gromgit/litebrew/internal/net"
	"os"
	"path/filepath"
)

func Update(path string) error {
	dirpath := filepath.Dir(path)
	if err := os.MkdirAll(dirpath, 0755); err != nil {
		return errors.New("cannot create directory " + dirpath)
	}
	fmt.Printf("Updating %s\n", path)
	if err := net.DownloadFile(path, config.JSON_URL); err != nil {
		return errors.New("cannot download " + config.JSON_URL)
	}
	return nil
}
