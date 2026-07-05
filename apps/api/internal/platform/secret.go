package platform

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

// GenerateToken renvoie un token de session crypto-aléatoire (URL-safe, 256 bits).
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

const base62Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GeneratePassword renvoie un mot de passe de 24 caractères base62 (~143 bits).
func GeneratePassword() (string, error) {
	const length = 24
	out := make([]byte, length)
	max := big.NewInt(int64(len(base62Alphabet)))
	for i := range out {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		out[i] = base62Alphabet[idx.Int64()]
	}
	return string(out), nil
}
