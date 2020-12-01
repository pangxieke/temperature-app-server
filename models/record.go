package models

import "time"

type Record struct {
	ID        uint    `gorm:"primary_key"`
	Num       float64 `gorm:"column:num"`
	SN        string  `gorm:"column:sn"`
	StoreId   uint    `gorm:"column:store_id"`
	UserId    uint    `gorm:"column:user_id"`
	FaceId    uint    `gorm:"column:face_id"`
	Type      uint    `gorm:"column:type"`
	FaceImage string  `gorm:"column:face_image"`
	CreatedAt int64   `gorm:"column:created_at"`
}

func (this *Record) TableName() string {
	return "t_record"
}

func (this *Record) Save() (err error) {
	this.CreatedAt = time.Now().Unix()
	return db.Save(this).Error
}

func (this *Record) Create() (err error) {

	return
}

func FindLastByFaceId(faceId uint) (item Record, err error) {
	s := db.Where("face_id = ?", faceId).Last(&item)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}

func FindListByFaceId(faceId uint) (items []Record, err error) {
	s := db.Where("face_id = ?", faceId).Order("id desc").Find(&items)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}

type RecordSN struct {
	SN       string `gorm:"column:sn"`
	Province uint `gorm:"column:province"`
	City     uint `gorm:"city"`
}

func FindListSN(faceId uint, startTime, endTime int) (items []RecordSN, err error) {
	s := db.Table("t_record").
		Select("t_record.sn, t_device.city, t_device.province").
		Where("face_id = ?", faceId).
		Joins("left join t_device on t_device.sn=t_record.sn")

	if startTime > 0 {
		s = s.Where("t_record.created_at >= ?", startTime)
	}
	if endTime != 0 {
		s = s.Where("t_record.created_at <= ?", endTime)
	}
	s = s.Group("t_device.city").Scan(&items)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}

func FindListBySN(faceId uint, city uint, startTime, endTime int) (items []Record, err error) {
	s := db.Table("t_record").
		//Select("t_record.sn, t_device.mac").
		Where("face_id = ?", faceId).
		Joins("inner join t_device on t_device.sn=t_record.sn").
		Where("city=?", city)
	if startTime > 0 {
		s = s.Where("t_record.created_at >= ?", startTime)
	}
	if endTime != 0 {
		s = s.Where("t_record.created_at <= ?", endTime)
	}
	s = s.Order("t_record.id desc").Find(&items)

	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return

}
