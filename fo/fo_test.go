package fo

import (
	"os"
	"testing"
)

func TestCreateDirectoryIfNotExist(t *testing.T) {
	path := "./this/is/a/test"
	defer os.RemoveAll("./this")

	if err := CreateDirectoryIfNotExist(path); err != nil {
		t.Fatalf("expected to pass, but got: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected path to exist, but got %v", err)
	}
}

// TODO: Find a smaller file to test with
func TestUnzip(t *testing.T) {
	url := "https://www.sample-videos.com/zip/10mb.zip"

	if err := DownloadFile("/tmp", url); err != nil {
		t.Fatalf("expected to pass, but got %v", err)
	}
	defer os.Remove("/tmp/test/zip")

	if _, err := os.Stat("/tmp/test.zip"); err != nil {
		t.Errorf("expected file to exist, but got %v", err)
	}

	if err := Unzip("/tmp/archive", "/tmp/test.zip"); err != nil {
		t.Errorf("expected to pass, but got %v", err)
	}
	defer os.RemoveAll("/tmp/archive")

	if _, err := os.Stat("/tmp/archive"); err != nil {
		t.Errorf("expected directory to exist, but got %v", err)
	}
}
