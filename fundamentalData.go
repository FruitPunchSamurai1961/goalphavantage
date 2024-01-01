package goalphavantage

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type State string

func (s State) Valid() bool {
	switch strings.ToLower(string(s)) {
	case "", "active", "delisted":
		return true
	default:
		return false
	}
}

func (l ListingStatusOptions) Valid() bool {
	if !l.State.Valid() {
		return false
	}

	if l.Date != "" {
		//Check for YYYY-MM-DD
		dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if !dateRegex.MatchString(l.Date) {
			return false
		}

		dateValue, err := time.Parse("2006-01-02", l.Date)
		if err != nil {
			return false
		}

		minimumDate := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
		if dateValue.Before(minimumDate) {
			return false
		}
	}

	return true
}

type ListingStatusOptions struct {
	Date  string `url:"date"`
	State State  `url:"state"`
}

type Listing struct {
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Exchange  string `json:"exchange"`
	AssetType string `json:"assetType"`
}

func (c *Client) GetListingStatus(ctx context.Context, options *ListingStatusOptions) (*[]Listing, error) {
	if options != nil && !options.Valid() {
		return nil, InValidInputError
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%sfunction=LISTING_STATUS&%s", c.BaseURL, c.buildQuery(options)), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	var res []Listing
	if err = c.doCSVRequest(req, &res); err != nil {
		return nil, fmt.Errorf("failed to get latest active listing status: %w", err)
	}
	return &res, nil
}
