package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Variables struct {
	Idno      string
	Variables []map[string]interface{} `json:"variables"`
	Vids      []string
}

func (c *Client) GetSurveyVars(ctx context.Context, idno string) (Variables, error) {

	//create a http request to the search endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiURL+"/"+idno+"/variables", nil)
	if err != nil {
		return Variables{}, AppErr{
			Message:    fmt.Errorf("failed to generate http request. %w", err).Error(),
			StatusCode: 1001,
		}
	}

	// make request and unmarshal response
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Variables{}, AppErr{
			Message:    fmt.Errorf("failed to complete http request. %w", err).Error(),
			StatusCode: 1001,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Variables{}, FetchErr{
			Message:    "non-200 status code from the API",
			StatusCode: resp.StatusCode,
		}
	}

	var vars Variables
	err = json.NewDecoder(resp.Body).Decode(&vars)
	if err != nil {
		return Variables{}, AppErr{
			Message:    fmt.Errorf("failed to unmarshal response. %w", err).Error(),
			StatusCode: 1001,
		}
	}
	vars.Idno = idno

	err = extractVids(&vars)
	if err != nil {
		return Variables{}, AppErr{
			Message:    fmt.Errorf("failed to extract variable ids from response. %w", err).Error(),
			StatusCode: 1001,
		}
	}

	return vars, nil

}

func extractVids(vars *Variables) error {

	for _, v := range vars.Variables {
		// extract the VID field from the map
		vid, ok := v["vid"].(string)
		if !ok {
			// create and return a new error
			return errors.New("VID field not found")
		}

		//append to vars
		vars.Vids = append(vars.Vids, vid)
	}
	return nil
}
