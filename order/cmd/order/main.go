package main

import (
	"github.com/Mostbesep/microservice-com-temp/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"time"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL" required:"true"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL" required:"true"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		r, err = order.NewPostgresqlRepository(cfg.DatabaseURL)
		if err != nil {
			println(err)
		}
		return err
	})
	defer r.Close()
	log.Println("Listening on port 8080...")
	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))
}
