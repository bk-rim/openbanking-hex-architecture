package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
)

func GenerateIdempotentKey() (string, error) {
	keyLengthStr := os.Getenv("IDEMPOTENT_KEY_LENGTH")
	if keyLengthStr == "" {
		return "", errors.New("IDEMPOTENT_KEY_LENGTH environment variable is not set")
	}

	keyLength, err := strconv.Atoi(keyLengthStr)
	if err != nil {
		return "", err
	}

	randomBytes := make([]byte, keyLength)

	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	key := base64.StdEncoding.EncodeToString(randomBytes)

	key = key[:keyLength]

	key = "JXJ" + key + "XXXZ"

	return key, nil
}
