package tl

import (
	"context"
	"dsc/inbrief/scraper/config"
	"dsc/inbrief/scraper/pkg/log"
	"fmt"
	"path/filepath"

	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
)

func InitClient(ctx context.Context, cfg config.Config) *client.Client {
	tdlibParameters := &client.SetTdlibParametersRequest{
		UseTestDc:           false,
		DatabaseDirectory:   filepath.Join(".tdlib", "database"),
		FilesDirectory:      filepath.Join(".tdlib", "files"),
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseMessageDatabase:  true,
		UseSecretChats:      false,
		ApiId:               cfg.Telegram.ApiId,
		ApiHash:             cfg.Telegram.ApiHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
	}
	// client authorizer
	authorizer := client.ClientAuthorizer(tdlibParameters)
	go client.CliInteractor(authorizer)

	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		log.L.Fatal("SetLogVerbosityLevel", zap.Error(err))
	}

	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		log.L.Fatal("NewClient error", zap.Error(err))
	}

	versionOption, err := client.GetOption(&client.GetOptionRequest{
		Name: "version",
	})
	if err != nil {
		log.L.Fatal("GetOption error", zap.Error(err))
	}

	commitOption, err := client.GetOption(&client.GetOptionRequest{
		Name: "commit_hash",
	})
	if err != nil {
		log.L.Fatal("GetOption", zap.Error(err))
	}

	me, err := tdlibClient.GetMe()
	if err != nil {
		log.L.Fatal("GetMe error", zap.Error(err))
	}

	log.L.Info("TDLib loaded",
		zap.String(
			"version",
			versionOption.(*client.OptionValueString).Value,
		),
		zap.String(
			"commit",
			commitOption.(*client.OptionValueString).Value,
		),
		zap.String(
			"me",
			fmt.Sprintf("%s %s", me.FirstName, me.LastName),
		),
	)

	return tdlibClient

}

type ClientService struct {
	client *client.Client
}
