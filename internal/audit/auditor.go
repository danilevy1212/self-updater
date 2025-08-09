package audit

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

type VerifySignatureFunc func(publicKey []byte, messageHex, signature string) (bool, error)

func verifyED22519Signature(publicKey []byte, messageHex, signature string) (bool, error) {
	if len(publicKey) == 0 {
		return false, errors.New("public key is empty")
	}

	if len(signature) == 0 {
		return false, errors.New("signature is empty")
	}

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return false, errors.New("failed to decode PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	rawSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode base64 signature: %w", err)
	}

	messageBytes, err := hex.DecodeString(messageHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode hex digest: %w", err)
	}

	switch pub := pubKey.(type) {
	case ed25519.PublicKey:
		return ed25519.Verify(pub, messageBytes, rawSignature), nil
	default:
		return false, fmt.Errorf("unsupported public key type: %T", pubKey)
	}
}

var VerifySignature VerifySignatureFunc = verifyED22519Signature
