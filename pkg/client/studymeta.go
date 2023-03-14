package client

import (
	"encoding/json"
	"log"
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
		log.Fatal(err)
	}

	// make request and unmarshal response
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// decode JSON response into SurveyMeta struct
	var meta SurveyMeta
	err = json.NewDecoder(resp.Body).Decode(&meta)
	if err != nil {
		return SurveyMeta{}, err
	}
	meta.Idno = idno

	return meta, nil
}
