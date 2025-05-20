package tl

import "github.com/zelenin/go-tdlib/client"
import tl "dsc/inbrief/scraper/pkg/tl"

func (s TlService) GetChannelsFromLink(link string) ([]tl.ChatId, error) {
	info, err := s.Client.CheckChatFolderInviteLink(
		&client.CheckChatFolderInviteLinkRequest{
			InviteLink: link,
		},
	)
	if err != nil {
		return nil, err
	}

	ids := make([]tl.ChatId, len(info.AddedChatIds))

	for i, id := range info.AddedChatIds {
		ids[i] = tl.ChatId(id)
	}

	return ids, nil
}

func (s TlService) GetChat(id tl.ChatId) (*client.Chat, error) {
	info, err := s.Client.GetChat(
		&client.GetChatRequest{
			ChatId: int64(id),
		},
	)
	if err != nil {
		return nil, err
	}

	return info, nil
}
