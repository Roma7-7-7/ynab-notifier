package budget

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Roma7-7-7/ynab-notifier/pkg/ynab"
)

const ten = 10
const tenFloat = float64(10)
const hundred = 100

type GeneralCategoryStatistic struct {
	Budgeted     int
	Activity     int
	Balance      int
	AvgSpent     int
	AvgSpentLeft int
	DaysLeft     int
}

func (s GeneralCategoryStatistic) BudgetedS() string {
	return FormatMoney(s.Budgeted)
}

func (s GeneralCategoryStatistic) ActivityS() string {
	return FormatMoney(s.Activity)
}

func (s GeneralCategoryStatistic) BalanceS() string {
	return FormatMoney(s.Balance)
}

func (s GeneralCategoryStatistic) AvgSpentS() string {
	return FormatMoney(s.AvgSpent)
}

func (s GeneralCategoryStatistic) AvgSpentLeftS() string {
	return FormatMoney(s.AvgSpentLeft)
}

func CalculateStatistic(c ynab.Category) GeneralCategoryStatistic {
	today := time.Now()
	return GeneralCategoryStatistic{
		Budgeted:     c.Budgeted,
		Activity:     c.Activity,
		Balance:      c.Balance,
		AvgSpent:     CalculateAvgSpent(c, today),
		AvgSpentLeft: CalculateAvgLeft(c, today),
		DaysLeft:     DaysLeft(today),
	}
}

func CalculateAvgSpent(c ynab.Category, date time.Time) int {
	if c.Activity == 0 {
		return 0
	}

	return c.Activity / date.Day()
}

func CalculateAvgLeft(c ynab.Category, date time.Time) int {
	return c.Balance / (daysInMonth(date) - date.Day() + 1)
}

// FormatMoney format money in cents to string like "123,456.78".
//
// Note: there are 3 cents at the end of value because of YNAB API.
func FormatMoney(money int) string {
	if money == 0 {
		return "0.00"
	}
	pref := ""
	if money < 0 {
		pref = "-"
		money = -money
	}
	if money < ten {
		return pref + "0.01"
	}
	money = int(math.Round(float64(money) / tenFloat))

	primaryStr := strconv.Itoa(money / hundred)
	cents := money % hundred

	// Add comma between thousands
	primary := ""
	for i := len(primaryStr) - 1; i >= 0; i-- {
		primary = string(primaryStr[i]) + primary
		if i != 0 && (len(primaryStr)-i)%3 == 0 {
			primary = "," + primary
		}
	}

	return fmt.Sprintf("%s%s.%02d", pref, primary, cents)
}

func DaysLeft(date time.Time) int {
	return daysInMonth(date) - date.Day()
}

func daysInMonth(date time.Time) int {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC).
		AddDate(0, 1, 0).AddDate(0, 0, -1).Day()
}
