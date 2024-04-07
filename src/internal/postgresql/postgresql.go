package postgresql

import (
	"banner/internal/config"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
)

type Postgres struct {
	Db *gorm.DB
}

type Banner struct {
	Feature  int
	Tags     pq.Int32Array `gorm:"type:integer[]"`
	Content  string
	IsActive bool `gorm:"default:true"`
}

func NewPostgresRepository(cfg config.DbConfig) *Postgres {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DbName, cfg.Port)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic("couldn't connect to database")
	}
	db.AutoMigrate(&Banner{})

	return &Postgres{db}
}

func (p *Postgres) Stop() error {
	val, err := p.Db.DB()
	if err != nil {
		return errors.New("failed to get database; error: " + err.Error())
	}
	return val.Close()
}

//func (p *Postgres) InsertBanner(banner *models.Banner) error {
//	res := models.Banner{}
//	for _, tag := range banner.TagIds {
//		p.Db.Where("feature = ? AND ? = ANY (tags)", banner.FeatureId, tag).First(&res)
//		if len(res.TagIds) > 0 {
//			return errors.New("no unique pairs tag-feature found")
//		}
//	}
//	return p.Db.Create(banner).Error
//}

func (p *Postgres) Fill() error {

	for i := 1; i < 500; i++ {
		for j := 1; j <= 20; j++ {
			count := rand.Intn(4) + 2
			banner := Banner{
				Feature: i,
			}
			for k := j; k < j+count; k++ {
				banner.Tags = append(banner.Tags, int32(k))
			}
			j += count
			res := p.Db.Create(&banner)
			if res.Error != nil {
				return res.Error
			}
		}
	}
	return nil
}
