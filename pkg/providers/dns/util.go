package dns

import (
	"path/filepath"

	"github.com/akamai/cli-terraform/pkg/tools"
)

// Utility method to create full import list file path
func createImportListFilename(resourceName string) string {

	return filepath.Join(tools.TFWorkPath, resourceName+"_resources.json")

}
