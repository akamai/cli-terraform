package imaging

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockProcessor struct {
	mock.Mock
}

func (m *mockProcessor) ProcessTemplates(i interface{}) error {
	args := m.Called(i)
	return args.Error(0)
}

func TestMain(m *testing.M) {
	if err := os.MkdirAll("./testdata/res", 0755); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := os.RemoveAll("./testdata/res"); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func TestCreateImaging(t *testing.T) {
	section := "test_section"
	tests := map[string]struct {
		init      func(*mockimaging, *mockProcessor)
		withError error
	}{
		"fetch policy set with given id and contract": {
			init: func(i *mockimaging, p *mockProcessor) {
				i.On("GetPolicySet", mock.Anything, imaging.GetPolicySetRequest{
					PolicySetID: "test_policyset_id",
					ContractID:  "ctr_123",
				}).Return(&imaging.PolicySet{
					ID:     "test_policyset_id",
					Name:   "some policy set",
					Region: "EMEA",
					Type:   "IMAGE",
				}, nil).Once()

				p.On("ProcessTemplates", TFImagingData{
					PolicySet: TFPolicySet{
						ID:         "test_policyset_id",
						ContractID: "ctr_123",
						Name:       "some policy set",
						Region:     "EMEA",
						Type:       "IMAGE",
					},
					Section: "test_section",
				}).Return(nil).Once()
			},
		},
		"error fetching policy set": {
			init: func(i *mockimaging, p *mockProcessor) {
				i.On("GetPolicySet", mock.Anything, i.On("GetPolicySet", mock.Anything, imaging.GetPolicySetRequest{
					PolicySetID: "test_policyset_id",
					ContractID:  "ctr_123",
				}).Return(nil, fmt.Errorf("oops"), nil).Once())
			},
			withError: ErrFetchingPolicySet,
		},
		"error processing template": {
			init: func(i *mockimaging, p *mockProcessor) {
				i.On("GetPolicySet", mock.Anything, imaging.GetPolicySetRequest{
					PolicySetID: "test_policyset_id",
					ContractID:  "ctr_123",
				}).Return(&imaging.PolicySet{
					ID:     "test_policyset_id",
					Name:   "some policy set",
					Region: "EMEA",
					Type:   "IMAGE",
				}, nil).Once()

				p.On("ProcessTemplates", TFImagingData{
					PolicySet: TFPolicySet{
						ID:         "test_policyset_id",
						ContractID: "ctr_123",
						Name:       "some policy set",
						Region:     "EMEA",
						Type:       "IMAGE",
					},
					Section: "test_section",
				}).Return(templates.ErrSavingFiles).Once()
			},
			withError: templates.ErrSavingFiles,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mi := new(mockimaging)
			mp := new(mockProcessor)
			test.init(mi, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createImaging(ctx, "ctr_123", "test_policyset_id", section, mi, mp)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			mi.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessPolicyTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    TFImagingData
		dir          string
		filesToCheck []string
	}{
		"policy with ER match rules and activations": {
			givenData: TFImagingData{
				PolicySet: TFPolicySet{
					ID:         "test_policyset_id",
					ContractID: "ctr_123",
					Name:       "some policy set",
					Region:     "EMEA",
					Type:       "IMAGE",
				},
				Section: "test_section",
			},
			dir:          "only_default_policy",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"imaging.tmpl":   fmt.Sprintf("./testdata/res/%s/imaging.tf", test.dir),
					"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
					"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
				},
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))

			for _, f := range test.filesToCheck {
				expected, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := ioutil.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
