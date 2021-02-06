package file

import (
	"github.com/h2non/filetype"
	"os"
)

func GetTypeFromPath(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	ft, err := filetype.MatchFile(path)
	if err != nil {
		return "", err
	} else {
		return ft.MIME.Value, nil
	}
}
