package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

// TFWorkPath is a target directory for generated terraform resources
var TFWorkPath = "./"

// CheckFiles verifies if all given files exist in filesystem
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
func CreateTFFilename(resourceName string) string {
	return filepath.Join(TFWorkPath, resourceName+".tf")
}
