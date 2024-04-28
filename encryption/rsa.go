package encryption

import (
	"crypto/rsa"
	"math/big"
)

type RSA struct {
	publicKey *rsa.PublicKey
}

func NewRSA(key *rsa.PublicKey) *RSA {
	return &RSA{publicKey: key}
}

func (r *RSA) Encrypt(value string) string {
	var (
		buffer string

		step  = 64
		start = 0
		end   = step
	)

	for start < len(value) {
		if end >= len(value) {
			end = len(value)
		}

		currentSubstring := value[start:end]

		valueBytes := []byte(currentSubstring)
		for i := len(currentSubstring); i < 64; i++ {
			valueBytes = append(valueBytes, byte(int(0)))
		}

		var valueBigInt big.Int
		encryptedResult := valueBigInt.Exp(valueBigInt.SetBytes(valueBytes), big.NewInt(int64(r.publicKey.E)), r.publicKey.N)
		resultInHex := encryptedResult.Text(16)

		if len(resultInHex)%2 == 1 {
			resultInHex = "0" + resultInHex
		}

		if len(resultInHex) != 128 {
			for i := 0; i < 128-len(resultInHex); i++ {
				resultInHex = "0" + resultInHex
			}
		}

		buffer += resultInHex

		start += step
		end += step
	}

	return buffer
}
