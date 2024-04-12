package postgresql

import (
	"banner/internal/config"
	"banner/models"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
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
	IsActive  bool           `gorm:"type:boolean;default:true;"`
	CreatedAt time.Time      `gorm:"autoUpdateTime:milli"`
	UpdatedAt time.Time      `gorm:"autoCreateTime"`
}

func NewPostgresRepository(cfg config.DbConfig) *Postgres {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DbName, cfg.Port)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	rawDB, _ := db.DB()
	rawDB.SetMaxOpenConns(2)
	if err != nil {
		panic("couldn't connect to database")
	}
	err = db.AutoMigrate(&Banner{}, &Data{})
	if err := db.AutoMigrate(&Banner{}, &Data{}); err != nil {
		panic("can't migrate databases")
	}
	db.Exec("DELETE  FROM data;")
	db.Exec("DELETE FROM banners;")
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

func (p *Postgres) Get(feature, tag int32) (data models.JSONMap, userAccess bool, found bool, err error) {
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
	if err := tx.Where("feature = ? AND tag = ?", feature, tag).First(&idToFind).Error; err != nil {
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

func (p *Postgres) GetMany(featureId int32, tagId int32, limit int32, offset int32) {
	//tx := p.Db.Begin()
	//defer func() {
	//	if r := recover(); r != nil {
	//		tx.Rollback()
	//	}
	//}()
	//res := make([]map[string]interface{}, 0)
	//elem := map[string]interface{}{
	//	"feature_id": 1,
	//	"tag_ids":    []int{},
	//	"is_active":  true,
	//	"updated_at": "2000-01-23T04:56:07.000+00:00",
	//	"banner_id":  0,
	//	"created_at": "2000-01-23T04:56:07.000+00:00",
	//	"content":    map[string]interface{}{},
	//}

}

func (p *Postgres) Fill() error {
	for i := 1; i < 1001; i++ {
		for j := 1; j <= 10; j++ {
			count := rand.Intn(4) + 2
			banner := models.InsertData{}
			val := int32(i)
			banner.Feature = val
			tags := make([]int32, 0, count)
			for k := int32(j); k < int32(j+count); k++ {
				tags = append(tags, k)
			}
			banner.TagIds = tags
			j += count
			_, err := p.Insert(&banner)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
