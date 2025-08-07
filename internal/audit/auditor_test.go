package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danilevy1212/self-updater/internal/audit/fixtures"
)

func Test_verifyED22519Signature(t *testing.T) {
	var publicKey = []byte(releaseFixture.PublicKey)
	var signature = []byte("mOnBz0kNnFWZfe/YyGND9y2C/2J0Z0sI8Y59HSTMP1tyagH+qrF6PzRLc1uKSn+Ks0DmPExYt1/FboiyT2r0DA==")
	var message = []byte("fad0f6e8b138bcf83c041db9ea83427c37f7cafc1efa7528860c68ded97770d1")

	t.Run("should return error if public key is empty", func(t *testing.T) {
		_, err := verifyED22519Signature([]byte{}, message, signature)
		assert.Error(t, err)
		assert.Equal(t, "public key is empty", err.Error())
	})

	t.Run("should return error if signature is empty", func(t *testing.T) {
		_, err := verifyED22519Signature(publicKey, message, []byte{})
		assert.Error(t, err)
		assert.Equal(t, "signature is empty", err.Error())
	})

	t.Run("should only accept valid public keys in pem format", func(t *testing.T) {
		_, err := verifyED22519Signature([]byte("aGVsbG8gd29ybGQgdGhpcyBpcyBub3QgYSB2YWxpZCBwdWJsaWMga2V5"), message, signature)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode PEM block")
	})

	t.Run("should return error if public key is not parseable", func(t *testing.T) {
		badPEM := []byte(`-----BEGIN PUBLIC KEY-----
aGVsbG8gd29ybGQgdGhpcyBpcyBub3QgYSB2YWxpZCBwdWJsaWMga2V5
-----END PUBLIC KEY-----`)

		_, err := verifyED22519Signature(badPEM, message, signature)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse public key")
	})

	t.Run("should only accept ED22519 public keys", func(t *testing.T) {
		_, err := verifyED22519Signature(fixtures.RSAPublicKey, message, signature)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported public key type")
	})

	t.Run("should error if signature is not base64 encoded", func(t *testing.T) {
		_, err := verifyED22519Signature(publicKey, message, []byte("invalidBase64Signature"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode base64 signature")
	})

	t.Run("should error if message is not hex encoded", func(t *testing.T) {
		_, err := verifyED22519Signature(publicKey, []byte("invalidHexMessage"), signature)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode hex digest")
	})

	t.Run("should return true if signature matches message and was signed with correct key", func(t *testing.T) {
		validSignature, err := verifyED22519Signature(publicKey, message, signature)
		assert.NoError(t, err, "should not return an error for valid signature")
		assert.True(t, validSignature, "should return true for valid signature")
	})
}
