package nadago

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNewClient(t *testing.T) {

	t.Run("test with httpclient works", func(*testing.T) {
		myClient := NewClient(
			"my-test-url",
			WithHTTPClient(&http.Client{
				Timeout: 1 * time.Second,
			}),
		)
		assert.Equal(t, "my-test-url", myClient.apiURL)
		assert.Equal(t, 1*time.Second, myClient.httpClient.Timeout)
	})

	t.Run("test with httpclient works", func(*testing.T) {
		baseURL := "my-test-url"
		username := "testuser"
		password := "testpassword"

		client := NewClient(baseURL, WithURLAuth(username, password))

		expectedURL := fmt.Sprintf("%s?username=%s&password=%s", baseURL, url.QueryEscape(username), url.QueryEscape(password))

		assert.Equal(t, expectedURL, client.apiURL, "apiURL should have the correct authentication parameters")
	})

}
