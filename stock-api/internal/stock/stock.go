package stock

import (
	"context"
	"sort"
	"time"
)

const TimeLayout = "2006-01-02"

type EncryptionResponse struct {
	Token   string `json:"token"`
	Payload string `json:"payload"`
}

type DecryptStockParams struct {
	Payload string `json:"payload"`
}

type TimeSeries struct {
	Date   time.Time `json:"date"`
	Open   string    `json:"open"`
	High   string    `json:"high"`
	Low    string    `json:"low"`
	Close  string    `json:"close"`
	Volume string    `json:"volume"`
}

// SliceSortedByDate retrieve a slice of TimeSeries sorted descending by date.
func SliceSortedByDate(timeSeries []TimeSeries) []TimeSeries {
	sort.Slice(timeSeries, func(i, j int) bool {
		return timeSeries[j].Date.Before(timeSeries[i].Date)
	})
	return timeSeries
}

type Service interface {
	GetStock(ctx context.Context, symbol string) (EncryptionResponse, error)
	DecryptStockData(ctx context.Context, data, token string) (TimeSeries, error)
}
