package remora

import (
	"fmt"
	"os"

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
	CheckEvery    string
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
func (r *Remora) LoadConfig(configpaths []string) {

	viper.SetConfigName("config")

	for _, configpath := range configpaths {
		viper.AddConfigPath(configpath)
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// some of our config names need aliasing to unmarshal correctly
	viper.RegisterAlias("acceptable-lag", "AcceptableLag")
	viper.RegisterAlias("check-every", "CheckEvery")
	viper.RegisterAlias("http-serve", "HTTPServe")

	umerr := viper.Unmarshal(&r.Config)
	if umerr != nil {
		fmt.Printf("unable to decode into struct, %v", umerr)
		os.Exit(1)
	}

}

// Run starts running checks and exposes the HTTP endpoint
func (r *Remora) Run() {

	fmt.Printf("%v\n", r.Config.MySQL.User)
}
