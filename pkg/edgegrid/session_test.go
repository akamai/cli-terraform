package edgegrid

import (
	"context"
	"flag"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestInitializeSession(t *testing.T) {
	tests := map[string]struct {
		edgercFile string
		section    string
		withError  bool
	}{
		"valid session initialized": {
			edgercFile: "./testdata/.edgerc",
			section:    "test_section",
		},
		"could not initialize session": {
			edgercFile: "./testdata/edgerc-invalid",
			withError:  true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			app := cli.NewApp()
			set := flag.NewFlagSet("test", 0)
			set.String("edgerc", test.edgercFile, "")
			set.String("section", test.section, "")
			cliCtx := cli.NewContext(app, set, nil)
			s, err := InitializeSession(cliCtx)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, s)
		})
	}
}

func TestWithSession(t *testing.T) {
	ctx := context.Background()
	s, err := session.New()
	require.NoError(t, err)
	ctx = WithSession(ctx, s)
	res := ctx.Value(sessionCtx)
	assert.Equal(t, s, res)
}

func TestGetSession(t *testing.T) {
	t.Run("session in context, should not panic", func(t *testing.T) {
		ctx := context.Background()
		s, err := session.New()
		require.NoError(t, err)
		ctx = WithSession(ctx, s)
		var res session.Session
		assert.NotPanics(t, func() {
			res = GetSession(ctx)
		})
		assert.Equal(t, s, res)
	})
	t.Run("session not in context, should panic", func(t *testing.T) {
		ctx := context.Background()
		assert.Panics(t, func() {
			GetSession(ctx)
		})
	})

}
