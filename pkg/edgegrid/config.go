// Package edgegrid contains code for manipulating edgegrid access settings
package edgegrid

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/urfave/cli/v2"
)

// GetEdgegridConfig gets configuration from .edgerc file
func GetEdgegridConfig(c *cli.Context) (*edgegrid.Config, error) {
	edgercOps := []edgegrid.Option{
		edgegrid.WithEnv(true),
		edgegrid.WithFile(GetEdgercPath(c)),
		edgegrid.WithSection(GetEdgercSection(c)),
	}
	config, err := edgegrid.New(edgercOps...)
	if err != nil {
		return nil, err
	}
	if c.IsSet("accountkey") {
		config.AccountKey = c.String("accountkey")
	}
	return config, nil
}

// GetEdgercPath returns the location of edgerc credential file or "~/.edgerc" if not found
func GetEdgercPath(c *cli.Context) string {
	edgercPath := c.String("edgerc")
	if edgercPath == "" {
		return edgegrid.DefaultConfigFile
	}
	return edgercPath
}

// GetEdgercSection returns the section in edgerc credential file or "default" if not found
func GetEdgercSection(c *cli.Context) string {
	edgercSection := c.String("section")
	if edgercSection == "" {
		return edgegrid.DefaultSection
	}
	return edgercSection
}

func getRetryConfig() (*session.RetryConfig, error) {
	retryDisabled, ok := os.LookupEnv("AKAMAI_RETRY_DISABLED")
	if ok {
		disabled, err := strconv.ParseBool(retryDisabled)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AKAMAI_RETRY_DISABLED environment variable: %w", err)
		}
		if disabled {
			return nil, nil
		}
	}
	conf := session.NewRetryConfig()
	max, ok := os.LookupEnv("AKAMAI_RETRY_MAX")
	if ok {
		v, err := strconv.Atoi(max)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AKAMAI_RETRY_MAX environment variable: %w", err)
		}
		conf.RetryMax = v
	}
	waitMin, ok := os.LookupEnv("AKAMAI_RETRY_WAIT_MIN")
	if ok {
		v, err := strconv.Atoi(waitMin)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AKAMAI_RETRY_WAIT_MIN environment variable: %w", err)
		}
		conf.RetryWaitMin = time.Duration(v) * time.Second
	}
	waitMax, ok := os.LookupEnv("AKAMAI_RETRY_WAIT_MAX")
	if ok {
		v, err := strconv.Atoi(waitMax)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AKAMAI_RETRY_WAIT_MAX environment variable: %w", err)
		}
		conf.RetryWaitMax = time.Duration(v) * time.Second
	}
	excludedEndpoints, ok := os.LookupEnv("AKAMAI_RETRY_EXCLUDED_ENDPOINTS")
	if ok {
		conf.ExcludedEndpoints = strings.Split(excludedEndpoints, ",")
	} else {
		conf.ExcludedEndpoints = append(conf.ExcludedEndpoints, "/identity-management/v3/user-admin/ui-identities/*")
	}

	return &conf, nil
}
