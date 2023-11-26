package goalphavantage

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const baseURL = "https://www.alphavantage.co/query?"

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

type statusErrorResponse struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

type apiErrorResponse struct {
	Information  string `json:"information,omitempty"`
	ErrorMessage string `json:"Error_Message,omitempty"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		}}
}

func (c *Client) doJSONRequest(req *http.Request, v interface{}) error {
	setUpHeaders(req)
	c.addAPIKey(req)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if err = checkStatusErrorResponse(res); err != nil {
		return err
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err = checkAPIResponseForErrorMessage(content); err != nil {
		return err
	}

	if err = json.Unmarshal(content, &v); err != nil {
		return err
	}
	return nil
}

func (c *Client) doCSVRequest(req *http.Request, v interface{}) error {
	setUpHeaders(req)
	c.addAPIKey(req)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if err = checkStatusErrorResponse(res); err != nil {
		return err
	}

	content, err := getContent(res)
	if err != nil {
		return err
	}

	contentType := res.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		return checkAPIResponseForErrorMessage(content)
	}

	csvReader := csv.NewReader(bytes.NewReader(content))
	if err := readCSV(csvReader, v); err != nil {
		return err
	}

	return nil
}

func readCSV(reader *csv.Reader, v interface{}) error {
	var records []map[string]string
	csvHeader, err := reader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if errors.Is(err, csv.ErrFieldCount) {
			continue
		} else if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		csvRow := make(map[string]string)
		for i, value := range record {
			csvRow[csvHeader[i]] = value
		}

		records = append(records, csvRow)
	}

	switch v := v.(type) {
	case *[]ActiveListing:
		for _, record := range records {
			listing := ActiveListing{
				Symbol:    record["symbol"],
				Name:      record["name"],
				Exchange:  record["exchange"],
				AssetType: record["assetType"],
			}

			if listing.Name == "" {
				listing.Name = listing.Symbol
			}

			*v = append(*v, listing)
		}
	default:
		return fmt.Errorf("unsupported type %T for v", v)
	}

	return nil
}

func (c *Client) buildQuery(options interface{}) string {
	queryParams := url.Values{}

	reflectValue := reflect.ValueOf(options)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	if reflectValue.Kind() != reflect.Struct {
		return ""
	}

	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Type().Field(i)
		fieldValue := reflectValue.Field(i)

		if fieldValue.IsZero() || field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("url")
		if tag == "" {
			tag = field.Name
		}

		tag = strings.Split(tag, ",")[0]

		switch fieldValue.Kind() {
		case reflect.Array, reflect.Slice:
			for j := 0; j < fieldValue.Len(); j++ {
				queryParams.Add(tag, fmt.Sprint(fieldValue.Index(j)))
			}
		default:
			queryParams.Add(tag, fmt.Sprint(fieldValue))
		}
	}

	return queryParams.Encode()
}

func checkAPIResponseForErrorMessage(content []byte) error {
	var apiErrRes apiErrorResponse
	if err := json.Unmarshal(content, &apiErrRes); err != nil {
		return err
	}

	if apiErrRes.ErrorMessage != "" {
		return &APIError{Message: fmt.Sprintf("alphvantage call error message: %s", apiErrRes.ErrorMessage)}
	}
	if apiErrRes.Information != "" {
		return &APIError{Message: fmt.Sprintf("alphvantage call error message: %s", apiErrRes.Information)}
	}

	return nil
}

func (c *Client) addAPIKey(req *http.Request) {
	q := req.URL.Query()
	q.Add("apikey", c.apiKey)
	req.URL.RawQuery = q.Encode()
}

func setUpHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
}

func getContent(res *http.Response) ([]byte, error) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func checkStatusErrorResponse(res *http.Response) error {
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes statusErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return &APIError{Message: fmt.Sprintf("alphvantage call error message: (%s), status code: %d", errRes.Detail, res.StatusCode)}
		} else {
			return fmt.Errorf("unknown error: %w, status code: %d", err, res.StatusCode)
		}
	}
	return nil
}
