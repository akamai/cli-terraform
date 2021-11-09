package tools

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/urfave/cli"
)

// GetEdgegridConfig gets configuration from .edgerc file
func GetEdgegridConfig(c *cli.Context) (*edgegrid.Config, error) {

	edgercOps := []edgegrid.Option{
		edgegrid.WithEnv(true),
		edgegrid.WithFile(getEdgercPath(c)),
		edgegrid.WithSection(getEdgercSection(c)),
	}
	config, err := edgegrid.New(edgercOps...)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getEdgercPath(c *cli.Context) string {
	edgercPath := c.GlobalString("edgerc")
	if edgercPath == "" {
		return edgegrid.DefaultConfigFile
	}
	return edgercPath
}

func getEdgercSection(c *cli.Context) string {
	edgercSection := c.GlobalString("section")
	if edgercSection == "" {
		return edgegrid.DefaultSection
	}
	return edgercSection
}
