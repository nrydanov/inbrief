package server

import (
	"context"
	pb "dsc/inbrief/scraper/pkg/proto"
	"dsc/inbrief/scraper/pkg/tl"
	"fmt"

	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
)

func (s *server) Fetch(
	ctx context.Context,
	req *pb.FetchRequest,
) (*pb.FetchResponse, error) {
	state := s.state

	info, err := state.TlClient.CheckChatFolderInviteLink(
		&client.CheckChatFolderInviteLinkRequest{
			InviteLink: req.ChatFolderLink,
		},
	)
	if err != nil {
		return nil, err
	}

	ids := make([]tl.ChatId, len(info.AddedChatIds))

	for i, id := range info.AddedChatIds {
		ids[i] = tl.ChatId(id)
	}

	zap.L().Debug("Scraping channels", zap.String("ids", fmt.Sprintf("%+v", ids)))

	for _, id := range ids {
		_, err := state.TlClient.GetChat(
			&client.GetChatRequest{
				ChatId: int64(id),
			},
		)
		if err != nil {
			return nil, err
		}

	}

	return &pb.FetchResponse{}, nil
}

// Health godoc
// @Summary      Health check
// @Description  Returns status ok
// @Tags         health
// @Produce      json
// @Success      200
// @Router       /health [get]
// func health(c *gin.Context) {
// 	c.JSON(200, gin.H{
// 		"status": "ok",
// 	})
// }

// Scrape godoc
// @Summary      Scrape request
// @Description  Handle scrape request with query parameters
// @Tags         scrape
// @Produce      json
// @Param        chat_folder_link  query     string    false  "Chat folder link"  default(https://t.me/addlist/grg7NStE6881MDE6)
// @Param        right_bound       query     string    true   "Right bound datetime"  format(date-time)  default(2025-05-20T15:00:00+04:00)
// @Param        left_bound        query     string    true   "Left bound datetime"   format(date-time)  default(2025-05-18T15:00:00+04:00)
// @Param        social            query     bool      false  "Social flag" default(false)
// @Success      200
// @Failure      400
// @Router       /scrape [get]
// func scrape(c *gin.Context, state *internal.AppState) {
// 	var req models.ScrapeRequest
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	var err error

// 	defer func() {
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error": err.Error(),
// 			})
// 		} else {
// 			c.JSON(http.StatusOK, gin.H{
// 				"status": "ok",
// 			})
// 		}
// 	}()

// 	info, err := state.TlClient.CheckChatFolderInviteLink(
// 		&client.CheckChatFolderInviteLinkRequest{
// 			InviteLink: req.ChatFolderLink,
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}

// 	ids := make([]tl.ChatId, len(info.AddedChatIds))

// 	for i, id := range info.AddedChatIds {
// 		ids[i] = tl.ChatId(id)
// 	}

// 	zap.L().Debug("Scraping channels", zap.String("ids", fmt.Sprintf("%+v", ids)))

// 	for _, id := range ids {
// 		chat, err := state.TlClient.GetChat(
// 			&client.GetChatRequest{
// 				ChatId: int64(id),
// 			},
// 		)
// 		if err != nil {
// 			return
// 		}

// 	}
// }
