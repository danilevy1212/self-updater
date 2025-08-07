package digest

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

type fileDigester func(path string) ([]byte, error)

func defaultFileDigester(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file in path `%s`: %w", path, err)
	}
	defer file.Close()

	hasher := sha256.New()
	buf := make([]byte, 4096)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("cannot read file `%s`: %w", path, err)
		}
		if n == 0 {
			break
		}
		hasher.Write(buf[:n])
	}

	return hasher.Sum([]byte{}), nil
}

var DigestFile fileDigester = defaultFileDigester
