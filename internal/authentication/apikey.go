package auth

import (
	"auditlog/internal/config"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type APIKey struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

func NewAPIKey() (*APIKey, error) {
	key := &APIKey{
		Expiry: time.Now().Add(config.KeyTTL),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	key.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	key.Hash = HashKey(key.Plaintext)

	return key, err
}

func HashKey(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}
