package goalphavantage

import (
	"context"
	"fmt"
	"net/http"
)

type ActiveListing struct {
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Exchange  string `json:"exchange"`
	AssetType string `json:"assetType"`
}

func (c *Client) GetLatestActiveListingStatus(ctx context.Context) (*[]ActiveListing, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%sfunction=LISTING_STATUS", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	var res []ActiveListing
	if err = c.doCSVRequest(req, &res); err != nil {
		return nil, fmt.Errorf("failed to get latest active listing status: %w", err)
	}
	return &res, nil
}
