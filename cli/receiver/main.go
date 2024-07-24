package main

import (
	"flag"

	"github.com/kuznetsovin/egts-protocol/cli/receiver/config"
	"github.com/kuznetsovin/egts-protocol/cli/receiver/server"
	"github.com/kuznetsovin/egts-protocol/cli/receiver/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfgFilePath := ""
	flag.StringVar(&cfgFilePath, "c", "", "Конфигурационный файл")
	flag.Parse()

	if cfgFilePath == "" {
		log.Fatalf("Не задан путь до конфига")
	}

	cfg, err := config.New(cfgFilePath)
	if err != nil {
		log.Fatalf("Ошибка парсинга конфига: %v", err)
	}

	ConfigureLogger(cfg.GetLogLevel())

	storages := storage.NewRepository()
	if err := storages.LoadStorages(cfg.Storage); err != nil {
		log.Errorf("ошибка загрузка хранилища: %v", err)

		// TODO: clear after test
		store := storage.LogConnector{}
		if err := store.Init(nil); err != nil {
			log.Fatal(err)
		}

		storages.AddStore(store)
		defer store.Close()
	}

	srv := server.New(cfg.GetListenAddress(), cfg.GetEmptyConnTTL(), storages)

	srv.Run()
}
