package service

import (
	"context"
	"encoding/json"

	"github.com/stock-cipher-api/stock-api/internal/stock"
)

type StockProvider interface {
	TimeSeriesDaily(ctx context.Context, symbol string) ([]stock.TimeSeries, error)
}

type SecurityProvider interface {
	EncryptData(ctx context.Context, data []byte) (stock.EncryptionResponse, error)
	DecryptData(ctx context.Context, data []byte, token string) (stock.TimeSeries, error)
}

type stockService struct {
	stockProvider    StockProvider
	securityProvider SecurityProvider
}

func NewStockService(stockProvider StockProvider, securityProvider SecurityProvider) *stockService {
	return &stockService{stockProvider: stockProvider, securityProvider: securityProvider}
}

// GetStock get stock information from an stock provider API call and encrypt data with encryption-api.
// returns stock payload encrypted with AES-256 CBC and a token that can be used to decrypt the payload.
func (srv stockService) GetStock(ctx context.Context, symbol string) (stock.EncryptionResponse, error) {
	timeSeries, err := srv.stockProvider.TimeSeriesDaily(ctx, symbol)
	if err != nil {
		return stock.EncryptionResponse{}, err
	}

	if len(timeSeries) == 0 {
		return stock.EncryptionResponse{}, stock.ErrTimeSeriesNotFound
	}

	timeSeries = stock.SliceSortedByDate(timeSeries)

	data, err := json.Marshal(&timeSeries[0])
	if err != nil {
		return stock.EncryptionResponse{}, err
	}

	return srv.securityProvider.EncryptData(ctx, data)
}

// DecryptStockData receives stock data encrypted and the token to decrypt with encryption-api.
func (srv stockService) DecryptStockData(ctx context.Context, data, token string) (stock.TimeSeries, error) {
	req := stock.DecryptStockParams{Payload: data}

	bytes, err := json.Marshal(req)
	if err != nil {
		return stock.TimeSeries{}, err
	}

	return srv.securityProvider.DecryptData(ctx, bytes, token)
}
