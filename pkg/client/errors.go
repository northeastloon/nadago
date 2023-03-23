package client

import (
	"fmt"
)

type FetchErr struct {
	StatusCode int
	Message    string
}

type CreateReqErr struct {
	StatusCode int
	Message    string
}

func (e FetchErr) Error() string {
	return fmt.Sprintf("failed to fetch response: %s with statuscode: %d", e.Message, e.StatusCode)
}

func (e CreateReqErr) Error() string {
	return fmt.Sprintf("failed to create http request: %s with statuscode: %d", e.Message, e.StatusCode)
}
