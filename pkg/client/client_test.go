package client

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {

	t.Run("test with httpclient works", func(*testing.T) {
		myClient := NewClient(
			"http://example.com/api",
			WithHTTPClient(&http.Client{
				Timeout: 1 * time.Second,
			}),
		)
		assert.Equal(t, "my-test-url", myClient.apiURL)
		assert.Equal(t, 1*time.Second, myClient.httpClient.Timeout)
	})

	t.Run("test with httpclient works", func(*testing.T) {
		baseURL := "http://example.com/api"
		username := "testuser"
		password := "testpassword"

		client := NewClient(baseURL, WithURLAuth(username, password))

		expectedURL := fmt.Sprintf("%s?username=%s&password=%s", baseURL, url.QueryEscape(username), url.QueryEscape(password))

		assert.Equal(t, expectedURL, client.apiURL, "apiURL should have the correct authentication parameters")
	})

}
