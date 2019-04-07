package userController

import (
	"../userEntitie"
	"../userService"
	"github.com/kataras/iris"
)

func SetEndpoints(app *iris.Application, configFile string) {
	InitService(configFile)

	app.Get("/", func(ctx iris.Context) {
		_, _ = ctx.JSON(nil)
	})

	app.Get("/user/{id:int min(1)}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		user, _ := userService.FindById(id)

		_, _ = ctx.JSON(user)
	})

	app.Get("/user/{username:string min(1)}/{password:string min(1)}", func(ctx iris.Context) {
		username := ctx.Params().GetString("username")
		password := ctx.Params().GetString("password")
		user, _ := userService.FindByUsernameAndPassword(username, password)

		if user.Id == 0 {
			_, _ = ctx.JSON(nil)
		} else {
			_, _ = ctx.JSON(user)
		}
	})

	app.Get("/user/", func(ctx iris.Context) {
		users, _ := userService.FindAll()
		_, _ = ctx.JSON(users)
	})

	app.Post("/user/", func(ctx iris.Context) {
		user := userEntitie.User{}
		err := ctx.ReadJSON(&user)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.WriteString(err.Error())
		} else {
			id, _ := userService.Create(user)
			_, _ = ctx.Writef("%v", id)
		}
	})

	app.Put("/user/", func(ctx iris.Context) {
		user := userEntitie.User{}
		err := ctx.ReadJSON(&user)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.WriteString(err.Error())
		} else {
			_ = userService.Update(user)
			_, _ = ctx.Writef("User Updated")
		}
	})

	app.Delete("/user/{id:int min(1)}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		_ = userService.Delete(id)
		_, _ = ctx.Writef("User Deleted")
	})
}
func InitService(configFile string) {
	userService.InitService(configFile)
}
