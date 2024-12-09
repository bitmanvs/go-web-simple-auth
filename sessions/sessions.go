package sessions

import (
	"errors"
)

type (
	Key struct {
		key []byte
	}

	Session []byte
)

var ErrNoSessionFound = errors.New("no active session in request context")

func KeyFromBytes(b []byte) Key {
	return Key{
		key: b,
	}
}

func (k Key) Bytes() []byte {
	return k.key
}
