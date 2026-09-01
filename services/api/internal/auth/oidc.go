package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	ErrMissingBearerToken = errors.New("missing bearer token")
	ErrInvalidBearerToken = errors.New("invalid bearer token")
)

type Verifier struct {
	provider *oidc.Provider
	tokens   *oidc.IDTokenVerifier
}

func NewVerifier(ctx context.Context, issuerURL, clientID string) (*Verifier, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}

	return &Verifier{
		provider: provider,
		tokens:   provider.Verifier(&oidc.Config{ClientID: clientID}),
	}, nil
}

func (v *Verifier) Subject(r *http.Request) (string, error) {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if authorization == "" {
		return "", ErrMissingBearerToken
	}

	scheme, rawToken, found := strings.Cut(authorization, " ")
	rawToken = strings.TrimSpace(rawToken)
	if !found || !strings.EqualFold(scheme, "Bearer") || rawToken == "" {
		return "", ErrInvalidBearerToken
	}

	verified, err := v.tokens.Verify(r.Context(), rawToken)
	if err != nil || strings.TrimSpace(verified.Subject) == "" {
		return "", ErrInvalidBearerToken
	}

	userInfo, err := v.provider.UserInfo(
		r.Context(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: rawToken, TokenType: "Bearer"}),
	)
	if err != nil || strings.TrimSpace(userInfo.Subject) == "" || userInfo.Subject != verified.Subject {
		return "", ErrInvalidBearerToken
	}

	return verified.Subject, nil
}
