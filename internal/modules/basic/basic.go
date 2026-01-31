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

	finalText := fmt.Sprintf(
		"🛰 𝗣𝗢𝗡𝗚 𝗦𝗧𝗔𝗧𝗨𝗦\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"⚡️ 𝗟𝗮𝘁𝗲𝗻𝗰𝘆  : %s\n"+
			"⚙️ 𝗦𝘁𝗮𝘁𝘂𝘀   : Safe\n"+
			"🤖 𝗘𝗻𝗴𝗶𝗻𝗲   : Zee-Ubot\n"+
			"━━━━━━━━━━━━━━━━━━━━━━",
		duration,
	)

	return c.Edit(finalText)
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
		b.WriteString("✨ 𝗭𝗘𝗘-𝗨𝗕𝗢𝗧 𝗠𝗘𝗡𝗨 ✨\n")
		b.WriteString("━━━━━━━━━━━━━━━━━━━━━━\n")
		b.WriteString("  [ 🗂 𝗠𝗔𝗜𝗡 𝗖𝗢𝗠𝗠𝗔𝗡𝗗𝗦 ]\n\n")

		for _, k := range keys {
			cmd := m.Commands[k]
			spacing := ""
			if len(cmd.Name) < maxLen {
				spacing = strings.Repeat(" ", maxLen-len(cmd.Name))
			}
			fmt.Fprintf(&b, "  ▫️ .%s %s ┆ %s\n", cmd.Name, spacing, cmd.Description)
		}

		b.WriteString("\n━━━━━━━━━━━━━━━━━━━━━━\n")
		b.WriteString("  ⋄ 𝗖𝗿𝗲𝗮𝘁𝗲𝗱 𝘄𝗶𝘁𝗵 ❤️ 𝗯𝘆 @zee")

		return c.Edit(b.String())
	}
}
