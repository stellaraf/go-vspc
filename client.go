package vspc

import (
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
)

func New(url, apiKey string) (*VeeamSPC, error) {
	auth, err := securityprovider.NewSecurityProviderBearerToken(apiKey)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	if !strings.HasSuffix(url, "/api/v3") {
		url += "/api/v3"
	}
	client, err := NewClient(url, WithRequestEditorFn(auth.Intercept))
	if err != nil {
		return nil, err
	}
	return client, nil
}
