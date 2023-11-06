package initializers

import (
	"fmt"

	"github.com/spf13/viper"
)

type Env struct {
	SpotifyClientID         string `mapstructure:"SPOTIFY_OAUTH_CLIENT_ID"`
	SpotifyClientSecret     string `mapstructure:"SPOTIFY_OAUTH_CLIENT_SECRET"`
	SpotifyOAuthRedirectURI string `mapstructure:"SPOTIFY_OAUTH_REDIRECT_URL"`
}

func LoadEnv() (env Env, err error) {
	viper.SetConfigFile(".env")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error env file: %s", err))
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshal env file: %s", err))
	}

	return
}
