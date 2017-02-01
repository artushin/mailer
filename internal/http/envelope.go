package http

import (
	"github.com/artushin/mailer/internal/pdf"
	"github.com/gin-gonic/gin"
	"log"
)

func envelopeHandler(c *gin.Context) {
	reqi, ok := c.Get("request")
	if !ok || reqi == nil {
		log.Printf("no request present\n")
		c.AbortWithStatus(500)
	}
	req := reqi.(*Request)

	c.Writer.Header().Set("content-type", "application/pdf")
	if err := pdf.GenerateEnvelope(&pdf.Envelope{
		From: req.Address,
		To:   req.OfficialAddress,
	}, c.Writer); err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.Status(200)
}
