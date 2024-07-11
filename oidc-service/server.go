package main

import (
	"github.com/mosheya/access-oidc-poc/oidc-service/config"
	"github.com/mosheya/access-oidc-poc/oidc-service/handler"
	"log"
	"net/http"
	"os"
)

var conf config.Config

func init() {
	conf.RouterUrl = os.Getenv("ROUTER_URL")
	conf.AccessServiceAdminToken = os.Getenv("JF_ACCESS_ADMIN_TOKEN")
	conf.ProviderUrl = "https://token.actions.githubusercontent.com"
	conf.Audience = "access-oidc-poc"
}

func main() {
	http.Handle("/token", handler.NewHandler(conf))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
