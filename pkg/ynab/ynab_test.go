package ynab_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Roma7-7-7/ynab-notifier/pkg/ynab"
)

func TestClient_GetCategory(t *testing.T) {
	type fields struct {
		token       string
		handlerFunc http.HandlerFunc
	}
	type args struct {
		budgetID   string
		categoryID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ynab.Category
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			fields: fields{
				token: "token",
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/v1/budgets/1234/categories/5678" {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("wrong url"))
						return
					}
					if r.Header.Get("Authorization") != "Bearer token" {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("wrong token"))
						return
					}

					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"data": {"category": {"id": "7733", "name": "testName", "budgeted": 100, "activity": -60, "balance": 40}}}`))
				},
			},
			args: args{
				budgetID:   "1234",
				categoryID: "5678",
			},
			want: &ynab.Category{
				ID:       "7733",
				Name:     "testName",
				Budgeted: 100,
				Activity: -60,
				Balance:  40,
			},
			wantErr: assert.NoError,
		},
		{
			name: "forbidden",
			fields: fields{
				token: "token",
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusForbidden)
				},
			},
			args: args{
				budgetID:   "1234",
				categoryID: "5678",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIsf(t, err, ynab.ErrForbidden, "GetCategory(ctx, %s, %s)", "1234", "5678")
			},
		},
		{
			name: "not_found",
			fields: fields{
				token: "token",
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				},
			},
			args: args{
				budgetID:   "1234",
				categoryID: "5678",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIsf(t, err, ynab.ErrNotFound, "GetCategory(ctx, %s, %s)", "1234", "5678")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.fields.handlerFunc)
			defer server.Close()

			c := ynab.NewClient(server.URL, tt.fields.token, zap.NewNop().Sugar())
			got, err := c.GetCategory(context.Background(), tt.args.budgetID, tt.args.categoryID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetCategory(ctx, %s, %s)", tt.args.budgetID, tt.args.categoryID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetCategory(ctx, %v, %v)", tt.args.budgetID, tt.args.categoryID)
		})
	}
}

func TestClient_Manual(t *testing.T) {
	t.Skipf("for manual run only")

	log, _ := zap.NewDevelopment()
	c := ynab.NewClient("https://api.ynab.com", os.Getenv("YNAB_ACCESS_TOKEN"), log.Sugar())

	t.Run("GetCategory", func(t *testing.T) {
		res, err := c.GetCategory(context.Background(), os.Getenv("YNAB_BUDGET_ID"), os.Getenv("YNAB_CATEGORY_ID"))
		require.NoError(t, err)
		t.Logf("%+v", res)
	})
}
