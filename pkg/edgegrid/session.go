package edgegrid

import (
	"context"
	"fmt"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/akamai/cli/pkg/log"
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

	retryConfig, err := getRetryConfig()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve retry configuration: %w", err)
	}

	options := []session.Option{
		session.WithSigner(edgerc),
		session.WithHTTPTracing(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED") == "true"),
		session.WithLog(log.FromContext(c.Context)),
	}
	if retryConfig != nil {
		// Exclude get user endpoint from retries as it may fail for some users returned by list users endpoint during `export-iam all` command.
		retryConfig.ExcludedEndpoints = []string{"/identity-management/v3/user-admin/ui-identities/*"}
		options = append(options, session.WithRetries(*retryConfig))
	}

	s, err := session.New(options...)
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
