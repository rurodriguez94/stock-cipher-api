package security

import (
	"context"
)

type Service interface {
	EncryptData(ctx context.Context, data interface{}) (Encryption, error)
	DecryptData(ctx context.Context, data, token string) (interface{}, error)
}

type Encryption struct {
	Token   string `json:"token"`
	Payload string `json:"payload"`
}

type DecryptReq struct {
	Payload string `json:"payload"`
}
