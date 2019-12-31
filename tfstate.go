// Copyright 2020. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type tfStateStruct struct {
	Version           int         `json:"version"`
	Terraform_version string      `json:"terraform_version"`
	Serial            int         `json:"serial"`
	Lineage           string      `json:"lineage"`
	Outputs           interface{} `json:"outputs,omitempty"`
	Resources         []*Resource `json:"resources"`
}

type Resource struct {
	Mode      string        `json:"mode"`
	Type      string        `json:"type"`
	Name      string        `json:"name"`
	Provider  string        `json:"provider"`
	Instances []interface{} `json:"instances,omitempty"`
}

var tfState *tfStateStruct

// Utility method to read in tfstate content
func readTfState() error {
	// tfWorkPath global var
	tfStateFilename := filepath.Join(tfWorkPath, "terraform.tfstate")
	if _, err := os.Stat(tfStateFilename); err != nil {
		return err
	}
	stateData, err := ioutil.ReadFile(tfStateFilename)
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

// check if resource present in state
func checkForResource(rtype string, name string) bool {

	if tfState == nil {
		if err := readTfState(); err != nil {
			// not differentiating between not exists and file error
			return false
		}
	}
	for _, r := range tfState.Resources {
		if r.Type == rtype && r.Name == name {
			return true
		}
	}

	return false
}
