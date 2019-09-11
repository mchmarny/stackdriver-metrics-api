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

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// routes
	r.GET("/", defaultRequestHandler)
	r.GET("/health", healthHandler)

	// api
	v1 := r.Group("/v1")
	{
		v1.GET("/counter/:metric", metricCounterHandler)
	}

	// server
	hostPost := net.JoinHostPort("0.0.0.0", port)
	logger.Printf("Server starting: %s \n", hostPost)
	if err := r.Run(hostPost); err != nil {
		logger.Fatal(err)
	}
}
