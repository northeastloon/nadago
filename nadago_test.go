package nadago_test

import (
	"context"
	"testing"
	"time"

	"github.com/northeastloon/nadago/nadago"
)

const (
	testTimeout = 30 * time.Second
)

var endpoints = []string{
	"https://catalog.ihsn.org/index.php/api/catalog",
	"https://microdata.worldbank.org/index.php/api/catalog",
	"https://www.ilo.org/surveyLib/index.php/api/catalog",
}

func TestEndpoints(t *testing.T) {
	for _, endpoint := range endpoints {
		testSearchForEndpoint(t, endpoint)
	}
}

func testSearchForEndpoint(t *testing.T, endpoint string) {
	client := nadago.NewClient(endpoint)
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	params := nadago.NewDefaultSearchParams()
	params.Ps = 2
	params.From = 2020
	params.To = 2020
	params.Country = "ZAF"

	surveys, err := client.Search(ctx, params)
	if err != nil {
		t.Fatalf("Error searching on endpoint %s: %v", endpoint, err)
	}

	if len(surveys) == 0 {
		t.Fatalf("Expected non-empty results from endpoint %s, but got none", endpoint)
	}

	// Check if Idno and Title fields are non-empty for each survey
	for _, survey := range surveys {
		if survey.Idno == "" {
			t.Errorf("Expected Idno to be non-empty for survey from endpoint %s", endpoint)
		}
		if survey.Title == "" {
			t.Errorf("Expected Title to be non-empty for survey from endpoint %s", endpoint)
		}
	}

	// Test GetSurveyMeta
	idno := surveys[0].Idno
	_, err = client.GetSurveyMeta(ctx, idno)
	if err != nil {
		t.Fatalf("Error fetching survey metadata on endpoint %s for Idno %s: %v", endpoint, idno, err)
	}

	// Test GetSurveyVars
	vars, err := client.GetSurveyVars(ctx, idno)
	if err != nil {
		t.Fatalf("Error fetching survey variables on endpoint %s for Idno %s: %v", endpoint, idno, err)
	}

	// Test GetVarMeta
	if len(vars.Vids) == 0 {
		t.Fatalf("Expected non-empty Vids from GetSurveyVars on endpoint %s for Idno %s", endpoint, idno)
	}

	vid := vars.Vids[0]
	_, err = client.GetVarMeta(ctx, idno, vid)
	if err != nil {
		t.Fatalf("Error fetching variable metadata on endpoint %s for Idno %s and Vid %s: %v", endpoint, idno, vid, err)
	}
}
