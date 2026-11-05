package main

import (
	"fmt"
	"log"

	"github.com/kenzosashida090/metro-math-go/db"
	"github.com/kenzosashida090/metro-math-go/internal/env"
	"github.com/kenzosashida090/metro-math-go/store"
)

//	@title		metro-math-go
//	@description	This is an API for a metronome.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

// @description
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":3030"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  env.GetString("MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxIdleConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	fmt.Println("DB CONNECTION ESTABLSIIHED", env.GetString("DB_ADDR", ""), env.GetString("ADDR", ":3030"))
	store := store.NewStorage(db)
	app := &application{
		config:  cfg,
		version: "0.0.1",
		env:     "DEV",
		store:   store,
	}
	mux := app.mount()
	err = app.run(mux)
	if err != nil {
		log.Panic(err)
	}
}
