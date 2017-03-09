package main

import (
	"os"
	"runtime"

	"github.com/devopsmakers/go-sakila-remora/remora"
	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	jww.SetStdoutThreshold(jww.LevelInfo)
	runtime.GOMAXPROCS(runtime.NumCPU())
	var r remora.Remora

	configpaths := []string{"/etc/remora/", "$HOME/.remora", ".", ".."}
	if err := r.LoadConfig(configpaths); err != nil {
		jww.FATAL.Fatalln(err)
	}

	if err := r.Serve(); err != nil {
		jww.FATAL.Fatalln(err)
	}

	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		os.Exit(-1)
	}
}
