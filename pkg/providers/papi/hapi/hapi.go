package hapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Certificate holds certificate information
type Certificate struct {
	CertificateID    string    `json:"certificateId"`
	CommonName       string    `json:"commonName"`
	SerialNumber     string    `json:"serialNumber"`
	SlotNumber       int       `json:"slotNumber"`
	ExpirationDate   time.Time `json:"expirationDate"`
	CertificateType  string    `json:"certificateType"`
	ValidationType   string    `json:"validationType"`
	Status           string    `json:"status"`
	AvailableDomains []string  `json:"availableDomains"`
}

// EdgeHostname represents edge hostname
type EdgeHostname struct {
	EdgeHostnameID    int    `json:"edgeHostnameId"`
	RecordName        string `json:"recordName"`
	DNSZone           string `json:"dnsZone"`
	SecurityType      string `json:"securityType"`
	UseDefaultTTL     bool   `json:"useDefaultTtl"`
	UseDefaultMap     bool   `json:"useDefaultMap"`
	IPVersionBehavior string `json:"ipVersionBehavior"`
	ProductID         string `json:"productId"`
	TTL               int    `json:"ttl"`
	Map               string `json:"map,omitempty"`
	SlotNumber        int    `json:"slotNumber,omitempty"`
	Comments          string `json:"comments"`
	SerialNumber      int    `json:"serialNumber,omitempty"`
	CustomTarget      string `json:"customTarget,omitempty"`
	ChinaCdn          struct {
		IsChinaCdn bool `json:"isChinaCdn"`
	} `json:"chinaCdn,omitempty"`
	IsEdgeIPBindingEnabled bool `json:"isEdgeIPBindingEnabled,omitempty"`
}

// ListEdgeHostnamesResponse holds a response from ListEdgeHostnames
type ListEdgeHostnamesResponse struct {
	EdgeHostnames []EdgeHostname `json:"edgeHostnames"`
}

// ChangeRequest represents a change response
type ChangeRequest struct {
	Action            string         `json:"action"`
	ChangeID          int            `json:"changeId"`
	Comments          string         `json:"comments"`
	Status            string         `json:"status"`
	StatusMessage     string         `json:"statusMessage"`
	StatusUpdateDate  time.Time      `json:"statusUpdateDate"`
	StatusUpdateEmail string         `json:"statusUpdateEmail"`
	SubmitDate        time.Time      `json:"submitDate"`
	Submitter         string         `json:"submitter"`
	SubmitterEmail    string         `json:"submitterEmail"`
	EdgeHostnames     []EdgeHostname `json:"edgeHostnames"`
}

// ListChangeRequestsResponse holds a response from ListChangeRequests
type ListChangeRequestsResponse struct {
	ChangeRequests []ChangeRequest `json:"changeRequests"`
}

// Patch contains json patch
type Patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// Products contains a response from ListProducts
type Products struct {
	ProductDisplayNames []struct {
		ProductID          string `json:"productId"`
		ProductDisplayName string `json:"productDisplayName"`
	} `json:"productDisplayNames"`
}

// ListEdgeHostnames return a list of all edge hostnames
func ListEdgeHostnames() (*ListEdgeHostnamesResponse, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/hapi/v1/edge-hostnames",
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ListEdgeHostnamesResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetEdgeHostname gets hostame by dnszone and recordname
func GetEdgeHostname(recordName string, dnsZone string) (*EdgeHostname, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s", dnsZone, recordName),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response EdgeHostname
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCertificate returns edge hostname certificate for given dnsZone and recordName
func GetCertificate(recordName string, dnsZone string) (*Certificate, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s/certificate", dnsZone, recordName),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response Certificate
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetEdgeHostnameByID returns edge hostname with given id
func GetEdgeHostnameByID(id string) (*EdgeHostname, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/edge-hostnames/%s", id),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response EdgeHostname
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteEdgeHostname deletes
func DeleteEdgeHostname(recordName string, dnsZone string) (*ChangeRequest, error) {
	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s", dnsZone, recordName),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ChangeRequest
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// PatchEdgeHostname updates edge hostname with given patch
func PatchEdgeHostname(recordName string, dnsZone string, patch []Patch) (*ChangeRequest, error) {

	r, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(
		Config,
		"PATCH",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s", dnsZone, recordName),
		bytes.NewReader(r),
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ChangeRequest
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListChangeRequestsByHostname returns a list of all ChangeRequests with given dnsZone and recordName
func ListChangeRequestsByHostname(recordName string, dnsZone string) (*ListChangeRequestsResponse, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/dns-zones/%s/edge-hostnames/%s/change-requests", dnsZone, recordName),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ListChangeRequestsResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListChangeRequests lists all ChangeRequests
func ListChangeRequests() (*ListChangeRequestsResponse, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/hapi/v1/change-requests",
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ListChangeRequestsResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetChangeRequest returns ChangeRequest with given id
func GetChangeRequest(id string) (*ChangeRequest, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/hapi/v1/change-requests/%s", id),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response ChangeRequest
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListProducts returns list of products
func ListProducts() (*Products, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/hapi/v1/products/display-names",
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response Products
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
