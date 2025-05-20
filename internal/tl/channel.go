package tl

import "github.com/zelenin/go-tdlib/client"
import tl "dsc/inbrief/scraper/pkg/tl"

func (s ClientService) GetChannelsFromLink(link string) ([]tl.ChatId, error) {
	info, err := s.client.CheckChatFolderInviteLink(
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
