package service

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"time"

	logs "github.com/stock-cipher-api/encryption-api/internal/common"
	"github.com/stock-cipher-api/encryption-api/internal/security"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

const keySize = 16

type securityService struct {
	TokenCache *cache.Cache
}

func NewSecurityService(c *cache.Cache) *securityService {
	return &securityService{TokenCache: c}
}

func (srv securityService) EncryptData(_ context.Context, data interface{}) (security.Encryption, error) {
	logs.Log().Info("encrypt data init")

	plaintext, err := json.Marshal(data)
	if err != nil {
		return security.Encryption{}, err
	}

	logs.Log().Info("creating new key")

	key, err := NewKey(keySize)
	if err != nil {
		return security.Encryption{}, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return security.Encryption{}, err
	}

	plaintext = PKCS5Padding(plaintext, block.BlockSize())

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return security.Encryption{}, err
	}

	logs.Log().Info("data encryption init")

	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	logs.Log().Info("data encryption success")

	token := uuid.New().String()

	_ = srv.TokenCache.Add(token, hex.EncodeToString(key), time.Minute*5)

	res := security.Encryption{
		Token:   token,
		Payload: hex.EncodeToString(ciphertext),
	}

	logs.Log().Info("encryption response", zap.String("token", res.Token), zap.String("payload", res.Payload))

	return res, nil
}

func (srv securityService) DecryptData(_ context.Context, data, token string) (interface{}, error) {
	value, found := srv.TokenCache.Get(token)
	if !found {
		return nil, security.ErrTokenNotFound
	}

	ciphertext, err := hex.DecodeString(data)
	if err != nil {
		return security.Encryption{}, err
	}

	key, err := hex.DecodeString(value.(string))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	plaintext := make([]byte, len(ciphertext))

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(plaintext, ciphertext)

	plaintext, err = PKCS5UnPadding(plaintext)
	if err != nil {
		return nil, err
	}

	var res json.RawMessage
	err = json.Unmarshal(plaintext, &res)
	if err != nil {
		return nil, err
	}

	srv.TokenCache.Delete(token)

	return res, nil
}

// NewKey creates a new random key of the given size.
func NewKey(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := io.ReadAtLeast(rand.Reader, b, size)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// PKCS5Padding paddings the src in relation to the block size.
func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// PKCS5Padding removes padding from src.
func PKCS5UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)], nil
}
