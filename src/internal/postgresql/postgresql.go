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
	if err != nil {
		panic("can't migrate databases")
	}
	return &Postgres{db}
}

func (p *Postgres) Stop() error {
	val, err := p.Db.DB()
	if err != nil {
		return errors.New("failed to get database; error: " + err.Error())
	}
	return val.Close()
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
	return tx.Commit().Error
}

func (p *Postgres) Get(feature, tag int32) (models.JSONMap, error) {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return nil, errors.New("can't start transaction; error: " + tx.Error.Error())
	}
	idToFind := Banner{
		Feature: feature,
		Tag:     tag,
	}
	err := p.Db.Model(Banner{
		Feature: feature,
		Tag:     tag,
	}).First(&idToFind).Error
	if err != nil {
		return nil, errors.New("can't find banner: " + err.Error())
	}
	result := Data{}
	err = p.Db.Model(Data{Id: idToFind.DataId}).First(&result).Error
	if err != nil {
		return nil, errors.New("can't find banner: " + err.Error())
	}
	return result.Content, tx.Commit().Error
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
