package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type fileDigester func(path string) (string, error)

func defaultFileDigester(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("could not open file in path `%s`: %w", path, err)
	}
	defer file.Close()

	hasher := sha256.New()
	buf := make([]byte, 4096)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("cannot read file `%s`: %w", path, err)
		}
		if n == 0 {
			break
		}
		hasher.Write(buf[:n])
	}

	digest := hasher.Sum([]byte{})
	return hex.EncodeToString(digest), nil
}

var DefaultFileDigester fileDigester = defaultFileDigester
