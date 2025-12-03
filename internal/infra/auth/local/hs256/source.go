package hs256

import (
	"context"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

// Source implements the portsauth.JwkSource interface.
type Source struct {
	secret []byte
	jwkSet jwk.Set
	key    jwk.Key
}

// NewSource creates a new Source.
func NewSource(secret []byte) (source *Source, err error) {
	key, err := jwk.Import(secret)
	if err != nil {
		return nil, err
	}

	key.Set(jwk.KeyIDKey, "local-hs256")
	key.Set(jwk.AlgorithmKey, "HS256")

	set := jwk.NewSet()
	err = set.AddKey(key)
	if err != nil {
		return nil, err
	}

	return &Source{
		secret: secret,
		jwkSet: set,
		key:    key,
	}, nil
}

func (s Source) Sign(token jwt.Token) (signed []byte, err error) {
	return jwt.Sign(token, jwt.WithKey(jwa.HS256(), s.key))
}

func (s Source) GetJwkSet(ctx context.Context) (set jwk.Set, err error) {
	return s.jwkSet, nil
}

func (s Source) Validate(signed []byte) (token jwt.Token, err error) {
	return jwt.Parse(signed, jwt.WithKeySet(s.jwkSet))
}
