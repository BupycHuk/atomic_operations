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
	"fmt"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	body, err := json.Marshal([]KeyValue{
		{Key: "test", Value: 15000},
	})

	body2, err := json.Marshal([]KeyValue{
		{Key: "test", Value: 15000},
		{Key: "test2", Value: -25000},
	})

	body3, err := json.Marshal([]KeyValue{
		{Key: "test", Value: -35000},
	})


	if err != nil {
		t.Fatal(err)
	}
	waitGroup := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {

		waitGroup.Add(1)
		go func() {
			send(bytes.NewBuffer(body), router, t)
			send(bytes.NewBuffer(body2), router, t)
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
	checkResult(router, t, 1500000000)
	for i := 0; i < 100000; i++ {

		waitGroup.Add(1)
		go func() {
			send(bytes.NewBuffer(body3), router, t)
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
	checkResult(router, t, 5000)
}
func checkResult(router *gin.Engine, t *testing.T, expected int) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/map/test", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)
	assert.Equal(t, fmt.Sprintf("{\"key\":\"test\",\"value\":%d}", expected), w.Body.String())
}

func send(b *bytes.Buffer, router *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/map", b)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}