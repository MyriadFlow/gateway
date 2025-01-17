package main

import (
	"fmt"

	"github.com/MyriadFlow/gateway/app"
	"github.com/MyriadFlow/gateway/config/dbconfig"
	"github.com/MyriadFlow/gateway/config/envconfig"
	"github.com/MyriadFlow/gateway/util/pkg/logwrapper"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app.Init()
	dbconfig.GetDb()
	logwrapper.Log.Info("Starting app")
	addr := fmt.Sprintf(":%d", envconfig.EnvVars.APP_PORT)
	err := app.GinApp.Run(addr)
	if err != nil {
		logwrapper.Log.Fatal(err)
	}
}
