package storage

import (
	"banner/internal/cashe"
	"banner/internal/config"
	"banner/internal/postgresql"
	"banner/models"
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

func (s *Storage) GetUserBanner(feature, tag int32) (models.JSONMap, bool, bool, error) {
	//fmt.Println(s.db.Get(feature, tag))
	return s.db.Get(feature, tag)

}

func (s *Storage) Update(feature, tag int32, newValue models.JSONMap) {
	fmt.Println(s.db.Update(feature, tag, newValue))

}

func (s *Storage) Delete(feature, tag int32) {
	fmt.Println(s.db.Delete(feature, tag))
}
