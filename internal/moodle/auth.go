package moodle

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Login exchanges a Moodle username and password for a web service token via
// the mobile app token service (POST /login/token.php).
// On success it returns the token string. On failure it returns a descriptive error.
func Login(ctx context.Context, baseURL, username, password string) (string, error) {
	endpoint := strings.TrimRight(baseURL, "/") + "/login/token.php"
	params := url.Values{
		"username": {username},
		"password": {password},
		"service":  {"moodle_mobile_app"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return "", fmt.Errorf("create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login http error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read login response: %w", err)
	}

	var result struct {
		Token string `json:"token"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("decode login response: %w", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("%s", result.Error)
	}
	if result.Token == "" {
		return "", fmt.Errorf("empty token in response")
	}
	return result.Token, nil
}
