package main

import (
	"flag"
	"fmt"
	stdLog "log"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/Roma7-7-7/ynab-notifier/internal/telegram"
	"github.com/Roma7-7-7/ynab-notifier/pkg/ynab"
)

func main() {
	var debug = flag.Bool("debug", false, "debug mode")

	flag.Parse()

	var err error
	var l *zap.Logger
	if *debug {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		stdLog.Fatalf("can't initialize  zap logger: %v", err)
	}
	log := l.Sugar()

	chatIDs, err := chatIDs()
	if err != nil {
		log.Fatalw("failed to parse chat ids", "error", err)
	}

	telebot, err := telegram.NewTelebot(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatalw("failed to create telebot", "error", err)
	}

	client := ynab.NewClient("https://api.ynab.com", os.Getenv("YNAB_ACCESS_TOKEN"), log)
	formatter, err := telegram.NewDefaultStatisticMessageFormatter()
	if err != nil {
		log.Fatalw("failed to create statistic message formatter", "error", err)
	}

	bot := telegram.NewBot(telegram.Dependencies{
		AllowedChats: chatIDs,
		YNAB: telegram.YNABDependencies{
			BudgetID:   os.Getenv("YNAB_BUDGET_ID"),
			CategoryID: os.Getenv("YNAB_CATEGORY_ID"),
			Client:     client,
		},
		StatisticMessageFormatter: formatter,
		Logger:                    log,
	})

	bot.Start(telebot)
}

func chatIDs() ([]int64, error) {
	res := make([]int64, 0)

	chatIDs := os.Getenv("TELEGRAM_CHAT_IDS")
	if chatIDs == "" {
		return nil, fmt.Errorf("TELEGRAM_CHAT_IDS environment variable is empty")
	}

	for _, chatID := range strings.Split(chatIDs, ",") {
		val, err := strconv.ParseInt(strings.TrimSpace(chatID), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse chat id %q: %w", chatID, err)
		}
		res = append(res, val)
	}

	return res, nil
}
