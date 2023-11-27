package goalphavantage

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Topic string
type Sort string

func (t Topic) Valid() bool {
	switch strings.ToLower(string(t)) {
	case "", "blockchain", "earnings", "ipo", "mergers_and_acquisitions", "financial_markets", "economy_fiscal", "economy_monetary", "economy_macro", "energy_transportation", "finance", "life_sciences", "manufacturing", "real_estate", "retail_wholesale", "technology":
		return true
	default:
		return false
	}
}

func (s Sort) Valid() bool {
	switch strings.ToLower(string(s)) {
	case "", "latest", "earliest", "relevance":
		return true
	default:
		return false
	}
}

func (n NewsSentimentOptions) Valid() bool {
	if !n.Sort.Valid() {
		return false
	}

	for _, topic := range n.Topics {
		if !topic.Valid() {
			return false
		}
	}

	return true
}

type NewsSentimentOptions struct {
	Tickers  []string `url:"tickers,omitempty"`
	Topics   []Topic  `url:"topics,omitempty"`
	TimeFrom string   `url:"time_from,omitempty"`
	TimeTo   string   `url:"time_to,omitempty"`
	Sort     Sort     `url:"sort,omitempty"`
	Limit    int      `url:"limit,omitempty"`
}

type TickerSentiment struct {
	Ticker               string `json:"ticker"`
	RelevanceScore       string `json:"relevance_score"`
	TickerSentimentScore string `json:"ticker_sentiment_score"`
	TickerSentimentLabel string `json:"ticker_sentiment_label"`
}

type TopicRelevance struct {
	Topic          string `json:"topic"`
	RelevanceScore string `json:"relevance_score"`
}

type NewsFeed struct {
	Title                 string            `json:"title"`
	URL                   string            `json:"url"`
	TimePublished         string            `json:"time_published"`
	Authors               []string          `json:"authors"`
	Summary               string            `json:"summary"`
	BannerImage           string            `json:"banner_image"`
	Source                string            `json:"source"`
	CategoryWithinSource  string            `json:"category_within_source"`
	SourceDomain          string            `json:"source_domain"`
	Topics                []TopicRelevance  `json:"topics"`
	OverallSentimentScore float64           `json:"overall_sentiment_score"`
	OverallSentimentLabel string            `json:"overall_sentiment_label"`
	TickerSentiment       []TickerSentiment `json:"ticker_sentiment"`
}
type NewsSentimentResponse struct {
	Items                    string     `json:"items"`
	SentimentScoreDefinition string     `json:"sentiment_score_definition"`
	RelevanceScoreDefinition string     `json:"relevance_score_definition"`
	Feed                     []NewsFeed `json:"feed"`
}

type RankedStock struct {
	Ticker           string `json:"ticker"`
	Price            string `json:"price"`
	ChangeAmount     string `json:"change_amount"`
	ChangePercentage string `json:"change_percentage"`
	Volume           string `json:"volume"`
}
type RankingResponse struct {
	Metadata           string        `json:"metadata"`
	LastUpdated        string        `json:"last_updated"`
	TopGainers         []RankedStock `json:"top_gainers"`
	TopLosers          []RankedStock `json:"top_losers"`
	MostActivelyTraded []RankedStock `json:"most_actively_traded"`
}

func (c *Client) GetNewsSentiment(ctx context.Context, options *NewsSentimentOptions) (*NewsSentimentResponse, error) {

	if !options.Valid() {
		return nil, InValidInputError
	}

	apiURL := fmt.Sprintf("%sfunction=NEWS_SENTIMENT&%s", c.BaseURL, c.buildQuery(options))

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	var res NewsSentimentResponse
	if err := c.doJSONRequest(req, &res); err != nil {
		return nil, fmt.Errorf("failed to get news sentiment: %w", err)
	}

	return &res, nil
}

func (c *Client) GetTopGainersLosers(ctx context.Context) (*RankingResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%sfunction=TOP_GAINERS_LOSERS", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	var res RankingResponse
	if err = c.doJSONRequest(req, &res); err != nil {
		return nil, fmt.Errorf("failed to get top gainers/losers: %w", err)
	}
	return &res, nil
}
