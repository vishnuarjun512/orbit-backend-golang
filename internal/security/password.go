package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemory      uint32 = 64 * 1024
	argonIterations  uint32 = 3
	argonParallelism uint8  = 2
	saltLength              = 16
	keyLength               = 32
)

// HashPassword generates a secure Argon2id password hash.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		keyLength,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	encodedPassword := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory,
		argonIterations,
		argonParallelism,
		encodedSalt,
		encodedHash,
	)

	return encodedPassword, nil
}

// VerifyPassword checks a password against its stored Argon2id hash.
func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 ||
		parts[1] != "argon2id" ||
		parts[2] != fmt.Sprintf("v=%d", argon2.Version) ||
		parts[3] != fmt.Sprintf(
			"m=%d,t=%d,p=%d",
			argonMemory,
			argonIterations,
			argonParallelism,
		) {
		return false, fmt.Errorf("invalid password hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) != saltLength {
		return false, fmt.Errorf("invalid password hash salt")
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(expectedHash) != keyLength {
		return false, fmt.Errorf("invalid password hash")
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		keyLength,
	)

	match := subtle.ConstantTimeCompare(expectedHash, actualHash) == 1

	return match, nil
}
