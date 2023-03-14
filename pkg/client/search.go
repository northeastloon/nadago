package client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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
	Idno    string
	Created time.Time
	Changed time.Time
	Data    interface{}
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

func (c *Client) Search(params *SearchParams) ([]Survey, error) {

	//create a http request to the search endpoint
	req, err := http.NewRequest("GET", c.apiURL+"/search", nil)
	if err != nil {
		log.Fatal(err)
	}

	//extract params into url.Values
	v, err := query.Values(params)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response SearchResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// extract response into slice of survey structs

	surveys, err := extractSurveys(&response.Result)
	if err != nil {
		log.Fatal(err)
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
