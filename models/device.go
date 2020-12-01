package models

type Device struct {
	ID       uint   `gorm:"primary_key"`
	SN       string `gorm:"column:sn"`
	MacAddr  string `gorm:"column:mac"`
	StoreId  uint   `gorm:"column:store_id"`
	Province uint   `gorm:"column:province"`
	City     uint   `gorm:"column:city"`
}

func (m *Device) TableName() string {
	return "t_device"
}

func (this *Device) FindDevByID(devId int) (res *Device, err error) {
	var u Device
	s := db.Where("id = ?", devId).First(&u)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}
	res = &u
	return
}

func (this *Device) FindDevBySN(sn string) (res *Device, err error) {
	var u Device
	s := db.Where("sn = ?", sn).First(&u)
	if s.RecordNotFound() {
		return
	}
	if err = s.Error; err != nil {
		return
	}
	res = &u
	return
}
