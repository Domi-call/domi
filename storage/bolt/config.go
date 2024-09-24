package bolt

import (
	"errors"
	"github.com/filebrowser/filebrowser/v2/storage/mysql"
	"gorm.io/gorm"

	"github.com/filebrowser/filebrowser/v2/settings"
)

type settingsBackend struct {
	//db *storm.DB
	db *gorm.DB
}

func (s settingsBackend) Get() (*settings.Settings, error) {
	set := &settings.Settings{}
	setTb := &settings.SettingsTb{}
	err := s.db.Table("settings").First(setTb).Error
	if err != nil {
		return nil, err
	}
	set, err = settings.SettingsTb2Settings(setTb)
	return set, err
	//return set, get(s.db, "settings", set)
}

func (s settingsBackend) Save(set *settings.Settings) error {
	//不存在记录则创建，存在则更新
	setTb := &settings.SettingsTb{}
	err := mysql.DB.Model(&settings.SettingsTb{}).First(setTb).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		tb, err := settings.Settings2SettingsTb(set)
		if err != nil {
			return err
		}
		return mysql.DB.Model(tb).Create(tb).Error

	}
	if err != nil {
		return err
	}
	tb, err := settings.Settings2SettingsTb(set)
	if err != nil {
		return err
	}
	return mysql.DB.Model(setTb).Updates(tb).Error

	//tb, err := settings.Settings2SettingsTb(set)
	//if err != nil {
	//	return err
	//}
	//return mysql.DB.Model(tb).Create(tb).Error
	//return save(s.db, "settings", set)
}

func (s settingsBackend) GetServer() (*settings.Server, error) {
	server := &settings.Server{}
	err := s.db.Model(server).First(server).Error
	if err != nil {
		return nil, err
	}
	return server, nil
	//return server, get(s.db, "server", server)
}

func (s settingsBackend) SaveServer(server *settings.Server) error {
	return mysql.DB.Model(server).Create(server).Error
	//return save(s.db, "server", server)
}
