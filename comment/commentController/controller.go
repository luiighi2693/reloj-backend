package commentController

import (
	"../commentEntitie"
	"../commentService"
	"github.com/kataras/iris"
)

func SetEndpoints(app *iris.Application, configFile string) {
	InitService(configFile)

	app.Get("/", func(ctx iris.Context) {
		_, _ = ctx.JSON(nil)
	})

	app.Get("/comment/{id:int min(1)}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		comment, _ := commentService.FindById(id)

		_, _ = ctx.JSON(comment)
	})

	app.Get("/comment/", func(ctx iris.Context) {
		comments, _ := commentService.FindAll()
		_, _ = ctx.JSON(comments)
	})

	app.Post("/comment/", func(ctx iris.Context) {
		comment := commentEntitie.Comment{}
		err := ctx.ReadJSON(&comment)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.WriteString(err.Error())
		} else {
			id, _ := commentService.Create(comment)
			_, _ = ctx.Writef("%v", id)
		}
	})

	app.Put("/comment/", func(ctx iris.Context) {
		comment := commentEntitie.Comment{}
		err := ctx.ReadJSON(&comment)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.WriteString(err.Error())
		} else {
			_ = commentService.Update(comment)
			_, _ = ctx.Writef("Comment Updated")
		}
	})

	app.Delete("/comment/{id:int min(1)}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		_ = commentService.Delete(id)
		_, _ = ctx.Writef("Comment Deleted")
	})
}
func InitService(configFile string) {
	commentService.InitService(configFile)
}
