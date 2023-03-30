package templates

import "github.com/stretchr/testify/mock"

// MockProcessor is a mock for TemplateProcessor
type MockProcessor struct {
	mock.Mock
}

// ProcessTemplates is a mocked version
func (m *MockProcessor) ProcessTemplates(i interface{}) error {
	args := m.Called(i)
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
