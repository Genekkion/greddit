package servicesauth

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"greddit/internal/domains/auth"

	portsauth "greddit/internal/ports/auth"
	dbportsauth "greddit/internal/ports/db/auth"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

// Service is the auth service.
type Service struct {
	logger    *slog.Logger
	jwkSource portsauth.JwkSource
	users     dbportsauth.UsersRepo
}

// NewService creates a new Service.
func NewService(logger *slog.Logger, jwkSource portsauth.JwkSource, users dbportsauth.UsersRepo) Service {
	return Service{
		logger:    logger,
		jwkSource: jwkSource,
		users:     users,
	}
}

const (
	tokenKeyUser = "user"
)

// TokenClaims is the claims for the JWT token.
type TokenClaims struct {
	UserId      auth.UserId
	Username    string
	DisplayName string
	Role        string
}

// Login logs in a user.
func (s Service) Login(ctx context.Context, username string) (signed []byte, err error) {
	user, err := s.users.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "auth.service :: Error getting user by username",
			"error", err,
		)
		return nil, err
	}

	claims := TokenClaims{
		UserId:      user.Id,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        string(user.Role),
	}
	b, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	exp := t.Add(30 * 24 * time.Hour)

	token, err := jwt.NewBuilder().
		IssuedAt(t).
		Expiration(exp).
		Claim(tokenKeyUser, string(b)).
		Build()
	if err != nil {
		s.logger.ErrorContext(ctx, "auth.service :: Error building JWT token",
			"error", err,
		)
		return nil, err
	}

	signed, err = s.jwkSource.Sign(token)
	if err != nil {
		s.logger.ErrorContext(ctx, "auth.service :: Error signing JWT token",
			"error", err,
		)
		return nil, err
	}

	return signed, nil
}

// ExtractClaims extracts the claims from the JWT token.
func (s Service) ExtractClaims(ctx context.Context, token []byte) (claims *TokenClaims, err error) {
	t, err := s.jwkSource.Validate(ctx, token)
	if err != nil {
		s.logger.ErrorContext(ctx, "auth.service :: Error validating JWT token",
			"error", err,
		)
		return nil, err
	}

	var cStr string
	err = t.Get(tokenKeyUser, &cStr)
	if err != nil {
		s.logger.ErrorContext(ctx, "auth.service :: Error extracting claims from JWT token",
			"error", err,
		)
		return nil, err
	}

	claims = &TokenClaims{}
	err = json.Unmarshal([]byte(cStr), claims)
	if err != nil {
		s.logger.ErrorContext(ctx, "auth.service :: Error unmarshalling claims from JWT token",
			"error", err,
		)
		return nil, err
	}

	return claims, nil
}
