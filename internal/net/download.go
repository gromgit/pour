package net

import (
	"crypto/sha256"
	"fmt"
	"github.com/gromgit/pour/internal/log"
	"io"
	"net/http"
	"os"
	"strings"
)

var ppfmt = "\r" + strings.Repeat(" ", 35) + "\rDownloading... %5.1f%% complete"

// Ref: https://golangcode.com/download-a-file-with-progress/
type PercentProgress struct {
	Expected, Total uint64
}

func (pp *PercentProgress) Write(p []byte) (int, error) {
	n := len(p)
	pp.Total += uint64(n)
	pp.PrintProgress()
	return n, nil
}

func (pp PercentProgress) PrintProgress() {
	pct := float64(pp.Total) * 100.0 / float64(pp.Expected)
	fmt.Fprintf(os.Stderr, ppfmt, pct)
}

func DownloadFile(filepath string, url string) error {

	log.Log("Downloading", url)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Create progress reporter and write the body to file
	pp := &PercentProgress{uint64(resp.ContentLength), 0}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, pp)); err != nil {
		return err
	}
	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Finalize the download
	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

func ChecksumFile(filepath string, sha string) error {

	log.Log("Verifying", filepath)

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
