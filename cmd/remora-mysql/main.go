package main

import (
	"os"
	"runtime"

	"github.com/devopsmakers/go-sakila-remora/mysql"
	"github.com/devopsmakers/go-sakila-remora/remora"

	jww "github.com/spf13/jwalterweatherman"
)

func main() {

	// Set this to the name of the service in the config file
	servicename := "mysql"

	// Set this to the name of the function that returns a remora.Result
	servicecheck := mysql.HealthCheck

	jww.SetStdoutThreshold(jww.LevelInfo)
	runtime.GOMAXPROCS(runtime.NumCPU())
	var r remora.Remora

	configpaths := []string{"/etc/remora/", "$HOME/.remora", ".", ".."}
	if err := r.LoadConfig(configpaths, servicename); err != nil {
		jww.FATAL.Fatalln(err)
	}

	if err := r.Serve(servicecheck); err != nil {
		jww.FATAL.Fatalln(err)
	}

	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		os.Exit(-1)
	}
}
