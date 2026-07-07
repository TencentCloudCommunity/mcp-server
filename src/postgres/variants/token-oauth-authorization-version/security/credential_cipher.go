package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

const CredentialEncryptionKeyEnv = "CREDENTIAL_ENCRYPTION_KEY"

type CredentialCipher struct {
	key []byte
}

func NewCredentialCipherFromEnv() (*CredentialCipher, error) {
	raw := strings.TrimSpace(os.Getenv(CredentialEncryptionKeyEnv))
	if raw == "" {
		return nil, fmt.Errorf("missing %s", CredentialEncryptionKeyEnv)
	}
	key, err := parseCredentialEncryptionKey(raw)
	if err != nil {
		return nil, err
	}
	return &CredentialCipher{key: key}, nil
}

func (c *CredentialCipher) EncryptString(plaintext string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("credential cipher is not initialized")
	}
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("create aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm cipher: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("read nonce: %w", err)
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

func (c *CredentialCipher) DecryptString(ciphertext string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("credential cipher is not initialized")
	}
	decoded, err := decodeCredentialBytes(strings.TrimSpace(ciphertext))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("create aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm cipher: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return "", fmt.Errorf("ciphertext is too short")
	}
	plaintext, err := gcm.Open(nil, decoded[:nonceSize], decoded[nonceSize:], nil)
	if err != nil {
		return "", fmt.Errorf("decrypt credential payload: %w", err)
	}
	return string(plaintext), nil
}

func parseCredentialEncryptionKey(raw string) ([]byte, error) {
	for _, decoder := range []func(string) ([]byte, error){
		base64.RawURLEncoding.DecodeString,
		base64.URLEncoding.DecodeString,
		base64.RawStdEncoding.DecodeString,
		base64.StdEncoding.DecodeString,
	} {
		decoded, err := decoder(raw)
		if err == nil && validAESKeySize(len(decoded)) {
			return decoded, nil
		}
	}
	if validAESKeySize(len(raw)) {
		return []byte(raw), nil
	}
	return nil, fmt.Errorf("%s must be a 16/24/32-byte raw string or base64-encoded AES key", CredentialEncryptionKeyEnv)
}

func decodeCredentialBytes(ciphertext string) ([]byte, error) {
	for _, decoder := range []func(string) ([]byte, error){
		base64.RawURLEncoding.DecodeString,
		base64.URLEncoding.DecodeString,
		base64.RawStdEncoding.DecodeString,
		base64.StdEncoding.DecodeString,
	} {
		decoded, err := decoder(ciphertext)
		if err == nil {
			return decoded, nil
		}
	}
	return nil, fmt.Errorf("decode encrypted credential payload failed")
}

func validAESKeySize(size int) bool {
	switch size {
	case 16, 24, 32:
		return true
	default:
		return false
	}
}
