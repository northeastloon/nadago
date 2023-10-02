package nadago

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSurveyVars(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(expectedVariablesResponse)))
		}))
		defer ts.Close()

		client := NewClient(ts.URL)

		ctx := context.Background()
		idno := "LBR_2020_FIES_v01_M_v01_A_OCS"

		expectedVids := []string{"V1", "V2", "V3", "V4", "V5", "V6", "V7", "V8", "V9", "V10", "V11", "V12", "V13", "V14", "V15", "V16", "V17", "V18", "V19", "V20", "V21", "V22"}

		variables, err := client.GetSurveyVars(ctx, idno)
		assert.NoError(t, err)
		assert.Equal(t, expectedVids, variables.Vids)
		assert.Equal(t, idno, variables.Idno)
		assert.NotNil(t, variables.Variables)
	})

	t.Run("bad request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
		}))
		defer ts.Close()

		client := NewClient(ts.URL)

		ctx := context.Background()
		idno := "LBR_2020_FIES_v01_M_v01_A_OCS"

		_, err := client.GetSurveyVars(ctx, idno)
		assert.Error(t, err)
		assert.Equal(t, FetchErr{
			Message:    "non-200 status code from the API",
			StatusCode: 404,
		}, err)
	})

	t.Run("failed to complete http request", func(t *testing.T) {

		failingClient := &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("failed to complete http request")
			}),
		}
		client := NewClient("http://invalid-url", WithHTTPClient(failingClient))

		ctx := context.Background()
		idno := "LBR_2020_FIES_v01_M_v01_A_OCS"

		_, err := client.GetSurveyVars(ctx, idno)
		assert.Error(t, err)

		appErr, ok := err.(AppErr)
		assert.True(t, ok, "error should be of type AppErr")
		assert.Equal(t, 1001, appErr.StatusCode, "error status code should match expected value")
	})

	t.Run("failed to unmarshal response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}))
		defer ts.Close()

		client := NewClient(ts.URL)

		ctx := context.Background()
		idno := "LBR_2020_FIES_v01_M_v01_A_OCS"

		_, err := client.GetSurveyVars(ctx, idno)
		assert.Error(t, err)
		assert.IsType(t, AppErr{}, err)
	})

}

var expectedVariablesResponse = `{"total":22,"variables":[{"uid":"19260146","sid":"10894","fid":"F1","vid":"V1","name":"Random_ID","labl":"Unique respondent identifier"},{"uid":"19260147","sid":"10894","fid":"F1","vid":"V2","name":"WORRIED","labl":"Worried you would not have enough food to eat because of a lack of money or other resources"},{"uid":"19260148","sid":"10894","fid":"F1","vid":"V3","name":"HEALTHY","labl":"Unable to eat healthy and nutritious food because of a lack of money or other resources"},{"uid":"19260149","sid":"10894","fid":"F1","vid":"V4","name":"FEWFOOD","labl":"Ate only a few kinds of foods because of a lack of money or other resources"},{"uid":"19260150","sid":"10894","fid":"F1","vid":"V5","name":"SKIPPED","labl":"Skipped a meal because there was not enough money or other resources to get food"},{"uid":"19260151","sid":"10894","fid":"F1","vid":"V6","name":"ATELESS","labl":"Ate less than you thought you should because of a lack of money or other resources"},{"uid":"19260152","sid":"10894","fid":"F1","vid":"V7","name":"RUNOUT","labl":"Household ran out of food because of a lack of money or other resources"},{"uid":"19260153","sid":"10894","fid":"F1","vid":"V8","name":"HUNGRY","labl":"Hungry but did not eat because there was not enough money or other resources for food?"},{"uid":"19260154","sid":"10894","fid":"F1","vid":"V9","name":"WHLDAY","labl":"Went without eating for a whole day because of a lack of money or other resources?"},{"uid":"19260155","sid":"10894","fid":"F1","vid":"V10","name":"wt","labl":"Post-stratification sampling weights"},{"uid":"19260156","sid":"10894","fid":"F1","vid":"V11","name":"year","labl":"Year when the study was administered in the country"},{"uid":"19260157","sid":"10894","fid":"F1","vid":"V12","name":"N_adults","labl":"Number of adults 15 years of age and above in household"},{"uid":"19260158","sid":"10894","fid":"F1","vid":"V13","name":"N_child","labl":"Number of children under 15 years of age in household"},{"uid":"19260159","sid":"10894","fid":"F1","vid":"V14","name":"Raw_score","labl":"Sum of Affirmative responses to FIES questions"},{"uid":"19260160","sid":"10894","fid":"F1","vid":"V15","name":"Raw_score_par","labl":"Estimated person parameters using the Rasch model"},{"uid":"19260161","sid":"10894","fid":"F1","vid":"V16","name":"Raw_score_par_error","labl":"Estimated person parameter errors using the Rasch model"},{"uid":"19260162","sid":"10894","fid":"F1","vid":"V17","name":"Prob_Mod_Sev","labl":"Probability of being moderately or severely food insecure"},{"uid":"19260163","sid":"10894","fid":"F1","vid":"V18","name":"Prob_sev","labl":"Probability of being severely food insecure"},{"uid":"19260164","sid":"10894","fid":"F1","vid":"V19","name":"Age","labl":"Age of the respondent"},{"uid":"19260165","sid":"10894","fid":"F1","vid":"V20","name":"Education","labl":"Education of the respondent"},{"uid":"19260166","sid":"10894","fid":"F1","vid":"V21","name":"Area","labl":"Area"},{"uid":"19260167","sid":"10894","fid":"F1","vid":"V22","name":"Gender","labl":"Gender of the respondent"}]}`
