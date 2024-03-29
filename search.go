package nadago

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

type SearchResults struct {
	Rows []map[string]interface{} `json:"rows"`
}

type SearchResponse struct {
	Result SearchResults `json:"result"`
}

type Survey struct {
	Idno     string    `json:"idno"`
	Title    string    `json:"title"`
	Nation   string    `json:"nation"`
	Start    int       `json:"year_start"`
	End      int       `json:"year_end"`
	Created  time.Time `json:"created"`
	Changed  time.Time `json:"changed"`
	Url      string    `json:"url"`
	Varcount int       `json:"varcount"`
	Data     interface{}
}

// define all search parameters for the search endpoint
type SearchParams struct {
	Keywords   string `url:"sk,omitempty"`
	From       int    `url:"from,omitempty"`
	To         int    `url:"to,omitempty"`
	Country    string `url:"country,omitempty"`
	Inc_iso    bool   `url:"inc_iso,omitempty"`
	Created    string `url:"created,omitempty"`
	Dtype      string `url:"dtype,omitempty"`
	Ps         int    `url:"ps,omitempty"`
	Page       int    `url:"page,omitempty"`
	Sort_by    string `url:"sort_by,omitempty"`
	Sort_order string `url:"sort_order,omitempty"`
	Format     string `url:"format,omitempty"`
}

// set the default parameters
func NewDefaultSearchParams() *SearchParams {
	return &SearchParams{
		Inc_iso:    true,
		Ps:         30,
		Page:       1,
		Sort_by:    "year",
		Sort_order: "asc",
		Format:     "json",
	}
}

func (s *Survey) UnmarshalJSON(data []byte) error {
	type Alias Survey
	aux := &struct {
		Start    interface{} `json:"year_start"`
		End      interface{} `json:"year_end"`
		Varcount interface{} `json:"varcount"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	s.Start, err = convertToInt(aux.Start)
	if err != nil {
		return err
	}

	s.End, err = convertToInt(aux.End)
	if err != nil {
		return err
	}

	s.Varcount, err = convertToInt(aux.Varcount)
	if err != nil {
		return err
	}

	return nil
}

func convertToInt(value interface{}) (int, error) {
	if value == nil {
		return 0, nil // Return zero value if the input is nil (null)
	}

	switch v := value.(type) {
	case float64:
		return int(v), nil
	case string:
		if v == "" {
			return 0, nil // Return zero if the string is empty
		}
		val, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("failed to convert string to int: %v", err)
		}
		return val, nil
	default:
		return 0, fmt.Errorf("unexpected type: %T", v)
	}
}

func (c *Client) Search(ctx context.Context, params *SearchParams) ([]Survey, error) {

	//create a http request to the search endpoint
	path := c.apiURL + "/search"
	parsedUrl, err := url.Parse(path)
	if err != nil {
		return []Survey{}, AppErr{
			Message:    fmt.Errorf("invalid URL format: %w", err).Error(),
			StatusCode: 1001,
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", parsedUrl.String(), nil)
	if err != nil {
		return []Survey{}, AppErr{
			Message:    fmt.Errorf("failed to generate http request. %w", err).Error(),
			StatusCode: 1001,
		}
	}

	//extract params into url.Values
	v, err := query.Values(params)
	if err != nil {
		return []Survey{}, fmt.Errorf("failed to query parameters: %w", err)
	}

	//add the parameters to the request url
	q := req.URL.Query()
	for param, value := range v {
		for _, val := range value {
			q.Add(param, val)
		}
	}

	//make request and unmarshal response
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(req.URL.String())
	if err != nil {
		return []Survey{}, AppErr{
			Message:    fmt.Errorf("failed to complete http request. %w", err).Error(),
			StatusCode: 1001,
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Survey{}, FetchErr{
			Message:    "non-200 status code from the API",
			StatusCode: resp.StatusCode,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Survey{}, AppErr{
			Message:    fmt.Errorf("failed to read respose %w", err).Error(),
			StatusCode: 1001,
		}
	}

	var response SearchResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []Survey{}, AppErr{
			Message:    fmt.Errorf("failed to unmarshal response. %w", err).Error(),
			StatusCode: 1001,
		}
	}

	// extract response into slice of survey structs

	surveys, err := extractSurveys(&response.Result)
	if err != nil {
		return []Survey{}, AppErr{
			Message:    fmt.Errorf("failed to unmarshal response into surveys slice. %w", err).Error(),
			StatusCode: 1001,
		}
	}

	return surveys, nil
}

func extractSurveys(search *SearchResults) ([]Survey, error) {
	surveys := make([]Survey, 0)
	for _, row := range search.Rows {

		rowBytes, err := json.Marshal(row)
		if err != nil {
			return []Survey{}, err
		}

		var survey Survey
		if err := json.Unmarshal(rowBytes, &survey); err != nil {
			return []Survey{}, err
		}
		survey.Data = row
		surveys = append(surveys, survey)
	}
	return surveys, nil
}
