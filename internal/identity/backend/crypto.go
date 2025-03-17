package backend

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

func randomBytes(n uint32) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func hashPassword(password string) (string, error) {
	// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
	var (
		t       uint32 = 2
		m       uint32 = 19 * 1024 // 19MB
		p       uint8  = 1
		keyLen  uint32 = 32
		saltLen uint32 = 16
	)

	salt, err := randomBytes(saltLen)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}

	hash := argon2.IDKey([]byte(password), salt, t, m, p, keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", m, t, p, b64Salt, b64Hash)
	return encoded, nil
}

func verifyPassword(password string, encodedHash string) (bool, error) {
	// https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	// parts[1] contains the algorithm
	if parts[1] != "argon2id" {
		return false, errors.New("invalid hash algorithm")
	}

	// parts[2] is version info, e.g., "v=19"
	if parts[2] != "v=19" {
		return false, errors.New("invalid hash version")
	}

	// parts[3] contains the parameters in the form "m=19456,t=1,p=4"
	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return false, errors.New("invalid parameter format")
	}

	var (
		m uint32
		t uint32
		p uint8
	)

	for _, param := range params {
		kv := strings.Split(param, "=")
		if len(kv) != 2 {
			return false, errors.New("invalid parameter format")
		}
		key, value := kv[0], kv[1]
		switch key {
		case "m":
			m64, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return false, errors.New("failed to parse memory")
			}
			m = uint32(m64)
		case "t":
			t64, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return false, err
			}
			t = uint32(t64)
		case "p":
			p64, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return false, err
			}
			p = uint8(p64)
		default:
			return false, errors.New("unknown parameter key")
		}
	}

	// Decode the Base64 encoded salt and hash.
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	keyLen := uint32(len(hash))

	// Compute the hash with the provided password using the extracted parameters and salt.
	computedHash := argon2.IDKey([]byte(password), salt, t, m, p, keyLen)

	// Compare the computed hash with the stored hash using constant-time comparison.
	if subtle.ConstantTimeCompare(hash, computedHash) == 1 {
		return true, nil
	}

	return false, nil
}
