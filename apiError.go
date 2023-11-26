package goalphavantage

type APIError struct {
	Message string
}

func (a *APIError) Error() string {
	return a.Message
}
