package api

import (
	"../db"
)

type Config struct {
	DbAPI db.DbAPI
}

type Api struct {
	DbAPI db.DbAPI
}

func New(cfg Config) Api {
	return Api{DbAPI: cfg.DbAPI}
}
