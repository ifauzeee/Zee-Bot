package afk

import (
	"fmt"
	"strings"
	"time"

	"zee-ubot/internal/database"
	"zee-ubot/internal/handlers"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

func init() {
	handlers.RegisterModule(Register)
}

func Register(m *handlers.Manager) {
	m.Register("afk", "Set status away (ex: .afk tidur)", afkHandler)
	m.RegisterHook("all", afkHook)
}

func afkHandler(c *handlers.Context) error {
	reason := "Away from keyboard"
	if len(c.Args) > 0 {
		reason = strings.Join(c.Args, " ")
	}

	startTime := time.Now().Unix()
	val := fmt.Sprintf("%d|%s", startTime, reason)

	if err := database.SetKV("afk_status", val); err != nil {
		return fmt.Errorf("failed to save afk status: %w", err)
	}

	style := c.NewStyle("𝗦𝗘𝗗𝗔𝗡𝗚 𝗔𝗙𝗞", "💤")
	style.AddRowWithIcon("📝", "Alasan", reason)
	style.AddRowWithIcon("⏰", "Waktu", time.Unix(startTime, 0).Format("15:04"))

	return c.Edit(style.Build())
}

func afkHook(c *handlers.Context, update tg.UpdateClass) error {
	if c.Msg == nil {
		return nil
	}

	if c.Msg.Out {
		return checkAndDisableAFK(c)
	}

	return handleIncomingAFK(c)
}

func checkAndDisableAFK(c *handlers.Context) error {
	text := c.Msg.Message
	if strings.HasPrefix(text, ".afk") || strings.HasPrefix(text, "!afk") {
		return nil
	}
	if strings.Contains(text, "💤") || strings.Contains(text, "✅") {
		return nil
	}

	status, _ := database.GetKV("afk_status")
	if status == "" {
		return nil
	}

	if err := database.DeleteKV("afk_status"); err != nil {
		c.Logger.Error("Failed to delete AFK status", zap.Error(err))
	}

	parts := strings.SplitN(status, "|", 2)
	if len(parts) == 2 {
		startTimeStr := parts[0]
		var startTime int64
		_, _ = fmt.Sscanf(startTimeStr, "%d", &startTime)
		duration := time.Since(time.Unix(startTime, 0)).Round(time.Second)

		style := c.NewStyle("𝗦𝗔𝗬𝗔 𝗞𝗘𝗠𝗕𝗔𝗟𝗜!", "✅")
		style.AddRowWithIcon("⏰", "Lama AFK", duration)

		_, err := c.Sender.To(c.Peer).Text(c.Ctx, style.Build())
		if err != nil {
			return fmt.Errorf("failed to send auto-resume message: %w", err)
		}
		return nil
	}
	return nil
}

func handleIncomingAFK(c *handlers.Context) error {
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
	if _, ok := c.Msg.PeerID.(*tg.PeerUser); ok {
		isPM = true
	}

	isMentioned := isPM || c.Msg.Mentioned

	if isMentioned {
		c.Logger.Info("AFK: Mentioned or PM received, sending reply",
			zap.Bool("isPM", isPM),
			zap.String("reason", reason))

		style := c.NewStyle("𝗦𝗘𝗗𝗔𝗡𝗚 𝗔𝗙𝗞", "💤")
		style.AddRowWithIcon("📝", "Alasan", reason)
		style.AddRowWithIcon("⏳", "Durasi", duration)

		err := c.Reply(style.Build())
		if err != nil {
			c.Logger.Error("AFK: Failed to send reply", zap.Error(err))
			return fmt.Errorf("failed to send auto-reply: %w", err)
		}
		return nil
	}

	return nil
}
