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
}

// Load returns an application configuration which is populated from the given environment variables.
func Load() (*Config, error) {
	// REST Server
	RESTServerPort := viper.GetInt("WS_REST_SERVER_PORT")

	// WS Server
	WsServerPort := viper.GetInt("WS_WEBSOCKET_SERVER_PORT")

	// default config
	config := Config{
		RESTServerPort: RESTServerPort,
		WsServerPort:   WsServerPort,
	}

	return &config, nil
}

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("WS_REST_SERVER_PORT", 8080)
	viper.SetDefault("WS_WEBSOCKET_SERVER_PORT", 7778)
}
