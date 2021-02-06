package net

import (
	"crypto/sha256"
	"fmt"
	cfg "github.com/gromgit/pour/internal/config"
	"io"
	"net/http"
	"os"
)

// Ref: https://golangcode.com/download-a-file-from-a-url/
func DownloadFile(filepath string, url string) error {

	cfg.Log("Downloading", url)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func ChecksumFile(filepath string, sha string) error {

	cfg.Log("Verifying", filepath)

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	fsha := fmt.Sprintf("%x", h.Sum(nil))
	if sha != fsha {
		return fmt.Errorf("SHA256 failed: expected %+v, got %+v", sha, fsha)
	}
	return nil

}

func DownloadAndVerify(filepath, url, sha string) error {

	if _, err := os.Stat(filepath); err != nil {
		// Download first
		if err := DownloadFile(filepath, url); err != nil {
			return fmt.Errorf("Download of %s failed: %+v", filepath, err)
		}
	}
	// Now do the checksum
	if err := ChecksumFile(filepath, sha); err != nil {
		return fmt.Errorf("Verification of %s failed: %+v", filepath, err)
	}
	return nil
}
