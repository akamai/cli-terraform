package edgegrid

import (
	"context"
	"fmt"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/urfave/cli/v2"
)

type ctxType string

var sessionCtx ctxType = "session"

// InitializeSession prepares a session.Session interface based on edgerc config
func InitializeSession(c *cli.Context) (session.Session, error) {
	edgerc, err := GetEdgegridConfig(c)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve edgegrid configuration: %s", err)
	}
	s, err := session.New(
		session.WithSigner(edgerc),
		session.WithHTTPTracing(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED") == "true"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize edgegrid session: %s", err)
	}
	return s, nil
}

// WithSession puts a session.Session in context
func WithSession(ctx context.Context, session session.Session) context.Context {
	return context.WithValue(ctx, sessionCtx, session)
}

// GetSession retrieves a session.Session from context
// It panics if session is not found, as we should ensure that session is always in context - if it is not, then it is an implementation error
func GetSession(ctx context.Context) session.Session {
	s, ok := ctx.Value(sessionCtx).(session.Session)
	if !ok {
		panic("context does not have an edgegrid session")
	}

	return s
}
