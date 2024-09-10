package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func CreateAndEncryptVaultKey(password string) (string, string, string) {
	// Generate a salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)

	}

	// Derive the key using PBKDF2 (could also use bcrypt or Argon2)
	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	// vaultKey for an extra layer of security

	// Encrypt the vault key using AES-GCM
	vaultKey, err := generateVaultKey()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)

	}
	encryptedKey, nonce, err := encryptAESGCM(vaultKey, key)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)

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

func EncryptSecretData(secretName, secretText string, vaultKey []byte) ([2]string, [2]string) {
	encryptedSecretName, nonceSecretName, err := encryptAESGCM([]byte(secretName), vaultKey)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	encryptedSecretText, nonceSecretText, err := encryptAESGCM([]byte(secretText), vaultKey)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)

	}

	encodedEncryptedSecretName := base64.StdEncoding.EncodeToString(encryptedSecretName)
	encodedNonceSecretName := base64.StdEncoding.EncodeToString(nonceSecretName)
	encodedEncryptedSecretText := base64.StdEncoding.EncodeToString(encryptedSecretText)
	encodedNonceSecretText := base64.StdEncoding.EncodeToString(nonceSecretText)

	// [2]{cipher, nonce}
	return [2]string{encodedEncryptedSecretName, encodedNonceSecretName}, [2]string{encodedEncryptedSecretText, encodedNonceSecretText}
}

func DecryptSecretData(encodedEncryptedName, encodedEncryptedText [2]string, vaultKey []byte) (string, string) {
	encodedEncryptedSecretName, encodedNonceSecretName := encodedEncryptedName[0], encodedEncryptedName[1]
	encodedEncryptedSecretText, encodedNonceSecretText := encodedEncryptedText[0], encodedEncryptedText[1]

	decodedEncryptedSecretName, _ := base64.StdEncoding.DecodeString(encodedEncryptedSecretName)
	decodedNonceSecretName, _ := base64.StdEncoding.DecodeString(encodedNonceSecretName)
	decodedEncryptedSecretText, _ := base64.StdEncoding.DecodeString(encodedEncryptedSecretText)
	decodedNonceSecretText, _ := base64.StdEncoding.DecodeString(encodedNonceSecretText)

	decryptedSecretName, err := decryptAESGCM(decodedEncryptedSecretName, vaultKey, decodedNonceSecretName)
	if err != nil {
		log.Fatal(err)
	}
	decryptedSecretText, err := decryptAESGCM(decodedEncryptedSecretText, vaultKey, decodedNonceSecretText)
	if err != nil {
		log.Fatal(err)
	}

	return string(decryptedSecretName[:]), string(decryptedSecretText[:])
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

func generateVaultKey() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes for AES-256
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
