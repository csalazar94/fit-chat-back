package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	time    = 1
	saltLen = 32
	memory  = 64 * 1024
	threads = 32
	keyLen  = 256
)

type ArgonParams struct {
	memory  uint32
	time    uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func randomBytes(saltLen uint32) ([]byte, error) {
	secret := make([]byte, saltLen)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func GenerateHash(password string, salt []byte, params *ArgonParams) (encodedHash string, err error) {
	if params == nil {
		params = &ArgonParams{
			memory:  memory,
			time:    time,
			threads: threads,
			keyLen:  keyLen,
			saltLen: saltLen,
		}
	}

	if len(salt) == 0 {
		salt, err = randomBytes(params.saltLen)
	}
	if err != nil {
		log.Printf("Error al generar el salt: %v", err)
		return encodedHash, err
	}

	hash := argon2.IDKey([]byte(password), salt, params.time, params.memory, params.threads, params.keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.memory, params.time, params.threads, b64Salt, b64Hash)

	return encodedHash, err
}

func CompareHash(encodedHash, password string) (bool, error) {
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, params.time, params.memory, params.threads, params.keyLen)
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (params *ArgonParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("the encoded hash is not in the correct format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version of argon2")
	}

	params = &ArgonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.memory, &params.time, &params.threads)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.saltLen = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.keyLen = uint32(len(hash))

	return params, salt, hash, nil
}
