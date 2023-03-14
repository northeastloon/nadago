package main

import (
	"fmt"

	"github.com/northeastloon/nadago/pkg/client"
)

func main() {

	//create new client without authentication
	c := client.NewClient("http://catalog.ihsn.org/index.php/api/catalog")

	// set search parameters
	p := client.NewDefaultSearchParams()
	p.Ps = 2
	p.From = 2020
	p.To = 2020
	p.Country = "ALB"

	res, err := c.Search(p)

	if err != nil {
		fmt.Println(err)
	}

	//print first item
	fmt.Println(res[0])

	//get metadata for first survey item
	idno := res[0].Idno

	meta, err := c.GetSurveyMeta(idno)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(meta)

	//get variables for survey item
	vars, err := c.GetSurveyVars(idno)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(vars.Variables[0])

	//get variable metadata
	vid := vars.Variables[0]["vid"]

	varmeta, err := c.GetVarMeta(idno, vid)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(varmeta.Data[0])

}
