package users

import (
	"encoding/json"
	"time"
)

type UserTb struct {
	Id           int             `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string          `json:"username" gorm:"column:username;type:varchar(255);not null"`
	Password     string          `json:"password" gorm:"column:password;type:varchar(255);not null"`
	Scope        string          `json:"scope" gorm:"column:scope;type:varchar(255);not null;default:'/'"`
	Locale       string          `json:"locale" gorm:"column:locale;type:varchar(10);not null;default:''"`
	LockPassword bool            `json:"lock_password" gorm:"column:lock_password;type:tinyint(1);not null;default:0"`
	ViewMode     string          `json:"view_mode" gorm:"column:view_mode;type:varchar(50);not null;default:''"`
	SingleClick  bool            `json:"single_click" gorm:"column:single_click;type:tinyint(1);not null;default:0"`
	PermAdmin    bool            `json:"perm_admin" gorm:"column:perm_admin;type:tinyint(1);not null;default:0"`
	PermExecute  bool            `json:"perm_execute" gorm:"column:perm_execute;type:tinyint(1);not null;default:0"`
	PermCreate   bool            `json:"perm_create" gorm:"column:perm_create;type:tinyint(1);not null;default:0"`
	PermRename   bool            `json:"perm_rename" gorm:"column:perm_rename;type:tinyint(1);not null;default:0"`
	PermModify   bool            `json:"perm_modify" gorm:"column:perm_modify;type:tinyint(1);not null;default:0"`
	PermDelete   bool            `json:"perm_delete" gorm:"column:perm_delete;type:tinyint(1);not null;default:0"`
	PermShare    bool            `json:"perm_share" gorm:"column:perm_share;type:tinyint(1);not null;default:0"`
	PermDownload bool            `json:"perm_download" gorm:"column:perm_download;type:tinyint(1);not null;default:0"`
	SortingBy    string          `json:"sorting_by" gorm:"column:sorting_by;type:varchar(50);not null;default:''"`
	SortingAsc   bool            `json:"sorting_asc" gorm:"column:sorting_asc;type:tinyint(1);not null;default:0"`
	HideDotfiles bool            `json:"hide_dotfiles" gorm:"column:hide_dotfiles;type:tinyint(1);not null;default:0"`
	DateFormat   bool            `json:"date_format" gorm:"column:date_format;type:tinyint(1);not null;default:0"`
	Commands     json.RawMessage `json:"commands" gorm:"column:commands;type:json;not null"`
	Rules        json.RawMessage `json:"rules" gorm:"column:rules;type:json;not null"`
	CreateTime   time.Time       `json:"create_time" gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime   time.Time       `json:"update_time" gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (UserTb) TableName() string {
	return "user"
}

func UserTb2User(userTb *UserTb) (*User, error) {
	user := &User{}
	user.ID = uint(userTb.Id)
	user.Username = userTb.Username
	user.Password = userTb.Password
	user.Scope = userTb.Scope
	user.Locale = userTb.Locale
	user.LockPassword = userTb.LockPassword
	user.ViewMode = ViewMode(userTb.ViewMode)
	user.SingleClick = userTb.SingleClick
	user.Perm.Admin = userTb.PermAdmin
	user.Perm.Execute = userTb.PermExecute
	user.Perm.Create = userTb.PermCreate
	user.Perm.Rename = userTb.PermRename
	user.Perm.Modify = userTb.PermModify
	user.Perm.Delete = userTb.PermDelete
	user.Perm.Share = userTb.PermShare
	user.Perm.Download = userTb.PermDownload
	user.Sorting.By = userTb.SortingBy
	user.Sorting.Asc = userTb.SortingAsc
	user.HideDotfiles = userTb.HideDotfiles
	user.DateFormat = userTb.DateFormat
	err := json.Unmarshal(userTb.Commands, &user.Commands)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(userTb.Rules, &user.Rules)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func User2UserTb(user *User) (*UserTb, error) {
	userTb := &UserTb{}
	userTb.Id = int(user.ID)
	userTb.Username = user.Username
	userTb.Password = user.Password
	userTb.Scope = user.Scope
	userTb.Locale = user.Locale
	userTb.LockPassword = user.LockPassword
	userTb.ViewMode = string(user.ViewMode)
	userTb.SingleClick = user.SingleClick
	userTb.PermAdmin = user.Perm.Admin
	userTb.PermExecute = user.Perm.Execute
	userTb.PermCreate = user.Perm.Create
	userTb.PermRename = user.Perm.Rename
	userTb.PermModify = user.Perm.Modify
	userTb.PermDelete = user.Perm.Delete
	userTb.PermShare = user.Perm.Share
	userTb.PermDownload = user.Perm.Download
	userTb.SortingBy = string(user.Sorting.By)
	userTb.SortingAsc = user.Sorting.Asc
	userTb.HideDotfiles = user.HideDotfiles
	userTb.DateFormat = user.DateFormat
	commands, err := json.Marshal(user.Commands)
	if err != nil {
		return nil, err
	}
	userTb.Commands = commands
	rules, err := json.Marshal(user.Rules)
	if err != nil {
		return nil, err
	}
	userTb.Rules = rules
	userTb.CreateTime = time.Now()
	userTb.UpdateTime = time.Now()
	return userTb, nil
}
