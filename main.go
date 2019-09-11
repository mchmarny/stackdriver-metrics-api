package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/mchmarny/gcputil/env"
	"github.com/mchmarny/gcputil/project"
)

const (
	defaultPort = "8080"
)

var (
	release   = env.MustGetEnvVar("RELEASE", "v0.0.1")
	port      = env.MustGetEnvVar("PORT", defaultPort)
	projectID = project.GetIDOrFail()
	logger    = log.New(os.Stdout, "", 0)
)

func setupRouter(debug bool) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if !debug {
		r.Use(gin.Recovery())
	}

	r.GET("/", defaultRequestHandler)
	r.GET("/health", healthHandler)

	v1 := r.Group("/v1")
	{
		v1.POST("/counter/:metric", metricCounterHandler)
	}

   return r
}

func main() {
	hostPost := net.JoinHostPort("0.0.0.0", port)
	logger.Printf("Server starting: %s \n", hostPost)
	if err := setupRouter(false).Run(hostPost); err != nil {
		logger.Fatal(err)
	}
}
