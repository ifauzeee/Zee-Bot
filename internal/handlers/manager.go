package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

var moduleRegistry []func(*Manager)

func RegisterModule(f func(*Manager)) {
	moduleRegistry = append(moduleRegistry, f)
}

func GetModules() []func(*Manager) {
	return moduleRegistry
}

type HandlerFunc func(c *Context) error

type Command struct {
	Name        string
	Description string
	Handler     HandlerFunc
}

type HookFunc func(c *Context, update tg.UpdateClass) error

type Manager struct {
	Commands map[string]Command
	Hooks    map[string][]HookFunc
	Logger   *zap.Logger
	SelfID   int64
}

type Context struct {
	Ctx    context.Context
	Sender *message.Sender
	Msg    *tg.Message
	Logger *zap.Logger
	Args   []string
	Peer   tg.InputPeerClass
	Raw    *tg.Client
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
		Hooks:    make(map[string][]HookFunc),
		Logger:   logger,
		SelfID:   0,
	}
}

func (m *Manager) Register(name, description string, handler HandlerFunc) {
	m.Commands[name] = Command{
		Name:        name,
		Description: description,
		Handler:     handler,
	}
}

func (m *Manager) RegisterHook(event string, handler HookFunc) {
	m.Hooks[event] = append(m.Hooks[event], handler)
}

func (m *Manager) HandleNewMessage(ctx context.Context, sender *message.Sender, raw *tg.Client, update tg.UpdateClass, entities tg.Entities) error {
	var msg *tg.Message
	var ok bool

	switch u := update.(type) {
	case *tg.UpdateNewMessage:
		msg, ok = u.Message.(*tg.Message)
	case *tg.UpdateNewChannelMessage:
		msg, ok = u.Message.(*tg.Message)
	}

	if !ok || msg == nil {
		return m.HandleUpdate(ctx, sender, raw, update, entities)
	}

	if err := m.HandleUpdate(ctx, sender, raw, update, entities); err != nil {
		return err
	}

	if msg.Out {
		return m.processMessage(ctx, sender, raw, msg, entities)
	}

	return nil

}

func (m *Manager) HandleUpdate(ctx context.Context, sender *message.Sender, raw *tg.Client, update tg.UpdateClass, entities tg.Entities) error {
	var peer tg.InputPeerClass
	var msg *tg.Message

	switch u := update.(type) {
	case *tg.UpdateNewMessage:
		if messageVal, ok := u.Message.(*tg.Message); ok {
			msg = messageVal
			peer, _ = m.resolvePeer(ctx, raw, msg.PeerID, entities)
		}
	case *tg.UpdateNewChannelMessage:
		if messageVal, ok := u.Message.(*tg.Message); ok {
			msg = messageVal
			peer, _ = m.resolvePeer(ctx, raw, msg.PeerID, entities)
		}
	}

	c := &Context{
		Ctx:    ctx,
		Sender: sender,
		Logger: m.Logger,
		Peer:   peer,
		Raw:    raw,
		Msg:    msg,
	}

	for _, h := range m.Hooks["all"] {
		if err := h(c, update); err != nil {
			m.Logger.Error("Update hook failed", zap.Error(err))
		}
	}

	switch u := update.(type) {
	case *tg.UpdateChatParticipants:
		for _, h := range m.Hooks["chat_participants"] {
			_ = h(c, u)
		}
	case *tg.UpdateChannel:
		for _, h := range m.Hooks["channel"] {
			_ = h(c, u)
		}
	}

	return nil
}

func (m *Manager) processMessage(ctx context.Context, sender *message.Sender, raw *tg.Client, msg *tg.Message, entities tg.Entities) error {
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
		peer, err := m.resolvePeer(ctx, raw, msg.PeerID, entities)
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
			Raw:    raw,
		}

		safeExec := func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic recovered: %v", r)
					m.Logger.Error("Panic in command handler", zap.String("command", cmdName), zap.Any("panic", r))
				}
			}()
			return cmd.Handler(c)
		}

		if err := safeExec(); err != nil {
			m.Logger.Error("Command execution failed", zap.Error(err))
		}
	}

	return nil
}

func (m *Manager) resolvePeer(ctx context.Context, raw *tg.Client, peer tg.PeerClass, entities tg.Entities) (tg.InputPeerClass, error) {
	switch p := peer.(type) {
	case *tg.PeerUser:
		if user, ok := entities.Users[p.UserID]; ok {
			return &tg.InputPeerUser{
				UserID:     user.ID,
				AccessHash: user.AccessHash,
			}, nil
		}
		if raw != nil {
			users, err := raw.UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUser{UserID: p.UserID}})
			if err == nil && len(users) > 0 {
				if user, ok := users[0].(*tg.User); ok {
					return &tg.InputPeerUser{
						UserID:     user.ID,
						AccessHash: user.AccessHash,
					}, nil
				}
			}
		}
		return nil, fmt.Errorf("user not found")
	case *tg.PeerChat:
		return &tg.InputPeerChat{
			ChatID: p.ChatID,
		}, nil
	case *tg.PeerChannel:
		if channel, ok := entities.Channels[p.ChannelID]; ok {
			return &tg.InputPeerChannel{
				ChannelID:  channel.ID,
				AccessHash: channel.AccessHash,
			}, nil
		}
		if raw != nil {
			channels, err := raw.ChannelsGetChannels(ctx, []tg.InputChannelClass{&tg.InputChannel{ChannelID: p.ChannelID}})
			if err == nil && len(channels.GetChats()) > 0 {
				if channel, ok := channels.GetChats()[0].(*tg.Channel); ok {
					return &tg.InputPeerChannel{
						ChannelID:  channel.ID,
						AccessHash: channel.AccessHash,
					}, nil
				}
			}
		}
		return nil, fmt.Errorf("channel not found")
	}
	return nil, nil
}
