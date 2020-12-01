package models

func Migrate() error {
	return db.AutoMigrate().Error

}
