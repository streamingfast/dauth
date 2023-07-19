package trust

import (
	"context"
	"testing"

	"github.com/streamingfast/dauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	Register()
	a, err := dauth.New("trust://?allowed=x-sf-something,x-sf-SoMeThingElse")
	require.NoError(t, err)

	p := a.(*trustPlugin)
	assert.Equal(t, map[string]bool{
		"x-sf-something":     true,
		"x-sf-somethingelse": true,
	},
		p.allowed)

}

func TestAuthenticateEmpty(t *testing.T) {

	p := &trustPlugin{}

	ctx, err := p.Authenticate(context.Background(), "", map[string][]string{
		"x-sf-something":     []string{"someval", "ignored"},
		"x-sf-somethingelse": []string{"someotherval", "ignored"},
	}, "10.0.0.1")
	require.NoError(t, err)

	auth := dauth.FromContext(ctx)
	assert.Equal(t, "someval", auth.Get("x-sf-something"))
	assert.Equal(t, "someotherval", auth.Get("x-sf-somethingelse"))
}

func TestAuthenticateAllowed(t *testing.T) {

	p := &trustPlugin{
		allowed: map[string]bool{
			"x-sf-something":     true,
			"x-sf-somethingelse": true,
		},
	}

	ctx, err := p.Authenticate(context.Background(), "", map[string][]string{
		"x-sf-something": []string{"someval", "ignored"},
		"x-sf-forbidden": []string{"forbiddenval", "ignored"},
	}, "10.0.0.1")
	require.NoError(t, err)

	auth := dauth.FromContext(ctx)
	assert.Equal(t, "someval", auth.Get("x-sf-something"))
	assert.Equal(t, "", auth.Get("x-sf-somethingelse"))
	assert.Equal(t, "", auth.Get("x-sf-forbidden"))
}
