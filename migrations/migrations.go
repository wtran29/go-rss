package migrations

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

func GenerateAPIKey() string {
	keyLength := 32 // You can change the length of the API key as per your requirement

	apiKey := make([]byte, keyLength)
	_, err := rand.Read(apiKey)
	if err != nil {
		log.Fatal("Error generating API key:", err)
	}

	sha256Hash := sha256.Sum256(apiKey)
	return hex.EncodeToString(sha256Hash[:])
}

func generateAddAPIKeySQL() string {
	// Use quote_literal to escape the result of the encode function
	// and ensure proper handling of parentheses
	return fmt.Sprintf("ALTER TABLE users ADD COLUMN api_key VARCHAR(64) DEFAULT %s NOT NULL UNIQUE", GenerateAPIKey())
}
