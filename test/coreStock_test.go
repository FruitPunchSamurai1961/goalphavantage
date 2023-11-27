package test

import (
	"context"
	"fmt"
	"github.com/FruitPunchSamurai1961/goalphavantage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMonthlyAdjustedCoreStockCall(t *testing.T) {
	apiKey, err := getApiKey()
	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotEmpty(t, apiKey, "API_KEY should not be empty")

	c := goalphavantage.NewClient(apiKey)
	ctx := context.Background()

	// Define test options
	options := &goalphavantage.CoreStockSharedInputOptions{
		Function:   goalphavantage.Function("TIME_SERIES_MONTHLY_ADJUSTED"),
		Symbol:     "AAPL",
		Datatype:   goalphavantage.DataType("json"),
		Adjusted:   goalphavantage.BoolString("true"),
		OutputSize: goalphavantage.OutputSize("compact"),
	}

	// Make API call
	response, err := c.GetTimeSeriesStockData(ctx, options)
	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotNil(t, response, "expecting non-nil response")

	assert.NotNil(t, response.MetaData, "expecting non-nil MetaData")
	assert.NotNil(t, response.MonthlyAdjustedTimeSeries, "expecting non-nil MonthlyAdjustedTimeSeries")

	assert.NotEmpty(t, response.MetaData.Symbol, "expecting non-empty Symbol field in MetaData")

}
