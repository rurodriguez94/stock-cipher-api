package provider

type TimeSeriesRes struct {
	TimeSeries map[string]TimeSeriesRaw `json:"Time Series (Daily)"`
}

type TimeSeriesRaw struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}
