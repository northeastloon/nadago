package client

import (
	"encoding/json"
	"log"
	"net/http"
)

type Variable struct {
	Idno string
	Vid  string
	Data interface{} `json:"variable"`
}

func (c *Client) GetVarMeta(idno string, vid string) (Variable, error) {

	//create a http request to the search endpoint
	req, err := http.NewRequest("GET", c.apiURL+"/"+idno+"/variables/"+vid, nil)
	if err != nil {
		log.Fatal(err)
	}

	// make request and unmarshal response
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var v Variable
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return Variable{}, err
	}
	v.Idno = idno
	v.Vid = vid

	return v, nil

}
