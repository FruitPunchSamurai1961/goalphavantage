package goalphavantage

import "errors"

type APIError struct {
	Message string
}

func (a *APIError) Error() string {
	return a.Message
}

func IsAPIError(err error) bool {
	var apiError *APIError
	return errors.As(err, &apiError)
}
