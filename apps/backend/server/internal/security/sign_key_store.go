package security

import (
	"context"
	"time"
)

type SigningKeyStore interface {
	GetActiveKey(
		ctx context.Context,
		tokenType TokenType,
		retainUntil time.Time,
	) (string, *RSAKeyPair, error)
	GetKeyPair(
		ctx context.Context,
		tokenType TokenType,
		kid string,
	) (*RSAKeyPair, error)
}
