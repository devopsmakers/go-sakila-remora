package main

import (
	"os"
	"runtime"

	"github.com/devopsmakers/go-sakila-remora/remora"
	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var r remora.Remora

	configpaths := []string{"/etc/remora/", "$HOME/.remora", "."}
	r.LoadConfig(configpaths)

	r.Run()

	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		os.Exit(-1)
	}
}
