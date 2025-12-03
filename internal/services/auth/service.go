package servicesauth

import portsauth "greddit/internal/ports/auth"

type Service struct {
	portsauth.JwkSource
}

func NewService(jwkSource portsauth.JwkSource) Service {
	return Service{
		JwkSource: jwkSource,
	}
}
