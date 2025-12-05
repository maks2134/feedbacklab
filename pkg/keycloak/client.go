package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

//TODO: переделать client на общую структуру

// Client represents a Keycloak Admin API client
type Client struct {
	baseURL      string
	realm        string
	adminUser    string
	adminPwd     string
	httpClient   *http.Client
	accessToken  string
	tokenExpires time.Time
}

// NewClient creates a new Keycloak Admin API client
func NewClient(baseURL, realm, adminUser, adminPwd string) *Client {
	return &Client{
		baseURL:    baseURL,
		realm:      realm,
		adminUser:  adminUser,
		adminPwd:   adminPwd,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// TokenResponse represents Keycloak token response
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
}

// getAccessToken retrieves or refreshes the admin access token
func (c *Client) getAccessToken() (string, error) {
	if c.accessToken != "" && time.Now().Before(c.tokenExpires.Add(-30*time.Second)) {
		return c.accessToken, nil
	}

	tokenURL := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", c.baseURL)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", "admin-cli")
	data.Set("username", c.adminUser)
	data.Set("password", c.adminPwd)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpires = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return c.accessToken, nil
}

// CreateUserRequest represents a user creation request
type CreateUserRequest struct {
	Username      string                     `json:"username"`
	Email         string                     `json:"email"`
	FirstName     string                     `json:"firstName,omitempty"`
	LastName      string                     `json:"lastName,omitempty"`
	Enabled       bool                       `json:"enabled"`
	EmailVerified bool                       `json:"emailVerified"`
	Credentials   []CredentialRepresentation `json:"credentials,omitempty"`
}

// CredentialRepresentation represents user credentials
type CredentialRepresentation struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type UserRepresentation struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Enabled          bool   `json:"enabled"`
	EmailVerified    bool   `json:"emailVerified"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
}

// CreateUser creates a new user in Keycloak
func (c *Client) CreateUser(req CreateUserRequest) (string, error) {
	token, err := c.getAccessToken()
	if err != nil {
		return "", err
	}

	userURL := fmt.Sprintf("%s/admin/realms/%s/users", c.baseURL, c.realm)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", userURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create user request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create user: status %d, body: %s", resp.StatusCode, string(body))
	}

	location := resp.Header.Get("Location")
	if location != "" {
		parts := location
		lastSlash := len(parts) - 1
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] == '/' {
				lastSlash = i
				break
			}
		}
		if lastSlash < len(parts)-1 {
			return parts[lastSlash+1:], nil
		}
	}

	return "", fmt.Errorf("user created but could not extract user ID")
}

// GetUserByEmail retrieves a user by email
func (c *Client) GetUserByEmail(email string) (*UserRepresentation, error) {
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	userURL := fmt.Sprintf("%s/admin/realms/%s/users?email=%s&exact=true", c.baseURL, c.realm, url.QueryEscape(email))

	httpReq, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user: status %d, body: %s", resp.StatusCode, string(body))
	}

	var users []UserRepresentation
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

// ResetPassword resets user password
func (c *Client) ResetPassword(userID string, newPassword string, temporary bool) error {
	token, err := c.getAccessToken()
	if err != nil {
		return err
	}

	passwordURL := fmt.Sprintf("%s/admin/realms/%s/users/%s/reset-password", c.baseURL, c.realm, userID)

	credential := CredentialRepresentation{
		Type:      "password",
		Value:     newPassword,
		Temporary: temporary,
	}

	reqBody, err := json.Marshal(credential)
	if err != nil {
		return fmt.Errorf("failed to marshal password request: %w", err)
	}

	httpReq, err := http.NewRequest("PUT", passwordURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create password request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to reset password: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// UpdatePassword updates user password (for user self-service)
func (c *Client) UpdatePassword(userID, currentPassword, newPassword string) error {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.baseURL, c.realm)

	user, err := c.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", "feedbacklab-api")
	data.Set("username", user.Username)
	data.Set("password", currentPassword)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("current password is incorrect")
	}

	return c.ResetPassword(userID, newPassword, false)
}

// GetUserByID retrieves a user by ID
func (c *Client) GetUserByID(userID string) (*UserRepresentation, error) {
	token, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	userURL := fmt.Sprintf("%s/admin/realms/%s/users/%s", c.baseURL, c.realm, userID)

	httpReq, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user: status %d, body: %s", resp.StatusCode, string(body))
	}

	var user UserRepresentation
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &user, nil
}
