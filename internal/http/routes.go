package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func Start(ctx context.Context, port int) {
	mapsKeyBytes, err := ioutil.ReadFile("./google-maps-key")
	if err != nil {
		log.Fatalf("could not read api key: %v", err)
	}
	mapsKey := strings.Trim(string(mapsKeyBytes), " \n")
	mapsClient, err := maps.NewClient(
		maps.WithAPIKey(string(mapsKey)),
		maps.WithRateLimit(100),
	)
	if err != nil {
		log.Fatalf("could not set up maps client: %v", err)
	}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		defer wg.Done()
		wg.Add(1)
		c.Next()
	})
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	api := r.Group("/pdfs", middleware(mapsClient))
	api.GET("/letter", letterHandler)
	api.GET("/envelope", envelopeHandler)

	go r.Run(fmt.Sprintf(":%d", port))
}

func Close() <-chan struct{} {
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		wg.Wait()
	}()
	return done
}
