package web

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type ObservabilityWeb struct {
}

func (o *ObservabilityWeb) RegisterRoutes(s *gin.Engine) {
	s.GET("/test/metric", func(context *gin.Context) {
		sleep := rand.Int31n(1000)
		time.Sleep(time.Millisecond * time.Duration(sleep))
		context.String(http.StatusOK, "ok")
	})
}
