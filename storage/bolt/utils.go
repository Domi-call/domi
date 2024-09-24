package bolt

import (
	"errors"
	"gorm.io/gorm"

	fbErrors "github.com/filebrowser/filebrowser/v2/errors"
)

//func get(db *storm.DB, name string, to interface{}) error {
//	err := db.Get("config", name, to)
//	if errors.Is(err, storm.ErrNotFound) {
//		return fbErrors.ErrNotExist
//	}
//
//	return err
//}

func get(db *gorm.DB, name string, to interface{}) error {
	err := db.Table(name).First(to).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fbErrors.ErrNotExist
	}
	return err
}

//	func save(db *storm.DB, name string, from interface{}) error {
//		return db.Set("config", name, from)
//	}
//func save(db *gorm.DB, name string, from interface{}) error {
//	return db.Set("config", name, from)
//}
