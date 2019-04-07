package commentService

import (
	"../commentEntitie"
	"../commentModel"
)

func InitService(configFile string) {
	commentModel.InitModel(configFile)
}

func FindById(id int) (commentEntitie.Comment, error) {
	var comment commentEntitie.Comment
	comment, err := commentModel.FindById(id)
	if err != nil {
		return comment, err
	}
	return comment, err
}

func FindAll() (comments []commentEntitie.Comment, err error) {
	comments, err = commentModel.FindAll()
	if err != nil {
		return comments, err
	}
	return comments, err
}

func Create(comment commentEntitie.Comment) (id int, err error) {
	return commentModel.Create(comment)
}

func Update(comment commentEntitie.Comment) (err error) {
	_ = commentModel.Update(comment)
	return
}

func Delete(id int) (err error) {
	_ = commentModel.Delete(id)
	return
}
