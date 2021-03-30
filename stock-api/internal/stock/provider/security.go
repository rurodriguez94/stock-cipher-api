package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	logs "github.com/stock-cipher-api/stock-api/internal/common"
	"github.com/stock-cipher-api/stock-api/internal/stock"
)

const (
	encryptionURL = "http://encryption-api:8081/security"
)

type securityProvider struct {
	client http.Client
}

func NewSecurityProvider(c http.Client) *securityProvider {
	return &securityProvider{client: c}
}

func (p securityProvider) EncryptData(ctx context.Context, data []byte) (stock.EncryptionResponse, error) {
	logs.Log().Info("EncryptData request init")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, encryptionURL+"/encrypt", bytes.NewBuffer(data))
	if err != nil {
		return stock.EncryptionResponse{}, err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return stock.EncryptionResponse{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return stock.EncryptionResponse{}, err
	}

	var encryptedRes stock.EncryptionResponse
	if err := json.Unmarshal(body, &encryptedRes); err != nil {
		return stock.EncryptionResponse{}, err
	}

	logs.Log().Info("EncryptData request success")

	return encryptedRes, nil
}

func (p securityProvider) DecryptData(ctx context.Context, data []byte, token string) (stock.TimeSeries, error) {
	logs.Log().Info("DecryptData request init")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, encryptionURL+"/decrypt/"+token, bytes.NewBuffer(data))
	if err != nil {
		return stock.TimeSeries{}, err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return stock.TimeSeries{}, err
	}

	if res.StatusCode == http.StatusNotFound {
		logs.Log().Error("token not found")
		return stock.TimeSeries{}, stock.ErrTokenNotFound
	}

	if res.StatusCode == http.StatusUnauthorized {
		logs.Log().Error("an error occurred handling payload encrypted")
		return stock.TimeSeries{}, stock.ErrCannotDecryptPayload
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return stock.TimeSeries{}, err
	}

	var decryptedRes stock.TimeSeries
	if err := json.Unmarshal(body, &decryptedRes); err != nil {
		return stock.TimeSeries{}, err
	}

	logs.Log().Info("DecryptData request success")

	return decryptedRes, nil
}
