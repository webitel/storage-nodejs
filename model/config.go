package model

import "net/http"

const (
	DEFAULT_LOCALE = "en"

	FILE_DRIVER_S3    = "s3"
	FILE_DRIVER_LOCAL = "local"

	DATABASE_DRIVER_POSTGRES = "postgres"
)

type LocalizationSettings struct {
	DefaultServerLocale *string
	DefaultClientLocale *string
	AvailableLocales    *string
}

func (s *LocalizationSettings) SetDefaults() {
	if s.DefaultServerLocale == nil {
		s.DefaultServerLocale = NewString(DEFAULT_LOCALE)
	}

	if s.DefaultClientLocale == nil {
		s.DefaultClientLocale = NewString(DEFAULT_LOCALE)
	}

	if s.AvailableLocales == nil {
		s.AvailableLocales = NewString("")
	}
}

type ServiceSettings struct {
	ListenAddress         *string
	ListenInternalAddress *string
	SessionCacheInMinutes *int
}

type SqlSettings struct {
	DriverName                  *string
	DataSource                  *string
	DataSourceReplicas          []string
	DataSourceSearchReplicas    []string
	MaxIdleConns                *int
	ConnMaxLifetimeMilliseconds *int
	MaxOpenConns                *int
	Trace                       bool
	AtRestEncryptKey            string
	QueryTimeout                *int
}

type NoSqlSettings struct {
	Host  *string
	Trace bool
}

type BrokerSettings struct {
	ConnectionString *string
}

type Config struct {
	NodeName               string
	DiscoverySettings      DiscoverySettings
	LocalizationSettings   LocalizationSettings
	ServiceSettings        ServiceSettings
	SqlSettings            SqlSettings
	NoSqlSettings          NoSqlSettings
	BrokerSettings         BrokerSettings
	MediaFileStoreSettings MediaFileStoreSettings
	ServerSettings         ServerSettings
}

type DiscoverySettings struct {
	Url string
}

type ServerSettings struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Network string `json:"network"`
}

type MediaFileStoreSettings struct {
	MaxSizeByte *int
	AllowMime   []string
	Directory   *string
	PathPattern *string
}

func (c *Config) IsValid() *AppError {

	if c.MediaFileStoreSettings.Directory == nil || len(*c.MediaFileStoreSettings.Directory) == 0 {
		return NewAppError("Config.IsValid", "model.config.is_valid.media_store_directory.app_error", nil, "", http.StatusInternalServerError)
	}
	return nil
}