package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckFiles verifies if all given files doesn't exist in filesystem
func CheckFiles(files ...string) error {
	for _, file := range files {
		_, err := os.Stat(file)
		if err == nil {
			return fmt.Errorf("file %s already exists", file)
		}
	}
	return nil
}

// CreateTFFilename creates full tf file path
func CreateTFFilename(resourceName string, tfWorkPath string) string {
	return filepath.Join(tfWorkPath, resourceName+".tf")
}
