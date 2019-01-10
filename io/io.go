package io

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// CreateDirectoryIfNotExist creates a directory or path of directories
// if it doesn't exist
func CreateDirectoryIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// DownloadFile downloads a file from the provided url and will save it under the
// provided name in the destination directory
func DownloadFile(destDir, filename, url string) error {
	// Make sure that the directory exists unless the current directory shortcut has been
	// given as an argument
	if destDir != "." && destDir != "./" {
		if err := CreateDirectoryIfNotExist(destDir); err != nil {
			return fmt.Errorf("couldn't create directory: %v", err)
		}
	}

	// Create the destination file
	fullPath := path.Join(destDir, filename)
	dest, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("couldn't create destination file: %v", err)
	}
	defer dest.Close()

	// Download the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download error: %v", err)
	}
	defer resp.Body.Close()

	// Write the body into the file
	_, err = io.Copy(dest, resp.Body)
	if err != nil {
		return fmt.Errorf("couldn't write content: %v", err)
	}

	return nil
}

// Unzip extracts the files of a .zip archive to the provided destination directory
func Unzip(destDir, archive string) error {
	// Make sure that the directory exists unless the current directory shortcut has been
	// given as an argument
	if destDir != "." && destDir != "./" {
		if err := CreateDirectoryIfNotExist(destDir); err != nil {
			return fmt.Errorf("couldn't create directory: %v", err)
		}
	}

	// Open the archive
	z, err := zip.OpenReader(archive)
	if err != nil {
		return fmt.Errorf("couldn't open zip file: %v", err)
	}
	defer z.Close()

	// Loop over the file in the archive
	for _, file := range z.File {
		// Open the file
		f, err := file.Open()
		if err != nil {
			return fmt.Errorf("error extracting %q: %v", file.Name, err)
		}

		// Create the destination file
		fullPath := path.Join(destDir, file.Name)
		dest, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("couldn't create destination file: %v", err)
		}
		defer dest.Close()

		// Write the extracted content into the file
		_, err = io.Copy(dest, f)
		if err != nil {
			return fmt.Errorf("couldn't write content: %v", err)
		}
	}
	return nil
}
