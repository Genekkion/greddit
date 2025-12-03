package portsauth

import (
	"context"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

// JwkSource is a source of JWKs.
type JwkSource interface {
	// Sign signs a token.
	Sign(token jwt.Token) (signed []byte, err error)

	// GetJwkSet returns the JWK set.
	GetJwkSet(ctx context.Context) (set jwk.Set, err error)

	// Validate validates a token.
	Validate(signed []byte) (token jwt.Token, err error)
}
