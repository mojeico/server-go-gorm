package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"myMod/entity"
)

type VideoRepository interface {
	Create(video entity.Video)
	Update(video entity.Video)
	Delete(video entity.Video)
	FindAll() []entity.Video
}

type database struct {
	connection *gorm.DB
}

func NewVideoRepository() VideoRepository {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=test  sslmode=disable password=postgres")
	if err != nil {
		fmt.Println(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, tableName string) string { return "test." + tableName }

	return &database{
		connection: db,
	}
}

func (d database) Create(video entity.Video) {
	d.connection.Create(&video)
}

func (d database) Update(video entity.Video) {
	d.connection.Save(&video)
}

func (d database) Delete(video entity.Video) {
	d.connection.Delete(&video)
}

func (d database) FindAll() []entity.Video {
	var videos []entity.Video
	d.connection.Set("gorm:auto_preload", true).Find(&videos, &entity.Person{})
	return videos
}
