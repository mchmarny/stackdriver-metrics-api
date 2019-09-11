package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

const (
	testFilePath = "sample/test.json"
)

func runRequest(r http.Handler, method, path string,
	body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestMetricCounterHandler(t *testing.T) {

	testFile, err := os.Open(testFilePath)
	if err != nil {
		t.Fatalf("test file not found: %s", testFilePath)
	}
	defer testFile.Close()

	content, _ := ioutil.ReadAll(testFile)

	router := setupRouter(true)
	w := runRequest(router, "POST", "/v1/counter/test", bytes.NewBuffer(content))

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
}
