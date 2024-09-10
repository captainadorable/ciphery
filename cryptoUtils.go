package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"

	"golang.org/x/crypto/pbkdf2"
)

func CreateAndEncryptVaultKey(password, vaultKey string) (string, string, string) {
	// Generate a salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal(err)
	}

	// Derive the key using PBKDF2 (could also use bcrypt or Argon2)
	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	// vaultKey for an extra layer of security

	// Encrypt the vault key using AES-GCM
	encryptedKey, nonce, err := encryptAESGCM([]byte(vaultKey), key)
	if err != nil {
		log.Fatal(err)
	}

	// Base64 encode everything for safe storage in JSON
	encodedEncryptedVaultKey := base64.StdEncoding.EncodeToString(encryptedKey)
	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	encodedNonce := base64.StdEncoding.EncodeToString(nonce)

	return encodedEncryptedVaultKey, encodedSalt, encodedNonce
}

func DecryptVaultKeyFromPassword(password, encodedSalt, encodedEncryptedVaultKey, encodedNonce string) ([]byte, bool) {
	// check if the given master password can decrypt vaultkey. if it can't return false meaning wrong master password.
	decodedSalt, _ := base64.StdEncoding.DecodeString(encodedSalt)
	decodedEncryptedVaultKey, _ := base64.StdEncoding.DecodeString(encodedEncryptedVaultKey)
	decodedNonce, _ := base64.StdEncoding.DecodeString(encodedNonce)

	derivedKey := pbkdf2.Key([]byte(password), decodedSalt, 4096, 32, sha256.New)

	decryptedVaultKey, err := decryptAESGCM(decodedEncryptedVaultKey, derivedKey, decodedNonce)
	auth := true
	if err != nil {
		auth = false
	}

	return decryptedVaultKey, auth
}

func encryptAESGCM(plaintext, key []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}
func decryptAESGCM(ciphertext, key, nonce []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}
