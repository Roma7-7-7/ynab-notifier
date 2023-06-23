package ynab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	ErrNotFound     = fmt.Errorf("not found")
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrForbidden    = fmt.Errorf("forbidden")
)

const getCategoryURL = "%s/v1/budgets/%s/categories/%s"

type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Budgeted int    `json:"budgeted"`
	Activity int    `json:"activity"`
}

type categoryResponse struct {
	Data struct {
		Category Category `json:"category"`
	} `json:"data"`
}

type Client struct {
	baseULR string
	token   string
	client  *http.Client
	log     Logger
}

func NewClient(baseURL, token string, log Logger) *Client {
	return &Client{
		baseULR: baseURL,
		token:   token,
		client:  &http.Client{},
		log:     log,
	}
}

func (c *Client) GetCategory(ctx context.Context, budgetID, categoryID string) (*Category, error) {
	c.log.Debugw("getting categoryID", "budgetID", budgetID, "categoryID", categoryID)

	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, fmt.Sprintf(getCategoryURL, c.baseULR, budgetID, categoryID), nil)
	if err != nil {
		c.log.Errorw("can't create request", "budgetID", budgetID, "categoryID", categoryID, "error", err)
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Errorw("can't do request", "budgetID", budgetID, "categoryID", categoryID, "error", err)
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp, budgetID, categoryID)
	}

	var res categoryResponse

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		c.log.Errorw("can't decode response", "budgetID", budgetID, "categoryID", categoryID, "error", err)
		return nil, fmt.Errorf("can't decode response: %w", err)
	}

	c.log.Debugw("got category", "budgetID", budgetID, "categoryID", categoryID)
	return &res.Data.Category, nil
}

func (c *Client) handleErrorResponse(resp *http.Response, budgetID string, categoryID string) (*Category, error) {
	if resp.StatusCode == http.StatusNotFound {
		c.log.Debugw("categoryID not found", "budgetID", budgetID, "categoryID", categoryID)
		return nil, ErrNotFound
	}
	if resp.StatusCode == http.StatusUnauthorized {
		c.log.Debugw("unauthorized", "budgetID", budgetID, "categoryID", categoryID)
		return nil, ErrUnauthorized
	}
	if resp.StatusCode == http.StatusForbidden {
		c.log.Debugw("forbidden", "budgetID", budgetID, "categoryID", categoryID)
		return nil, ErrForbidden
	}

	c.log.Warnw("unexpected status code",
		"budgetID", budgetID, "categoryID", categoryID, "statusCode", resp.StatusCode, "payload", resp.Body,
	)

	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
