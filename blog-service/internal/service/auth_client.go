package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthClient struct {
	baseURL string
	http    *http.Client
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 4 * time.Second,
		},
	}
}

func (c *AuthClient) GetUsernameByID(ctx context.Context, id int64) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/users/%d", c.baseURL, id), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("auth returned %d", resp.StatusCode)
	}

	var out struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.Username, nil
}
