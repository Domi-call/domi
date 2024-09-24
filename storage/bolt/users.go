package bolt

import (
	"errors"
	"github.com/filebrowser/filebrowser/v2/storage/mysql"
	"gorm.io/gorm"

	fbErrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/users"
)

type usersBackend struct {
	//db *storm.DB
	db *gorm.DB
}

func (st usersBackend) GetBy(i interface{}) (user *users.User, err error) {
	userTb := &users.UserTb{}

	var arg string
	switch i.(type) {
	case uint:
		//arg = "ID"
		arg = "id"
	case string:
		//arg = "Username"
		arg = "username"
	default:
		return nil, fbErrors.ErrInvalidDataType
	}
	err = mysql.DB.Model(userTb).Where(arg+" = ?", i).First(userTb).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fbErrors.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	user, err = users.UserTb2User(userTb)
	return user, err

	//err = st.db.One(arg, i, user)
	//if err != nil {
	//	if errors.Is(err, storm.ErrNotFound) {
	//		return nil, fbErrors.ErrNotExist
	//	}
	//	return nil, err
	//}

	//return
}

func (st usersBackend) Gets() ([]*users.User, error) {
	var allUserTbs []*users.UserTb
	err := mysql.DB.Model(&users.UserTb{}).Find(&allUserTbs).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fbErrors.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	var allUsers []*users.User
	for _, userTb := range allUserTbs {
		user, err := users.UserTb2User(userTb)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, user)
	}
	return allUsers, err

	//var allUsers []*users.User
	//err := st.db.All(&allUsers)
	//if errors.Is(err, storm.ErrNotFound) {
	//	return nil, fbErrors.ErrNotExist
	//}
	//
	//if err != nil {
	//	return allUsers, err
	//}
	//
	//return allUsers, err
}

func (st usersBackend) Update(user *users.User, fields ...string) error {
	userTb, err := users.User2UserTb(user)
	if err != nil {
		return err
	}
	tx := mysql.DB.Begin()
	if len(fields) == 0 {
		err := tx.Model(userTb).Updates(userTb).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, field := range fields {
		err := tx.Model(userTb).Update(field, userTb).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil

	//if len(fields) == 0 {
	//	return st.Save(user)
	//}
	//
	//for _, field := range fields {
	//	userField := reflect.ValueOf(user).Elem().FieldByName(field)
	//	if !userField.IsValid() {
	//		return fmt.Errorf("invalid field: %s", field)
	//	}
	//	val := userField.Interface()
	//	if err := st.db.UpdateField(user, field, val); err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (st usersBackend) Save(user *users.User) error {
	userTb, err := users.User2UserTb(user)
	if err != nil {
		return err
	}
	//检查username是否存在
	var userTbCheck users.UserTb
	err = mysql.DB.Model(&users.UserTb{}).Where("username = ?", userTb.Username).First(&userTbCheck).Error
	if err == nil {
		return fbErrors.ErrExist
	}
	err = mysql.DB.Create(userTb).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
	//err := st.db.Save(user)
	//if errors.Is(err, storm.ErrAlreadyExists) {
	//	return fbErrors.ErrExist
	//}
	//return err
}

func (st usersBackend) DeleteByID(id uint) error {
	user, err := st.GetBy(id)
	if err != nil {
		return err
	}
	return st.db.Model(&users.UserTb{}).Delete(user).Error
	//return st.db.DeleteStruct(&users.User{ID: id})
}

func (st usersBackend) DeleteByUsername(username string) error {
	user, err := st.GetBy(username)
	if err != nil {
		return err
	}
	return st.db.Model(&users.UserTb{}).Delete(user).Error
	//return st.db.DeleteStruct(user)
}
