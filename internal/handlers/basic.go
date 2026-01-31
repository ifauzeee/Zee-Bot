package handlers

import (
	"fmt"
	"time"
)

func PingHandler(c *Context) error {
	start := time.Now()

	if err := c.Edit("Pong! ..."); err != nil {
		return err
	}

	duration := time.Since(start)
	finalText := fmt.Sprintf("Pong! 🏓\nLatency: %s", duration.Round(time.Millisecond))

	return c.Edit(finalText)
}

func HelpHandler(commands map[string]Command) HandlerFunc {
	return func(c *Context) error {
		var helpText string = "**Available Commands:**\n\n"
		for _, cmd := range commands {
			helpText += fmt.Sprintf("• `.%s`: %s\n", cmd.Name, cmd.Description)
		}
		return c.Edit(helpText)
	}
}
