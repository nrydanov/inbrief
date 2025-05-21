package streaming

import "github.com/nrydanov/inbrief/internal"

func SubscribeToChat(chId int64, state *internal.AppState) error {
	tlClient := state.TlClient

	for {
		// updates, err := tlClient.GetListener().Updates(
		// 	&tlClient.GetChatUpdatesRequest{
		// 		ChatId: chId,
		// 	},
		// )

		// if err != nil {
		// 	zap.L().Error("GetChatUpdates", zap.Error(err))
		// 	return err
		// }

		// for _, update := range updates.Updates {
		// 	switch msg := update.Update.(type) {
		// 	case *tlClient.UpdateChatMessage:
		// 		zap.L().Debug(msg.Message.Content.GetMessageText())
		// 	}
		// }
	}
}
