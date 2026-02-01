package mocks

import (
	"go.uber.org/zap"
)

type MockSender struct {
	LastMessage string
	LastEdit    string
	LastReply   string
}

func (m *MockSender) Reset() {
	m.LastMessage = ""
	m.LastEdit = ""
	m.LastReply = ""
}

func NewTestContext(sender *MockSender, logger *zap.Logger) interface{} {
	return nil
}
