//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type profileResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type taskResponse struct {
	Title string `json:"title"`
}

func baseURL(t *testing.T) string {
	t.Helper()
	if v := os.Getenv("BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

func mustDo(t *testing.T, req *http.Request) *http.Response {
	t.Helper()
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func readAllAndClose(t *testing.T, rc io.ReadCloser) []byte {
	t.Helper()
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return b
}

func postForm(t *testing.T, urlStr string, values url.Values) (*http.Response, []byte) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBufferString(values.Encode()))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := mustDo(t, req)
	body := readAllAndClose(t, resp.Body)
	return resp, body
}

func TestE2E_HealthAndAuthFlow(t *testing.T) {
	b := baseURL(t)

	// Health
	{
		req, err := http.NewRequest(http.MethodGet, b+"/health", nil)
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		resp := mustDo(t, req)
		_ = readAllAndClose(t, resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("health status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	}

	unique := fmt.Sprintf("%d", time.Now().UnixNano())
	name := "ci-user-" + unique
	email := "ci-" + unique + "@example.com"
	password := "P@ssw0rd-" + unique

	// Signup
	var signup tokenResponse
	{
		resp, body := postForm(t, b+"/v1/signup", url.Values{
			"name":     []string{name},
			"email":    []string{email},
			"password": []string{password},
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("signup status code = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, string(body))
		}
		if err := json.Unmarshal(body, &signup); err != nil {
			t.Fatalf("unmarshal signup response: %v; body=%s", err, string(body))
		}
		if signup.AccessToken == "" {
			t.Fatalf("signup accessToken is empty")
		}
	}

	// Login
	var login tokenResponse
	{
		resp, body := postForm(t, b+"/v1/login", url.Values{
			"email":    []string{email},
			"password": []string{password},
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("login status code = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, string(body))
		}
		if err := json.Unmarshal(body, &login); err != nil {
			t.Fatalf("unmarshal login response: %v; body=%s", err, string(body))
		}
		if login.AccessToken == "" {
			t.Fatalf("login accessToken is empty")
		}
	}

	authHeader := "Bearer " + login.AccessToken

	// Profile (protected)
	{
		req, err := http.NewRequest(http.MethodGet, b+"/v1/profile", nil)
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		req.Header.Set("Authorization", authHeader)
		resp := mustDo(t, req)
		body := readAllAndClose(t, resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("profile status code = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, string(body))
		}
		var profile profileResponse
		if err := json.Unmarshal(body, &profile); err != nil {
			t.Fatalf("unmarshal profile response: %v; body=%s", err, string(body))
		}
		if profile.Email != email {
			t.Fatalf("profile.email = %q, want %q", profile.Email, email)
		}
	}

	// Create task (protected)
	{
		values := url.Values{"title": []string{"task-" + unique}}
		req, err := http.NewRequest(http.MethodPost, b+"/v1/task", bytes.NewBufferString(values.Encode()))
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", authHeader)
		resp := mustDo(t, req)
		body := readAllAndClose(t, resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("create task status code = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, string(body))
		}
	}

	// Fetch tasks (protected)
	{
		req, err := http.NewRequest(http.MethodGet, b+"/v1/task", nil)
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		req.Header.Set("Authorization", authHeader)
		resp := mustDo(t, req)
		body := readAllAndClose(t, resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("fetch tasks status code = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, string(body))
		}
		var tasks []taskResponse
		if err := json.Unmarshal(body, &tasks); err != nil {
			t.Fatalf("unmarshal tasks response: %v; body=%s", err, string(body))
		}
		if len(tasks) < 1 {
			t.Fatalf("expected at least 1 task")
		}
	}
}
