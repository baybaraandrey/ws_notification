package config

import (
	"github.com/spf13/viper"
)

// Config represents an application configuration.
type Config struct {
	// the REST server port. Defaults to 8080
	RESTServerPort int
	// the Websocket server port
	WsServerPort int
	// the user name for connecting to the database.
	DbUser string
	// database password
	DbPass string
	// database default name
	DbName string
	// database host
	DbHost string
	// database port
	DbPort int

	// media
	CoursesMediaRoot string
}

// Load returns an application configuration which is populated from the given environment variables.
func Load() (*Config, error) {
	// REST Server
	RESTServerPort := viper.GetInt("REST_SERVER_PORT")

	// WS Server
	WsServerPort := viper.GetInt("WS_SERVER_PORT")

	// Database
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASS")
	dbName := viper.GetString("DB_NAME")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetInt("DB_PORT")

	// default config
	config := Config{
		RESTServerPort: RESTServerPort,
		WsServerPort:   WsServerPort,
		DbUser:         dbUser,
		DbPass:         dbPass,
		DbName:         dbName,
		DbHost:         dbHost,
		DbPort:         dbPort,
	}

	return &config, nil
}

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("REST_SERVER_PORT", 8080)
	viper.SetDefault("WS_SERVER_PORT", 7778)
	viper.SetDefault("DB_USER", "courses")
	viper.SetDefault("DB_PASS", "courses")
	viper.SetDefault("DB_NAME", "courses")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 27017)
}
