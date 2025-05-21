package tl

import (
	"github.com/zelenin/go-tdlib/client"
)

type ChatId int64

func ExtractChatIds(info *client.ChatFolderInviteLinkInfo) []ChatId {
	ids := make([]ChatId, len(info.AddedChatIds))

	for i, id := range info.AddedChatIds {
		ids[i] = ChatId(id)
	}

	return ids
}
