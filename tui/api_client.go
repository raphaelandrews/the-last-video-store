package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thelastvideostore/internal/models"
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Session    *SessionState
}

func NewAPIClient(baseURL string, session *SessionState) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Session: session,
	}
}

func (c *APIClient) doRequest(method, path string, body interface{}, target interface{}, retryOn401 bool) error {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Session.IsLoggedIn {
		req.Header.Set("Authorization", "Bearer "+c.Session.AccessToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 && retryOn401 && c.Session.RefreshToken != "" {
		if err := c.refreshToken(); err != nil {
			return err
		}
		return c.doRequest(method, path, body, target, false)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
			Code  int    `json:"code"`
		}
		json.Unmarshal(respBody, &errResp)
		if errResp.Error != "" {
			return fmt.Errorf("%s (HTTP %d)", errResp.Error, resp.StatusCode)
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if target != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, target); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
	}

	return nil
}

func (c *APIClient) refreshToken() error {
	var resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	body := map[string]string{"refresh_token": c.Session.RefreshToken}
	if err := c.doRequest("POST", "/api/v1/auth/refresh", body, &resp, false); err != nil {
		c.Session.Logout()
		return err
	}

	c.Session.AccessToken = resp.AccessToken
	c.Session.RefreshToken = resp.RefreshToken
	return nil
}

func (c *APIClient) Login(username, password string) (*models.UserResponse, error) {
	body := map[string]string{"username": username, "password": password}
	var resp struct {
		AccessToken  string              `json:"access_token"`
		RefreshToken string              `json:"refresh_token"`
		User         models.UserResponse `json:"user"`
	}
	if err := c.doRequest("POST", "/api/v1/auth/login", body, &resp, true); err != nil {
		return nil, err
	}
	c.Session.Login(resp.AccessToken, resp.RefreshToken, &resp.User)
	return &resp.User, nil
}

func (c *APIClient) Register(username, password string) (*models.UserResponse, error) {
	body := map[string]string{"username": username, "password": password}
	var user models.UserResponse
	if err := c.doRequest("POST", "/api/v1/auth/register", body, &user, false); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *APIClient) SearchMovies(query string) ([]models.MovieResponse, error) {
	var movies []models.MovieResponse
	path := fmt.Sprintf("/api/v1/movies/search?q=%s", query)
	if err := c.doRequest("GET", path, nil, &movies, true); err != nil {
		return nil, err
	}
	return movies, nil
}

func (c *APIClient) GetMovies(genre string, page int) ([]models.MovieResponse, int, error) {
	path := fmt.Sprintf("/api/v1/movies?page=%d&page_size=20", page)
	if genre != "" {
		path += fmt.Sprintf("&genre=%s", genre)
	}
	var resp struct {
		Movies []models.MovieResponse `json:"movies"`
		Total  int                    `json:"total"`
	}
	if err := c.doRequest("GET", path, nil, &resp, true); err != nil {
		return nil, 0, err
	}
	return resp.Movies, resp.Total, nil
}

func (c *APIClient) GetMovie(id string) (*models.MovieResponse, error) {
	var movie models.MovieResponse
	if err := c.doRequest("GET", "/api/v1/movies/"+id, nil, &movie, true); err != nil {
		return nil, err
	}
	return &movie, nil
}

func (c *APIClient) RentMovie(movieID string) (*models.RentalResponse, error) {
	var rental models.RentalResponse
	body := map[string]string{"movie_id": movieID}
	if err := c.doRequest("POST", "/api/v1/rentals/rent", body, &rental, true); err != nil {
		return nil, err
	}
	return &rental, nil
}

func (c *APIClient) ReturnMovie(rentalID string) (*models.RentalResponse, error) {
	var rental models.RentalResponse
	body := map[string]string{"rental_id": rentalID}
	if err := c.doRequest("POST", "/api/v1/rentals/return", body, &rental, true); err != nil {
		return nil, err
	}
	return &rental, nil
}

func (c *APIClient) GetRentalHistory() ([]models.RentalResponse, error) {
	var rentals []models.RentalResponse
	if err := c.doRequest("GET", "/api/v1/rentals/history", nil, &rentals, true); err != nil {
		return nil, err
	}
	return rentals, nil
}

func (c *APIClient) GetUsers() ([]models.UserResponse, error) {
	var users []models.UserResponse
	if err := c.doRequest("GET", "/api/v1/users", nil, &users, true); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *APIClient) UpdateUser(id string, tier string, banned *bool) (*models.UserResponse, error) {
	body := map[string]interface{}{}
	if tier != "" {
		body["tier"] = tier
	}
	if banned != nil {
		body["banned"] = *banned
	}
	var user models.UserResponse
	if err := c.doRequest("PUT", "/api/v1/users/"+id, body, &user, true); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *APIClient) GetAuditEntries() ([]interface{}, error) {
	var entries []interface{}
	if err := c.doRequest("GET", "/api/v1/audit", nil, &entries, true); err != nil {
		return nil, err
	}
	return entries, nil
}

func (c *APIClient) GetWishlist() ([]map[string]interface{}, error) {
	var items []map[string]interface{}
	if err := c.doRequest("GET", "/api/v1/wishlist", nil, &items, true); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *APIClient) AddToWishlist(movieID string) error {
	body := map[string]string{"movie_id": movieID}
	return c.doRequest("POST", "/api/v1/wishlist", body, nil, true)
}

func (c *APIClient) RemoveFromWishlist(movieID string) error {
	return c.doRequest("DELETE", "/api/v1/wishlist/"+movieID, nil, nil, true)
}

func (c *APIClient) GetStaffPicks() ([]models.MovieResponse, error) {
	var movies []models.MovieResponse
	if err := c.doRequest("GET", "/api/v1/movies/staff-picks", nil, &movies, true); err != nil {
		return nil, err
	}
	return movies, nil
}

func (c *APIClient) GetLastChance() ([]models.MovieResponse, error) {
	var movies []models.MovieResponse
	if err := c.doRequest("GET", "/api/v1/movies/last-chance", nil, &movies, true); err != nil {
		return nil, err
	}
	return movies, nil
}

func (c *APIClient) TOTPSetup(userID string, enabled bool) (string, string, error) {
	body := map[string]bool{"enabled": enabled}
	var resp struct {
		Secret string `json:"secret"`
		URL    string `json:"url"`
	}
	if err := c.doRequest("POST", "/api/v1/users/"+userID+"/totp", body, &resp, true); err != nil {
		return "", "", err
	}
	return resp.Secret, resp.URL, nil
}

func (c *APIClient) LoginTOTP(tempToken, code string) (*models.UserResponse, error) {
	body := map[string]string{"code": code}
	c.Session.AccessToken = tempToken
	var resp struct {
		AccessToken  string              `json:"access_token"`
		RefreshToken string              `json:"refresh_token"`
		User         models.UserResponse `json:"user"`
	}
	if err := c.doRequest("POST", "/api/v1/auth/login/totp", body, &resp, true); err != nil {
		c.Session.AccessToken = ""
		return nil, err
	}
	c.Session.Login(resp.AccessToken, resp.RefreshToken, &resp.User)
	return &resp.User, nil
}
