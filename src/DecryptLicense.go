package main

import (
	"crypto/aes"
	"encoding/base64"
)

const (
	blockSize = 16
)

var aesxor = []byte{5, 29, 252, 145, 197, 55, 82, 1, 188, 93, 221, 222, 110, 39, 134, 162}

func base64Decode(input string) []byte {
	data, _ := base64.StdEncoding.DecodeString(input)
	return data
}

func decryptAes128Ecb(data, key []byte) []byte {
	cipher, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))

	for bs, be := 0, blockSize; bs < len(data); bs, be = bs+blockSize, be+blockSize {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

func genAesKey(key string) []byte {
	for len(key) < blockSize {
		key += key
	}

	aeskey := make([]byte, blockSize)
	for i := 0; i < blockSize; i++ {
		aeskey[i] = aesxor[i] ^ key[i]
	}

	return aeskey
}

// LicenseDecrypt : license tier decryptor
func LicenseDecrypt(encrypted string, key string) []byte {
	aeskey := genAesKey(key)
	data := base64Decode(encrypted)
	decrypted := decryptAes128Ecb(data, aeskey)

	return decrypted
}
