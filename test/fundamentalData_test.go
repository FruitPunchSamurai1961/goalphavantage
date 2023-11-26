package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"goalphavantage.FruitPunchSamurai1961.net"
	"testing"
)

func TestGetLatestActiveListingStatus(t *testing.T) {
	apiKey, err := getApiKey()
	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotEmpty(t, apiKey, "API_KEY should not be empty")

	c := goalphavantage.NewClient(apiKey)
	ctx := context.Background()
	listings, err := c.GetLatestActiveListingStatus(ctx)

	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotNil(t, listings, "expecting non-nil result")

	for idx, listing := range *listings {
		assert.NotEmpty(t, listing.Symbol, fmt.Sprintf("Symbol should not be empty. ActiveListing: %+v\n, Index: %d", listing, idx))
		assert.NotEmpty(t, listing.Name, fmt.Sprintf("Name should not be empty. ActiveListing: %+v\n, Index: %d", listing, idx))
		assert.NotEmpty(t, listing.AssetType, fmt.Sprintf("AssetType should not be empty. ActiveListing: %+v\n, Index: %d", listing, idx))
		assert.NotEmpty(t, listing.Exchange, fmt.Sprintf("Exchange should not be empty. ActiveListing: %+v\n, Index: %d", listing, idx))
	}
}
