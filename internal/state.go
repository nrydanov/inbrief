package internal

import "github.com/zelenin/go-tdlib/client"

type AppState struct {
	TlClient *client.Client
	Server
}
