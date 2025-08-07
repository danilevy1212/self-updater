package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/danilevy1212/self-updater/internal/digest"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: sign <private.pem> <file-to-sign>")
		os.Exit(1)
	}

	keyPath := os.Args[1]
	filePath := os.Args[2]

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		panic("invalid PEM block")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	privKey := priv.(ed25519.PrivateKey)

	payload, err := digest.DigestFile(filePath)
	if err != nil {
		panic(err)
	}

	sig := ed25519.Sign(privKey, []byte(payload))
	fmt.Println(base64.StdEncoding.EncodeToString(sig))
}
