package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gmc, _ := cipher.NewGCM(block)
	nonce := make([]byte, gmc.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gmc.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, _ := aes.NewCipher(key)
	gmc, _ := cipher.NewGCM(block)
	nonceSize := gmc.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, _ := gmc.Open(nil, nonce, ciphertext, nil)
	return plaintext
}
