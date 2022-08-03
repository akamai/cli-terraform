package dns

import (
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fileutilsmock struct {
	mock.Mock
	createModuleArg string
	appendRootArg   string
}

func (m *fileutilsmock) createModuleTF(_ context.Context, modName string, content string, tfWorkPath string) error {
	m.createModuleArg = content
	args := m.Called(modName, content, tfWorkPath)
	return args.Error(0)
}
func (m *fileutilsmock) appendRootModuleTF(configText string) error {
	m.appendRootArg = configText
	args := m.Called(configText)
	return args.Error(0)
}

func assertFileWithContent(t *testing.T, expectedPath, actual string) {
	expectedResult, err := ioutil.ReadFile(expectedPath)
	if err != nil {
		fmt.Print("incorrect expected file")
		return
	}
	expected := strings.ReplaceAll(string(expectedResult), " ", "")
	actual = strings.ReplaceAll(actual, " ", "")
	expectedSplitted := strings.Split(expected, "\n")
	actualSplitted := strings.Split(actual, "\n")
	sort.Strings(expectedSplitted)
	sort.Strings(actualSplitted)
	assert.Equal(t, expectedSplitted, actualSplitted)
}
