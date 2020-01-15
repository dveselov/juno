package providers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUCallerSendSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintln(w, `
			{
				"status": true,
				"ucaller_id": 103000,
				"phone": 79991234567,
				"code": 7777
			}
		`)
	}))
	defer testServer.Close()
	provider := UCallerProvider{
		ServiceID:  "service_id",
		SecretKey:  "secret_key",
		BaseAPIUrl: testServer.URL,
	}
	err := provider.Send("79818246403", "1337")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUCallerSendError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintln(w, `
			{
				"status": false,
				"error": "Test case error message"
			}
		`)
	}))
	defer testServer.Close()
	provider := UCallerProvider{
		ServiceID:  "service_id",
		SecretKey:  "secret_key",
		BaseAPIUrl: testServer.URL,
	}
	err := provider.Send("79818246403", "1337")
	if err.Error() != "Test case error message" {
		t.Fatal(err)
	}
}
