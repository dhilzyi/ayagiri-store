package auth

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/alexedwards/argon2id"
)

func VerifyHashPassword(password, hash string) (bool, error) {
	result, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return result, nil
}

func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil // 64 char hex string
}

func HashPassword(password string) (string, error) {
	params := argon2id.Params{
		Memory:      65536,
		Iterations:  3,
		Parallelism: 1,
		KeyLength:   16,
		SaltLength:  16,
	}
	hashed, err := argon2id.CreateHash(password, &params)
	if err != nil {
		return "", err
	}

	return hashed, nil
}
