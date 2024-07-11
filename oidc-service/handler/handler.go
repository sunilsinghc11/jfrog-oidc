package handler

import (
	"encoding/json"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mosheya/access-oidc-poc/oidc-service/client"
	"github.com/mosheya/access-oidc-poc/oidc-service/config"
	"github.com/mosheya/access-oidc-poc/oidc-service/token"
	"io"
	"log"
	"net/http"
)

type tokenHandler struct {
	accessClient client.AccessClient
	githubClient client.GithubClient
}

func NewHandler(config config.Config) http.Handler {
	accessClient := client.NewAccessClient(config.RouterUrl, config.AccessServiceAdminToken)
	oidcClient := client.NewGithubClient(config.ProviderUrl)
	return tokenHandler{
		accessClient: accessClient,
		githubClient: oidcClient,
	}
}

func (h tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read body")
		w.WriteHeader(400)
		return
	}
	githubReq, err := getOidcRequest(payload)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	keyFunc, err := h.getKeyFunc()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	oidcClaim, err := token.VerifyToken(githubReq.Token, keyFunc)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("token verification failed"))
		return
	}
	createToken, err := h.accessClient.CreateToken(oidcClaim.Subject)
	if err != nil {
		log.Println("failed to crate access token")
		w.WriteHeader(400)
		return
	}
	fmt.Println(createToken)
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("{\"access_token\":\"%s\"}", createToken)))
}

type GithubOidcRequest struct {
	Token string `json:"token,omitempty"`
}

func getOidcRequest(payload []byte) (*GithubOidcRequest, error) {
	log.Println("request body: " + string(payload))
	var githubReq = GithubOidcRequest{}
	err := json.Unmarshal(payload, &githubReq)
	if err != nil {
		log.Println("failed to parse github request: " + string(payload))
		return nil, err
	}
	return &githubReq, nil
}

func (h tokenHandler) getKeyFunc() (func(token2 *jwt.Token) (interface{}, error), error) {
	certificate := h.githubClient.GetCertificate()
	marshal, err := json.Marshal(certificate)
	if err != nil {
		log.Println("failed to marshal jwks")
		return nil, err
	}
	jwks, err := keyfunc.NewJSON(marshal)
	if err != nil {
		log.Println("failed to parse jwks")
		return nil, err
	}
	return jwks.Keyfunc, nil
}
