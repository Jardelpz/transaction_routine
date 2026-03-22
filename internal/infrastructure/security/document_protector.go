package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

type DocumentProtector struct {
	encryptionKey []byte
	hashSalt      string
}

func NewDocumentProtector() (*DocumentProtector, error) {
	encryptionKey := os.Getenv("DOCUMENT_ENCRYPTION_KEY")
	hashSalt := os.Getenv("DOCUMENT_HASH_SALT")
	if len(encryptionKey) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}

	return &DocumentProtector{
		encryptionKey: []byte(encryptionKey),
		hashSalt:      hashSalt,
	}, nil
}

func (p *DocumentProtector) Hash(document string) string {
	sum := sha256.Sum256([]byte(document + p.hashSalt))
	return hex.EncodeToString(sum[:])
}

func (p *DocumentProtector) Encrypt(document string) (string, error) {
	block, err := aes.NewCipher(p.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(document), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (p *DocumentProtector) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(p.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid encrypted payload")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
