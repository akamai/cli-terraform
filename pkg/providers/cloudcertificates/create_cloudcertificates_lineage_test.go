package cloudcertificates

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	edgegridcloudcertificates "github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/cli-terraform/v3/pkg/templates"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testRSASignedCertificatePEM = `-----BEGIN CERTIFICATE-----
rsa-signed-certificate
-----END CERTIFICATE-----
`
	testECDSASignedCertificatePEM = `-----BEGIN CERTIFICATE-----
ecdsa-signed-certificate
-----END CERTIFICATE-----
`
	testTrustChainPEM = `-----BEGIN CERTIFICATE-----
test-trust-chain
-----END CERTIFICATE-----
`
	testLineageID              int64 = 123456
	testHeadGenerationID       int64 = 567890
	testStagingGenerationID    int64 = 123400
	testProductionGenerationID int64 = 123401
	testStagingActivationID    int64 = 1001
	testProductionActivationID int64 = 1002
)

func TestCreateCloudCertificate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		dir                 string
		acknowledgeWarnings bool
		init                func(*edgegridcloudcertificates.Mock, string)
		expectedError       string
	}{
		{
			name: "RSA only no upload",
			dir:  "rsa_only_no_upload",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}}), nil),
		},
		{
			name: "RSA only no upload with common name subject",
			dir:  "rsa_only_no_upload_common_name_subject",
			init: setupCloudCertificate(testLineageWithSubject(
				[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
				edgegridcloudcertificates.Subject{CommonName: "example.com"}), nil),
		},
		{
			name: "exact lineage name follows substring result",
			dir:  "rsa_only_no_upload",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}}), nil,
				[]*edgegridcloudcertificates.ListLineagesResponse{
					{Lineages: []edgegridcloudcertificates.Lineage{
						{LineageID: 123455, LineageName: "rsa-only-no-upload-previous"},
						{LineageID: testLineageID, LineageName: "rsa-only-no-upload"},
					}},
				}),
		},
		{
			name: "exact lineage name on next page after substring results",
			dir:  "rsa_only_no_upload",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}}), nil,
				[]*edgegridcloudcertificates.ListLineagesResponse{
					{
						Lineages:   substringLineages("rsa-only-no-upload", edgegridcloudcertificates.DefaultListLineagesPageSize),
						NextCursor: new("next-page"),
					},
					{Lineages: []edgegridcloudcertificates.Lineage{
						{LineageID: testLineageID, LineageName: "rsa-only-no-upload"},
					}},
				}),
		},
		{
			name: "ECDSA only no upload",
			dir:  "ecdsa_only_no_upload",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "ECDSA", KeySize: "P-256"}}), nil),
		},
		{
			name: "RSA and ECDSA no upload",
			dir:  "rsa_and_ecdsa_no_upload",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{
				{KeyType: "RSA", KeySize: "4096"},
				{KeyType: "ECDSA", KeySize: "P-384"},
			}), nil),
		},
		{
			name: "RSA only with signed certificate and trust chain",
			dir:  "rsa_only_with_upload",
			init: setupCloudCertificate(testLineage(
				[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
				[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}), nil),
		},
		{
			name: "RSA only with signed certificate and no trust chain",
			dir:  "rsa_only_with_signed_certificate_no_trust_chain",
			init: setupCloudCertificate(testLineage(
				[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
				[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, nil)}), nil),
		},
		{
			name: "ECDSA only with signed certificate and no trust chain",
			dir:  "ecdsa_only_with_upload",
			init: setupCloudCertificate(testLineage(
				[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "ECDSA", KeySize: "P-256"}},
				[]edgegridcloudcertificates.Algorithm{readyAlgorithm("ECDSA", testECDSASignedCertificatePEM, nil)}), nil),
		},
		{
			name: "ECDSA only with signed certificate and trust chain",
			dir:  "ecdsa_only_with_signed_certificate_and_trust_chain",
			init: setupCloudCertificate(testLineage(
				[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "ECDSA", KeySize: "P-256"}},
				[]edgegridcloudcertificates.Algorithm{readyAlgorithm("ECDSA", testECDSASignedCertificatePEM, new(testTrustChainPEM))}), nil),
		},
		{
			name: "RSA and ECDSA with signed certificates and mixed trust chains",
			dir:  "rsa_and_ecdsa_both_with_upload",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{
				{KeyType: "RSA", KeySize: "4096"},
				{KeyType: "ECDSA", KeySize: "P-384"},
			}, []edgegridcloudcertificates.Algorithm{
				readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM)),
				readyAlgorithm("ECDSA", testECDSASignedCertificatePEM, nil),
			}), nil),
		},
		{
			name: "RSA and ECDSA with signed certificates and no trust chain",
			dir:  "rsa_and_ecdsa_with_signed_certificates_no_trust_chain",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{
				{KeyType: "RSA", KeySize: "4096"},
				{KeyType: "ECDSA", KeySize: "P-384"},
			}, []edgegridcloudcertificates.Algorithm{
				readyAlgorithm("RSA", testRSASignedCertificatePEM, nil),
				readyAlgorithm("ECDSA", testECDSASignedCertificatePEM, nil),
			}), nil),
		},
		{
			name: "RSA and ECDSA with signed certificates and trust chains",
			dir:  "rsa_and_ecdsa_with_signed_certificates_and_trust_chains",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{
				{KeyType: "RSA", KeySize: "4096"},
				{KeyType: "ECDSA", KeySize: "P-384"},
			}, []edgegridcloudcertificates.Algorithm{
				readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM)),
				readyAlgorithm("ECDSA", testECDSASignedCertificatePEM, new(testTrustChainPEM)),
			}), nil),
		},
		{
			name: "RSA and ECDSA with only RSA upload",
			dir:  "rsa_and_ecdsa_only_rsa_with_upload",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{
				{KeyType: "RSA", KeySize: "4096"},
				{KeyType: "ECDSA", KeySize: "P-384"},
			}, []edgegridcloudcertificates.Algorithm{
				readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM)),
				{CertificateStatus: edgegridcloudcertificates.CertificateStatusCSRReady, KeyType: "ECDSA"},
			}), nil),
		},
		{
			name: "RSA and ECDSA with READY_FOR_USE ECDSA certificate missing",
			dir:  "rsa_and_ecdsa_ready_without_ecdsa_certificate",
			init: setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{
				{KeyType: "RSA", KeySize: "4096"},
				{KeyType: "ECDSA", KeySize: "P-384"},
			}, []edgegridcloudcertificates.Algorithm{
				readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM)),
				{
					CertificateStatus:    edgegridcloudcertificates.CertificateStatusReadyForUse,
					KeyType:              "ECDSA",
					SignedCertificatePEM: nil,
				},
			}), nil),
		},
		{
			name: "RSA only with upload and acknowledgement true",
			dir:  "rsa_only_with_upload_acknowledgement_true",
			init: setupCloudCertificate(testLineage(
				[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
				[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}), nil),
			acknowledgeWarnings: true,
		},
		{
			name: "RSA only with upload without acknowledgement",
			dir:  "rsa_only_with_upload_acknowledgement_false",
			init: setupCloudCertificate(testLineage(
				[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
				[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}), nil),
			acknowledgeWarnings: false,
		},
		{
			name:                "RSA only without upload and acknowledgement true",
			dir:                 "rsa_only_without_upload_acknowledgement_true",
			init:                setupCloudCertificate(testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}}), nil),
			acknowledgeWarnings: true,
		},
		{
			name: "RSA with staging activation",
			dir:  "rsa_with_staging_activation",
			init: setupCloudCertificate(lineageWithActivations(
				testLineage(
					[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
					[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}),
				int64Ptr(testStagingGenerationID), nil), testActivations(int64Ptr(testStagingGenerationID), nil)),
		},
		{
			name: "RSA with production activation",
			dir:  "rsa_with_production_activation",
			init: setupCloudCertificate(lineageWithActivations(
				testLineage(
					[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
					[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}),
				nil, int64Ptr(testProductionGenerationID)), testActivations(nil, int64Ptr(testProductionGenerationID))),
		},
		{
			name: "RSA with staging and production activations",
			dir:  "rsa_with_staging_and_production_activations",
			init: setupCloudCertificate(lineageWithActivations(
				testLineage(
					[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
					[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}),
				int64Ptr(testStagingGenerationID), int64Ptr(testProductionGenerationID)),
				testActivations(int64Ptr(testStagingGenerationID), int64Ptr(testProductionGenerationID))),
		},
		{
			name: "RSA with latest generation on staging and production",
			dir:  "rsa_with_latest_generation_on_staging_and_production",
			init: setupCloudCertificate(lineageWithActivations(
				uploadFromCurrentProductionLineage(),
				int64Ptr(testProductionGenerationID), int64Ptr(testProductionGenerationID)),
				testActivations(int64Ptr(testProductionGenerationID), int64Ptr(testProductionGenerationID))),
		},
		{
			name: "RSA with incomplete latest staging and production activations",
			dir:  "rsa_with_staging_and_production_activations",
			init: setupCloudCertificate(lineageWithActivations(
				testLineage(
					[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
					[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}),
				int64Ptr(testStagingGenerationID), int64Ptr(testProductionGenerationID)),
				testActivationsWithIncompleteLatest()),
		},
		{
			name: "RSA with in-progress staging and production activations",
			dir:  "rsa_with_in_progress_staging_and_production_activations",
			init: setupCloudCertificate(
				lineageWithActivations(
					testLineage(
						[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
						[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}),
					int64Ptr(testStagingGenerationID), int64Ptr(testProductionGenerationID)),
				testActivationsWithStatus("IN_PROGRESS")),
		},
		{
			name: "RSA with in-progress staging and completed production activations",
			dir:  "rsa_with_in_progress_staging_and_completed_production_activation",
			init: setupCloudCertificate(
				lineageWithActivations(
					testLineage(
						[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
						[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}),
					int64Ptr(testStagingGenerationID), int64Ptr(testProductionGenerationID)),
				testActivationsWithInProgressStagingAndCompletedProduction()),
		},
		{
			name: "RSA with partially successful staging and production activations",
			dir:  "rsa_with_staging_and_production_activations",
			init: setupCloudCertificate(
				lineageWithActivations(
					testLineage(
						[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
						[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}),
					int64Ptr(testStagingGenerationID), int64Ptr(testProductionGenerationID)),
				testActivationsWithStatus("PARTIAL_SUCCESS")),
		},
		{
			name: "activation pagination through next cursor",
			dir:  "rsa_with_staging_activation",
			init: setupCloudCertificateWithActivationPages(
				lineageWithActivations(
					testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
						[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))}),
					int64Ptr(testStagingGenerationID), nil),
				&edgegridcloudcertificates.ListActivationsResponse{
					NextCursor: new("next-page"),
					Items: []edgegridcloudcertificates.GetActivationStatusResponse{{
						ActivationID:      2001,
						ActivationStatus:  "COMPLETE",
						GenerationID:      testHeadGenerationID,
						TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkStaging),
					}},
				},
				testActivations(int64Ptr(testStagingGenerationID), nil)),
		},
		{
			name: "RSA with upload from current production on staging and production when head is nil",
			dir:  "rsa_with_upload_from_current_production",
			init: setupCloudCertificate(
				lineageWithActivations(
					uploadFromCurrentProductionLineage(),
					int64Ptr(testProductionGenerationID),
					int64Ptr(testProductionGenerationID),
				),
				testActivations(int64Ptr(testProductionGenerationID), int64Ptr(testProductionGenerationID))),
		},
		{
			name:          "lineage not found",
			dir:           "lineage-not-found",
			expectedError: `error fetching certificate lineage: lineage with name "lineage-not-found" not found`,
			init: setupCloudCertificate(edgegridcloudcertificates.GetLineageResponse{}, nil, []*edgegridcloudcertificates.ListLineagesResponse{
				{Lineages: []edgegridcloudcertificates.Lineage{}},
			}),
		},
		{
			name:          "only substring lineage name result",
			dir:           "lineage-not-found",
			expectedError: `error fetching certificate lineage: lineage with name "lineage-not-found" not found`,
			init: setupCloudCertificate(edgegridcloudcertificates.GetLineageResponse{}, nil, []*edgegridcloudcertificates.ListLineagesResponse{
				{Lineages: []edgegridcloudcertificates.Lineage{
					{LineageID: 123455, LineageName: "lineage-not-found-previous"},
				}},
			}),
		},
		{
			name:          "list lineages returns an error",
			dir:           "list-lineages-error",
			expectedError: "error fetching certificate lineage: list lineages failed",
			init:          setupListLineagesError(fmt.Errorf("list lineages failed")),
		},
		{
			name:          "get lineage returns an error",
			dir:           "get-lineage-error",
			expectedError: "error fetching certificate lineage: get lineage failed",
			init: setupGetLineageError(testLineage([]edgegridcloudcertificates.KeySpecResponse{
				{KeyType: "RSA", KeySize: "2048"},
			}), fmt.Errorf("get lineage failed")),
		},
		{
			name:          "list activations returns an error",
			dir:           "list-activations-error",
			expectedError: "error fetching certificate lineage activations: list activations failed",
			init: setupListActivationsError(lineageWithActivations(
				testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}}),
				int64Ptr(testStagingGenerationID), nil), fmt.Errorf("list activations failed")),
		},
		{
			name:          "no matching activation found",
			dir:           "no-matching-activation",
			expectedError: "error fetching certificate lineage activations: no activation found for a current network generation on STAGING",
			init: setupCloudCertificate(lineageWithActivations(
				testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}}),
				int64Ptr(testStagingGenerationID), nil), &edgegridcloudcertificates.ListActivationsResponse{}),
		},
		{
			// Only a theoretical case as if generation is active on network, it should always have activation in the list
			name:          "activation exists but not for the current generation",
			dir:           "activation-for-wrong-generation",
			expectedError: "error fetching certificate lineage activations: no activation found for a current network generation on STAGING",
			init: setupCloudCertificate(lineageWithActivations(
				testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}}),
				int64Ptr(testStagingGenerationID), nil), &edgegridcloudcertificates.ListActivationsResponse{
				Items: []edgegridcloudcertificates.GetActivationStatusResponse{{
					ActivationID:      999,
					ActivationStatus:  "COMPLETE",
					GenerationID:      testHeadGenerationID,
					TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkStaging),
				}},
			}),
		},
		{
			name:          "no matching staging and production activations found",
			dir:           "no-matching-activations",
			expectedError: "error fetching certificate lineage activations: no activation found for a current network generation on STAGING and PRODUCTION",
			init: setupCloudCertificate(lineageWithActivations(
				testLineage([]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}}),
				int64Ptr(testStagingGenerationID), int64Ptr(testProductionGenerationID)), &edgegridcloudcertificates.ListActivationsResponse{}),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &edgegridcloudcertificates.Mock{}
			lineageName := strings.ReplaceAll(test.dir, "_", "-")
			test.init(client, lineageName)

			tempDir := t.TempDir()
			params := createCloudCertificateParams{
				lineageName:         lineageName,
				acknowledgeWarnings: test.acknowledgeWarnings,
				edgercPath:          "~/.edgerc",
				configSection:       "default",
				client:              client,
				templateProcessor: templates.FSTemplateProcessor{
					TemplatesFS: templateFiles,
					TemplateTargets: map[string]string{
						"cloudcertificate.tmpl": filepath.Join(tempDir, "cloudcertificate.tf"),
						"variables.tmpl":        filepath.Join(tempDir, "variables.tf"),
						"import.tmpl":           filepath.Join(tempDir, "import.sh"),
					},
					AdditionalFuncs: additionalFunctions,
				},
			}

			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createCloudCertificate(ctx, params)

			client.AssertExpectations(t)
			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
				return
			}
			require.NoError(t, err)

			if test.dir != "" {
				for _, filename := range []string{"cloudcertificate.tf", "import.sh", "variables.tf"} {
					expected, err := os.ReadFile(filepath.Join("testdata", test.dir, filename))
					require.NoError(t, err)
					actual, err := os.ReadFile(filepath.Join(tempDir, filename))
					require.NoError(t, err)
					assert.Equal(t, string(expected), string(actual))
				}
			}

		})
	}
}

func setupCloudCertificate(lineage edgegridcloudcertificates.GetLineageResponse, activations *edgegridcloudcertificates.ListActivationsResponse,
	listLineages ...[]*edgegridcloudcertificates.ListLineagesResponse) func(*edgegridcloudcertificates.Mock, string) {
	return func(client *edgegridcloudcertificates.Mock, lineageName string) {
		lineage.LineageName = lineageName
		responses := []*edgegridcloudcertificates.ListLineagesResponse{
			{Lineages: []edgegridcloudcertificates.Lineage{edgegridcloudcertificates.Lineage(lineage)}},
		}
		if len(listLineages) != 0 {
			responses = listLineages[0]
		}

		cursor := ""
		for _, response := range responses {
			client.On("ListLineages", mock.Anything, edgegridcloudcertificates.ListLineagesRequest{
				LineageName: lineage.LineageName,
				After:       cursor,
			}).Return(response, nil).Once()
			if response.NextCursor != nil {
				cursor = *response.NextCursor
			}
		}
		if lineage.LineageID != 0 {
			client.On("GetLineage", mock.Anything, edgegridcloudcertificates.GetLineageRequest{
				LineageID: lineage.LineageID,
				ExpandGenerations: []edgegridcloudcertificates.ExpandGenerations{
					edgegridcloudcertificates.ExpandGenerationsHead,
					edgegridcloudcertificates.ExpandGenerationsCurrentProduction,
					edgegridcloudcertificates.ExpandGenerationsCurrentStaging,
				},
			}).Return(&lineage, nil).Once()
		}
		if activations != nil {
			client.On("ListActivations", mock.Anything, edgegridcloudcertificates.ListActivationsRequest{
				LineageID: lineage.LineageID,
			}).Return(activations, nil).Once()
		}
	}
}

func setupCloudCertificateWithActivationPages(lineage edgegridcloudcertificates.GetLineageResponse,
	firstPage, secondPage *edgegridcloudcertificates.ListActivationsResponse) func(*edgegridcloudcertificates.Mock, string) {
	return func(client *edgegridcloudcertificates.Mock, lineageName string) {
		lineage.LineageName = lineageName
		client.On("ListLineages", mock.Anything, edgegridcloudcertificates.ListLineagesRequest{
			LineageName: lineageName,
		}).Return(&edgegridcloudcertificates.ListLineagesResponse{
			Lineages: []edgegridcloudcertificates.Lineage{edgegridcloudcertificates.Lineage(lineage)},
		}, nil).Once()
		client.On("GetLineage", mock.Anything, edgegridcloudcertificates.GetLineageRequest{
			LineageID: lineage.LineageID,
			ExpandGenerations: []edgegridcloudcertificates.ExpandGenerations{
				edgegridcloudcertificates.ExpandGenerationsHead,
				edgegridcloudcertificates.ExpandGenerationsCurrentProduction,
				edgegridcloudcertificates.ExpandGenerationsCurrentStaging,
			},
		}).Return(&lineage, nil).Once()
		client.On("ListActivations", mock.Anything, edgegridcloudcertificates.ListActivationsRequest{
			LineageID: lineage.LineageID,
		}).Return(firstPage, nil).Once()
		client.On("ListActivations", mock.Anything, edgegridcloudcertificates.ListActivationsRequest{
			LineageID: lineage.LineageID,
			Cursor:    "next-page",
		}).Return(secondPage, nil).Once()
	}
}

func setupListLineagesError(err error) func(*edgegridcloudcertificates.Mock, string) {
	return func(client *edgegridcloudcertificates.Mock, lineageName string) {
		client.On("ListLineages", mock.Anything, edgegridcloudcertificates.ListLineagesRequest{
			LineageName: lineageName,
		}).Return(nil, err).Once()
	}
}

func setupGetLineageError(lineage edgegridcloudcertificates.GetLineageResponse, err error) func(*edgegridcloudcertificates.Mock, string) {
	return func(client *edgegridcloudcertificates.Mock, lineageName string) {
		lineage.LineageName = lineageName
		client.On("ListLineages", mock.Anything, edgegridcloudcertificates.ListLineagesRequest{
			LineageName: lineageName,
		}).Return(&edgegridcloudcertificates.ListLineagesResponse{
			Lineages: []edgegridcloudcertificates.Lineage{edgegridcloudcertificates.Lineage(lineage)},
		}, nil).Once()
		client.On("GetLineage", mock.Anything, edgegridcloudcertificates.GetLineageRequest{
			LineageID: lineage.LineageID,
			ExpandGenerations: []edgegridcloudcertificates.ExpandGenerations{
				edgegridcloudcertificates.ExpandGenerationsHead,
				edgegridcloudcertificates.ExpandGenerationsCurrentProduction,
				edgegridcloudcertificates.ExpandGenerationsCurrentStaging,
			},
		}).Return(nil, err).Once()
	}
}

func setupListActivationsError(lineage edgegridcloudcertificates.GetLineageResponse, err error) func(*edgegridcloudcertificates.Mock, string) {
	return func(client *edgegridcloudcertificates.Mock, lineageName string) {
		lineage.LineageName = lineageName
		client.On("ListLineages", mock.Anything, edgegridcloudcertificates.ListLineagesRequest{
			LineageName: lineageName,
		}).Return(&edgegridcloudcertificates.ListLineagesResponse{
			Lineages: []edgegridcloudcertificates.Lineage{edgegridcloudcertificates.Lineage(lineage)},
		}, nil).Once()
		client.On("GetLineage", mock.Anything, edgegridcloudcertificates.GetLineageRequest{
			LineageID: lineage.LineageID,
			ExpandGenerations: []edgegridcloudcertificates.ExpandGenerations{
				edgegridcloudcertificates.ExpandGenerationsHead,
				edgegridcloudcertificates.ExpandGenerationsCurrentProduction,
				edgegridcloudcertificates.ExpandGenerationsCurrentStaging,
			},
		}).Return(&lineage, nil).Once()
		client.On("ListActivations", mock.Anything, edgegridcloudcertificates.ListActivationsRequest{
			LineageID: lineage.LineageID,
		}).Return(nil, err).Once()
	}
}

func testLineage(keySpecs []edgegridcloudcertificates.KeySpecResponse, algorithms ...[]edgegridcloudcertificates.Algorithm) edgegridcloudcertificates.GetLineageResponse {
	return testLineageWithSubject(keySpecs, edgegridcloudcertificates.Subject{
		CommonName:         "example.com",
		Organization:       "Example, Inc.",
		OrganizationalUnit: "Security",
		Country:            "US",
		State:              "Massachusetts",
		Locality:           "Cambridge",
	}, algorithms...)
}

func testLineageWithSubject(keySpecs []edgegridcloudcertificates.KeySpecResponse, subject edgegridcloudcertificates.Subject, algorithms ...[]edgegridcloudcertificates.Algorithm) edgegridcloudcertificates.GetLineageResponse {
	headAlgorithms := csrReadyAlgorithms(keySpecs)
	if len(algorithms) != 0 {
		headAlgorithms = algorithms[0]
	}
	return edgegridcloudcertificates.GetLineageResponse{
		ContractID:    "ctr_1-ABC",
		GeoClass:      "STANDARD_WORLDWIDE",
		GroupID:       123,
		KeySpecs:      keySpecs,
		LineageID:     testLineageID,
		LineageType:   "MULTIPLE_GENERATION",
		SANs:          []string{"example.com", "*.example.com"},
		SecureNetwork: "ENHANCED_TLS",
		Subject:       subject,
		Head: &edgegridcloudcertificates.HeadGeneration{
			Generation: edgegridcloudcertificates.Generation{
				Algorithms: headAlgorithms,
			},
			HeadGenerationID: testHeadGenerationID,
		},
	}
}

func uploadFromCurrentProductionLineage() edgegridcloudcertificates.GetLineageResponse {
	lineage := testLineage(
		[]edgegridcloudcertificates.KeySpecResponse{{KeyType: "RSA", KeySize: "2048"}},
		[]edgegridcloudcertificates.Algorithm{readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM))},
	)
	lineage.Head = nil
	lineage.CurrentProduction = &edgegridcloudcertificates.ProductionGeneration{
		Generation: edgegridcloudcertificates.Generation{
			Algorithms: []edgegridcloudcertificates.Algorithm{
				readyAlgorithm("RSA", testRSASignedCertificatePEM, new(testTrustChainPEM)),
			},
		},
		ProductionGenerationID: testProductionGenerationID,
	}
	return lineage
}

func testActivations(stagingGenerationID, productionGenerationID *int64) *edgegridcloudcertificates.ListActivationsResponse {
	activations := &edgegridcloudcertificates.ListActivationsResponse{}
	if productionGenerationID != nil {
		activations.Items = append(activations.Items, edgegridcloudcertificates.GetActivationStatusResponse{
			ActivationID:      testProductionActivationID,
			ActivationStatus:  "COMPLETE",
			GenerationID:      *productionGenerationID,
			TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkProduction),
		})
	}
	if stagingGenerationID != nil {
		activations.Items = append(activations.Items, edgegridcloudcertificates.GetActivationStatusResponse{
			ActivationID:      testStagingActivationID,
			ActivationStatus:  "COMPLETE",
			GenerationID:      *stagingGenerationID,
			TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkStaging),
		})
	}
	return activations
}

func testActivationsWithIncompleteLatest() *edgegridcloudcertificates.ListActivationsResponse {
	return &edgegridcloudcertificates.ListActivationsResponse{
		Items: []edgegridcloudcertificates.GetActivationStatusResponse{
			{
				ActivationID:      2001,
				ActivationStatus:  "ABORTED",
				GenerationID:      testHeadGenerationID,
				TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkProduction),
			},
			{
				ActivationID:      2002,
				ActivationStatus:  "ABORTED",
				GenerationID:      testHeadGenerationID,
				TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkStaging),
			},
			{
				ActivationID:      testProductionActivationID,
				ActivationStatus:  "COMPLETE",
				GenerationID:      testProductionGenerationID,
				TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkProduction),
			},
			{
				ActivationID:      testStagingActivationID,
				ActivationStatus:  "COMPLETE",
				GenerationID:      testStagingGenerationID,
				TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkStaging),
			},
		},
	}
}

func testActivationsWithStatus(status string) *edgegridcloudcertificates.ListActivationsResponse {
	return &edgegridcloudcertificates.ListActivationsResponse{
		Items: []edgegridcloudcertificates.GetActivationStatusResponse{
			{
				ActivationID:      testProductionActivationID,
				ActivationStatus:  status,
				GenerationID:      testProductionGenerationID,
				TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkProduction),
			},
			{
				ActivationID:      testStagingActivationID,
				ActivationStatus:  status,
				GenerationID:      testStagingGenerationID,
				TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkStaging),
			},
		},
	}
}

func testActivationsWithInProgressStagingAndCompletedProduction() *edgegridcloudcertificates.ListActivationsResponse {
	return &edgegridcloudcertificates.ListActivationsResponse{
		Items: []edgegridcloudcertificates.GetActivationStatusResponse{
			{
				ActivationID:      testProductionActivationID,
				ActivationStatus:  "COMPLETE",
				GenerationID:      testProductionGenerationID,
				TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkProduction),
			},
			{
				ActivationID:      testStagingActivationID,
				ActivationStatus:  "IN_PROGRESS",
				GenerationID:      testStagingGenerationID,
				TargetEnvironment: string(edgegridcloudcertificates.TargetNetworkStaging),
			},
		},
	}
}

func lineageWithActivations(lineage edgegridcloudcertificates.GetLineageResponse, stagingGenerationID, productionGenerationID *int64) edgegridcloudcertificates.GetLineageResponse {
	if stagingGenerationID != nil {
		lineage.CurrentStaging = &edgegridcloudcertificates.StagingGeneration{
			StagingGenerationID: *stagingGenerationID,
		}
	}
	if productionGenerationID != nil {
		if lineage.CurrentProduction == nil {
			lineage.CurrentProduction = &edgegridcloudcertificates.ProductionGeneration{}
		}
		lineage.CurrentProduction.ProductionGenerationID = *productionGenerationID
	}

	return lineage
}

func int64Ptr(value int64) *int64 {
	return &value
}

func csrReadyAlgorithms(keySpecs []edgegridcloudcertificates.KeySpecResponse) []edgegridcloudcertificates.Algorithm {
	algorithms := make([]edgegridcloudcertificates.Algorithm, 0, len(keySpecs))
	for _, keySpec := range keySpecs {
		algorithms = append(algorithms, edgegridcloudcertificates.Algorithm{
			CertificateStatus: edgegridcloudcertificates.CertificateStatusCSRReady,
			KeyType:           keySpec.KeyType,
		})
	}
	return algorithms
}

func readyAlgorithm(keyType string, signedCertificatePEM string, trustChainPEM *string) edgegridcloudcertificates.Algorithm {
	return edgegridcloudcertificates.Algorithm{
		CertificateStatus:    edgegridcloudcertificates.CertificateStatusReadyForUse,
		KeyType:              keyType,
		SignedCertificatePEM: new(signedCertificatePEM),
		TrustChainPEM:        trustChainPEM,
	}
}

func substringLineages(lineageName string, count int) []edgegridcloudcertificates.Lineage {
	lineages := make([]edgegridcloudcertificates.Lineage, count)
	for i := range lineages {
		lineages[i] = edgegridcloudcertificates.Lineage{
			LineageID:   int64(i + 1),
			LineageName: fmt.Sprintf("%s-substring-%d", lineageName, i+1),
		}
	}

	return lineages
}
