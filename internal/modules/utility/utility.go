package utility

import "zee-ubot/internal/handlers"

func init() {
	handlers.RegisterModule(Register)
}

func Register(m *handlers.Manager) {
}
