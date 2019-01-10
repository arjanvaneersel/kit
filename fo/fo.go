package fo

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// SplitPathAndFile takes a full path (windows and unix) or an URL splits the path and file
// The last element is considered to be the file and is returned as such
// When an URL is provided the prodocol indicator will be removed and the domain name will become a part of the path
// So http://example.com/files/file.zip will return example.com/files as the path
func SplitPathAndFile(p string) (string, string, error) {
	sep := "/"
	if strings.Contains(p, "://") {
		parts := strings.Split(p, "://")
		p = strings.Join(parts[1:], "/")
	} else if strings.Contains(p, "\\") {
		sep = "\\"
	} else if !strings.Contains(p, "/") {
		return "", "", fmt.Errorf("invalid path")
	}

	parts := strings.Split(p, sep)
	l := len(parts)
	if l <= 1 {
		return "", "", fmt.Errorf("couldn't split %v", p)
	}

	dir := path.Join(append([]string{sep}, parts[:l-2]...)...)
	filename := parts[l-1]

	return dir, filename, nil
}

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
func DownloadFile(destDir, url string) (string, error) {
	// Make sure that the directory exists unless the current directory shortcut has been
	// given as an argument
	if destDir != "." && destDir != "./" {
		if err := CreateDirectoryIfNotExist(destDir); err != nil {
			return "", fmt.Errorf("couldn't create directory: %v", err)
		}
	}

	// Get the filename from the URL
	_, filename, err := SplitPathAndFile(url)
	if err != nil {
		return "", fmt.Errorf("couldn't get filename: %v", err)
	}

	// Create the destination file
	fullPath := path.Join(destDir, filename)
	dest, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("couldn't create destination file: %v", err)
	}
	defer dest.Close()

	// Download the data
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("download error: %v", err)
	}
	defer resp.Body.Close()

	// Write the body into the file
	_, err = io.Copy(dest, resp.Body)
	if err != nil {
		return "", fmt.Errorf("couldn't write content: %v", err)
	}

	return fullPath, nil
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
