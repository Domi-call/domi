package bolt

import (
	"errors"
	"gorm.io/gorm"

	fbErrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/share"
)

type shareBackend struct {
	//db *storm.DB
	db *gorm.DB
}

func (s shareBackend) All() ([]*share.Link, error) {
	var v []*share.Link
	err := s.db.Model(&share.Link{}).Find(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v, fbErrors.ErrNotExist
	}
	return v, err

	//var v []*share.Link
	//err := s.db.All(&v)
	//if errors.Is(err, storm.ErrNotFound) {
	//	return v, fbErrors.ErrNotExist
	//}
	//
	//return v, err
}

func (s shareBackend) FindByUserID(id uint) ([]*share.Link, error) {
	var v []*share.Link
	err := s.db.Model(&share.Link{}).Where("user_id = ?", id).Find(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v, fbErrors.ErrNotExist
	}
	return v, err

	//var v []*share.Link
	//err := s.db.Select(q.Eq("UserID", id)).Find(&v)
	//if errors.Is(err, storm.ErrNotFound) {
	//	return v, fbErrors.ErrNotExist
	//}
	//
	//return v, err
}

func (s shareBackend) GetByHash(hash string) (*share.Link, error) {
	var v share.Link
	err := s.db.Model(&share.Link{}).Where("hash = ?", hash).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fbErrors.ErrNotExist
	}
	return &v, err

	//var v share.Link
	//err := s.db.One("Hash", hash, &v)
	//if errors.Is(err, storm.ErrNotFound) {
	//	return nil, fbErrors.ErrNotExist
	//}
	//
	//return &v, err
}

func (s shareBackend) GetPermanent(path string, id uint) (*share.Link, error) {
	var v share.Link
	err := s.db.Model(&share.Link{}).Where("path = ? AND expire = 0 AND user_id = ?", path, id).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fbErrors.ErrNotExist
	}
	return &v, err
	//var v share.Link
	//err := s.db.Select(q.Eq("Path", path), q.Eq("Expire", 0), q.Eq("UserID", id)).First(&v)
	//if errors.Is(err, storm.ErrNotFound) {
	//	return nil, fbErrors.ErrNotExist
	//}
	//
	//return &v, err
}

func (s shareBackend) Gets(path string, id uint) ([]*share.Link, error) {
	var v []*share.Link
	err := s.db.Model(&share.Link{}).Where("path = ? AND user_id = ?", path, id).Find(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v, fbErrors.ErrNotExist
	}
	return v, err

	//var v []*share.Link
	//err := s.db.Select(q.Eq("Path", path), q.Eq("UserID", id)).Find(&v)
	//if errors.Is(err, storm.ErrNotFound) {
	//	return v, fbErrors.ErrNotExist
	//}
	//
	//return v, err
}

func (s shareBackend) Save(l *share.Link) error {
	return s.db.Save(l).Error
	//return s.db.Save(l)
}

func (s shareBackend) Delete(hash string) error {
	return s.db.Where("hash = ?", hash).Delete(&share.Link{}).Error
	//err := s.db.DeleteStruct(&share.Link{Hash: hash})
	//if errors.Is(err, storm.ErrNotFound) {
	//	return nil
	//}
	//return err
}
