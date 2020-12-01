package models

import (
	"fmt"
	"math/rand"
	"time"
)

type Face struct {
	ID        uint   `gorm:"primary_key"`
	BaiDuID   string `gorm:"column:baidu_user_id"`
	UserInfo  string `gorm:"column:user_info"`
	CreatedAt int64  `gorm:"column:created_at"`
}

func (this *Face) TableName() string {
	return "t_face"
}

func (this *Face) Save() (err error) {
	this.CreatedAt = time.Now().Unix()
	return db.Save(this).Error
}

func (this *Face) Create() (err error) {

	return
}

// 查询，或创建
func (this *Face) Register() (item Face,  err error) {
	// 查询
	if this.BaiDuID == "" {
		return
	}

	//find
	item, err = FindFaceByBaiDuId(this.BaiDuID)
	if err != nil {
		return
	}
	fmt.Println(item)
	if item.ID != 0 {
		return
	}

	item.BaiDuID = this.BaiDuID
	item.UserInfo = this.UserInfo
	err = item.Save()
	if err != nil {
		return
	}

	return
}

func FindFaceByBaiDuId(baiDuId string) (item Face, err error) {
	s := db.Where("baidu_user_id = ?", baiDuId).First(&item)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}

//TODO
func CreatRandBaiDuId() (string) {
	l := 10
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
