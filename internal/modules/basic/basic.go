package basic

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"zee-ubot/internal/handlers"
)

func init() {
	handlers.RegisterModule(Register)
}

func Register(m *handlers.Manager) {
	m.Register("ping", "Check bot latency", pingHandler)
	m.Register("help", "Show this help message", helpHandler(m))
}

func pingHandler(c *handlers.Context) error {
	start := time.Now()
	if err := c.EditStatus("⚡️", "Analyzing latency..."); err != nil {
		return fmt.Errorf("ping: failed to set status: %w", err)
	}
	duration := time.Since(start).Round(time.Millisecond)

	style := c.NewStyle("Pong Status", "🛰")
	style.AddRowWithIcon("⚡️", "Latency", duration)
	style.AddRowWithIcon("⚙️", "Status", "Safe")
	style.AddRowWithIcon("🤖", "Engine", "Zee-Ubot")

	if err := c.Edit(style.Build()); err != nil {
		return fmt.Errorf("ping: failed to send result: %w", err)
	}
	return nil
}

func helpHandler(m *handlers.Manager) handlers.HandlerFunc {
	return func(c *handlers.Context) error {
		keys := make([]string, 0, len(m.Commands))
		maxLen := 0
		for k := range m.Commands {
			keys = append(keys, k)
			if len(k) > maxLen {
				maxLen = len(k)
			}
		}
		sort.Strings(keys)

		var b strings.Builder
		for _, k := range keys {
			cmd := m.Commands[k]
			spacing := ""
			if len(cmd.Name) < maxLen {
				spacing = strings.Repeat(" ", maxLen-len(cmd.Name))
			}
			fmt.Fprintf(&b, "%s .%s %s %s %s\n", handlers.RowPrefix, cmd.Name, spacing, handlers.Arrow, cmd.Description)
		}

		style := c.NewStyle("Zee-Ubot Menu", "✨")

		style.AddRawRow(b.String())

		style.SetFooter("⋄ 𝗖𝗿𝗲𝗮𝘁𝗲𝗱 𝘄𝗶𝘁𝗵 ❤️ 𝗯𝘆 @zee")

		if err := c.Edit(style.Build()); err != nil {
			return fmt.Errorf("help: failed to send menu: %w", err)
		}
		return nil
	}
}
