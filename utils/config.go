package utils

import (
	"github.com/webitel/storage/model"
	"os"
	"path/filepath"
)

var (
	commonBaseSearchPaths = []string{
		".",
		"..",
		"../..",
		"../../..",
	}
)

func LoadConfig(fileName string) (*model.Config, string, map[string]interface{}, *model.AppError) {
	var envConfig = make(map[string]interface{})
	dbDatasource := "postgres://webitel:webitel@localhost:5432/webitel?fallback_application_name=storage&sslmode=disable&connect_timeout=10&search_path=storage"
	dbDriverName := "postgres"
	maxIdleConns := 5
	maxOpenConns := 5
	connMaxLifetimeMilliseconds := 3600000

	return &model.Config{
		LocalizationSettings: model.LocalizationSettings{
			DefaultClientLocale: model.NewString(model.DEFAULT_LOCALE),
			DefaultServerLocale: model.NewString(model.DEFAULT_LOCALE),
			AvailableLocales:    model.NewString(model.DEFAULT_LOCALE),
		},
		ServiceSettings: model.ServiceSettings{
			ListenAddress:         model.NewString(":10023"),
			ListenInternalAddress: model.NewString(":10021"),
		},
		MediaFileStoreSettings: model.MediaFileStoreSettings{
			MaxSizeByte: model.NewInt(50 * 1000000),
			Directory:   model.NewString("/tmp/media_storage"),
			PathPattern: model.NewString("$DOMAIN"),
			AllowMime:   []string{"video/mp4", "audio/mp3", "audio/wav"},
		},
		SqlSettings: model.SqlSettings{
			DriverName: &dbDriverName,
			DataSource: &dbDatasource,
			//DataSourceReplicas:          []string{"postgres://webitel:webitel@10.10.10.25:5432/webitel?sslmode=disable&connect_timeout=10&search_path=storage"},
			MaxIdleConns:                &maxIdleConns,
			MaxOpenConns:                &maxOpenConns,
			ConnMaxLifetimeMilliseconds: &connMaxLifetimeMilliseconds,
			Trace:                       false,
		},
		NoSqlSettings: model.NoSqlSettings{
			Host:  model.NewString("http://10.10.10.200:9200"),
			Trace: true,
		},
		BrokerSettings: model.BrokerSettings{
			ConnectionString: model.NewString("amqp://webitel:webitel@cloud-ua2.webitel.com:5672?heartbeat=0"),
		},
		DiscoverySettings: model.DiscoverySettings{
			Url: "192.168.177.199:8500",
		},
		ServerSettings: model.ServerSettings{
			Address: "",
			Port:    8041,
			Network: "tcp",
		},
	}, "", envConfig, nil
}

func FindPath(path string, baseSearchPaths []string, filter func(os.FileInfo) bool) string {
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path
		}

		return ""
	}

	searchPaths := []string{}
	for _, baseSearchPath := range baseSearchPaths {
		searchPaths = append(searchPaths, baseSearchPath)
	}

	// Additionally attempt to search relative to the location of the running binary.
	var binaryDir string
	if exe, err := os.Executable(); err == nil {
		if exe, err = filepath.EvalSymlinks(exe); err == nil {
			if exe, err = filepath.Abs(exe); err == nil {
				binaryDir = filepath.Dir(exe)
			}
		}
	}
	if binaryDir != "" {
		for _, baseSearchPath := range baseSearchPaths {
			searchPaths = append(
				searchPaths,
				filepath.Join(binaryDir, baseSearchPath),
			)
		}
	}

	for _, parent := range searchPaths {
		found, err := filepath.Abs(filepath.Join(parent, path))
		if err != nil {
			continue
		} else if fileInfo, err := os.Stat(found); err == nil {
			if filter != nil {
				if filter(fileInfo) {
					return found
				}
			} else {
				return found
			}
		}
	}

	return ""
}

// FindDir looks for the given directory in nearby ancestors relative to the current working
// directory as well as the directory of the executable, falling back to `./` if not found.
func FindDir(dir string) (string, bool) {
	found := FindPath(dir, commonBaseSearchPaths, func(fileInfo os.FileInfo) bool {
		return fileInfo.IsDir()
	})
	if found == "" {
		return "./", false
	}

	return found, true
}
