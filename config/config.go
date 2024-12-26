package config

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Domain         string   `mapstructure:"domain"`
	APIPrefix      string   `mapstructure:"api_prefix"`
	ImagePrefix    string   `mapstructure:"image_prefix"`
	TrustedProxies []string `mapstructure:"trusted_proxies"`
	OpenAICIDRs    []string `mapstructure:"openai_cidrs"`
}

var (
	cfg *Config
)

func init() {
	viper.SetEnvPrefix("DETECTIVE")
	viper.EnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigName("config")                   // name of config file (without extension)
	viper.SetConfigType("yaml")                     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/fake-ai-detective/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.fake-ai-detective") // call multiple times to add many search paths
	viper.AddConfigPath("./config")                 // optionally look for config in the working directory
	err := viper.ReadInConfig()                     // Find and read the config file
	if err != nil {                                 // Handle errors reading the config file
		panic(fmt.Errorf("fatal error while reading config file: %w", err))
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error while unmarshaling config file: %w", err))
	}

	cfg.APIPrefix = path.Clean("/" + cfg.APIPrefix)
	cfg.ImagePrefix = path.Clean("/" + cfg.ImagePrefix)

	log.Printf("Config: %+v", cfg)
}

func GetConfig() *Config {
	return cfg
}
