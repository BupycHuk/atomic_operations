package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"encoding/json"
	"bytes"
	"sync"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	body, err := json.Marshal(KeyValue{Key:"test", Value: 15000})
	if err != nil {
		t.Fatal(err)
	}
	waitGroup := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {

		waitGroup.Add(1)
		go func() {
			b := bytes.NewBuffer(body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/map", b)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/map/test", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)

	assert.Equal(t, "{\"key\":\"test\",\"value\":1500000000}", w.Body.String())
}