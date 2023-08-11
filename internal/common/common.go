package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	config "fourclover.org/internal/config"
)

// File struct
type File struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	Permission   string `json:"permission"`
	LastModified string `json:"last_modified"`

	Blake2b_256 string `json:"blake2b-256,omitempty"`
	Blake2b_512 string `json:"blake2b-512,omitempty"`
	Blake3      string `json:"blake3,omitempty"`
	CRC32       string `json:"crc32,omitempty"`
	MD5         string `json:"md5,omitempty"`
	SHA1        string `json:"sha1,omitempty"`
	SHA3_224    string `json:"sha3-224,omitempty"`
	SHA3_256    string `json:"sha3-256,omitempty"`
	SHA3_384    string `json:"sha3-384,omitempty"`
	SHA3_512    string `json:"sha3-512,omitempty"`
	SHA256      string `json:"sha256"`
	SHA512      string `json:"sha512,omitempty"`
}

// GetConfigDataInBytes converts the config data to bytes
func GetConfigDataInBytes(configData config.YamlConfig) ([]byte, error) {
	configDataInBytes, err := json.Marshal(configData)
	if err != nil {
		return nil, err
	}
	return configDataInBytes, nil
}

// Check if cli argument has certain key. Return true if it has, false if it doesn't.
func CheckCliArg(key string) bool {
	for _, arg := range os.Args {
		if arg == key {
			return true
		}
	}
	return false
}

func Encrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func Decrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = fmt.Errorf("ciphertext too short")
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}

func ValidateDirectory(dir string) error {
	// Check if the directory path is empty
	if dir == "" {
		return fmt.Errorf("directory path is empty")
	}

	// Check if the directory name is "." or ".."
	if dir == "." || dir == ".." {
		return fmt.Errorf("'%s' is not a valid input. Expecting a parent directory", dir)
	}

	// Check if the directory name starts with "./" or ".\"
	if strings.HasPrefix(dir, "./") || strings.HasPrefix(dir, ".\\") {
		dir = dir[2:] // remove the "./" or ".\" prefix
	}

	// Check if the directory name ends with a path separator
	if strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir = dir[:len(dir)-1] // remove the trailing path separator
	}

	// Check if the directory name contains any path separators
	if strings.ContainsAny(dir, `/\`) {
		return fmt.Errorf("'%s' is not a valid input. Expecting a parent directory", dir)
	}

	return nil
}

func ValidateFile(file string) error {
	// Check if the file path is empty
	if file == "" {
		return fmt.Errorf("file path is empty")
	}

	// Check if the file name is "." or ".."
	if file == "." || file == ".." {
		return fmt.Errorf("'%s' is not a valid input. Expecting a file", file)
	}

	return nil
}
