package dns

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatingImportingScript(t *testing.T) {
	zoneConfigMap := map[string]Types{"a": {"b", "c", "d"}, "e": {"f", "g", "h"}}
	importScript, err := buildZoneImportScript("some-zone", zoneConfigMap, "resource_name")
	require.NoError(t, err)
	assertFileWithContent(t, "./testdata/import_script/import.sh", importScript)
}
