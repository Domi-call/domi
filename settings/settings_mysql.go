package settings

import (
	"encoding/json"
	"github.com/filebrowser/filebrowser/v2/users"
	"time"
)

type SettingsTb struct {
	Id                            int             `json:"id" gorm:"primaryKey;autoIncrement"`
	Key                           string          `json:"key" gorm:"column:key;type:text;not null"`
	Signup                        bool            `json:"signup" gorm:"column:signup;type:tinyint(1);not null;default:0"`
	CreateUserDir                 bool            `json:"create_user_dir" gorm:"column:create_user_dir;type:tinyint(1);not null;default:0"`
	UserHomeBasePath              string          `json:"user_home_base_path" gorm:"column:user_home_base_path;type:varchar(255);not null;default:'/users'"`
	DefaultsScope                 string          `json:"defaults_scope" gorm:"column:defaults_scope;type:varchar(255);not null;default:'.'"`
	DefaultsLocale                string          `json:"defaults_locale" gorm:"column:defaults_locale;type:varchar(10);not null;default:'en'"`
	DefaultsViewMode              string          `json:"defaults_view_mode" gorm:"column:defaults_view_mode;type:varchar(50);not null;default:''"`
	DefaultsSingleClick           bool            `json:"defaults_single_click" gorm:"column:defaults_single_click;type:tinyint(1);not null;default:0"`
	DefaultsSortingBy             string          `json:"defaults_sorting_by" gorm:"column:defaults_sorting_by;type:varchar(50);not null;default:''"`
	DefaultsSortingAsc            bool            `json:"defaults_sorting_asc" gorm:"column:defaults_sorting_asc;type:tinyint(1);not null;default:0"`
	DefaultsPermAdmin             bool            `json:"defaults_perm_admin" gorm:"column:defaults_perm_admin;type:tinyint(1);not null;default:0"`
	DefaultsPermExecute           bool            `json:"defaults_perm_execute" gorm:"column:defaults_perm_execute;type:tinyint(1);not null;default:1"`
	DefaultsPermCreate            bool            `json:"defaults_perm_create" gorm:"column:defaults_perm_create;type:tinyint(1);not null;default:1"`
	DefaultsPermRename            bool            `json:"defaults_perm_rename" gorm:"column:defaults_perm_rename;type:tinyint(1);not null;default:1"`
	DefaultsPermModify            bool            `json:"defaults_perm_modify" gorm:"column:defaults_perm_modify;type:tinyint(1);not null;default:1"`
	DefaultsPermDelete            bool            `json:"defaults_perm_delete" gorm:"column:defaults_perm_delete;type:tinyint(1);not null;default:1"`
	DefaultsPermShare             bool            `json:"defaults_perm_share" gorm:"column:defaults_perm_share;type:tinyint(1);not null;default:1"`
	DefaultsPermDownload          bool            `json:"defaults_perm_download" gorm:"column:defaults_perm_download;type:tinyint(1);not null;default:1"`
	DefaultsCommands              json.RawMessage `json:"defaults_commands" gorm:"column:defaults_commands;type:json;not null"`
	DefaultsHideDotfiles          bool            `json:"defaults_hide_dotfiles" gorm:"column:defaults_hide_dotfiles;type:tinyint(1);not null;default:0"`
	DefaultsDateFormat            bool            `json:"defaults_date_format" gorm:"column:defaults_date_format;type:tinyint(1);not null;default:0"`
	AuthMethod                    json.RawMessage `json:"auth_method" gorm:"column:auth_method;type:varchar(50);not null;default:'json'"`
	GpfsAPI                       string          `json:"gpfs_api" gorm:"column:gpfs_api;type:varchar(255);not null;default:'https://xxx.xxx.xxx.xxx:443'"`
	GpfsApiUsername               string          `json:"gpfs_api_username" gorm:"column:gpfs_api_username;type:varchar(255);not null;default:''"`
	GpfsApiPassword               string          `json:"gpfs_api_password" gorm:"column:gpfs_api_password;type:varchar(255);not null;default:''"`
	GpfsFileSystemName            string          `json:"gpfs_filesystem_name" gorm:"column:gpfs_filesystem_name;type:varchar(255);not null;default:'gpfsdata'"`
	GpfsRootPath                  string          `json:"gpfs_root_path" gorm:"column:gpfs_root_path;type:varchar(255);not null;default:''"`
	GpfsQuotaLimit                int             `json:"gpfs_quota_limit" gorm:"column:gpfs_quota_limit;type:int;not null;default:50"`
	GpfsQuotaMax                  int             `json:"gpfs_quota_max" gorm:"column:gpfs_quota_max;type:int;not null;default:70"`
	GpfsServer                    string          `json:"gpfs_server" gorm:"column:gpfs_server;type:varchar(50);not null;default:''"`
	GpfsSysuser                   string          `json:"gpfs_sysuser" gorm:"column:gpfs_sysuser;type:varchar(255);not null;default:''"`
	GpfsPrivateKey                string          `json:"gpfs_private_key" gorm:"column:gpfs_private_key;type:text;not null"`
	BrandingName                  string          `json:"branding_name" gorm:"column:branding_name;type:varchar(255);not null;default:''"`
	BrandingDisableExternal       bool            `json:"branding_disable_external" gorm:"column:branding_disable_external;type:tinyint(1);not null;default:0"`
	BrandingDisableUsedPercentage bool            `json:"branding_disable_used_percentage" gorm:"column:branding_disable_used_percentage;type:tinyint(1);not null;default:0"`
	BrandingFiles                 string          `json:"branding_files" gorm:"column:branding_files;type:text;not null"`
	BrandingTheme                 string          `json:"branding_theme" gorm:"column:branding_theme;type:varchar(255);not null;default:''"`
	BrandingColor                 string          `json:"branding_color" gorm:"column:branding_color;type:varchar(255);not null;default:''"`
	TusChunkSize                  int             `json:"tus_chunk_size" gorm:"column:tus_chunk_size;type:int;not null;default:10485760"`
	TusRetryCount                 int             `json:"tus_retry_count" gorm:"column:tus_retry_count;type:int;not null;default:5"`
	CommandsAfterCopy             json.RawMessage `json:"commands_after_copy" gorm:"column:commands_after_copy;type:json;not null"`
	CommandsAfterDelete           json.RawMessage `json:"commands_after_delete" gorm:"column:commands_after_delete;type:json;not null"`
	CommandsAfterRename           json.RawMessage `json:"commands_after_rename" gorm:"column:commands_after_rename;type:json;not null"`
	CommandsAfterSave             json.RawMessage `json:"commands_after_save" gorm:"column:commands_after_save;type:json;not null"`
	CommandsAfterUpload           json.RawMessage `json:"commands_after_upload" gorm:"column:commands_after_upload;type:json;not null"`
	CommandsBeforeCopy            json.RawMessage `json:"commands_before_copy" gorm:"column:commands_before_copy;type:json;not null"`
	CommandsBeforeDelete          json.RawMessage `json:"commands_before_delete" gorm:"column:commands_before_delete;type:json;not null"`
	CommandsBeforeRename          json.RawMessage `json:"commands_before_rename" gorm:"column:commands_before_rename;type:json;not null"`
	CommandsBeforeSave            json.RawMessage `json:"commands_before_save" gorm:"column:commands_before_save;type:json;not null"`
	CommandsBeforeUpload          json.RawMessage `json:"commands_before_upload" gorm:"column:commands_before_upload;type:json;not null"`
	Shell                         json.RawMessage `json:"shell" gorm:"column:shell;type:json;not null"`
	Rules                         json.RawMessage `json:"rules" gorm:"column:rules;type:json;not null"`
	CreateTime                    time.Time       `json:"create_time" gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime                    time.Time       `json:"update_time" gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (SettingsTb) TableName() string {
	return "setting"
}

func SettingsTb2Settings(tb *SettingsTb) (*Settings, error) {
	set := &Settings{}
	set.Key = []byte(tb.Key)
	set.Signup = tb.Signup
	set.CreateUserDir = tb.CreateUserDir
	set.UserHomeBasePath = tb.UserHomeBasePath
	set.Defaults.Scope = tb.DefaultsScope
	set.Defaults.Locale = tb.DefaultsLocale
	set.Defaults.ViewMode = users.ViewMode(tb.DefaultsViewMode)
	set.Defaults.SingleClick = tb.DefaultsSingleClick
	set.Defaults.Sorting.By = tb.DefaultsSortingBy
	set.Defaults.Sorting.Asc = tb.DefaultsSortingAsc
	set.Defaults.Perm.Admin = tb.DefaultsPermAdmin
	set.Defaults.Perm.Execute = tb.DefaultsPermExecute
	set.Defaults.Perm.Create = tb.DefaultsPermCreate
	set.Defaults.Perm.Rename = tb.DefaultsPermRename
	set.Defaults.Perm.Modify = tb.DefaultsPermModify
	set.Defaults.Perm.Delete = tb.DefaultsPermDelete
	set.Defaults.Perm.Share = tb.DefaultsPermShare
	set.Defaults.Perm.Download = tb.DefaultsPermDownload
	err := json.Unmarshal(tb.DefaultsCommands, &set.Defaults.Commands)
	if err != nil {
		return nil, err
	}
	set.Defaults.HideDotfiles = tb.DefaultsHideDotfiles
	set.Defaults.DateFormat = tb.DefaultsDateFormat
	set.AuthMethod = AuthMethod(tb.AuthMethod)
	set.Gpfs.Api = tb.GpfsAPI
	set.Gpfs.Username = tb.GpfsApiUsername
	set.Gpfs.Password = tb.GpfsApiPassword
	set.Gpfs.FileSystemName = tb.GpfsFileSystemName
	set.Gpfs.RootPath = tb.GpfsRootPath
	set.Gpfs.QuotaLimmit = uint16(tb.GpfsQuotaLimit)
	set.Gpfs.QuotaMax = uint16(tb.GpfsQuotaMax)
	set.Gpfs.Server = tb.GpfsServer
	set.Gpfs.Sysuser = tb.GpfsSysuser
	set.Gpfs.PrivateKey = tb.GpfsPrivateKey
	set.Branding.Name = tb.BrandingName
	set.Branding.DisableExternal = tb.BrandingDisableExternal
	set.Branding.DisableUsedPercentage = tb.BrandingDisableUsedPercentage
	set.Branding.Files = tb.BrandingFiles
	set.Branding.Theme = tb.BrandingTheme
	set.Branding.Color = tb.BrandingColor
	set.Tus.ChunkSize = uint64(tb.TusChunkSize)
	set.Tus.RetryCount = uint16(tb.TusRetryCount)
	set.Commands = make(map[string][]string)
	set.Commands["after_copy"] = []string{}
	set.Commands["after_delete"] = []string{}
	set.Commands["after_rename"] = []string{}
	set.Commands["after_save"] = []string{}
	set.Commands["after_upload"] = []string{}
	set.Commands["before_copy"] = []string{}
	set.Commands["before_delete"] = []string{}
	set.Commands["before_rename"] = []string{}
	set.Commands["before_save"] = []string{}
	set.Commands["before_upload"] = []string{}
	afterCopy := []string{}
	err = json.Unmarshal(tb.CommandsAfterCopy, &afterCopy)
	if err != nil {
		return nil, err
	}
	set.Commands["after_copy"] = afterCopy
	afterDelete := []string{}
	err = json.Unmarshal(tb.CommandsAfterDelete, &afterDelete)
	if err != nil {
		return nil, err
	}
	set.Commands["after_delete"] = afterDelete
	afterRename := []string{}
	err = json.Unmarshal(tb.CommandsAfterRename, &afterRename)
	if err != nil {
		return nil, err
	}
	set.Commands["after_rename"] = afterRename
	afterSave := []string{}
	err = json.Unmarshal(tb.CommandsAfterSave, &afterSave)
	if err != nil {
		return nil, err
	}
	set.Commands["after_save"] = afterSave
	afterUpload := []string{}
	err = json.Unmarshal(tb.CommandsAfterUpload, &afterUpload)
	if err != nil {
		return nil, err
	}
	set.Commands["after_upload"] = afterUpload
	beforeCopy := []string{}
	err = json.Unmarshal(tb.CommandsBeforeCopy, &beforeCopy)
	if err != nil {
		return nil, err
	}
	set.Commands["before_copy"] = beforeCopy
	beforeDelete := []string{}
	err = json.Unmarshal(tb.CommandsBeforeDelete, &beforeDelete)
	if err != nil {
		return nil, err
	}
	set.Commands["before_delete"] = beforeDelete
	beforeRename := []string{}
	err = json.Unmarshal(tb.CommandsBeforeRename, &beforeRename)
	if err != nil {
		return nil, err
	}
	set.Commands["before_rename"] = beforeRename
	beforeSave := []string{}
	err = json.Unmarshal(tb.CommandsBeforeSave, &beforeSave)
	if err != nil {
		return nil, err
	}
	set.Commands["before_save"] = beforeSave
	beforeUpload := []string{}
	err = json.Unmarshal(tb.CommandsBeforeUpload, &beforeUpload)
	if err != nil {
		return nil, err
	}
	set.Commands["before_upload"] = beforeUpload
	err = json.Unmarshal(tb.Shell, &set.Shell)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tb.Rules, &set.Rules)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func Settings2SettingsTb(set *Settings) (*SettingsTb, error) {
	setTb := &SettingsTb{}
	setTb.Key = string(set.Key)
	setTb.Signup = set.Signup
	setTb.CreateUserDir = set.CreateUserDir
	setTb.UserHomeBasePath = set.UserHomeBasePath
	setTb.DefaultsScope = set.Defaults.Scope
	setTb.DefaultsLocale = set.Defaults.Locale
	setTb.DefaultsViewMode = string(set.Defaults.ViewMode)
	setTb.DefaultsSingleClick = set.Defaults.SingleClick
	setTb.DefaultsSortingBy = set.Defaults.Sorting.By
	setTb.DefaultsSortingAsc = set.Defaults.Sorting.Asc
	setTb.DefaultsPermAdmin = set.Defaults.Perm.Admin
	setTb.DefaultsPermExecute = set.Defaults.Perm.Execute
	setTb.DefaultsPermCreate = set.Defaults.Perm.Create
	setTb.DefaultsPermRename = set.Defaults.Perm.Rename
	setTb.DefaultsPermModify = set.Defaults.Perm.Modify
	setTb.DefaultsPermDelete = set.Defaults.Perm.Delete
	setTb.DefaultsPermShare = set.Defaults.Perm.Share
	setTb.DefaultsPermDownload = set.Defaults.Perm.Download
	commands, err := json.Marshal(set.Defaults.Commands)
	if err != nil {
		return nil, err
	}
	setTb.DefaultsCommands = commands
	setTb.DefaultsHideDotfiles = set.Defaults.HideDotfiles
	setTb.DefaultsDateFormat = set.Defaults.DateFormat
	setTb.AuthMethod = json.RawMessage(string(set.AuthMethod))
	setTb.GpfsAPI = set.Gpfs.Api
	setTb.GpfsApiUsername = set.Gpfs.Username
	setTb.GpfsApiPassword = set.Gpfs.Password
	setTb.GpfsFileSystemName = set.Gpfs.FileSystemName
	setTb.GpfsRootPath = set.Gpfs.RootPath
	setTb.GpfsQuotaLimit = int(set.Gpfs.QuotaLimmit)
	setTb.GpfsQuotaMax = int(set.Gpfs.QuotaMax)
	setTb.GpfsServer = set.Gpfs.Server
	setTb.GpfsSysuser = set.Gpfs.Sysuser
	setTb.GpfsPrivateKey = set.Gpfs.PrivateKey
	setTb.BrandingName = set.Branding.Name
	setTb.BrandingDisableExternal = set.Branding.DisableExternal
	setTb.BrandingDisableUsedPercentage = set.Branding.DisableUsedPercentage
	setTb.BrandingFiles = set.Branding.Files
	setTb.BrandingTheme = set.Branding.Theme
	setTb.BrandingColor = set.Branding.Color
	setTb.TusChunkSize = int(set.Tus.ChunkSize)
	setTb.TusRetryCount = int(set.Tus.RetryCount)
	afterCopy, err := json.Marshal(set.Commands["after_copy"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsAfterCopy = afterCopy
	afterDelete, err := json.Marshal(set.Commands["after_delete"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsAfterDelete = afterDelete
	afterRename, err := json.Marshal(set.Commands["after_rename"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsAfterRename = afterRename
	afterSave, err := json.Marshal(set.Commands["after_save"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsAfterSave = afterSave
	afterUpload, err := json.Marshal(set.Commands["after_upload"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsAfterUpload = afterUpload
	beforeCopy, err := json.Marshal(set.Commands["before_copy"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsBeforeCopy = beforeCopy
	beforeDelete, err := json.Marshal(set.Commands["before_delete"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsBeforeDelete = beforeDelete
	beforeRename, err := json.Marshal(set.Commands["before_rename"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsBeforeRename = beforeRename
	beforeSave, err := json.Marshal(set.Commands["before_save"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsBeforeSave = beforeSave
	beforeUpload, err := json.Marshal(set.Commands["before_upload"])
	if err != nil {
		return nil, err
	}
	setTb.CommandsBeforeUpload = beforeUpload
	shell, err := json.Marshal(set.Shell)
	if err != nil {
		return nil, err
	}
	setTb.Shell = shell
	rules, err := json.Marshal(set.Rules)

	if err != nil {
		return nil, err
	}
	setTb.Rules = rules
	return setTb, nil
}
