package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SurveyMeta struct {
	Idno string
	Data interface{} `json:"dataset"`
}

func (c *Client) GetSurveyMeta(idno string) (SurveyMeta, error) {

	//create a http request to the search endpoint
	req, err := http.NewRequest("GET", c.apiURL+"/"+idno, nil)
	if err != nil {
		return []Survey, CreateReqErr{
			Message:    err.Error(),
			StatusCode: 1001,
		}
	}

	// make request and unmarshal response
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []Survey, FetchErr{
			Message:    err.Error(),
			StatusCode: 1002,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Survey, FetchErr{
			Message:    "non-200 status code from the API",
			StatusCode: resp.StatusCode,
		}
	}

	// decode JSON response into SurveyMeta struct
	var meta SurveyMeta
	err = json.NewDecoder(resp.Body).Decode(&meta)
	if err != nil {
		return SurveyMeta{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	meta.Idno = idno

	return meta, nil
}
