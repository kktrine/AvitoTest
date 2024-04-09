package postgresql

import (
	"banner/internal/config"
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
	id      int `gorm:"foreignKey:id;references:id"`
	Feature int `gorm:"index"`
	Tag     int `gorm:"type:integer[]"`
}

type Data struct {
	id       int    `gorm:"primary_key;auto_increment"`
	context  string `gorm:"type:text;default:'this is json';not null"`
	isActive bool   `gorm:"type:boolean;default:true;"`
}

func NewPostgresRepository(cfg config.DbConfig) *Postgres {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DbName, cfg.Port)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic("couldn't connect to database")
	}
	err = db.AutoMigrate(&Banner{})
	if err != nil {
		panic("can't migrate database")
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

func (p *Postgres) InsertErrorHandler(b *Banner) error {
	err := p.Db.Create(b).Error
	if err != nil {

		return errors.New("can't insert banner into existing records; error: " + err.Error())
	}
	return nil
}

func (p *Postgres) Fill() error {

	for i := 1; i < 100; i++ {
		for j := 1; j <= 20; j++ {
			count := rand.Intn(4) + 2
			banner := Banner{
				Feature: i,
			}
			for k := j; k < j+count; k++ {
				banner.Tags = append(banner.Tags, int32(k))
			}
			j += count
			err := p.InsertErrorHandler(&banner)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
