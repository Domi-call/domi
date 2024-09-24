package bolt

import (
	"encoding/json"
	"gorm.io/gorm"

	"github.com/filebrowser/filebrowser/v2/auth"
	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/settings"
)

type authBackend struct {
	//db *storm.DB
	db *gorm.DB
}

type AutherTb struct {
	Id        int             `gorm:"primaryKey"`
	JsonAuth  json.RawMessage `gorm:"type:json"`
	ProxyAuth json.RawMessage `gorm:"type:json"`
	HookAuth  json.RawMessage `gorm:"type:json"`
}

func (s authBackend) Get(t settings.AuthMethod) (auth.Auther, error) {
	var auther auth.Auther

	autherTb := &AutherTb{}
	err := s.db.Table("auther").First(autherTb).Error
	if err != nil {
		return nil, errors.ErrNotExist
	}
	switch t {
	case auth.MethodJSONAuth:
		reCaptcha := &auth.ReCaptcha{}
		err = json.Unmarshal(autherTb.JsonAuth, reCaptcha)
		if err != nil {
			return nil, err
		}
		auther = &auth.JSONAuth{
			ReCaptcha: reCaptcha,
		}
	case auth.MethodProxyAuth:
		proxyAuth := &auth.ProxyAuth{}
		err = json.Unmarshal(autherTb.ProxyAuth, proxyAuth)
		if err != nil {
			return nil, err
		}
		auther = proxyAuth
	case auth.MethodHookAuth:
		//auther = &auth.HookAuth{}
		//return nil, errors.ErrNotExist
		hookAuth := &auth.HookAuth{}
		err = json.Unmarshal(autherTb.HookAuth, hookAuth)
		if err != nil {
			return nil, err
		}
		auther = hookAuth
	case auth.MethodNoAuth:
		auther = &auth.NoAuth{}
	default:
		return nil, errors.ErrInvalidAuthMethod
	}
	return auther, err
	//return auther, get(s.db, "auther", auther)
}

func (s authBackend) Save(a auth.Auther) error {
	//switch a.(type) {
	//case auth.JSONAuth:
	//
	//case auth.ProxyAuth:
	//
	//case auth.HookAuth:
	//
	//case auth.NoAuth:
	//
	//default:
	//	return errors.ErrInvalidAuthMethod
	//}
	//return save(s.db, "auther", a)
	return nil
}
