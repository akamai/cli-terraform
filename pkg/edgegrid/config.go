package edgegrid

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
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
