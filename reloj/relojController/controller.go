package relojController

import (
	"../relojEntitie"
	"../relojService"
	"github.com/kataras/iris"
)

func SetEndpoints(app *iris.Application, configFile string) {
	InitService(configFile)

	app.Get("/", func(ctx iris.Context) {
		_, _ = ctx.JSON(nil)
	})

	//app.Get("/user/{id:int min(1)}", func(ctx iris.Context) {
	//	id, _ := ctx.Params().GetInt("id")
	//	user, _ := relojService.FindById(id)
	//
	//	_, _ = ctx.JSON(user)
	//})
	//
	//app.Get("/user/{username:string min(1)}/{password:string min(1)}", func(ctx iris.Context) {
	//	username := ctx.Params().GetString("username")
	//	password := ctx.Params().GetString("password")
	//	user, _ := relojService.FindByUsernameAndPassword(username, password)
	//
	//	if user.Id == 0 {
	//		_, _ = ctx.JSON(nil)
	//	} else {
	//		_, _ = ctx.JSON(user)
	//	}
	//})

	app.Get("/reloj/", func(ctx iris.Context) {
		users, _ := relojService.FindAll()
		_, _ = ctx.JSON(users)
	})

	app.Post("/reloj/dni/verification", func(ctx iris.Context) {
		user := relojEntitie.User{}
		err := ctx.ReadJSON(&user)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.WriteString(err.Error())
		} else {
			hasError, err := relojService.VerifyDni(user.NroDoc)
			if !hasError {
				_, _ = ctx.JSON(iris.Map{"hasError": hasError, "error": nil})
			} else {
				_, _ = ctx.JSON(iris.Map{"hasError": hasError, "error": err.Error()})
			}
		}
	})

	app.Post("/reloj/dni-password/verification", func(ctx iris.Context) {
		user := relojEntitie.User{}
		err := ctx.ReadJSON(&user)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.WriteString(err.Error())
		} else {
			hasError, err := relojService.VerifyDniAndPassword(user)
			if !hasError {
				_, _ = ctx.JSON(iris.Map{"hasError": hasError, "error": nil})
			} else {
				_, _ = ctx.JSON(iris.Map{"hasError": hasError, "error": err.Error()})
			}
		}
	})

	app.Post("/reloj/marcacion/verification", func(ctx iris.Context) {
		user := relojEntitie.User{}
		err := ctx.ReadJSON(&user)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.WriteString(err.Error())
		} else {
			message := ""
			hasError, err, message := relojService.VerifyMarkByOperation(user.NroDoc, user.Operation) //I - O
			if !hasError {
				_, _ = ctx.JSON(iris.Map{"hasError": hasError, "error": nil, "message": message})
			} else {
				_, _ = ctx.JSON(iris.Map{"hasError": hasError, "error": err.Error(), "message": message})
			}
		}
	})
	//app.Post("/user/", func(ctx iris.Context) {
	//	user := relojEntitie.User{}
	//	err := ctx.ReadJSON(&user)
	//	if err != nil {
	//		ctx.StatusCode(iris.StatusBadRequest)
	//		_, _ = ctx.WriteString(err.Error())
	//	} else {
	//		id, _ := relojService.Create(user)
	//		_, _ = ctx.Writef("%v", id)
	//	}
	//})
	//
	//app.Put("/user/", func(ctx iris.Context) {
	//	user := relojEntitie.User{}
	//	err := ctx.ReadJSON(&user)
	//	if err != nil {
	//		ctx.StatusCode(iris.StatusBadRequest)
	//		_, _ = ctx.WriteString(err.Error())
	//	} else {
	//		_ = relojService.Update(user)
	//		_, _ = ctx.Writef("User Updated")
	//	}
	//})
	//
	//app.Delete("/user/{id:int min(1)}", func(ctx iris.Context) {
	//	id, _ := ctx.Params().GetInt("id")
	//	_ = relojService.Delete(id)
	//	_, _ = ctx.Writef("User Deleted")
	//})
}
func InitService(configFile string) {
	relojService.InitService(configFile)
}
