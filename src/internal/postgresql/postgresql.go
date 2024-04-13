package postgresql

import (
	"banner/internal/config"
	"banner/models"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Postgres struct {
	Db *gorm.DB
}

type Banner struct {
	DataId  int32 `gorm:"foreignKey:id;references:id"`
	Feature int32 `gorm:"uniqueIndex:idx_banner_feature_tag"`
	Tag     int32 `gorm:"uniqueIndex:idx_banner_feature_tag"`
}

type Data struct {
	Id        int32          `gorm:"primary_key;auto_increment"`
	Content   models.JSONMap `gorm:"type:json;default:'{\"key\": \"value\"}';not null"`
	IsActive  bool           `gorm:"type:boolean;default:false;"`
	CreatedAt time.Time      `gorm:"autoUpdateTime:milli"`
	UpdatedAt time.Time      `gorm:"autoCreateTime"`
}

func NewPostgresRepository(cfg config.DbConfig) *Postgres {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DbName, cfg.Port)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	rawDB, _ := db.DB()
	rawDB.SetMaxOpenConns(128)
	rawDB.SetMaxIdleConns(256)

	if err != nil {
		panic("couldn't connect to database")
	}
	err = db.AutoMigrate(&Banner{}, &Data{})
	if err := db.AutoMigrate(&Banner{}, &Data{}); err != nil {
		panic("can't migrate databases")
	}
	return &Postgres{db}
}

func (p *Postgres) Stop() error {
	val, err := p.Db.DB()
	if err != nil {
		return errors.New("failed to get database; error: " + err.Error())
	}
	if err := val.Close(); err != nil {
		return errors.New("failed to close database connection; error: " + err.Error())
	}

	return nil
}

func (p *Postgres) Insert(record *models.InsertData) (int32, error) {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return 0, errors.New("can't start transaction; error: " + tx.Error.Error())
	}

	d := Data{
		Content:  record.Content,
		IsActive: record.IsActive,
	}
	if err := tx.Create(&d).Error; err != nil {
		tx.Rollback()
		return 0, errors.New("can't insert data: " + err.Error())
	}
	banners := make([]Banner, 0, len(record.TagIds))
	for _, i := range record.TagIds {
		banners = append(banners, Banner{DataId: d.Id, Feature: record.Feature, Tag: i})
	}
	if err := tx.Create(&banners).Error; err != nil {
		tx.Rollback()
		return 0, errors.New("can't insert banner: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return 0, errors.New("can't commit transaction: " + err.Error())
	}

	return banners[0].DataId, nil
}

func (p *Postgres) Get(feature, tag int32) (data models.JSONMap, id int32, userAccess bool, found bool, err error) {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	data = nil
	userAccess = false
	found = false
	err = nil

	if tx.Error != nil {
		err = fmt.Errorf("can't start transaction: %w", tx.Error)
		return
	}

	id, errId := p.findId(feature, tag, tx)
	if errId != nil {
		return
	}

	var result Data
	if err = p.Db.Where("id = ?", id).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		err = fmt.Errorf("failed to find banner: %w", err)
		return
	}
	found = true
	data = result.Content
	userAccess = result.IsActive
	if err = tx.Commit().Error; err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
		return
	}

	return
}

func (p *Postgres) findId(feature, tag int32, tx *gorm.DB) (int32, error) {
	var idToFind Banner
	if err := tx.Model(&Banner{}).Where("feature = ? AND tag = ?", feature, tag).First(&idToFind).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("banner with feature %d and tag %d not found", feature, tag)
		}
		return 0, fmt.Errorf("failed to find banner: %w", err)
	}
	return idToFind.DataId, nil
}

func (p *Postgres) Update(id int32, newValue *models.InsertData) (bool, error) {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var count int64
	tx.Model(&Data{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		tx.Rollback()
		return false, nil
	}
	if len(newValue.TagIds) > 0 || newValue.Feature > 0 {
		var deletedBanners []Banner
		tx.Model(&Banner{}).Where("data_id = ?", id).Find(&deletedBanners)
		tx.Model(&Banner{}).Where("data_id = ?", id).Delete(&deletedBanners)
		if newValue.Feature != 0 {
			for i := range deletedBanners {
				deletedBanners[i].Feature = newValue.Feature
			}
		}

		numOfTags := len(newValue.TagIds)
		if numOfTags > 0 {
			for i := 0; i < numOfTags && i < len(deletedBanners); i++ {
				deletedBanners[i].Tag = newValue.TagIds[i]
			}
			for i, feature := len(deletedBanners), deletedBanners[0].Feature; i < numOfTags; i++ {
				deletedBanners = append(deletedBanners, Banner{DataId: id, Feature: feature, Tag: newValue.TagIds[i]})
			}
		}
		err := tx.Model(&Banner{}).Create(deletedBanners)
		if err.Error != nil {
			tx.Rollback()
			return true, errors.New("can't update banner: " + err.Error.Error())
		}

	}
	newData := Data{
		Content:  newValue.Content,
		IsActive: newValue.IsActive,
	}
	errUpd := tx.Model(&Data{}).Where("id = ?", id).Updates(&newData)
	if errUpd.Error != nil {
		tx.Rollback()
		return true, errors.New("can't update banner: " + errUpd.Error.Error())
	}

	return true, tx.Commit().Error
}

func (p *Postgres) Delete(id int32) (bool, error) {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := tx.Delete(&Data{}, id)
	if err.Error != nil {
		tx.Rollback()
		return err.RowsAffected > 0, errors.New("can't delete data: " + err.Error.Error())
	}
	if err.RowsAffected == 0 {
		return false, nil
	}
	err = tx.Where("data_id = ?", id).Delete(&Banner{})
	if err.Error != nil {
		tx.Rollback()
		return err.RowsAffected > 0, errors.New("can't delete banner: " + err.Error.Error())
	}
	if tx.Error != nil {
		tx.Rollback()
		return err.RowsAffected > 0, errors.New("something went wrong: " + tx.Error.Error())
	}
	return err.RowsAffected > 0, tx.Commit().Error
}

func (p *Postgres) GetMany(featureId int32, tagId int32, limit int32, offset int32) ([]map[string]interface{}, error) {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if limit == 0 {
		limit = -1
	}
	if offset == 0 {
		offset = -1
	}
	var ids []int32
	bannersIds := make(map[int32]string)
	var resData []Data
	var resBanners []Banner
	if featureId > 0 && tagId > 0 {
		tx.Model(&Banner{}).Where("feature = ? AND tag = ?", featureId, tagId).Pluck("data_id", &ids)
	} else {
		tx.Model(&Banner{}).Where("feature = ? OR tag = ?", featureId, tagId).Pluck("data_id", &ids)
	}
	tx.Model(&Banner{}).Where("data_id IN (?)", ids).Find(&resBanners)
	tx.Model(&Data{}).Limit(int(limit)).Offset(int(offset)).Where("id IN (?)", ids).Distinct().Find(&resData)
	res := make([]map[string]interface{}, 0, len(resData))
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("something went wrong: " + tx.Error.Error())
	}
	bannerGroups := make(map[string]struct {
		DataID  int32
		Feature int32
		Tags    []int32
	})
	for _, banner := range resBanners {
		key := fmt.Sprintf("%d_%d", banner.DataId, banner.Feature)
		group, ok := bannerGroups[key]
		if !ok {
			group = struct {
				DataID  int32
				Feature int32
				Tags    []int32
			}{
				DataID:  banner.DataId,
				Feature: banner.Feature,
				Tags:    []int32{}}
			bannersIds[banner.DataId] = key
		}
		group.Tags = append(group.Tags, banner.Tag)
		bannerGroups[key] = group

	}

	for _, i := range resData {
		elem := map[string]interface{}{
			"feature_id": 0,
			"tag_ids":    []int{},
			"is_active":  true,
			"updated_at": "",
			"banner_id":  0,
			"created_at": "",
			"content":    map[string]interface{}{},
		}
		elem["content"] = i.Content
		elem["is_active"] = i.IsActive
		elem["updated_at"] = i.UpdatedAt
		elem["created_at"] = i.CreatedAt
		elem["banner_id"] = i.Id
		elem["tag_ids"] = bannerGroups[bannersIds[i.Id]].Tags
		elem["feature_id"] = bannerGroups[bannersIds[i.Id]].Feature
		res = append(res, elem)
	}

	return res, tx.Commit().Error
}
