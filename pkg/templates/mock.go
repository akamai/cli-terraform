package templates

import "github.com/stretchr/testify/mock"

// MockProcessor is a mock for TemplateProcessor
type MockProcessor struct {
	mock.Mock
}

// MockMultiTargetProcessor is a mock for MultiTargetProcessor
type MockMultiTargetProcessor struct {
	mock.Mock
}

// ProcessTemplates is a mocked version
func (m *MockProcessor) ProcessTemplates(i interface{}, filterFuncs ...func([]string) ([]string, error)) error {
	var args mock.Arguments
	switch len(filterFuncs) {
	case 0:
		args = m.Called(i)
	case 1:
		args = m.Called(i, filterFuncs[0])
	}
	return args.Error(0)
}

// AddTemplateTarget is a mocked version
func (m *MockProcessor) AddTemplateTarget(templateName, targetPath string) {
	m.Called(templateName, targetPath)
}

// TemplateExists is a mocked version
func (m *MockProcessor) TemplateExists(fileName string) bool {
	args := m.Called(fileName)
	return args.Bool(0)
}

// ProcessTemplates is a mocked version
func (m *MockMultiTargetProcessor) ProcessTemplates(i MultiTargetData, filterFuncs ...func([]string) ([]string, error)) error {
	var args mock.Arguments
	switch len(filterFuncs) {
	case 0:
		args = m.Called(i)
	case 1:
		args = m.Called(i, filterFuncs[0])
	}
	return args.Error(0)
}
