package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

var (
	ErrMissingBearerToken = errors.New("missing bearer token")
	ErrInvalidBearerToken = errors.New("invalid bearer token")
)

type Verifier struct {
	idTokens *oidc.IDTokenVerifier
}

func NewVerifier(ctx context.Context, issuerURL, clientID string) (*Verifier, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}

	return &Verifier{
		idTokens: provider.Verifier(&oidc.Config{ClientID: clientID}),
	}, nil
}

func (v *Verifier) Subject(r *http.Request) (string, error) {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if authorization == "" {
		return "", ErrMissingBearerToken
	}

	scheme, rawToken, found := strings.Cut(authorization, " ")
	if !found || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(rawToken) == "" {
		return "", ErrInvalidBearerToken
	}

	idToken, err := v.idTokens.Verify(r.Context(), strings.TrimSpace(rawToken))
	if err != nil {
		return "", ErrInvalidBearerToken
	}
	if strings.TrimSpace(idToken.Subject) == "" {
		return "", ErrInvalidBearerToken
	}

	return idToken.Subject, nil
}
