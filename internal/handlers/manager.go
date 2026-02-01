package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"
	"zee-ubot/internal/database"

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
		return nil
	}

	if msg.Out {
		isAfkCmd := strings.HasPrefix(msg.Message, ".afk") || strings.HasPrefix(msg.Message, "!afk")
		isAfkStatus := strings.Contains(msg.Message, "💤") || strings.Contains(msg.Message, "✅")

		if !isAfkCmd && !isAfkStatus {
			m.checkAndDisableAFK(ctx, sender, raw, msg, entities)
		}

		return m.processMessage(ctx, sender, raw, msg, entities)
	}

	return m.handleIncomingAFK(ctx, sender, raw, msg, entities)
}

func (m *Manager) checkAndDisableAFK(ctx context.Context, sender *message.Sender, raw *tg.Client, msg *tg.Message, entities tg.Entities) {
	status, _ := database.GetKV("afk_status")
	if status != "" {
		_ = database.DeleteKV("afk_status")
		parts := strings.SplitN(status, "|", 2)
		if len(parts) == 2 {
			startTimeStr := parts[0]
			var startTime int64
			_, _ = fmt.Sscanf(startTimeStr, "%d", &startTime)
			duration := time.Since(time.Unix(startTime, 0)).Round(time.Second)

			peer, _ := m.resolvePeer(ctx, raw, msg.PeerID, entities)
			if peer != nil {
				style := (&Context{Ctx: ctx, Sender: sender, Peer: peer, Raw: raw}).NewStyle("𝗦𝗔𝗬𝗔 𝗞𝗘𝗠𝗕𝗔𝗟𝗜!", "✨")
				style.AddRowWithIcon("⏰", "Lama AFK", duration)
				_, _ = sender.To(peer).Text(ctx, style.Build())
			}
		}
	}
}

func (m *Manager) handleIncomingAFK(ctx context.Context, sender *message.Sender, raw *tg.Client, msg *tg.Message, entities tg.Entities) error {
	status, _ := database.GetKV("afk_status")
	if status == "" {
		return nil
	}

	parts := strings.SplitN(status, "|", 2)
	if len(parts) < 2 {
		return nil
	}
	startTimeStr, reason := parts[0], parts[1]

	var startTime int64
	_, _ = fmt.Sscanf(startTimeStr, "%d", &startTime)
	duration := time.Since(time.Unix(startTime, 0)).Round(time.Second)

	isPM := false
	if _, ok := msg.PeerID.(*tg.PeerUser); ok {
		isPM = true
	}

	isMentioned := false
	if isPM {
		isMentioned = true
	} else {
		isMentioned = msg.Mentioned
	}

	if isMentioned {
		m.Logger.Info("AFK: Mentioned or PM received, sending reply",
			zap.Bool("isPM", isPM),
			zap.String("reason", reason))

		peer, err := m.resolvePeer(ctx, raw, msg.PeerID, entities)
		if err != nil {
			m.Logger.Warn("AFK: Could not resolve peer for reply", zap.Error(err))
			return nil
		}

		style := (&Context{Ctx: ctx, Sender: sender, Peer: peer, Raw: raw}).NewStyle("𝗦𝗘𝗗𝗔𝗡𝗚 𝗔𝗙𝗞", "💤")
		style.AddRowWithIcon("📝", "Alasan", reason)
		style.AddRowWithIcon("⏳", "Durasi", duration)
		_, err = sender.To(peer).Reply(msg.ID).Text(ctx, style.Build())
		if err != nil {
			m.Logger.Error("AFK: Failed to send reply", zap.Error(err))
		}
		return err
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

		if err := cmd.Handler(c); err != nil {
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
