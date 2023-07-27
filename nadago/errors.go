package nadago

import (
	"fmt"
)

type FetchErr struct {
	StatusCode int
	Message    string
}

type AppErr struct {
	StatusCode int
	Message    string
}

func (e FetchErr) Error() string {
	return fmt.Sprintf("failed to fetch response: %s with statuscode: %d", e.Message, e.StatusCode)
}

func (e AppErr) Error() string {
	return fmt.Sprintf("Application side error: %s with statuscode: %d", e.Message, e.StatusCode)
}
