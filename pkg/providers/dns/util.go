package dns

import (
	"path/filepath"
)

// Utility method to create full import list file path
func createImportListFilename(resourceName, tfWorkPath string) string {

	return filepath.Join(tfWorkPath, resourceName+"_resources.json")

}
