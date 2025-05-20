package server

import (
	"fmt"

	"dsc/inbrief/scraper/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "dsc/inbrief/scraper/docs"

)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
}

func New(cfg *config.Config) *Server {
	r := gin.Default()
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r.GET("/health", health)

	r.GET("/scrape", scrape)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	return &Server{
		cfg:    cfg,
		router: r,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	return s.router.Run(addr)
}
