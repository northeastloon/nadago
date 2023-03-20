package client

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Variables struct {
	Idno      string
	Variables []map[string]interface{} `json:"variables"`
	Vids      []string
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

	err = extractVids(&vars)
	if err != nil {
		return Variables{}, err
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
