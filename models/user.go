package models

import "time"

type User struct {
	ID         uint   `gorm:"primary_key"`
	UserNO     string `gorm:"column:user_no"`
	FaceId     uint   `gorm:"column:face_id"`
	Name       string `gorm:"column:name"`
	Tel        string `gorm:"column:tel"`
	IdNum      string `gorm:"column:id_num"`
	Age        uint   `gorm:"column:age"`
	CompanyId  uint   `gorm:"column:company_id"`
	Company    string `gorm:"column:company"`
	Department string `gorm:"column:department"`
	FaceImage  string `gorm:"column:face_image"`
	AppId      uint   `gorm:"column:app_id"`
	Group      string `gorm:"column:group"`
	CreatedAt  int64  `gorm:"column:created_at"`
}

func (this *User) TableName() string {
	return "t_user"
}

func (this *User) Save() (err error) {
	this.CreatedAt = time.Now().Unix()
	return db.Save(this).Error
}

func (this *User) Create() (err error) {

	return
}

func FindUserByNOAndCompany(userNO string, companyId uint) (item User, err error) {
	s := db.Where("user_no = ?", userNO).Where("company_id = ?", companyId).First(&item)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}

func FindUserByFaceId(faceId uint) (item User, err error) {
	s := db.Where("face_id = ?", faceId).First(&item)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}

func FindUserByFaceIdAndCompany(faceId uint, companyId uint) (item User, err error) {
	if faceId == 0 || companyId == 0 {
		return
	}
	s := db.Where("face_id = ?", faceId).Where("company_id = ?", companyId).First(&item)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}

func FindUserByBaiDuId(baiDuId string) (item User, err error) {
	if baiDuId == "" {
		return
	}
	s := db.Table("t_user").
		Joins("inner join t_face on t_user.face_id=t_face.id").
		Where("t_face.baidu_user_id = ?", baiDuId).Last(&item)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}

func FindUserByMobile(mobile string) (item User, err error) {
	s := db.Where("tel = ?", mobile).Order("id desc").First(&item)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}

	return
}
