package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSurveyMeta(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		expectedIdno := "ARG_2021_HFS-Q1Q2_v01_M"

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"idno": "%s", "dataset": %s}`, expectedIdno, expectedSurveymetaResponse)))
		}))
		defer ts.Close()

		client := NewClient(ts.URL)

		ctx := context.Background()
		idno := "ARG_2021_HFS-Q1Q2_v01_M"

		meta, err := client.GetSurveyMeta(ctx, idno)
		assert.NoError(t, err)
		assert.Equal(t, expectedIdno, meta.Idno)
		assert.NotNil(t, meta.Data)
	})

	t.Run("bad request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
		}))
		defer ts.Close()

		client := NewClient(ts.URL)

		ctx := context.Background()
		idno := "ARG_2021_HFS-Q1Q2_v01_M"

		_, err := client.GetSurveyMeta(ctx, idno)
		assert.Error(t, err)
		assert.Equal(t, AppErr{
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
		idno := "ARG_2021_HFS-Q1Q2_v01_M"

		_, err := client.GetSurveyMeta(ctx, idno)
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
		idno := "ARG_2021_HFS-Q1Q2_v01_M"

		_, err := client.GetSurveyMeta(ctx, idno)
		assert.Error(t, err)
		assert.IsType(t, AppErr{}, err)
	})
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

var expectedSurveymetaResponse = `{"status":"success","dataset":{"id":"10226","repositoryid":"central","type":"survey","idno":"ARG_2021_HFS-Q1Q2_v01_M","title":
"High Frequency Survey 2021, Quarter 1 and 2","year_start":"2021","year_end":"2021","nation":"Argentina","authoring_entity":"UN Refugee Agency (UNHCR)","published":
"1","created":"2022-05-03T12:46:00+00:00","changed":"2022-05-05T18:45:52+00:00","varcount":"485","total_views":"5641","total_downloads":"0","formid":"5","data_access_type":
"remote","remote_data_url":"https:\/\/microdata.unhcr.org\/index.php\/catalog\/492","data_class_id":null,"data_class_code":null,"data_class_title":null,"thumbnail":null,
"link_study":"https:\/\/microdata.unhcr.org\/index.php\/catalog\/492","link_indicator":null,"link_report":null,"metadata":{"doc_desc":{"title":"ARG_2021_HFS-Q1Q2_v01_M",
"idno":"DDI_ARG_2021_HFS-Q1Q2_v01_M","producers":[{"name":"UN Refugee Agency","abbreviation":"UNHCR","affiliation":"UN","role":"Documentation of the study"},
{"name":"Development Economics Data Group","abbreviation":"DECDG","affiliation":"The World Bank","role":"Metadata adapted for Microdata Library"}],"prod_date":"2021-10-11",
"version_statement":{"version":"Version 01: This metadata was downloaded from the UNHCR Microdata Library catalog (https:\/\/microdata.unhcr.org\/index.php).
 The following two metadata fields were edited - Document and Survey ID."}},"study_desc":{"title_statement":{"idno":"ARG_2021_HFS-Q1Q2_v01_M","title":"High Frequency Survey 2021",
 "sub_title":"Quarter 1 and 2","alt_title":"HFS-Q1Q2 2021"},"authoring_entity":[{"name":"UN Refugee Agency (UNHCR)","affiliation":"UN"}],"distribution_statement":
 {"contact":[{"name":"Curation team","affiliation":"UNHCR","email":"microdata@unhcr.org","uri":"https:\/\/microdata.unhcr.org"}]},"series_statement":{"series_name":
 "Other Household Survey [hh\/oth]"},"version_statement":{"version":"Edited, cleaned and anonymised data.","version_date":"2021-10-11"},"study_info":{"keywords":[{"keyword":
 "HFS","vocab":"","uri":""},{"keyword":"High Frequency Survey","vocab":"","uri":""},{"keyword":"Covid","vocab":"","uri":""}],"topics":[{"topic":"Health","vocab":"","uri":""},
 {"topic":"Protection","vocab":"","uri":""},{"topic":"Livelihood and Social cohesion","vocab":"","uri":""},{"topic":"Transportation","vocab":"","uri":""},{"topic":"Basic Needs",
 "vocab":"","uri":""}],"abstract":"The data was collected using the High Frequency Survey (HFS), the new regional data collection tool & methodology launched in the Americas.
  The survey allowed for better reaching populations of interest with new remote modalities (phone interviews and self-administered surveys online) and improved sampling guidance
   and strategies. It includes a set of standardized regional core questions while allowing for operation-specific customizations. The core questions revolve around populations
    of interest's demographic profile, difficulties during their journey, specific protection needs, access to documentation & regularization, health access, coverage of basic needs
	, coping capacity & negative mechanisms used, and well-being & local integration. The data collected has been used by countries in their protection monitoring analysis and
	 vulnerability analysis.","coll_dates":[{"start":"2021-01-01","end":"2021-06-30","cycle":""}],"nation":[{"name":"Argentina","abbreviation":"ARG"}],"geog_coverage":"National
	  coverage","analysis_unit":"Household","universe":"All people of concern.","data_kind":"Sample survey data [ssd]","notes":"The scope includes:\n- household characteristics\n
	  - demographic profile\n- journey\n- protection needs\n- access to documentation\n- health\n- basic needs\n- coping capacity\n- local integration\n- Covid impact"},"method":
	  {"data_collection":{"data_collectors":[{"name":"UN Refugee Agency","abbreviation":"UNHCR","affiliation":"UN"}],"sampling_procedure":"In the absence of a well-developed
	   sampling-frame for forcibly displaced populations in the Americas, the High Frequency Survey employed a multi-frame sampling strategy where respondents entered the sample
	    through one of three channels: (i) those who opt-in to complete an online self-administered version of the questionnaire which was widely circulated through refugee
		 social media; (ii) persons identified through UNHCR and partner databases who were remotely-interviewed by phone; and (iii) random selection from the cases approaching
		  UNHCR for registration or assistance. The total sample size was 406 households.","coll_mode":["Other [oth]"],"research_instrument":"The questionnaire contained the
		   following sections: journey, family composition, vulnerability, basic Needs, coping capacity, well-being, COVID-19 Impact.","coll_situation":"Data collection
		    modalities include phone interviews (to accommodate for mobility restrictions due to the COVID-19 pandemic) and self-administered online surveys. Enumerators
			 were trained at the local level, keeping in line with the regional standards while reflecting the contextual nuances for each country."}},"data_access":
			 {"dataset_use":{"cit_req":"UNHCR (2021). Argentina: High Frequency Survey - Q1Q2 2021. Accessed from: https:\/\/microdata.unhcr.org","disclaimer":"The user of
			  the data acknowledges that the original collector of the data, the authorized distributor of the data, and the relevant funding agency bear no responsibility
			   for use of the data or for interpretations or inferences based upon such uses."}}},"schematype":"survey"}}}}`
