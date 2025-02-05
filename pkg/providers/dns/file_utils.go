package dns

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/akamai/cli/v2/pkg/terminal"
)

type fileUtils interface {
	createModuleTF(ctx context.Context, modName string, content string, tfWorkPath string) error
	appendRootModuleTF(configText string) error
}

type fileUtilsProcessor struct{}

// Work routine to create module TF file
func (fileUtilsProcessor) createModuleTF(ctx context.Context, modName, content, tfWorkPath string) (err error) {
	term := terminal.Get(ctx)
	term.Printf("Creating zone name %s module configuration file...", modName)
	namedModulePath := createNamedModulePath(modName, tfWorkPath)
	if !createDirectory(namedModulePath) {
		return fmt.Errorf("failed to create name module folder: %s", namedModulePath)
	}
	moduleFilename := filepath.Join(namedModulePath, normalizeResourceName(modName)+".tf")
	if _, err := os.Stat(moduleFilename); err == nil {
		// File exists.
		return fmt.Errorf("module configuration file already exists: %s", moduleFilename)
	}
	f, err := os.Create(moduleFilename)
	if err != nil {
		return fmt.Errorf("failed to create name module configuration file: %s", namedModulePath)
	}
	defer func(f *os.File) {
		if e := f.Close(); e != nil {
			err = e
		}
	}(f)
	_, err = f.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write name module configuration: %s", namedModulePath)
	}
	err = f.Sync()

	return err
}

// Flush string to root module TF file
func (fileUtilsProcessor) appendRootModuleTF(configText string) error {
	// save top level Zone TF config
	_, err := zoneTFFileHandle.Write([]byte(configText))
	if err != nil {
		return fmt.Errorf("failed to save zone configuration file")
	}
	err = zoneTFFileHandle.Sync()
	if err != nil {
		return err
	}

	return nil
}
