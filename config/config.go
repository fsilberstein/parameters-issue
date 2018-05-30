package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/fsilberstein/parameters-issue/logger"
	"github.com/spf13/viper"
)

// All ENVs
var (
	Port                int
	ZipkinHost          string
	ElasticIndex        string
	ElasticSniff        bool
	ElasticHost         string
	ElasticResponseSize int
	ElasticDebug        bool
)

func init() {
	InitConfig()
}

// InitConfig : gets the service configuration
func InitConfig() {

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("ELASTIC_RESPONSE_SIZE", 10000)
	viper.SetDefault("ELASTIC_DEBUG", false)

	if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("ENVIRONMENT") == "DEV" {
		_, dirname, _, _ := runtime.Caller(0)
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(filepath.Dir(dirname))
		err := viper.ReadInConfig()
		if err != nil {
			logger.LogStdErr.Error(err)
		}
	} else {
		viper.AutomaticEnv()
	}

	// Assign env variables value to global variables
	Port = viper.GetInt("APP_PORT")
	// Elastic configuration
	ElasticIndex = viper.GetString("ELASTIC_INDEX")
	ElasticSniff = viper.GetBool("ELASTIC_SNIFF")
	ElasticHost = viper.GetString("ELASTIC_HOST")
	ElasticResponseSize = viper.GetInt("ELASTIC_RESPONSE_SIZE")
	ElasticDebug = viper.GetBool("ELASTIC_DEBUG")
}
