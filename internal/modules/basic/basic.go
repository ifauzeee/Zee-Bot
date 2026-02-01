package basic

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"zee-ubot/internal/handlers"
)

func Register(m *handlers.Manager) {
	m.Register("ping", "Check bot latency", pingHandler)
	m.Register("help", "Show this help message", helpHandler(m))
}

func pingHandler(c *handlers.Context) error {
	start := time.Now()
	if err := c.Edit("⚡️ `Analyzing latency...` "); err != nil {
		return err
	}
	duration := time.Since(start).Round(time.Millisecond)

	style := c.NewStyle("Pong Status", "🛰")
	style.AddRowWithIcon("⚡️", "Latency", duration)
	style.AddRowWithIcon("⚙️", "Status", "Safe")
	style.AddRowWithIcon("🤖", "Engine", "Zee-Ubot")

	return c.Edit(style.Build())
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
			fmt.Fprintf(&b, "▫️ .%s %s ┆ %s\n", cmd.Name, spacing, cmd.Description)
		}

		style := c.NewStyle("Zee-Ubot Menu", "✨")
		style.SetFooter("⋄ 𝗖𝗿𝗲𝗮𝘁𝗲𝗱 𝘄𝗶𝘁𝗵 ❤️ 𝗯𝘆 @zee")

		output := style.Build()
		output = strings.Replace(output, handlers.Divider, handlers.Divider+"\n  [ 🗂 𝗠𝗔𝗜𝗡 𝗖𝗢𝗠𝗠𝗔𝗡𝗗𝗦 ]\n\n"+b.String()+"\n"+handlers.Divider, 1)
		output = strings.Replace(output, "🏷 **Title/Name** : `Unknown`"+"\n", "", 1)

		return c.Edit(output)
	}
}
