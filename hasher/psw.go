package hasher

import (
	"golang.org/x/crypto/argon2"
)

func GenerateHash(psw, salt []byte) []byte {
	return argon2.IDKey([]byte(psw), []byte(salt), 1, 64*1024, 4, 32)
}
