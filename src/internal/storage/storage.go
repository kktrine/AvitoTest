package storage

import (
	"banner/internal/cashe"
	"banner/internal/config"
	"banner/internal/postgresql"
	"banner/models"
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

func (s *Storage) Insert(record *models.InsertData) (int32, error) {
	return s.db.Insert(record)
}

func (s *Storage) GetUserBanner(feature, tag int32) (models.JSONMap, bool, bool, error) {
	return s.db.Get(feature, tag)
}

func (s *Storage) Update(id int32, record *models.InsertData) (bool, error) {
	//fmt.Println(s.db.Update(feature, tag, newValue))
	return s.db.Update(id, record)
}

func (s *Storage) Delete(id int32) (bool, error) {
	return s.db.Delete(id)
}

func (s *Storage) GetMany(featureId int32, tagId int32, limit int32, offset int32) ([]map[string]interface{}, error) {
	return s.db.GetMany(featureId, tagId, limit, offset)
}
