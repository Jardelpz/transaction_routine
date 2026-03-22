package security

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testKey32 = "12345678901234567890123456789012"
const testSalt = "test-salt"

func setupProtector(t *testing.T) *DocumentProtector {
	t.Helper()
	t.Setenv("DOCUMENT_ENCRYPTION_KEY", testKey32)
	t.Setenv("DOCUMENT_HASH_SALT", testSalt)
	p, err := NewDocumentProtector()
	require.NoError(t, err)
	return p
}

func TestNewDocumentProtector(t *testing.T) {
	t.Run("success with valid 32-byte key", func(t *testing.T) {
		t.Setenv("DOCUMENT_ENCRYPTION_KEY", testKey32)
		t.Setenv("DOCUMENT_HASH_SALT", testSalt)

		p, err := NewDocumentProtector()
		require.NoError(t, err)
		assert.NotNil(t, p)
	})

	t.Run("error when key is not 32 bytes - too short", func(t *testing.T) {
		t.Setenv("DOCUMENT_ENCRYPTION_KEY", "short")
		t.Setenv("DOCUMENT_HASH_SALT", testSalt)

		p, err := NewDocumentProtector()
		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Contains(t, err.Error(), "32 bytes")
	})

	t.Run("error when key is not 32 bytes - too long", func(t *testing.T) {
		t.Setenv("DOCUMENT_ENCRYPTION_KEY", strings.Repeat("a", 40))
		t.Setenv("DOCUMENT_HASH_SALT", testSalt)

		p, err := NewDocumentProtector()
		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestDocumentProtector_Hash(t *testing.T) {
	p := setupProtector(t)

	t.Run("deterministic - same input produces same hash", func(t *testing.T) {
		doc := "12345678901"
		h1 := p.Hash(doc)
		h2 := p.Hash(doc)
		assert.Equal(t, h1, h2)
	})

	t.Run("hash is 64 hex characters", func(t *testing.T) {
		h := p.Hash("12345678901")
		assert.Len(t, h, 64)
		for _, c := range h {
			assert.Contains(t, "0123456789abcdef", string(c))
		}
	})

	t.Run("different documents produce different hashes", func(t *testing.T) {
		h1 := p.Hash("12345678901")
		h2 := p.Hash("12345678902")
		assert.NotEqual(t, h1, h2)
	})

	t.Run("different salt produces different hash", func(t *testing.T) {
		t.Setenv("DOCUMENT_ENCRYPTION_KEY", testKey32)
		t.Setenv("DOCUMENT_HASH_SALT", "salt-a")
		p1, _ := NewDocumentProtector()

		t.Setenv("DOCUMENT_HASH_SALT", "salt-b")
		p2, _ := NewDocumentProtector()

		h1 := p1.Hash("12345678901")
		h2 := p2.Hash("12345678901")
		assert.NotEqual(t, h1, h2)
	})

	t.Run("empty string produces valid hash", func(t *testing.T) {
		h := p.Hash("")
		assert.Len(t, h, 64)
	})
}

func TestDocumentProtector_EncryptDecrypt(t *testing.T) {
	p := setupProtector(t)

	t.Run("encrypt then decrypt returns original", func(t *testing.T) {
		original := "12345678901"
		encrypted, err := p.Encrypt(original)
		require.NoError(t, err)

		decrypted, err := p.Decrypt(encrypted)
		require.NoError(t, err)
		assert.Equal(t, original, decrypted)
	})

	t.Run("same plaintext produces different ciphertext each time", func(t *testing.T) {
		plain := "12345678901"
		enc1, err1 := p.Encrypt(plain)
		enc2, err2 := p.Encrypt(plain)
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, enc1, enc2)
	})

	t.Run("decrypt with wrong key fails", func(t *testing.T) {
		encrypted, err := p.Encrypt("12345678901")
		require.NoError(t, err)

		t.Setenv("DOCUMENT_ENCRYPTION_KEY", "09876543210987654321098765432109")
		pWrong, _ := NewDocumentProtector()

		_, err = pWrong.Decrypt(encrypted)
		assert.Error(t, err)
	})

	t.Run("decrypt with invalid base64 fails", func(t *testing.T) {
		_, err := p.Decrypt("not-valid-base64!!!")
		assert.Error(t, err)
	})

	t.Run("decrypt with tampered ciphertext fails", func(t *testing.T) {
		encrypted, err := p.Encrypt("12345678901")
		require.NoError(t, err)

		data, _ := base64.StdEncoding.DecodeString(encrypted)
		data[len(data)-1] ^= 0xFF // flip last byte
		tampered := base64.StdEncoding.EncodeToString(data)

		_, err = p.Decrypt(tampered)
		assert.Error(t, err)
	})

	t.Run("decrypt with too short payload fails", func(t *testing.T) {
		shortPayload := base64.StdEncoding.EncodeToString([]byte("x"))
		_, err := p.Decrypt(shortPayload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid encrypted payload")
	})

	t.Run("empty string encrypt and decrypt", func(t *testing.T) {
		encrypted, err := p.Encrypt("")
		require.NoError(t, err)
		decrypted, err := p.Decrypt(encrypted)
		require.NoError(t, err)
		assert.Equal(t, "", decrypted)
	})

	t.Run("encryption produces valid base64 output", func(t *testing.T) {
		encrypted, err := p.Encrypt("12345678901")
		require.NoError(t, err)
		_, err = base64.StdEncoding.DecodeString(encrypted)
		assert.NoError(t, err)
	})
}

func TestDocumentProtector_Encrypt_WithDifferentKeys(t *testing.T) {
	t.Setenv("DOCUMENT_ENCRYPTION_KEY", testKey32)
	t.Setenv("DOCUMENT_HASH_SALT", testSalt)
	p1, _ := NewDocumentProtector()

	t.Setenv("DOCUMENT_ENCRYPTION_KEY", "09876543210987654321098765432109")
	p2, _ := NewDocumentProtector()

	plain := "12345678901"
	enc1, _ := p1.Encrypt(plain)
	enc2, _ := p2.Encrypt(plain)

	assert.NotEqual(t, enc1, enc2)

	dec1, _ := p1.Decrypt(enc1)
	dec2, _ := p2.Decrypt(enc2)
	assert.Equal(t, plain, dec1)
	assert.Equal(t, plain, dec2)
}
