package stock

import "errors"

var (
	ErrStockSymbolNotFound  = errors.New("stock symbol not found")
	ErrStockSymbolMalformed = errors.New("stock symbol malformed")
	ErrTokenNotFound        = errors.New("token not found")
	ErrCannotDecryptPayload = errors.New("cannot decrypt payload")
	ErrTimeSeriesNotFound   = errors.New("time series not found")
)
