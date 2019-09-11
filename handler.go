package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func metricCounterHandler(c *gin.Context) {

	metric := c.Param("metric")
	logger.Printf("Message: %s", metric)
	if metric == "" {
		logger.Println("Error on nil metric parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Null Argument",
			"status":  http.StatusBadRequest,
		})
		return
	}

	var content []map[string]interface{}
	if err := c.BindJSON(&content); err != nil {
		logger.Printf("Error while parsing posted content: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Content",
			"status":  http.StatusBadRequest,
		})
		return
	}

	resp := &ResponseObject{
		RequestID: newID(),
		Timestamp: time.Now().UTC(),
		Metric:    metric,
		Result:    int64(len(content)),
	}

	if err := post(c.Request.Context(),
		resp.Metric, resp.Timestamp, resp.Result); err != nil {
		logger.Printf("Error posting metrics: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error processing metric",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, resp)

}

func defaultRequestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"counter": "POST /v1/counter/:metric",
		"status":  http.StatusOK,
	})
}

// ResponseObject represents body of the request response
type ResponseObject struct {
	RequestID string    `json:"id"`
	Timestamp time.Time `json:"ts"`
	Metric    string    `json:"metric"`
	Result    int64     `json:"result,omitempty"`
}

func newID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		logger.Fatalf("Error while getting id: %v\n", err)
	}
	return id.String()
}
