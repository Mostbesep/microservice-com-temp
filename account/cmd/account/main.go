package main

import (
	"github.com/Mostbesep/microservice-com-temp/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"time"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}
	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(i int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
			return err
		}
		defer r.Close()
		log.Println("server is start at ", time.Now(), " Listening on port 8080...")
		s := account.NewAccountService(r)
		log.Fatal(account.ListenGRPC(s, 8080))
		return
	})
}
