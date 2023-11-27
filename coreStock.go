package goalphavantage

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Function string
type DataType string
type Interval string
type OutputSize string

type BoolString string

func (f Function) Valid() bool {
	switch strings.ToUpper(string(f)) {
	case "TIME_SERIES_INTRADAY", "TIME_SERIES_DAILY", "TIME_SERIES_DAILY_ADJUSTED", "TIME_SERIES_WEEKLY", "TIME_SERIES_WEEKLY_ADJUSTED", "TIME_SERIES_MONTHLY", "TIME_SERIES_MONTHLY_ADJUSTED", "GLOBAL_QUOTE":
		return true
	default:
		return false
	}
}

func (d DataType) Valid() bool {
	switch strings.ToLower(string(d)) {
	case "", "json", "csv":
		return true
	default:
		return false
	}
}

func (i Interval) Valid() bool {
	switch strings.ToLower(string(i)) {
	case "1min", "5min", "15min", "30min", "60min":
		return true
	default:
		return false
	}
}

func (o OutputSize) Valid() bool {
	switch strings.ToLower(string(o)) {
	case "", "compact", "full":
		return true
	default:
		return false
	}
}

func (b BoolString) Valid() bool {
	switch strings.ToLower(string(b)) {
	case "", "true", "false":
		return true
	default:
		return false
	}
}

func (c CoreStockSharedInputOptions) Valid() bool {
	if !c.Function.Valid() || c.Symbol == "" {
		return false
	}

	if strings.ToUpper(string(c.Function)) == "TIME_SERIES_INTRADAY" && !c.Interval.Valid() {
		return false
	}

	if !c.Datatype.Valid() || !c.Adjusted.Valid() || !c.ExtendedHours.Valid() || !c.OutputSize.Valid() {
		return false
	}
	return true
}

type CoreStockSharedInputOptions struct {
	Function      Function   `url:"function"`
	Symbol        string     `url:"symbol"`
	Interval      Interval   `url:"interval"`
	Datatype      DataType   `url:"datatype, omitempty"`
	Adjusted      BoolString `url:"adjusted, omitempty"`
	ExtendedHours BoolString `url:"extended_hours, omitempty"`
	Month         string     `url:"month, omitempty"`
	OutputSize    OutputSize `url:"outputsize, omitempty"`
}

type MetaData struct {
	Information   *string `json:"1. Information,omitempty"`
	Symbol        *string `json:"2. Symbol,omitempty"`
	LastRefreshed *string `json:"3. Last Refreshed,omitempty"`
	Interval      *string `json:"4. Interval,omitempty"`
	OutputSize    *string `json:"5. Output Size,omitempty"`
	TimeZone      *string `json:"6. Time Zone,omitempty"`
}

type CoreStockData struct {
	Open                  *string `json:"1. open"`
	High                  *string `json:"2. high"`
	Low                   *string `json:"3. low"`
	Close                 *string `json:"4. close"`
	Volume                *string `json:"5. volume,omitempty"`
	AdjustedClose         *string `json:"5. adjusted close,omitempty"`
	VolumeForAdjustedCall *string `json:"6. volume,omitempty"`
	DividendAmount        *string `json:"7. dividend amount,omitempty"`
}

type GlobalQuote struct {
	Symbol           *string `json:"01. symbol"`
	Open             *string `json:"02. open"`
	High             *string `json:"03. high"`
	Low              *string `json:"04. low"`
	Price            *string `json:"05. price"`
	Volume           *string `json:"06. volume"`
	LatestTradingDay *string `json:"07. latest trading day"`
	PreviousClose    *string `json:"08. previous close"`
	Change           *string `json:"09. change"`
	ChangePercent    *string `json:"10. change percent"`
}

type CoreStockResponse struct {
	MetaData                  *MetaData                 `json:"Meta Data,omitempty"`
	LatestQuote               *GlobalQuote              `json:"Global Quote,omitempty"`
	MonthlyAdjustedTimeSeries map[string]*CoreStockData `json:"Monthly Adjusted Time Series,omitempty"`
	MonthlyTimeSeries         map[string]*CoreStockData `json:"Monthly Time Series,omitempty"`
	OneMinTimeSeries          map[string]*CoreStockData `json:"Time Series (1min),omitempty"`
	FiveMinTimeSeries         map[string]*CoreStockData `json:"Time Series (5min),omitempty"`
	FifteenMinTimeSeries      map[string]*CoreStockData `json:"Time Series (15min),omitempty"`
	ThirtyMinTimeSeries       map[string]*CoreStockData `json:"Time Series (30min),omitempty"`
	HourTimeSeries            map[string]*CoreStockData `json:"Time Series (60min),omitempty"`
	DailyTimeSeries           map[string]*CoreStockData `json:"Time Series (Daily),omitempty"`
	WeeklyTimeSeries          map[string]*CoreStockData `json:"Weekly Time Series,omitempty"`
	WeeklyAdjustedTimeSeries  map[string]*CoreStockData `json:"Weekly Adjusted Time Series,omitempty"`
}

func (c *Client) GetTimeSeriesStockData(ctx context.Context, options *CoreStockSharedInputOptions) (*CoreStockResponse, error) {
	if !options.Valid() {
		return nil, InValidInputError
	}

	apiURL := fmt.Sprintf("%s%s", c.BaseURL, c.buildQuery(options))

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	var res CoreStockResponse
	if strings.ToLower(string(options.Datatype)) == "csv" {
		if err := c.doCSVRequest(req, &res); err != nil {
			return nil, fmt.Errorf("failed to get time sereis stock data: %w", err)
		}
	} else {
		if err := c.doJSONRequest(req, &res); err != nil {
			return nil, fmt.Errorf("failed to get time sereis stock data: %w", err)
		}
	}

	return &res, nil
}
