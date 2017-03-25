package main

import (
	"os"
	"runtime"

	"github.com/devopsmakers/go-sakila-remora/mysql"
	"github.com/devopsmakers/go-sakila-remora/remora"

	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	jww.SetStdoutThreshold(jww.LevelInfo)
	runtime.GOMAXPROCS(runtime.NumCPU())
	var r remora.Remora

	servicename := "mysql"
	configpaths := []string{"/etc/remora/", "$HOME/.remora", ".", ".."}
	if err := r.LoadConfig(configpaths, servicename); err != nil {
		jww.FATAL.Fatalln(err)
	}

	if err := r.Serve(mysql.HealthCheck); err != nil {
		jww.FATAL.Fatalln(err)
	}

	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		os.Exit(-1)
	}
}
