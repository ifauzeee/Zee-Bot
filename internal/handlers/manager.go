package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type HandlerFunc func(c *Context) error

type Command struct {
	Name        string
	Description string
	Handler     HandlerFunc
}

type Manager struct {
	Commands map[string]Command
	Logger   *zap.Logger
}

type Context struct {
	Ctx    context.Context
	Sender *message.Sender
	Msg    *tg.Message
	Logger *zap.Logger
	Args   []string
	Peer   tg.InputPeerClass
}

func (c *Context) Edit(text string) error {
	_, err := c.Sender.To(c.Peer).Edit(c.Msg.ID).Text(c.Ctx, text)
	return err
}

func (c *Context) Reply(text string) error {
	_, err := c.Sender.To(c.Peer).Reply(c.Msg.ID).Text(c.Ctx, text)
	return err
}

func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		Commands: make(map[string]Command),
		Logger:   logger,
	}
}

func (m *Manager) Register(name, description string, handler HandlerFunc) {
	m.Commands[name] = Command{
		Name:        name,
		Description: description,
		Handler:     handler,
	}
}

func (m *Manager) HandleNewMessage(ctx context.Context, sender *message.Sender, update tg.UpdateClass, entities tg.Entities) error {
	switch u := update.(type) {
	case *tg.UpdateNewMessage:
		msg, ok := u.Message.(*tg.Message)
		if !ok {
			return nil
		}
		if msg.Out {
			return m.processMessage(ctx, sender, msg, entities)
		}
	case *tg.UpdateNewChannelMessage:
		msg, ok := u.Message.(*tg.Message)
		if !ok {
			return nil
		}
		if msg.Out {
			return m.processMessage(ctx, sender, msg, entities)
		}
	}
	return nil
}

func (m *Manager) processMessage(ctx context.Context, sender *message.Sender, msg *tg.Message, entities tg.Entities) error {
	text := msg.Message
	if len(text) < 2 {
		return nil
	}

	prefix := text[0]
	if prefix != '.' && prefix != '!' {
		return nil
	}

	parts := strings.Fields(text[1:])
	if len(parts) == 0 {
		return nil
	}
	cmdName := strings.ToLower(parts[0])

	if cmd, exists := m.Commands[cmdName]; exists {
		peer, err := m.resolvePeer(msg.PeerID, entities)
		if err != nil {
			m.Logger.Warn("Could not resolve peer", zap.Error(err))
			return nil
		}

		m.Logger.Info("Executing command", zap.String("command", cmdName))

		c := &Context{
			Ctx:    ctx,
			Sender: sender,
			Msg:    msg,
			Logger: m.Logger,
			Args:   parts[1:],
			Peer:   peer,
		}

		if err := cmd.Handler(c); err != nil {
			m.Logger.Error("Command execution failed", zap.Error(err))
		}
	}

	return nil
}

func (m *Manager) resolvePeer(peer tg.PeerClass, entities tg.Entities) (tg.InputPeerClass, error) {
	switch p := peer.(type) {
	case *tg.PeerUser:
		user, ok := entities.Users[p.UserID]
		if !ok {
			return nil, fmt.Errorf("user not found")
		}
		return &tg.InputPeerUser{
			UserID:     user.ID,
			AccessHash: user.AccessHash,
		}, nil
	case *tg.PeerChat:
		return &tg.InputPeerChat{
			ChatID: p.ChatID,
		}, nil
	case *tg.PeerChannel:
		channel, ok := entities.Channels[p.ChannelID]
		if !ok {
			return nil, fmt.Errorf("channel not found")
		}
		return &tg.InputPeerChannel{
			ChannelID:  channel.ID,
			AccessHash: channel.AccessHash,
		}, nil
	}
	return nil, nil
}
