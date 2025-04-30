package main

import (
	"github.com/SwanHtetAungPhyo/auth/cmd/config"
	"github.com/SwanHtetAungPhyo/auth/cmd/server"
	"os"
	"os/signal"
)

func init() {
	cfg := config.ConfigState{}
	cfg.GetConfig()
	config.RegisterToConsul()
}
func main() {

	app := server.NewApp()
	go func() {
		app.Start()
	}()

	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	<-osChan
	app.Stop()

}
