package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoToken struct {
	paseto      *paseto.V2
	symetricKey []byte
}

func NewPasetoToken(symetricKey string) (TokenManager, error) {
	if len(symetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d character", chacha20poly1305.KeySize)
	}

	pasetoToken := &PasetoToken{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}

	return pasetoToken, nil
}

// CreateToken implements TokenManager.
func (manage *PasetoToken) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return manage.paseto.Encrypt(manage.symetricKey, payload, nil)
}

// VerifyToken implements TokenManager.
func (manager *PasetoToken) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := manager.paseto.Decrypt(token, manager.symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
