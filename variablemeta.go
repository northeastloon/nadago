package nadago

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Variable struct {
	Idno string
	Vid  string
	Data interface{} `json:"variable"`
}

func (c *Client) GetVarMeta(ctx context.Context, idno string, vid string) (Variable, error) {

	//create a http request to the search endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiURL+"/"+idno+"/variables/"+vid, nil)
	if err != nil {
		return Variable{}, AppErr{
			Message:    fmt.Errorf("failed to generate http request. %w", err).Error(),
			StatusCode: 1001,
		}
	}

	// make request and unmarshal response
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Variable{}, AppErr{
			Message:    fmt.Errorf("failed to complete http request. %w", err).Error(),
			StatusCode: 1001,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Variable{}, FetchErr{
			Message:    "non-200 status code from the API",
			StatusCode: resp.StatusCode,
		}
	}

	var v Variable
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return Variable{}, AppErr{
			Message:    fmt.Errorf("failed to unmarshal response. %w", err).Error(),
			StatusCode: 1001,
		}
	}
	v.Idno = idno
	v.Vid = vid

	return v, nil

}
