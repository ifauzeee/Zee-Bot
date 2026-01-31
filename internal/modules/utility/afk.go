package utility

import (
	"fmt"
	"strings"
	"time"

	"zee-ubot/internal/database"
	"zee-ubot/internal/handlers"
)

func Register(m *handlers.Manager) {
	m.Register("afk", "Set AFK status", afkHandler)
}

func afkHandler(c *handlers.Context) error {
	reason := "Sedang tidak di tempat."
	if len(c.Args) > 0 {
		reason = strings.Join(c.Args, " ")
	}

	startTime := time.Now().Unix()
	afkData := fmt.Sprintf("%d|%s", startTime, reason)

	if err := database.SetKV("afk_status", afkData); err != nil {
		return c.Edit("❌ Gagal mengaktifkan AFK.")
	}

	return c.Edit(fmt.Sprintf(
		"💤 𝗔𝗙𝗞 𝗗𝗜𝗔𝗞𝗧𝗜𝗙𝗞𝗔𝗡 💤\n"+
			"━━━━━━━━━━━━━━━━━━━━━━\n"+
			"📝 𝗔𝗹𝗮𝘀𝗮𝗻 : %s\n"+
			"━━━━━━━━━━━━━━━━━━━━━━",
		reason,
	))
}
