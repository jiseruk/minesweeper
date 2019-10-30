package config

import (
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
	"github.com/spf13/viper"
)

// Config is global object that holds all application level variables.
var Config appConfig

type appConfig struct {
	// the shared DB ORM object
	DB *gorm.DB
	// the error thrown be GORM when using DB ORM object
	DBErr error
	// the server port. Defaults to 8080
	ServerPort int `mapstructure:"server_port"`
	// the data source name (DSN) for connecting to the database. required.
	DSN string `mapstructure:"dsn"`

	Dialect string `mapstructure:"dialect"`
}

func init() {
	if err := LoadConfig("./config", "../config"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}

	Config.DB, Config.DBErr = gorm.Open(Config.Dialect, Config.DSN)
	if Config.DBErr != nil {
		panic(Config.DBErr)
	}

	Config.DB.AutoMigrate(&models.Game{}) // This is needed for generation of schema for postgres image.

	fmt.Println(fmt.Sprintf("Successfully connected to :%v", Config.DSN))

	//BoardService = services.NewBoardService(&daos.BoardDAOImpl{})
}

func LoadConfig(configPaths ...string) error {
	v := viper.New()
	if os.Getenv("GO_ENVIRONMENT") == "local" {
		v.SetConfigName("local")
	} else {
		v.SetConfigName("now")
	}
	v.SetConfigType("yaml")
	v.SetEnvPrefix("minesweeper")
	v.AutomaticEnv()

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read the configuration file: %s", err)
	}
	return v.Unmarshal(&Config)
}
