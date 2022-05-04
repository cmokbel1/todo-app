package crypto

import (
	"crypto/rand"

	"github.com/alexedwards/argon2id"
	"github.com/cmokbel1/todo-app/backend/todo"
)

const alphabet = "`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const passwordByteLen = 16

// RandomString creates a random 16 byte string using crypto/rand.
func RandomString() string {
	bytes := make([]byte, passwordByteLen)

	// Note that err == nil only if we read len(b) bytes.
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	for i, b := range bytes {
		bytes[i] = alphabet[b%byte(len(alphabet))]
	}

	return string(bytes)
}

// CreateHash creates a hash from a string
func CreateHash(pw string) (string, error) {
	hash, err := argon2id.CreateHash(pw, argon2id.DefaultParams)
	if err != nil {
		return "", todo.Err(todo.EINTERNAL, "failed to hash password: %v", err)
	}
	return hash, nil
}

// ComparePasswordAndHash does a constant time comparison between a plaintex password and a hash.
func ComparePasswordAndHash(password string, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return match, todo.Err(todo.EINTERNAL, "failed to compare and hash passwords: %v", err)
	} else if !match {

	}
	return match, nil
}
