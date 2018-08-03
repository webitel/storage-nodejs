package model

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
	Host *string
}

type BrokerSettings struct {
	ConnectionString *string
}

type Config struct {
	LocalizationSettings   LocalizationSettings
	ServiceSettings        ServiceSettings
	SqlSettings            SqlSettings
	NoSqlSettings          NoSqlSettings
	BrokerSettings         BrokerSettings
	MediaFileStoreSettings MediaFileStoreSettings
}

type MediaFileStoreSettings struct {
	MaxMb       *int
	Directory   *string
	PathPattern *string
}
