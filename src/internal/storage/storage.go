package storage

import (
	"banner/internal/cashe"
	"banner/internal/config"
	"banner/internal/postgresql"
	"fmt"
	"time"
)

type Storage struct {
	cfg   *config.Config
	db    *postgresql.Postgres
	cache *cashe.Cache
}

func NewStorage() *Storage {
	cfg := config.MustLoad()
	db := postgresql.NewPostgresRepository(cfg.DbConfig)
	cache := cashe.NewCache(5*time.Minute, 5*time.Minute+30*time.Second)
	return &Storage{
		cfg:   cfg,
		db:    db,
		cache: cache,
	}
}

func (s *Storage) Fill() error {
	return s.db.Fill()
}

func (s *Storage) Get(feature, tag int32) {
	fmt.Println(s.db.Get(feature, tag))
}
