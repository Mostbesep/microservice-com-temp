package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL" required:"true"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL" required:"true"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL" required:"true"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	// use deprecated NewDefaultServer instead New reason: use playground option in browser
	// for handle err in response:
	// [{"message":"transport not supported"}],"data":null}
	http.Handle("/graphql", handler.NewDefaultServer(s.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
