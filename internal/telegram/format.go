package telegram

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Roma7-7-7/ynab-notifier/internal/budget"
)

func NewDefaultStatisticMessageFormatter() (StatisticMessageFormatter, error) {
	t, err := template.New("defaultStatisticMessageFormatter").
		Parse(`Залишок:                               {{.BalanceS}} грн.
Середньо витрачається:    {{.AvgSpentS}} грн. в день

Середньо залишилося:       {{.AvgSpentLeftS}} грн. в день
`)
	if err != nil {
		return nil, fmt.Errorf("parsing defaultStatisticMessageFormatter template: %w", err)
	}

	return func(cat budget.GeneralCategoryStatistic) (string, error) {
		var buff bytes.Buffer
		if err = t.Execute(&buff, cat); err != nil {
			return "", fmt.Errorf("executing defaultStatisticMessageFormatter template: %w", err)
		}
		return buff.String(), nil
	}, nil
}
