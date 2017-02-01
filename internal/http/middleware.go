package http

import (
	"encoding/json"
	"fmt"
	"github.com/artushin/mailer/internal/address"
	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Name            string
	Address         []string
	OfficialAddress []string
}

type SunlightResponse struct {
	Results []*SunlightResult `json:"results"`
}

type SunlightResult struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Office    string `json:"office"`
}

func middleware(client *maps.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.ParseForm()
		name := c.Request.Form.Get("name")
		if len(name) == 0 {
			c.AbortWithStatus(400)
			return
		}

		addr := c.Request.Form.Get("address")
		if len(addr) == 0 {
			c.AbortWithStatus(400)
			return
		}

		a, err := address.Get(c, name, addr, client)
		if err != nil {
			if err == address.ErrNoAddress || err == address.ErrNotUS {
				c.AbortWithStatus(400)
				return
			}
			c.AbortWithStatus(500)
			return
		}

		resp, err := http.Get(fmt.Sprintf("https://congress.api.sunlightfoundation.com/legislators/locate?latitude=%f&longitude=%f", a.Latitude, a.Longitude))
		if err != nil {
			log.Printf("unable to get from sunlight: %v\n", err)
			c.AbortWithStatus(500)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("got invalid status from sunlight: %v\n", err)
			c.AbortWithStatus(resp.StatusCode)
			return
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("unable to read sunlight response: %v\n", err)
			c.AbortWithStatus(500)
			return
		}

		var sunlightResponse *SunlightResponse
		if err := json.Unmarshal(b, &sunlightResponse); err != nil {
			log.Printf("unable to read sunlight response: %v\n", err)
			c.AbortWithStatus(500)
			return
		}

		if len(sunlightResponse.Results) == 0 {
			log.Printf("no results from sunlight\n")
			c.AbortWithStatus(400)
		}

		oa, err := address.Get(c, fmt.Sprintf("%s %s", sunlightResponse.Results[0].FirstName, sunlightResponse.Results[0].LastName), sunlightResponse.Results[0].Office, client)
		if err != nil {
			if err == address.ErrNoAddress || err == address.ErrNotUS {
				c.AbortWithStatus(400)
				return
			}
			c.AbortWithStatus(500)
			return
		}

		c.Set("request", &Request{
			Name:            name,
			Address:         a.Address,
			OfficialAddress: oa.Address,
		})
	}
}
