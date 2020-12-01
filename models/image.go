package models

import (
	"errors"
	"fmt"
	"time"
)

type Image struct {
	ID        uint   `gorm:"primary_key"`
	StoreId   uint   `gorm:"column:store_id" json:"store_id"`
	SN        string `gorm:"column:sn" json:"sn"`
	Url       string `gorm:"column:url" json:"url"`
	CreatedAt int64  `gorm:"column:created_at"`
}

func (this *Image) TableName() string {
	return "t_image"
}

func (this *Image) Save() (err error) {
	this.CreatedAt = time.Now().Unix()
	return db.Save(this).Error
}

func (this *Image) Find(id uint) (item Image, err error) {
	s := db.First(&item, id)
	if s.RecordNotFound() {
		err = fmt.Errorf("not found, id=%d", id)
		return
	}
	if err = s.Error; err != nil {
		return
	}
	return
}

func (this *Image) GetInfoBySN(SN string) (item Image, err error) {
	s := db.Where("sn = ?", SN).First(&item)
	if s.RecordNotFound() {
		err = errors.New("record not found")
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}
