package test

import (
	"context"
	"fmt"
	"github.com/FruitPunchSamurai1961/goalphavantage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertNonEmptyFields(t *testing.T, listings *[]goalphavantage.Listing) {
	for idx, listing := range *listings {
		assert.NotEmpty(t, listing.Symbol, fmt.Sprintf("Symbol should not be empty. Listing: %+v\n, Index: %d", listing, idx))
		assert.NotEmpty(t, listing.AssetType, fmt.Sprintf("AssetType should not be empty. Listing: %+v\n, Index: %d", listing, idx))
		assert.NotEmpty(t, listing.Exchange, fmt.Sprintf("Exchange should not be empty. Listing: %+v\n, Index: %d", listing, idx))
	}
}

func assertInvalidInputError(t *testing.T, err error) {
	assert.NotNil(t, err, "expecting not-nil error")
	assert.ErrorIs(t, err, goalphavantage.InValidInputError, "expecting error to be type InvalidInputError")
}

func TestListingStatusCall(t *testing.T) {
	apiKey, err := getApiKey()
	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotEmpty(t, apiKey, "API_KEY should not be empty")

	c := goalphavantage.NewClient(apiKey)
	ctx := context.Background()

	//Test default call
	listings, err := c.GetListingStatus(ctx, nil)

	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotNil(t, listings, "expecting non-nil result")
	assertNonEmptyFields(t, listings)

	//Test call to get delisted stocks
	options := goalphavantage.ListingStatusOptions{
		Date:  "",
		State: "delisted",
	}
	listings, err = c.GetListingStatus(ctx, &options)

	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotNil(t, listings, "expecting non-nil result")
	assertNonEmptyFields(t, listings)

	//Test Invalid State Option
	options.State = "invalid-State"
	listings, err = c.GetListingStatus(ctx, &options)
	assert.Nil(t, listings, fmt.Sprintf("expecting null listings, got listings: %v", listings))
	assertInvalidInputError(t, err)

	//Test Invalid Date Option
	options.State = ""
	options.Date = "1234-56-78"
	listings, err = c.GetListingStatus(ctx, &options)
	assert.Nil(t, listings, fmt.Sprintf("expecting null listings, got listings: %v", listings))
	assertInvalidInputError(t, err)

	//Test Valid Date Option
	options.State = ""
	options.Date = "2010-01-02"
	listings, err = c.GetListingStatus(ctx, &options)

	assert.Nil(t, err, fmt.Sprintf("expecting nil error, got error: %v", err))
	assert.NotNil(t, listings, "expecting non-nil result")
	assertNonEmptyFields(t, listings)
}
