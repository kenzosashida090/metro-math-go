package main

import "github.com/kenzosashida090/metro-math-go/store"

type application struct {
	config config
	store  store.Storage
	env    string
}

type config struct {
	addt        string
	frontEndUrl string
	db          dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
