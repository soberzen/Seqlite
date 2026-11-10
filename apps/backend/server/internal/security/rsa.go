package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type RSAKeyPair struct {
	PrivatePEM string
	PublicPEM  string
}

func GenerateRSAPair() (*RSAKeyPair, error) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generate rsa key failed: %w", err)
	}

	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private),
	})
	publicBytes, err := x509.MarshalPKIXPublicKey(&private.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key failed: %w", err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	})
	return &RSAKeyPair{
		PrivatePEM: string(privatePEM),
		PublicPEM:  string(publicPEM),
	}, nil
}
