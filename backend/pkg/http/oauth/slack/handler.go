package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"lunch/pkg/http/auth"
	"lunch/pkg/jwt"
	"lunch/pkg/users"
	service_users "lunch/pkg/users/service"
)

func Handler(cfg *Configuration, jwtService *jwt.Service, usersService *service_users.Service) http.HandlerFunc {
	type request struct {
		Code        string `json:"code"`
		RedirectURI string `json:"redirectUri"`
	}

	type slackAuthedUser struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
	}
	type slackOAuthResponse struct {
		OK         bool             `json:"ok"`
		Error      string           `json:"error"`
		AuthedUser *slackAuthedUser `json:"authed_user"`
	}

	type slackUser struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	type slackIdentityResponse struct {
		OK    bool       `json:"ok"`
		Error string     `json:"error"`
		User  *slackUser `json:"user"`
	}

	client := http.Client{}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[ERROR] failed to decode request: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		form := url.Values{}
		form.Set("client_id", cfg.ClientID)
		form.Set("client_secret", cfg.ClientSecret)
		form.Set("code", req.Code)
		form.Set("redirect_uri", req.RedirectURI)
		form.Set("grant_type", "authorization_code")

		resp, err := client.Post(
			"https://slack.com/api/oauth.v2.access",
			"application/x-www-form-urlencoded",
			strings.NewReader(form.Encode()),
		)
		if err != nil {
			log.Printf("[ERROR] failed to post request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var oauthResponse slackOAuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&oauthResponse); err != nil {
			log.Printf("[ERROR] failed to decode oauth response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !oauthResponse.OK {
			log.Printf("[ERROR] failed to get access token: %s", oauthResponse.Error)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		identityRequest, err := http.NewRequest(http.MethodGet, "https://slack.com/api/users.identity", nil)
		if err != nil {
			log.Printf("[ERROR] failed to create request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		identityRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauthResponse.AuthedUser.AccessToken))

		identityResponse, err := client.Do(identityRequest.WithContext(r.Context()))
		if err != nil {
			log.Printf("[ERROR] failed to get user identity: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer identityResponse.Body.Close()

		var identityResponseBody slackIdentityResponse
		if err := json.NewDecoder(identityResponse.Body).Decode(&identityResponseBody); err != nil {
			log.Printf("[ERROR] failed to decode identity response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !identityResponseBody.OK {
			log.Printf("[ERROR] failed to get user identity: %s", identityResponseBody.Error)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user := &users.User{
			ID:   identityResponseBody.User.ID,
			Name: identityResponseBody.User.Name,
		}

		if err := usersService.Create(r.Context(), user); err != nil {
			log.Printf("[ERROR] failed to create user: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token, err := jwtService.NewToken(r.Context(), user)
		if err != nil {
			log.Printf("[ERROR] failed to generate token: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		secure := r.TLS != nil
		auth.SetCookie(w, token, secure)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Printf("[ERROR] failed to encode response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
