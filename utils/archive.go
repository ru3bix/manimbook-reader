package utils

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"archive/zip"
	"github.com/spf13/afero"
	"os"
)

func Unzip(src string, fs afero.Fs, destPath string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("cannot open zip file: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		// Construct the full destination path.
		fpath := filepath.Join(destPath, f.Name)

		// Prevent Zip Slip (directory traversal attack).
		if !strings.HasPrefix(filepath.Clean(fpath), filepath.Clean(destPath)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := fs.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}
		if err := fs.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		outFile, err := fs.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			inFile.Close()
			return err
		}
		if _, err := io.Copy(outFile, inFile); err != nil {
			inFile.Close()
			outFile.Close()
			return err
		}

		inFile.Close()
		outFile.Close()
	}
	return nil
}
