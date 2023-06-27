package telegram

import (
	"context"
	"time"

	tb "gopkg.in/telebot.v3"

	"github.com/Roma7-7-7/ynab-notifier/internal/budget"
	"github.com/Roma7-7-7/ynab-notifier/pkg/ynab"
)

type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

type YNABClient interface {
	GetCategory(ctx context.Context, budgetID, categoryID string) (*ynab.Category, error)
}

type StatisticMessageFormatter func(cat budget.GeneralCategoryStatistic) (string, error)

type Bot struct {
	chatIDs map[int64]struct{}

	ynabBudgetID   string
	ynabCategoryID string
	ynabClient     YNABClient
	msgFormatter   StatisticMessageFormatter

	log Logger
}

type YNABDependencies struct {
	BudgetID   string
	CategoryID string
	Client     YNABClient
}

type Dependencies struct {
	AllowedChats              []int64
	YNAB                      YNABDependencies
	StatisticMessageFormatter StatisticMessageFormatter
	Logger                    Logger
}

func NewBot(deps Dependencies) *Bot {
	chatIDs := make(map[int64]struct{}, len(deps.AllowedChats))
	for _, chatID := range deps.AllowedChats {
		chatIDs[chatID] = struct{}{}
	}

	return &Bot{
		chatIDs: chatIDs,

		ynabBudgetID:   deps.YNAB.BudgetID,
		ynabCategoryID: deps.YNAB.CategoryID,
		ynabClient:     deps.YNAB.Client,

		msgFormatter: deps.StatisticMessageFormatter,

		log: deps.Logger,
	}
}

func (b *Bot) Start(bot *tb.Bot) {
	bot.Use(AllowedChatsMiddleware(b.chatIDs, b.log))

	bot.Handle("/start", b.stateHandler)
	bot.Handle("/state", b.stateHandler)

	bot.Start()
}

func (b *Bot) stateHandler(c tb.Context) error {
	b.log.Infow("status handler", "chatID", c.Chat().ID)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFunc()
	cat, err := b.ynabClient.GetCategory(ctx, b.ynabBudgetID, b.ynabCategoryID)
	if err != nil {
		b.log.Errorw("failed to get category", "error", err)
		return b.sendWithErrorLogging(c, "Unexpected error occurred. You know whom to call")
	}

	if cat == nil {
		b.log.Warnw("category is nil", "chatID", c.Chat().ID, "budgetID", b.ynabBudgetID, "categoryID", b.ynabCategoryID)
		return b.sendWithErrorLogging(c, "Unexpected error occurred. You know whom to call")
	}

	msg, err := b.msgFormatter(budget.CalculateStatistic(*cat))
	if err != nil {
		b.log.Errorw("failed to format message", "chatID", c.Chat().ID, "error", err)
		return b.sendWithErrorLogging(c, "Unexpected error occurred. You know whom to call")
	}

	return b.sendWithErrorLogging(c, msg)
}

func (b *Bot) sendWithErrorLogging(c tb.Context, msg string) error {
	if err := c.Send(msg); err != nil {
		b.log.Errorw("failed to send message", "chatID", c.Chat().ID, "error", err)
		return err
	}

	return nil
}

func AllowedChatsMiddleware(chatIDs map[int64]struct{}, log Logger) func(next tb.HandlerFunc) tb.HandlerFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			if _, ok := chatIDs[c.Chat().ID]; !ok {
				log.Warnw("chat is not allowed", "chatID", c.Chat().ID)
				if err := c.Send("You are not allowed to use this bot"); err != nil {
					log.Errorw("failed to send message", "chatID", c.Chat().ID, "error", err)
				}

				return nil
			}

			return next(c)
		}
	}
}
