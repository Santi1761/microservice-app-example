package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Role      string `json:"role"`
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type UserService struct {
	Client            HTTPDoer
	UserAPIAddress    string
	AllowedUserHashes map[string]interface{}
}

func (s *UserService) Login(ctx context.Context, username, password string) (User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return User{}, ErrWrongCredentials
	}

	// Credenciales demo (user_pass)
	if _, ok := s.AllowedUserHashes[fmt.Sprintf("%s_%s", username, password)]; !ok {
		return User{}, ErrWrongCredentials
	}

	return s.getUser(ctx, username)
}

func (s *UserService) getUser(ctx context.Context, username string) (User, error) {
	var user User

	// Base desde env; si viene vacía, default al servicio docker
	base := strings.TrimSpace(s.UserAPIAddress)
	if base == "" {
		base = "http://users-api:8083"
	}
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}

	u, err := url.Parse(base)
	if err != nil {
		return user, fmt.Errorf("invalid users api base: %w", err)
	}
	u.Path = path.Join(u.Path, "/users", url.PathEscape(username))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return user, err
	}

	// Adjunta JWT para users-api (usa el mismo jwtSecret que lee de env JWT_SECRET)
	if tok, err := s.getUserAPIToken(username); err == nil && tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return user, fmt.Errorf("users-api %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &user); err != nil {
		return user, err
	}
	return user, nil
}

func (s *UserService) getUserAPIToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["scope"] = "read"
	return token.SignedString([]byte(jwtSecret))
}
