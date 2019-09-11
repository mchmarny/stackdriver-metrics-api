package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFilePath = "sample/test.json"
)

func runRequest(r http.Handler, method, path string,
	body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestMetricCounterHandler(t *testing.T) {

	testJSON := `[
					{
						"type": "User",
						"site_admin": false,
						"id": 1
					},
					{
						"type": "User",
						"site_admin": false,
						"id": 2
					},
					{
						"type": "User",
						"site_admin": false,
						"id": 3
					}
				]`

	router := setupRouter(true)
	w := runRequest(router, "POST", "/v1/counter/test",
		bytes.NewBuffer([]byte(testJSON)))

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("Response: %v", w.Body)
}
