package client

import (
	"encoding/json"
	"log"
	"net/http"
)

type Variables struct {
	Idno      string
	Variables []map[string]interface{} `json:"variables"`
}

func (c *Client) GetSurveyVars(idno string) (Variables, error) {

	//create a http request to the search endpoint
	req, err := http.NewRequest("GET", c.apiURL+"/"+idno+"/variables", nil)
	if err != nil {
		log.Fatal(err)
	}

	// make request and unmarshal response
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var vars Variables
	err = json.NewDecoder(resp.Body).Decode(&vars)
	if err != nil {
		return Variables{}, err
	}
	vars.Idno = idno

	return vars, nil

}
