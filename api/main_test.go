package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {

	testcases := map[string]struct {
		code int
		// message string
		endpoint string
		httpVerb string
	}{
		"Doesn't exist": {httpVerb: "GET", code: http.StatusNotFound, endpoint: "non-existant"},
		"Races":         {httpVerb: "GET", code: http.StatusOK, endpoint: "racing/2"},
		"Wrong Verb":    {httpVerb: "PUT", code: http.StatusMethodNotAllowed, endpoint: "racing"},
		"Tags":          {httpVerb: "GET", code: http.StatusOK, endpoint: "tags"},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest(tc.httpVerb, "http://localhost/"+tc.endpoint, nil)
			log.Printf("URL: %s", request.URL)
			response := httptest.NewRecorder()
			Routes().ServeHTTP(response, request)
			log.Printf("Response: %#+v", response)
			log.Printf("Body: %s", response.Body)
			assert.Equal(t, tc.code, response.Code, "Unexpected response code")
		})
	}
}
