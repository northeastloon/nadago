package nadago

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(expectedSearchResponse)))
		}))
		defer ts.Close()

		client := NewClient(ts.URL)

		ctx := context.Background()

		params := NewDefaultSearchParams()

		surveys, err := client.Search(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, len(surveys), 5)

		idno := "ALB_2020_ES-COVID19-R1_v01_M"
		created, _ := time.Parse(time.RFC3339, "2022-05-11T11:14:45+00:00")
		changed, _ := time.Parse(time.RFC3339, "2022-05-11T11:14:46+00:00")

		assert.Equal(t, idno, surveys[0].Idno)
		assert.Equal(t, created, surveys[0].Created)
		assert.Equal(t, changed, surveys[0].Changed)
		assert.NotNil(t, surveys[0].Data)

	})

	t.Run("bad request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
		}))
		defer ts.Close()

		client := NewClient(ts.URL)

		ctx := context.Background()

		params := NewDefaultSearchParams()

		_, err := client.Search(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, FetchErr{
			Message:    "non-200 status code from the API",
			StatusCode: 404,
		}, err)
	})

	t.Run("parameters created correctly", func(t *testing.T) {

		params := NewDefaultSearchParams()

		ptype := reflect.TypeOf(params)
		if ptype.Kind() == reflect.Ptr {
			ptype = ptype.Elem()
		}
		numFields := ptype.NumField()

		assert.Equal(t, 12, numFields)
		assert.Equal(t, 30, params.Ps)

	})

	t.Run("failed to complete http request", func(t *testing.T) {

		failingClient := &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("failed to complete http request")
			}),
		}
		client := NewClient("http://invalid-url", WithHTTPClient(failingClient))

		ctx := context.Background()
		params := NewDefaultSearchParams()

		_, err := client.Search(ctx, params)
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
		params := NewDefaultSearchParams()

		_, err := client.Search(ctx, params)
		assert.Error(t, err)
		assert.IsType(t, AppErr{}, err)
	})
}

var expectedSearchResponse = `{"result":{"rows":[{"idno":"ALB_2020_ES-COVID19-R1_v01_M","formid":5,"form_model":"remote","title":"Enterprise Survey Follow-up on COVID-19 2020, Round 1","nation":"Albania","year_start":2020,"year_end":2020,"repositoryid":"central","created":"2022-05-11T11:14:45+00:00","changed":"2022-05-11T11:14:46+00:00","varcount":85,"total_views":125,"authoring_entity":"World Bank Group","total_downloads":6,"rank":1,"type":"survey","id":10252,"url":"https:\/\/catalog.ihsn.org\/catalog\/10252"},{"idno":"WLD_2020_CTIS_v01_M","formid":5,"form_model":"remote","title":"COVID-19 Trends and Impact Survey (2020-Ongoing)","nation":"Afghanistan, Albania, Algeria, Angola, Argentina, Armenia, Australia, Austria, Azerbaijan, Banglades","year_start":2020,"year_end":2021,"repositoryid":"central","created":"2021-11-03T19:32:10+00:00","changed":"2021-11-03T19:34:21+00:00","varcount":0,"total_views":293,"authoring_entity":"Facebook Data for Good, Carnegie Mellon University, University of Maryland","total_downloads":0,"rank":1,"type":"survey","id":9884,"url":"https:\/\/catalog.ihsn.org\/catalog\/9884"},{"idno":"WLD_2020_FBS_v01_M","formid":5,"form_model":"remote","title":"Future of Business Survey 2020","nation":"Albania, Algeria, American Samoa...and 176 more","year_start":2020,"year_end":2020,"repositoryid":"central","created":"2021-12-08T21:42:48+00:00","changed":"2022-06-14T13:16:49+00:00","varcount":0,"total_views":160,"authoring_entity":"Facebook, The Organisation for Economic Co-operation and Development (OECD), World Bank","total_downloads":5,"rank":1,"type":"survey","id":9891,"url":"https:\/\/catalog.ihsn.org\/catalog\/9891"},{"idno":"ALB_2020_WBCS_v01_M","formid":5,"form_model":"remote","title":"World Bank Group Country Survey 2020","nation":"Albania","year_start":2020,"year_end":2020,"repositoryid":"central","created":"2021-01-19T01:55:01+00:00","changed":"2021-01-19T01:55:01+00:00","varcount":289,"total_views":394,"authoring_entity":"Public Opinion Research Group","total_downloads":31,"rank":1,"type":"survey","id":9523,"url":"https:\/\/catalog.ihsn.org\/catalog\/9523"},{"idno":"ALB_2020_FIES_v01_M_v01_A_OCS","formid":5,"title":"Food Insecurity Experience Scale 2020","nation":"Albania","authoring_entity":"FAO Statistics Division","form_model":"remote","year_start":2020,"year_end":2020,"repositoryid":"central","link_da":"https:\/\/microdata.fao.org\/index.php\/catalog\/1921","created":"2023-01-25T16:22:45+00:00","changed":"2023-01-25T16:22:45+00:00","varcount":0,"total_views":0,"total_downloads":0,"rank":1,"type":"survey","id":10987,"url":"https:\/\/catalog.ihsn.org\/catalog\/10987"}],"found":5,"total":10174,"limit":15,"offset":0,"search_counts_by_type":{"survey":5},"page":1}}`
