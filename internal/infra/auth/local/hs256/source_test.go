package hs256

import (
	"crypto/rand"
	"testing"
	"time"

	"greddit/internal/test"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

func newTestSource(t *testing.T) Source {
	t.Helper()

	b := make([]byte, 32)
	_, err := rand.Read(b)
	test.NilErr(t, err)

	source, err := NewSource(b)
	test.NilErr(t, err)

	return *source
}

func TestSource_GetJwkSet(t *testing.T) {
	t.Parallel()

	source := newTestSource(t)
	ctx := t.Context()

	set, err := source.GetJwkSet(ctx)
	test.NilErr(t, err)
	test.Assert(t, "Expected non-nil key set", set != nil)
}

func TestSource_SignValidate(t *testing.T) {
	t.Parallel()

	source := newTestSource(t)

	tn := time.Now()
	te := tn.Add(time.Hour)

	k := "hello"
	v := "world"

	token, err := jwt.NewBuilder().
		IssuedAt(tn).
		Expiration(te).
		Claim(k, v).
		Build()
	test.NilErr(t, err)

	b, err := source.Sign(token)
	test.NilErr(t, err)

	token, err = source.Validate(b)
	test.NilErr(t, err)

	tt, ok := token.IssuedAt()
	test.Assert(t, "Expected issue date to be set", ok)
	test.Assert(t, "Unexpected issue date", tn.Sub(tt) < time.Second)

	tt, ok = token.Expiration()
	test.Assert(t, "Expected expiration date to be set", ok)
	test.Assert(t, "Unexpected expiration date", te.Sub(tt) < time.Second)

	var vv string
	err = token.Get(k, &vv)
	test.NilErr(t, err)
	test.AssertEqual(t, "Unexpected value", v, vv)
}
