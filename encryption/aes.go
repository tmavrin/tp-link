package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AES struct {
	Key []byte
	IV  []byte
}

func NewAES(key, iv []byte) *AES {
	return &AES{Key: key, IV: iv}
}

func (a *AES) Encrypt(src string) (string, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCEncrypter(block, []byte(a.IV))

	content := pkcs5Padding([]byte(src), block.BlockSize())
	encrypted := make([]byte, len(content))
	ecb.CryptBlocks(encrypted, content)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (a *AES) Decrypt(src string) (string, error) {
	value, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCDecrypter(block, []byte(a.IV))
	decrypted := make([]byte, len(value))
	ecb.CryptBlocks(decrypted, value)

	return string(pkcs5Trimming(decrypted)), nil
}

func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	paddedText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(src, paddedText...)
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
