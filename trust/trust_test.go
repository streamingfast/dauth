package trust

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/streamingfast/dauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	Register()
	a, err := dauth.New("trust://?allowed=x-something,x-SoMeThingElse", zap.NewNop())
	require.NoError(t, err)

	p := a.(*trustPlugin)
	assert.Equal(t, map[string]bool{
		"x-something":     true,
		"x-somethingelse": true,
	},
		p.allowed)

}

func TestAuthenticateEmpty(t *testing.T) {

	p := &trustPlugin{}

	ctx, err := p.Authenticate(context.Background(), "", map[string][]string{
		"x-something":     []string{"someval", "ignored"},
		"x-somethingelse": []string{"someotherval", "ignored"},
	}, "10.0.0.1")
	require.NoError(t, err)

	auth := dauth.FromContext(ctx)
	assert.Equal(t, "someval", auth.Get("x-something"))
	assert.Equal(t, "someotherval", auth.Get("x-somethingelse"))
}

func TestAuthenticateAllowed(t *testing.T) {

	p := &trustPlugin{
		allowed: map[string]bool{
			"x-something":     true,
			"x-somethingelse": true,
		},
	}

	ctx, err := p.Authenticate(context.Background(), "", map[string][]string{
		"x-something": []string{"someval", "ignored"},
		"x-forbidden": []string{"forbiddenval", "ignored"},
	}, "10.0.0.1")
	require.NoError(t, err)

	auth := dauth.FromContext(ctx)
	assert.Equal(t, "someval", auth.Get("x-something"))
	assert.Equal(t, "", auth.Get("x-somethingelse"))
	assert.Equal(t, "", auth.Get("x-forbidden"))
}
