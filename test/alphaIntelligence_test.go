package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	goalphavantage "goalphavantage.FruitPunchSamurai1961.net"
	"testing"
)

func TestGetNewsSentimentWithTickerOption(t *testing.T) {
	apiKey, err := getApiKey()
	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotEmpty(t, apiKey, "API_KEY should not be empty")

	c := goalphavantage.NewClient(apiKey)
	ctx := context.Background()

	// Test with valid tickers
	options := goalphavantage.NewsSentimentOptions{
		Tickers: []string{"AAPL", "IBM"},
	}
	res, err := c.GetNewsSentiment(ctx, &options)
	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotNil(t, res, "expecting non-nil result")

	// Test with empty tickers
	optionsEmpty := goalphavantage.NewsSentimentOptions{
		Tickers: []string{},
	}
	resEmpty, errEmpty := c.GetNewsSentiment(ctx, &optionsEmpty)
	assert.Nil(t, errEmpty, fmt.Sprintf("expecting nil error, got error: %v", errEmpty))
	assert.NotNil(t, resEmpty, "expecting non-nil result")

	// Test with invalid tickers
	optionsInvalid := goalphavantage.NewsSentimentOptions{
		Tickers: []string{"INVALID_TICKER"},
	}
	resInvalid, errInvalid := c.GetNewsSentiment(ctx, &optionsInvalid)
	assert.NotNil(t, errInvalid, "expecting non-nil error for invalid tickers")

	var apiError *goalphavantage.APIError
	ok := errors.As(errInvalid, &apiError)
	assert.True(t, ok, fmt.Sprintf("expecting APIError got: %v", errInvalid))

	assert.Nil(t, resInvalid, "expecting nil result for invalid tickers")
}

func TestGetTopLosersGainersStock(t *testing.T) {
	apiKey, err := getApiKey()
	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotEmpty(t, apiKey, "API_KEY should not be empty")

	c := goalphavantage.NewClient(apiKey)
	ctx := context.Background()

	res, err := c.GetTopGainersLosers(ctx)
	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotNil(t, res, "expecting non-nil result")

	// Additional assertions
	assert.Equal(t, 20, len(res.TopGainers), "Expecting 20 top gainers")
	assert.Equal(t, 20, len(res.TopLosers), "Expecting 20 top losers")
	assert.Equal(t, 20, len(res.MostActivelyTraded), "Expecting 20 most actively traded stocks")

	// Check the structure and data of each RankedStock item
	checkRankedStockList(t, res.TopGainers, "top gainer")
	checkRankedStockList(t, res.TopLosers, "top loser")
	checkRankedStockList(t, res.MostActivelyTraded, "most actively traded")
}

func checkRankedStockList(t *testing.T, stocks []goalphavantage.RankedStock, listType string) {
	for i, stock := range stocks {
		assert.NotEmpty(t, stock.Ticker, fmt.Sprintf("Expecting non-empty ticker for %s #%d", listType, i+1))
		assert.NotEmpty(t, stock.Price, fmt.Sprintf("Expecting non-empty price for %s #%d", listType, i+1))
		assert.NotEmpty(t, stock.ChangeAmount, fmt.Sprintf("Expecting non-empty change amount for %s #%d", listType, i+1))
		assert.NotEmpty(t, stock.ChangePercentage, fmt.Sprintf("Expecting non-empty change percentage for %s #%d", listType, i+1))
		assert.NotEmpty(t, stock.Volume, fmt.Sprintf("Expecting non-empty volume for %s #%d", listType, i+1))
	}
}
