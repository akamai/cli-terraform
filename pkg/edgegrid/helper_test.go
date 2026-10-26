package edgegrid

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func setupTestCommand(t *testing.T, values map[string]string) *cli.Command {
	t.Helper()

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "edgerc"},
			&cli.StringFlag{Name: "section"},
			&cli.StringFlag{Name: "accountkey"},
		},
	}
	for name, value := range values {
		require.NoError(t, cmd.Set(name, value))
	}

	return cmd
}
