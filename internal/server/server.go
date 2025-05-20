package server

import (
	"time"

	"dsc/inbrief/scraper/config"
	"dsc/inbrief/scraper/pkg/log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "dsc/inbrief/scraper/docs"

	ginzap "github.com/gin-contrib/zap"
)

func NewServer(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(ginzap.Ginzap(log.L, time.RFC3339, true))
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	{
		r.GET("/health", health)

		r.GET("/scrape", scrape)

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.GET("/docs", func(c *gin.Context) {
			c.Redirect(302, "/swagger/index.html")
		})
	}

	return r
}
