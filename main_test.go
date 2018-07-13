package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"encoding/json"
	"bytes"
	"sync"
	"github.com/gin-gonic/gin"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	body, err := json.Marshal([]KeyValue{
		{Key: "test", Value: 15000000},
	})

	body2, err := json.Marshal([]KeyValue{
		{Key: "test", Value: 15000},
		{Key: "test2", Value: -25000},
	})
	if err != nil {
		t.Fatal(err)
	}
	waitGroup := sync.WaitGroup{}
	for i := 0; i < 100; i++ {

		waitGroup.Add(1)
		go func() {
			send(bytes.NewBuffer(body), router, t)
			send(bytes.NewBuffer(body2), router, t)
			waitGroup.Done()
		}()
	}
	for i := 0; i < 100; i++ {

		waitGroup.Add(1)
		go func() {
			send(bytes.NewBuffer(body2), router, t)
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
func send(b *bytes.Buffer, router *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/map", b)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}