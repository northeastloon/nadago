# Nadago

Nadago is a wrapper for the NADA consumer API, allowing users to access a subset of endpoints using Go, namely catalog search and fetching of study and variable metadata. The NADA API is implemented by survey metadata catalogs such as those of the International Household Survey Network, World Bank Microdata Library, and International Labour Organization (ILO) Data Catalog.

## Installation

 ```
 go get github.com/northeastloon/nadago
 ```

 ## Usage

 ```
package main

import (
	"context"
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

	ctx := context.Background()

	res, err := c.Search(ctx, p)

	if err != nil {
		fmt.Println(err)
	}

	//print first item
	fmt.Println(res[0])

	//get metadata for first survey item
	idno := res[0].Idno

	meta, err := c.GetSurveyMeta(ctx, idno)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(meta)

	//get variables for survey item
	vars, err := c.GetSurveyVars(ctx, idno)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(vars.Vids)

	//get variable metadata
	vid := vars.Vids[0]
	fmt.Println(vid)

	varmeta, err := c.GetVarMeta(ctx, idno, vid)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(varmeta.Data)

}
 
 ```
