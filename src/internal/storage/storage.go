package storage

import (
	"banner/internal/cashe"
	"banner/internal/postgresql"
	"banner/models"
)

type Storage struct {
	db    *postgresql.Postgres
	cache *cashe.Cache
}

func NewStorage() *Storage {
	db := postgresql.NewPostgresRepository()
	cache := cashe.NewCache()
	return &Storage{
		db:    db,
		cache: cache,
	}
}

func (s *Storage) Insert(record *models.InsertData) (int32, error) {
	id, err := s.db.Insert(record)
	if err != nil {
		return 0, err
	}
	s.cache.AddOne(cashe.Item{
		BannerID:  id,
		FeatureID: record.Feature,
		TagIDs:    record.TagIds,
		IsActive:  record.IsActive,
		Content:   record.Content,
	})
	return id, nil
}

func (s *Storage) GetUserBanner(feature, tag int32, fromBD bool) (models.JSONMap, bool, bool, error) {
	if !fromBD {
		content, userAccess := s.cache.Get(feature, tag)
		if content != nil {
			//fmt.Printf("from cache: feature: %d, tag: %d!\n", feature, tag)
			return content, userAccess, true, nil
		}
	}
	content, id, userAccess, found, err := s.db.Get(feature, tag)
	if err != nil {
		return nil, false, false, err
	}
	s.cache.AddOne(cashe.Item{
		BannerID:  id,
		FeatureID: feature,
		TagIDs:    []int32{tag},
		IsActive:  userAccess,
		Content:   content,
	})
	return content, userAccess, found, err
}

func (s *Storage) Update(id int32, record *models.InsertData) (bool, error) {
	return s.db.Update(id, record)
}

func (s *Storage) Delete(id int32) (bool, error) {
	return s.db.Delete(id)
}

func (s *Storage) GetMany(featureId int32, tagId int32, limit int32, offset int32) ([]map[string]interface{}, error) {
	return s.db.GetMany(featureId, tagId, limit, offset)
}

func (s *Storage) Stop() error {
	return s.db.Stop()
}
