package telegram

import (
	"time"

	tb "gopkg.in/telebot.v3"
)

func NewTelebot(token string) (*tb.Bot, error) {
	return tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 60 * time.Second}, //nolint: gomnd // 60 seconds
	})
}
