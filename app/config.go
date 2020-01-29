package app

import (
	"flag"
	"fmt"
	"github.com/webitel/storage/model"
)

var (
	consulHost            = flag.String("consul", "consul:8500", "Host to consul")
	dataSource            = flag.String("data_source", "postgres://opensips:webitel@postgres:5432/webitel?fallback_application_name=storage&sslmode=disable&connect_timeout=10&search_path=storage", "Data source")
	amqpSource            = flag.String("amqp", "amqp://webitel:webitel@rabbit:5672?heartbeat=10", "AMQP connection")
	elasticSource         = flag.String("elastic", "http://10.10.10.200:9200", "Elastic endpoint")
	grpcServerPort        = flag.Int("grpc_port", 0, "GRPC port")
	dev                   = flag.Bool("dev", false, "enable dev mode")
	internalServerAddress = flag.String("internal_address", ":10021", "Internal server address")
	publicServerAddress   = flag.String("public_address", ":10023", "Public server address")
)

func loadConfig(fileName string) (*model.Config, *model.AppError) {
	flag.Parse()

	return &model.Config{
		NodeName: fmt.Sprintf("%s-%s", model.APP_SERVICE_NAME, model.NewId()),
		IsDev:    *dev,
		LocalizationSettings: model.LocalizationSettings{
			DefaultClientLocale: model.NewString(model.DEFAULT_LOCALE),
			DefaultServerLocale: model.NewString(model.DEFAULT_LOCALE),
			AvailableLocales:    model.NewString(model.DEFAULT_LOCALE),
		},
		ServiceSettings: model.ServiceSettings{
			ListenAddress:         publicServerAddress,
			ListenInternalAddress: internalServerAddress,
		},
		MediaFileStoreSettings: model.MediaFileStoreSettings{
			MaxSizeByte: model.NewInt(100 * 1000000),
			Directory:   model.NewString("/tmp/media_storage"),
			PathPattern: model.NewString("$DOMAIN/$Y"),
			AllowMime:   []string{"video/mp4", "audio/mp3", "audio/wav", "audio/mpeg", "video/x-matroska", "video/mpeg"},
		},
		SqlSettings: model.SqlSettings{
			DriverName:                  model.NewString("postgres"),
			DataSource:                  dataSource,
			MaxIdleConns:                model.NewInt(5),
			MaxOpenConns:                model.NewInt(5),
			ConnMaxLifetimeMilliseconds: model.NewInt(3600000),
			Trace:                       false,
		},
		NoSqlSettings: model.NoSqlSettings{
			Host:  elasticSource,
			Trace: true,
		},
		BrokerSettings: model.BrokerSettings{
			ConnectionString: amqpSource,
		},
		DiscoverySettings: model.DiscoverySettings{
			Url: *consulHost,
		},
		ServerSettings: model.ServerSettings{
			Address: "",
			Port:    *grpcServerPort,
			Network: "tcp",
		},
	}, nil
}

func (a *App) Config() *model.Config {
	if cfg := a.config.Load(); cfg != nil {
		return cfg.(*model.Config)
	}
	return &model.Config{}
}

func (a *App) LoadConfig(configFile string) *model.AppError {
	cfg, err := loadConfig(configFile)
	if err != nil {
		return err
	}

	if err = cfg.IsValid(); err != nil {
		return err
	}
	a.configFile = configFile
	a.id = &cfg.NodeName

	a.config.Store(cfg)
	return nil
}
