package postgresql

import (
	"banner/internal/config"
	"banner/models"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
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
	Id       int32          `gorm:"primary_key;auto_increment"`
	Content  models.JSONMap `gorm:"type:json;default:'{\"key\": \"value\"}';not null"`
	IsActive bool           `gorm:"type:boolean;default:true;"`
}

func NewPostgresRepository(cfg config.DbConfig) *Postgres {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DbName, cfg.Port)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
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

func (p *Postgres) Insert(record *models.InsertData) error {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return errors.New("can't start transaction; error: " + tx.Error.Error())
	}

	d := Data{
		Content:  record.Content,
		IsActive: record.IsActive,
	}
	if err := tx.Create(&d).Error; err != nil {
		tx.Rollback()
		return errors.New("can't insert data: " + err.Error())
	}
	var banners []Banner
	for _, i := range record.TagIds {
		banners = append(banners, Banner{d.Id, record.Feature, i})
	}
	if err := tx.Create(&banners).Error; err != nil {
		tx.Rollback()
		return errors.New("can't insert banner: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("can't commit transaction: " + err.Error())
	}

	return nil
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

func (p *Postgres) Update(feature, tag int32, newValue models.JSONMap) error {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	id, err := p.findId(feature, tag, tx)
	if err != nil {
		return err
	}
	errUpd := tx.Model(&Data{}).Where("id = ?", id).Updates(map[string]interface{}{"content": newValue})
	if errUpd.Error != nil {
		tx.Rollback()
		return errors.New("can't update banner: " + errUpd.Error.Error())
	}
	return tx.Commit().Error
}

func (p *Postgres) Delete(feature, tag int32) error {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	id, err := p.findId(feature, tag, tx)
	if err != nil {
		return err
	}
	if err := tx.Delete(&Data{}, id).Error; err != nil {
		tx.Rollback()
		return errors.New("can't delete data: " + err.Error())
	}
	if err := tx.Where("data_id = ?", id).Delete(&Banner{}).Error; err != nil {
		tx.Rollback()
		return errors.New("can't delete banner: " + err.Error())
	}
	if tx.Error != nil {
		return errors.New("something went wrong: " + tx.Error.Error())
	}
	return tx.Commit().Error
}

func (p *Postgres) Fill() error {
	for i := 1; i < 10; i++ {
		for j := 1; j <= 20; j++ {
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
			err := p.Insert(&banner)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
