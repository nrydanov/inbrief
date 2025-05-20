package server

import (
	"dsc/inbrief/scraper/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health godoc
// @Summary      Health check
// @Description  Returns status ok
// @Tags         health
// @Produce      json
// @Success      200
// @Router       /health [get]
func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

// Scrape godoc
// @Summary      Scrape request
// @Description  Handle scrape request with query parameters
// @Tags         scrape
// @Produce      json
// @Param        chat_folder_link  query     string    false  "Chat folder link"  default(https://t.me/addlist/4eTLcGIrIx9hNzUy)
// @Param        right_bound       query     string    true   "Right bound datetime"  format(date-time)  default(2025-05-20T15:00:00+04:00)
// @Param        left_bound        query     string    true   "Left bound datetime"   format(date-time)  default(2025-05-18T15:00:00+04:00)
// @Param        social            query     bool      false  "Social flag" default(false)
// @Success      200
// @Failure      400
// @Router       /scrape [get]
func scrape(c *gin.Context) {
	var req models.ScrapeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
