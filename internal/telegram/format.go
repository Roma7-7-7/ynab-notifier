package telegram

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Roma7-7-7/ynab-notifier/internal/budget"
)

type extendedStatistic struct {
	budget.GeneralCategoryStatistic
}

func (s extendedStatistic) DaysLeftS() string {
	switch s.DaysLeft {
	case 1, 21, 31: //nolint: gomnd // 1, 21, 31 день
		return fmt.Sprintf("%d день", s.DaysLeft)
	case 2, 3, 4, 22, 23, 24: //nolint: gomnd // 2, 3, 4, 22, 23, 24 дні
		return fmt.Sprintf("%d дні", s.DaysLeft)
	default:
		return fmt.Sprintf("%d днів", s.DaysLeft)
	}
}

func NewDefaultStatisticMessageFormatter() (StatisticMessageFormatter, error) {
	t, err := template.New("defaultStatisticMessageFormatter").
		Parse(`Статистика: 🔴 {{.AvgSpentS}} грн. в день

Залишок:      🟢 {{.BalanceS}} грн. / {{.DaysLeftS}}
В день:          🟡 {{.AvgSpentLeftS}} грн. в день
`)
	if err != nil {
		return nil, fmt.Errorf("parsing defaultStatisticMessageFormatter template: %w", err)
	}

	return func(cat budget.GeneralCategoryStatistic) (string, error) {
		var buff bytes.Buffer
		if err = t.Execute(&buff, extendedStatistic{cat}); err != nil {
			return "", fmt.Errorf("executing defaultStatisticMessageFormatter template: %w", err)
		}
		return buff.String(), nil
	}, nil
}
