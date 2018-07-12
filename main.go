package main

import (
	"github.com/gin-gonic/gin"
	"math/big"
)

var db = Map{
	m: map[string]**big.Int{},
}

type KeyValue struct {
	Key string `json:"key" binding:"required"`
	Value int64 `json:"value" binding:"required"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	gin.DisableConsoleColor()
	r := gin.New()


	// Ping test
	r.GET("/map/:key", func(c *gin.Context) {
		key := c.Params.ByName("key")
		value, ok := db.m[key]
		if ok {
			c.JSON(200, gin.H{"key": key, "value": value })
		} else {
			c.JSON(200, gin.H{"key": key, "status": "no value"})
		}
	})

	r.POST("/map", func(c *gin.Context) {
		// Parse JSON
		var json KeyValue
		err := c.Bind(&json)
		if err != nil {
			c.JSON(400, gin.H{"status": "bad request", "err": err.Error()})
		}

		result := db.AddBigInt(json.Key, *big.NewInt(json.Value))

		c.JSON(200, gin.H{"status": "ok", "result": result.Text(10)})
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
