package main

import (
	"github.com/kataras/iris"

	"../user/userController"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	"../util/config"
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
)

func NewApp(configFile string) *iris.Application {
	app := iris.New()

	enableApiDoc(app)

	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	userController.SetEndpoints(app, configFile)

	return app
}

func enableApiDoc(app *iris.Application) {
	yaag.Init(&yaag.Config{
		On:       true,
		DocTitle: "Iris",
		DocPath:  "apidoc-user.html",
		BaseUrls: map[string]string{"Production": "", "Staging": ""},
	})
	app.Use(irisyaag.New())
}

func main() {
	configFileName := "config.json"
	app := NewApp(configFileName)
	_ = app.Run(iris.Addr(config.GetPort(config.LoadConfiguration(configFileName).UserPort)), iris.WithoutServerError(iris.ErrServerClosed))
}
