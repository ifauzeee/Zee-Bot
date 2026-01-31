package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"zee-ubot/internal/config"
	"zee-ubot/internal/database"
	"zee-ubot/internal/handlers"
	"zee-ubot/internal/middleware"
	"zee-ubot/internal/modules/admin"
	"zee-ubot/internal/modules/basic"
	"zee-ubot/internal/modules/utility"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

func prompt(msg string) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		fmt.Print(msg)
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(text), nil
	}
}

type TermAuth struct {
	PhoneFn func(context.Context) (string, error)
	CodeFn  func(context.Context) (string, error)
	PassFn  func(context.Context) (string, error)
}

func (t TermAuth) Phone(ctx context.Context) (string, error)                    { return t.PhoneFn(ctx) }
func (t TermAuth) Password(ctx context.Context) (string, error)                 { return t.PassFn(ctx) }
func (t TermAuth) Code(ctx context.Context, _ *tg.AuthSentCode) (string, error) { return t.CodeFn(ctx) }
func (t TermAuth) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("signup not supported")
}
func (t TermAuth) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	if err := database.Init(); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	manager := handlers.NewManager(logger)
	basic.Register(manager)
	admin.Register(manager)
	utility.Register(manager)

	var storage session.Storage

	if cfg.SessionString != "" {
		logger.Info("Using Session String from Environment")
		data, err := base64.StdEncoding.DecodeString(cfg.SessionString)
		if err != nil {
			logger.Fatal("Invalid Session String", zap.Error(err))
		}

		storage = &StaticMemoryStorage{Data: data}
	} else {
		sessionDir := "session"
		if err := os.MkdirAll(sessionDir, 0700); err != nil {
			logger.Fatal("Failed to create session dir", zap.Error(err))
		}
		sessionPath := filepath.Join(sessionDir, "session.json")
		storage = &telegram.FileSessionStorage{Path: sessionPath}
	}

	var api *tg.Client

	dispatcher := tg.NewUpdateDispatcher()
	dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		if api == nil {
			return nil
		}
		sender := message.NewSender(api)
		return manager.HandleNewMessage(ctx, sender, api, update, e)
	})
	dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		if api == nil {
			return nil
		}
		sender := message.NewSender(api)
		return manager.HandleNewMessage(ctx, sender, api, update, e)
	})

	gaps := updates.New(updates.Config{
		Handler: dispatcher,
		Logger:  logger.Named("updates"),
	})

	client := telegram.NewClient(cfg.AppID, cfg.AppHash, telegram.Options{
		Logger:         logger,
		SessionStorage: storage,
		UpdateHandler:  gaps,
		Middlewares: []telegram.Middleware{
			middleware.FloodWait(logger),
		},
	})

	api = client.API()

	flow := auth.NewFlow(
		TermAuth{
			PhoneFn: prompt("Enter Phone Number: "),
			CodeFn:  prompt("Enter Code: "),
			PassFn:  prompt("Enter Password (if known): "),
		},
		auth.SendCodeOptions{},
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger.Info("Starting Zee-Ubot...")

	err = client.Run(ctx, func(ctx context.Context) error {
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		self, err := client.Self(ctx)
		if err != nil {
			return err
		}
		logger.Info("Logged in!",
			zap.String("first_name", self.FirstName),
			zap.String("username", self.Username),
			zap.Int64("id", self.ID),
		)

		return gaps.Run(ctx, client.API(), self.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				logger.Info("Update manager started")
			},
		})
	})

	if err != nil {
		logger.Fatal("Error running client", zap.Error(err))
	}
}

type StaticMemoryStorage struct {
	Data []byte
}

func (s *StaticMemoryStorage) LoadSession(ctx context.Context) ([]byte, error) {
	if len(s.Data) == 0 {
		return nil, session.ErrNotFound
	}
	c := make([]byte, len(s.Data))
	copy(c, s.Data)
	return c, nil
}

func (s *StaticMemoryStorage) StoreSession(ctx context.Context, data []byte) error {
	s.Data = data
	return nil
}
