package client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type GithubClient interface {
	GetOpenIdConfig() Config
	GetCertificate() Jwks
}

type githubClient struct {
	baseUrl string
	client  http.Client
}

func (g githubClient) GetOpenIdConfig() Config {
	req, _ := http.NewRequest(http.MethodGet, g.baseUrl+"/.well-known/openid-configuration", nil)
	response, err := g.client.Do(req)
	if err != nil {
		log.Println("request failed")
		return Config{}
	}
	bytes, _ := io.ReadAll(response.Body)
	var config = Config{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Println("failed to parse request body")
		return Config{}
	}
	return config
}

func (g githubClient) GetCertificate() Jwks {
	openIdConfig := g.GetOpenIdConfig()
	req, _ := http.NewRequest(http.MethodGet, openIdConfig.JwksUri, nil)
	response, err := g.client.Do(req)
	if err != nil {
		log.Println("request failed")
		return Jwks{}
	}
	bytes, _ := io.ReadAll(response.Body)
	var jwks = Jwks{}
	err = json.Unmarshal(bytes, &jwks)
	if err != nil {
		log.Println("failed to parse request body")
		return Jwks{}
	}
	return jwks
}

func NewGithubClient(baseUrl string) GithubClient {
	return &githubClient{
		baseUrl: baseUrl,
		client: http.Client{
			Timeout:   time.Second * 30,
			Transport: http.DefaultTransport,
		},
	}
}

type Config struct {
	Issuer                           string   `json:"issuer"`
	JwksUri                          string   `json:"jwks_uri"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	ClaimsSupported                  []string `json:"claims_supported"`
	IdTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported                  []string `json:"scopes_supported"`
}

type Jwks struct {
	Keys []struct {
		N   string   `json:"n"`
		Kty string   `json:"kty"`
		Kid string   `json:"kid"`
		Alg string   `json:"alg"`
		E   string   `json:"e"`
		Use string   `json:"use"`
		X5C []string `json:"x5c"`
		X5T string   `json:"x5t"`
	} `json:"keys"`
}
