package dns

import (
	"path/filepath"
)

// Utility method to create full import list file path
func createImportListFilename(resourceName string) string {

	return filepath.Join(tfWorkPath, resourceName+"_resources.json")

}
