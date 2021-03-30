package provider

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	logs "github.com/stock-cipher-api/stock-api/internal/common"
	"github.com/stock-cipher-api/stock-api/internal/stock"

	"go.uber.org/zap"
)

const (
	alphaURL = "https://www.alphavantage.co/query"
)

type stockProvider struct {
	APIKey string
}

func NewStockProvider(apikey string) *stockProvider {
	return &stockProvider{APIKey: apikey}
}

// TimeSeriesDaily returns stock time series daily from alphavantage.
func (p stockProvider) TimeSeriesDaily(ctx context.Context, symbol string) ([]stock.TimeSeries, error) {
	logs.Log().Info("TimeSeriesDaily request init")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, alphaURL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("function", "TIME_SERIES_DAILY")
	query.Add("symbol", symbol)
	query.Add("apikey", p.APIKey)

	req.URL.RawQuery = query.Encode()

	c := http.Client{}

	res, err := c.Do(req)
	if err != nil {
		logs.Log().Error("TimeSeriesDaily request failed")
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		logs.Log().Error("TimeSeriesDaily request failed", zap.Error(stock.ErrStockSymbolNotFound))
		return nil, stock.ErrStockSymbolNotFound
	}

	if res.StatusCode == http.StatusBadRequest {
		logs.Log().Error("TimeSeriesDaily request failed", zap.Error(stock.ErrStockSymbolMalformed))
		return nil, stock.ErrStockSymbolMalformed
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil
	}

	var ts TimeSeriesRes
	if err := json.Unmarshal(body, &ts); err != nil {
		return nil, err
	}

	var stockTimeSeries []stock.TimeSeries
	for date, timeSerie := range ts.TimeSeries {
		t, err := time.Parse(stock.TimeLayout, date)
		if err != nil {
			return nil, err
		}

		newTs := stock.TimeSeries{
			Date:   t,
			Open:   timeSerie.Open,
			High:   timeSerie.High,
			Low:    timeSerie.Low,
			Close:  timeSerie.Close,
			Volume: timeSerie.Volume,
		}

		stockTimeSeries = append(stockTimeSeries, newTs)
	}

	logs.Log().Info("TimeSeriesDaily request success")

	return stockTimeSeries, nil
}
