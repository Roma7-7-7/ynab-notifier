package ynab

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
		want    *Category
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
					w.Write([]byte(`{"data": {"category": {"id": "7733", "name": "testName", "budgeted": 100, "activity": -50}}}`))
				},
			},
			args: args{
				budgetID:   "1234",
				categoryID: "5678",
			},
			want: &Category{
				ID:       "7733",
				Name:     "testName",
				Budgeted: 100,
				Activity: -50,
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
				return assert.ErrorIsf(t, err, ErrForbidden, "GetCategory(ctx, %s, %s)", "1234", "5678")
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
				return assert.ErrorIsf(t, err, ErrNotFound, "GetCategory(ctx, %s, %s)", "1234", "5678")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.fields.handlerFunc)
			defer server.Close()

			c := NewClient(server.URL, tt.fields.token, zap.NewNop().Sugar())
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
	c := NewClient("https://api.ynab.com", os.Getenv("YNAB_ACCESS_TOKEN"), log.Sugar())

	t.Run("GetCategory", func(t *testing.T) {
		res, err := c.GetCategory(context.Background(), os.Getenv("YNAB_BUDGET_ID"), os.Getenv("YNAB_CATEGORY_ID"))
		require.NoError(t, err)
		t.Logf("%+v", res)
	})
}
