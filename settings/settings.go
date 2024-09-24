package settings

import (
	"crypto/rand"
	"log"
	"strings"
	"time"

	"github.com/filebrowser/filebrowser/v2/rules"
)

const DefaultUsersHomeBasePath = "/users"

// AuthMethod describes an authentication method.
type AuthMethod string

// Settings contain the main settings of the application.
type Settings struct {
	Key              []byte              `json:"key"`
	Signup           bool                `json:"signup"`
	CreateUserDir    bool                `json:"createUserDir"`
	UserHomeBasePath string              `json:"userHomeBasePath"`
	Defaults         UserDefaults        `json:"defaults"`
	AuthMethod       AuthMethod          `json:"authMethod"`
	Gpfs             Gpfs                `json:"gpfs"`
	Branding         Branding            `json:"branding"`
	Tus              Tus                 `json:"tus"`
	Commands         map[string][]string `json:"commands"`
	Shell            []string            `json:"shell"`
	Rules            []rules.Rule        `json:"rules"`
}

// GetRules implements rules.Provider.
func (s *Settings) GetRules() []rules.Rule {
	return s.Rules
}

// Server specific settings.
type Server struct {
	Id                    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Root                  string `json:"root" gorm:"column:root;type:varchar(255);not null"`
	BaseURL               string `json:"baseURL" gorm:"column:base_url;type:varchar(255)"`
	Socket                string `json:"socket" gorm:"column:socket;type:varchar(255)"`
	TLSKey                string `json:"tlsKey" gorm:"column:tls_key;type:varchar(255)"`
	TLSCert               string `json:"tlsCert" gorm:"column:tls_cert;type:varchar(255)"`
	Port                  string `json:"port" gorm:"column:port;type:varchar(10);not null"`
	Address               string `json:"address" gorm:"column:address;type:varchar(255);not null"`
	Log                   string `json:"log" gorm:"column:log;type:varchar(255);not null"`
	EnableThumbnails      bool   `json:"enableThumbnails" gorm:"column:enable_thumbnails;type:boolean;default:false"`
	ResizePreview         bool   `json:"resizePreview" gorm:"column:resize_preview;type:boolean;default:false"`
	EnableExec            bool   `json:"enableExec" gorm:"column:enable_exec;type:boolean;default:false"`
	TypeDetectionByHeader bool   `json:"typeDetectionByHeader" gorm:"column:type_detection_by_header;type:boolean;default:false"`
	AuthHook              string `json:"authHook" gorm:"column:auth_hook;type:varchar(255)"`
	TokenExpirationTime   string `json:"tokenExpirationTime" gorm:"column:token_expiration_time;type:varchar(255)"`
}

func (Server) TableName() string {
	return "server"
}

// Clean cleans any variables that might need cleaning.
func (s *Server) Clean() {
	s.BaseURL = strings.TrimSuffix(s.BaseURL, "/")
}

func (s *Server) GetTokenExpirationTime(fallback time.Duration) time.Duration {
	if s.TokenExpirationTime == "" {
		return fallback
	}

	duration, err := time.ParseDuration(s.TokenExpirationTime)
	if err != nil {
		log.Printf("[WARN] Failed to parse tokenExpirationTime: %v", err)
		return fallback
	}
	return duration
}

// GenerateKey generates a key of 512 bits.
func GenerateKey() ([]byte, error) {
	b := make([]byte, 64) //nolint:gomnd
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
