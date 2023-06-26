package telegram

import (
	"time"

	"gopkg.in/telebot.v3"
	tb "gopkg.in/telebot.v3"
)

func NewTelebot(token string) (*telebot.Bot, error) {
	return tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 60 * time.Second},
	})
}
