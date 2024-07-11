package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type authedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.wrapped.RoundTrip(req)
}

type AccessClient interface {
	CreateToken(subject string) (string, error)
}

type accessClient struct {
	routerUrl  string
	httpClient *http.Client
}

func (a *accessClient) CreateToken(subject string) (string, error) {
	payload := TokenRequest{
		GrantType:   "client_credentials",
		Scope:       getScope(subject),
		Refreshable: false,
		ExpiresIn:   1600,
		Audience:    "jfrt@*",
		Issuer:      "oidc@poc",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("failed to parse request " + err.Error())
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, a.routerUrl+"/access/api/v1/oauth/token", strings.NewReader(string(body)))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		fmt.Println("failed to create request " + err.Error())
		return "", err
	}
	resp, err := a.httpClient.Do(request)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("request failed with status code: " + string(rune(resp.StatusCode)))
		return "", fmt.Errorf("request failed")
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to read body " + err.Error())
		return "", err
	}
	var tokenRes = &TokenResponse{}
	err = json.Unmarshal(all, tokenRes)
	if err != nil {
		fmt.Println("failed to parse response")
	}
	return tokenRes.AccessToken, nil
}

func getScope(subject string) string {
	return federatedCreds[subject]
}

var federatedCreds = map[string]string{
	"repo:mosheya/access-oidc-poc:ref:refs/heads/dev":  "applied-permissions/groups:oidc-poc",
	"repo:mosheya/access-oidc-poc:ref:refs/heads/main": "applied-permissions/groups:oidc-poc",
}

func NewAccessClient(routerUrl string, token string) AccessClient {
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &authedTransport{
			token:   token,
			wrapped: http.DefaultTransport,
		},
	}
	return &accessClient{
		routerUrl:  routerUrl,
		httpClient: client,
	}
}

type TokenRequest struct {
	GrantType   string `json:"grant_type"`
	Scope       string `json:"scope,omitempty"`
	Refreshable bool   `json:"refreshable"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	Audience    string `json:"audience,omitempty"`
	Issuer      string `json:"issuer,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
