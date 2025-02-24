package vspc_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.stellar.af/go-utils/environment"
	"go.stellar.af/go-vspc"
)

type EnvType struct {
	APIKey string `env:"VSPC_API_KEY"`
	URL    string `env:"VSPC_URL"`
}

var Env EnvType

func init() {
	err := environment.Load(&Env, &environment.EnvironmentOptions{
		DotEnv: true,
	})
	if err != nil {
		panic(err)
	}
}

func TestClient(t *testing.T) {
	client, err := vspc.New(Env.URL, Env.APIKey)
	require.NoError(t, err)
	ctx := context.Background()
	res, err := client.GetAboutInformation(ctx, &vspc.GetAboutInformationParams{})
	require.NoError(t, err)
	about, err := vspc.ParseGetAboutInformationRes(res)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, about.StatusCode())
	require.NotNil(t, about.JSON200)
	require.NotNil(t, about.JSON200.Data)
	ver := *about.JSON200.Data.ServerVersion
	assert.NotEmpty(t, ver)
	t.Logf("version=%s", ver)
}
