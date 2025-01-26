package tests

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"
)

func verifySignature(httpMethod, endpointUrl, lowercaseHash, timeStamp, encodedSignature string, publicKeyPEM string) (bool, error) {
	// Step 1: Build the string to sign
	stringToSign := fmt.Sprintf("%s:%s:%s:%s:", httpMethod, endpointUrl, lowercaseHash, timeStamp)

	// Step 2: Decode the base64-encoded signature
	signature, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	// Step 3: Hash the stringToSign using SHA256
	hashed := sha256.Sum256([]byte(stringToSign))

	// Step 4: Parse the public key from PEM format
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return false, errors.New("failed to parse public key PEM")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("not a valid RSA public key")
	}

	// Step 5: Verify the signature
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, fmt.Errorf("signature verification failed: %w", err)
	}

	return true, nil
}

func TestTOtlah124(t *testing.T) {
	// Example inputs
	httpMethod := "POST"
	endpointUrl := "/api/v1/resource"
	lowercaseHash := "abc123hashedvalue" // Replace with actual hash
	timeStamp := "2025-01-19T03:25:53+07:00"
	encodedSignature := "DW9PY8zxnnPA6lmfvzBUVfJSc8RzhYhFqo4rXhyr9SpFjy5ZzFl7g7ByDnPkfUO6o8FtbAzpY5MAtExx92NrgJKXkiiEHkXN10swopBgCLG894FK8O6idMW2VHUOx" // Replace with actual signature

	// Public key in PEM format (replace with your actual public key)
	publicKeyPEM := `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCDr5Uxwl1a4dtnaOqUpdEsFFw4
ekTtYt9M+t0ycbzZDRNWx8Hfot5Q+s2oU9FW+wOBowHF8AH4/GBo+nQsE7hHnKJq
JkKeUH0WFn8I7jMUdHkevvV7T6kkG+u5NQ+B8eBLRRK/09BqNVxRQeO8o0HMKw9y
5IXC/HlP9hvkhFkftwIDAQAB
-----END PUBLIC KEY-----
`

	// Verify the signature
	isValid, err := verifySignature(httpMethod, endpointUrl, lowercaseHash, timeStamp, encodedSignature, publicKeyPEM)
	if err != nil {
		fmt.Println("Error:", err)
	} else if isValid {
		fmt.Println("Signature is valid!")
	} else {
		fmt.Println("Signature is invalid.")
	}
}
