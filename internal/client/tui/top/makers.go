// Package top содержит функцию prepareMakers, которая инициализирует карту создателей экранов для TUI-приложения.
package top

import (
	"github.com/romanp1989/gophkeeper/internal/client/grpc"
	"github.com/romanp1989/gophkeeper/internal/client/tui"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens/auth"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens/blobs"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens/cards"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens/credentials"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens/remotes"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens/secrets"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens/storage"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens/texts"
)

func prepareMakers(client grpc.ClientGRPCInterface) map[tui.Screen]tui.ScreenMaker {
	return map[tui.Screen]tui.ScreenMaker{
		tui.BlobEditScreen:       &blobs.BlobEditScreen{},
		tui.CardEditScreen:       &cards.CardEditScreen{},
		tui.CredentialEditScreen: &credentials.CredentialEditScreen{},
		tui.FilePickScreen:       &blobs.FilePickScreen{},
		tui.LoginScreen:          &auth.AuthenticateScreen{},
		tui.RemoteOpenScreen:     &remotes.RemoteOpenScreenMaker{Client: client},
		tui.SecretTypeScreen:     &secrets.SecretTypeScreen{},
		tui.StorageBrowseScreen:  &storage.BrowseStorageScreen{},
		tui.TextEditScreen:       &texts.TextEditScreen{},
	}
}
