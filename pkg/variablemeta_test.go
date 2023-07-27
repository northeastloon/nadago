package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVarMeta(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(expectedVarMetaResponse)))
		}))
		defer ts.Close()

		client := NewClient(ts.URL)

		ctx := context.Background()
		idno := "LBR_2020_FIES_v01_M_v01_A_OCS"
		vid := "V1"

		varmeta, err := client.GetVarMeta(ctx, idno, vid)
		assert.NoError(t, err)
		assert.Equal(t, idno, varmeta.Idno)
		assert.Equal(t, vid, varmeta.Vid)
		assert.NotNil(t, varmeta.Data)
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
		vid := "V1"

		_, err := client.GetVarMeta(ctx, idno, vid)
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
		vid := "V1"

		_, err := client.GetVarMeta(ctx, idno, vid)
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
		vid := "V1"

		_, err := client.GetVarMeta(ctx, idno, vid)
		assert.Error(t, err)
		assert.IsType(t, AppErr{}, err)
	})

}

var expectedVarMetaResponse = `{"variable":{"uid":"19260146","sid":"10894","fid":"F1","vid":"V1","name":"Random_ID","labl":"Unique respondent identifier","qstn":null,"catgry":null,"metadata":{"file_id":"F1","vid":"V1","name":"Random_ID","var_intrvl":"contin","var_dcml":"0","var_wgt":null,"var_is_wgt":null,"loc_start_pos":null,"loc_end_pos":null,"loc_width":null,"loc_rec_seg_no":null,"labl":"Unique respondent identifier","var_imputation":null,"var_security":null,"var_resp_unit":null,"var_analysis_unit":null,"var_qstn_preqtxt":null,"var_qstn_qstnlit":null,"var_qstn_postqtxt":null,"var_qstn_ivuinstr":null,"var_universe":null,"var_universe_clusion":null,"var_sumstat":[{"value":"0","type":"vald","wgtd":null},{"value":"0","type":"invd","wgtd":null}],"var_txt":null,"var_catgry":[],"var_codinstr":null,"var_concept":[],"var_format":{"type":"numeric","schema":"other","category":null,"name":null},"var_notes":null,"var_val_range":{"min":null,"max":null},"fid":"F1"},"keywords":null}}`
