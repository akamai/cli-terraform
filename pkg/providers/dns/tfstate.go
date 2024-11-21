package dns

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type tfStateStruct struct {
	Version          int         `json:"version"`
	TerraformVersion string      `json:"terraform_version"`
	Serial           int         `json:"serial"`
	Lineage          string      `json:"lineage"`
	Outputs          interface{} `json:"outputs"`
	Resources        []*Resource `json:"resources"`
}

// Resource describes tf resource
type Resource struct {
	Mode      string        `json:"mode"`
	Type      string        `json:"type"`
	Name      string        `json:"name"`
	Provider  string        `json:"provider"`
	Instances []interface{} `json:"instances"`
}

var tfState *tfStateStruct

// Utility method to read in tf state content
func readTfState(tfWorkPath string) error {
	// TFWorkPath global var
	tfStateFileName := filepath.Join(tfWorkPath, "terraform.tfstate")
	if _, err := os.Stat(tfStateFileName); err != nil {
		return err
	}
	stateData, err := ioutil.ReadFile(tfStateFileName)
	if err != nil {
		return err
	}
	tfState = &tfStateStruct{}
	err = json.Unmarshal(stateData, tfState)
	if err != nil {
		return err
	}

	return nil
}
