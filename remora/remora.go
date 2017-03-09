package remora

import (
	"fmt"

	"github.com/spf13/viper"
)

// Remora interface
type Remora struct {
	Config *Config
}

// Config basic remora configuration
type Config struct {
	MySQL         Connection
	AcceptableLag string
	CacheTTL      string
	HTTPServe     int
}

// Connection detail for MySQL
type Connection struct {
	Port int
	Host string
	User string
	Pass string
}

type status struct {
}

// LoadConfig inits the config file and reads the default config information
// into Remora.Config. For testability it accepts an array containing dirs to
// search for a config file.
func (r *Remora) LoadConfig(configpaths []string) error {

	// Explicitly reset viper - helps for testing errors
	viper.Reset()

	viper.SetConfigName("config")

	for _, configpath := range configpaths {
		viper.AddConfigPath(configpath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// some of our config names need aliasing to unmarshal correctly
	viper.RegisterAlias("acceptable-lag", "AcceptableLag")
	viper.RegisterAlias("cache-ttl", "CacheTTL")
	viper.RegisterAlias("http-serve", "HTTPServe")

	if err := viper.Unmarshal(&r.Config); err != nil {
		return err
	}

	return nil
}

// Run starts running checks and exposes the HTTP endpoint
func (r *Remora) Run() {

	fmt.Printf("%v\n", r.Config.MySQL.User)
}
