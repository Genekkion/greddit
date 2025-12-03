package hs256

import "crypto/rand"

func NewSecret() (b []byte, err error) {
	b = make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
