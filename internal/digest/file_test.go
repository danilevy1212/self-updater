package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_defaultFileDigester(t *testing.T) {
	content := []byte("hello, world!")
	expected := sha256.Sum256(content)
	expectedHex := hex.EncodeToString(expected[:])

	tmpFile, err := os.CreateTemp("", "digest-test-*")
	assert.NoError(t, err, "error creating temporary file")
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(content)
	assert.NoError(t, err, "error writing to temporary file")
	assert.NoError(t, tmpFile.Close(), "error closing temporary file")

	d, err := defaultFileDigester(tmpFile.Name())
	assert.NoError(t, err, "error during file digestion")
	assert.Equal(t, expectedHex, d, "hex sha256 values don't match")
}

