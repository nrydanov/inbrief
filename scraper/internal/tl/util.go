package tl

import (
	"errors"
	"fmt"

	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
)

type ChatId int64

func ExtractChatIds(info *client.ChatFolderInviteLinkInfo) []ChatId {
	ids := make([]ChatId, len(info.AddedChatIds))

	for i, id := range info.AddedChatIds {
		ids[i] = ChatId(id)
	}

	return ids
}

func ExtractUsername(c *client.Client, chat *client.Chat) (string, error) {
	switch e := chat.Type.(type) {
	case *client.ChatTypeSupergroup:
		group, err := c.GetSupergroup(&client.GetSupergroupRequest{
			SupergroupId: e.SupergroupId,
		})
		if err != nil {
			zap.L().Debug("Unable to convert chat to supergroup", zap.Error(err))
			return "", err
		}
		if group.Usernames == nil {
			return fmt.Sprintf("%d", chat.Id), nil
		}
		return group.Usernames.ActiveUsernames[0], nil
	default:
		return "", errors.New("unsupported chat type")
	}
}
