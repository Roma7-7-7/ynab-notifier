package budget_test

import (
	"testing"
	"time"

	"github.com/Roma7-7-7/ynab-notifier/internal/budget"
	"github.com/Roma7-7-7/ynab-notifier/pkg/ynab"
)

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		name         string
		arg          int
		wantPositive string
		wantNegative string
	}{
		{
			name:         "0.00",
			arg:          0,
			wantPositive: "0.00",
			wantNegative: "0.00",
		},
		{
			name:         "0.01_1",
			arg:          1,
			wantPositive: "0.01",
			wantNegative: "-0.01",
		},
		{
			name:         "0.01_2",
			arg:          12,
			wantPositive: "0.01",
			wantNegative: "-0.01",
		},
		{
			name:         "0.12",
			arg:          123,
			wantPositive: "0.12",
			wantNegative: "-0.12",
		},
		{
			name:         "1.23",
			arg:          1234,
			wantPositive: "1.23",
			wantNegative: "-1.23",
		},
		{
			name:         "12.35",
			arg:          12345,
			wantPositive: "12.35",
			wantNegative: "-12.35",
		},
		{
			name:         "123.46",
			arg:          123456,
			wantPositive: "123.46",
			wantNegative: "-123.46",
		},
		{
			name:         "1,234.57",
			arg:          1234567,
			wantPositive: "1,234.57",
			wantNegative: "-1,234.57",
		},
		{
			name:         "12,345.68",
			arg:          12345678,
			wantPositive: "12,345.68",
			wantNegative: "-12,345.68",
		},
		{
			name:         "123,456.79",
			arg:          123456789,
			wantPositive: "123,456.79",
			wantNegative: "-123,456.79",
		},
		{
			name:         "1,234,567.89",
			arg:          1234567890,
			wantPositive: "1,234,567.89",
			wantNegative: "-1,234,567.89",
		},
		{
			name:         "12,345,678.90",
			arg:          12345678901,
			wantPositive: "12,345,678.90",
			wantNegative: "-12,345,678.90",
		},
		{
			name:         "123,456,789.01",
			arg:          123456789012,
			wantPositive: "123,456,789.01",
			wantNegative: "-123,456,789.01",
		},
		{
			name:         "1,234,567,890.12",
			arg:          1234567890123,
			wantPositive: "1,234,567,890.12",
			wantNegative: "-1,234,567,890.12",
		},
		{
			name:         "999,999.99",
			arg:          999999994,
			wantPositive: "999,999.99",
			wantNegative: "-999,999.99",
		},
		{
			name:         "1,000,000.00_1",
			arg:          999999995,
			wantPositive: "1,000,000.00",
			wantNegative: "-1,000,000.00",
		},
		{
			name:         "1,000,000.00_2",
			arg:          999999999,
			wantPositive: "1,000,000.00",
			wantNegative: "-1,000,000.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := budget.FormatMoney(tt.arg); got != tt.wantPositive {
				t.Errorf("FormatMoney() = %v, want positive %v", got, tt.wantPositive)
			}
			if got := budget.FormatMoney(-tt.arg); got != tt.wantNegative {
				t.Errorf("FormatMoney() = %v, want negative %v", got, tt.wantNegative)
			}
		})
	}
}

func TestCalculateAvgSpentAndCalculateAvgLeft(t *testing.T) {
	type args struct {
		c    ynab.Category
		time time.Time
	}
	tests := []struct {
		name      string
		args      args
		wantSpent int
		wantLeft  int
	}{
		{
			name: "0_activity_at_first_day",
			args: args{
				c: ynab.Category{
					Activity: 0,
					Budgeted: 1000000,
					Balance:  1000000,
				},
				time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantSpent: 0,
			wantLeft:  32258,
		},
		{
			name: "0_activity_at_middle_day",
			args: args{
				c: ynab.Category{
					Activity: 0,
					Budgeted: 1000000,
					Balance:  1000000,
				},
				time: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			wantSpent: 0,
			wantLeft:  58823,
		},
		{
			name: "0_activity_at_last_day",
			args: args{
				c: ynab.Category{
					Activity: 0,
					Budgeted: 1000000,
					Balance:  1000000,
				},
				time: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			},
			wantSpent: 0,
			wantLeft:  1000000,
		},
		{
			name: "200_activity_at_first_day",
			args: args{
				c: ynab.Category{
					Activity: 200000,
					Budgeted: 1000000,
					Balance:  800000,
				},
				time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantSpent: 200000,
			wantLeft:  25806,
		},
		{
			name: "200_activity_at_middle_day",
			args: args{
				c: ynab.Category{
					Activity: 200000,
					Budgeted: 1000000,
					Balance:  800000,
				},
				time: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			wantSpent: 13333,
			wantLeft:  47058,
		},
		{
			name: "200_activity_at_last_day",
			args: args{
				c: ynab.Category{
					Activity: 200000,
					Budgeted: 1000000,
					Balance:  800000,
				},
				time: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			},
			wantSpent: 6451,
			wantLeft:  800000,
		},
		{
			name: "500_activity_at_first_day_other_month",
			args: args{
				c: ynab.Category{
					Activity: 500000,
					Budgeted: 1000000,
					Balance:  800000,
				},
				time: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			wantSpent: 500000,
			wantLeft:  28571,
		},
		{
			name: "500_activity_at_middle_day_other_month",
			args: args{
				c: ynab.Category{
					Activity: 500000,
					Budgeted: 1000000,
					Balance:  800000,
				},
				time: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
			},
			wantSpent: 33333,
			wantLeft:  57142,
		},
		{
			name: "500_activity_at_last_day_other_month",
			args: args{
				c: ynab.Category{
					Activity: 500000,
					Budgeted: 1000000,
					Balance:  800000,
				},
				time: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			},
			wantSpent: 17857,
			wantLeft:  800000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := budget.CalculateAvgSpent(tt.args.c, tt.args.time); got != tt.wantSpent {
				t.Errorf("CalculateAvgSpent() = %v, wantSpent %v", got, tt.wantSpent)
			}
			if got := budget.CalculateAvgLeft(tt.args.c, tt.args.time); got != tt.wantLeft {
				t.Errorf("CalculateAvgLeft() = %v, want %v", got, tt.wantLeft)
			}
		})
	}
}
