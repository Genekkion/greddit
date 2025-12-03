package servicesauth

import (
	"context"
	"time"

	"greddit/internal/domains/auth"

	portsauth "greddit/internal/ports/auth"
	dbportsauth "greddit/internal/ports/db/auth"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

const (
	tokenKeyUser = "user"
)

type Service struct {
	jwkSource portsauth.JwkSource
	users     dbportsauth.UsersRepo
}

func NewService(jwkSource portsauth.JwkSource, users dbportsauth.UsersRepo) Service {
	return Service{
		jwkSource: jwkSource,
		users:     users,
	}
}

type TokenClaims struct {
	UserId      auth.UserId
	Username    string
	DisplayName string
	Role        auth.Role
}

func (s Service) Login(ctx context.Context, username string) (signed []byte, err error) {
	user, err := s.users.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	claims := TokenClaims{
		UserId:      user.Id,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role,
	}

	t := time.Now()
	exp := t.Add(30 * 24 * time.Hour)

	token, err := jwt.NewBuilder().
		IssuedAt(t).
		Expiration(exp).
		Claim(tokenKeyUser, claims).
		Build()
	if err != nil {
		return nil, err
	}

	signed, err = s.jwkSource.Sign(token)
	if err != nil {
		return nil, err
	}

	return signed, nil
}

func (s Service) ExtractClaims(ctx context.Context, token []byte) (claims *TokenClaims, err error) {
	t, err := s.jwkSource.Validate(ctx, token)
	if err != nil {
		return nil, err
	}

	claims = &TokenClaims{}
	err = t.Get(tokenKeyUser, &claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
