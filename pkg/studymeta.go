package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SurveyMeta struct {
	Idno string
	Data interface{} `json:"dataset"`
}

func (c *Client) GetSurveyMeta(ctx context.Context, idno string) (SurveyMeta, error) {

	//create a http request to the search endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiURL+"/"+idno, nil)
	if err != nil {
		return SurveyMeta{}, AppErr{
			Message:    fmt.Errorf("failed to generate http request. %w", err).Error(),
			StatusCode: 1001,
		}
	}

	// make request and unmarshal response
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SurveyMeta{}, AppErr{
			Message:    fmt.Errorf("failed to complete http request. %w", err).Error(),
			StatusCode: 1001,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SurveyMeta{}, FetchErr{
			Message:    "non-200 status code from the API",
			StatusCode: resp.StatusCode,
		}
	}

	// decode JSON response into SurveyMeta struct
	var meta SurveyMeta
	err = json.NewDecoder(resp.Body).Decode(&meta)
	if err != nil {
		return SurveyMeta{}, AppErr{
			Message:    fmt.Errorf("failed to unmarshal response. %w", err).Error(),
			StatusCode: 1001,
		}
	}
	meta.Idno = idno

	return meta, nil
}
