package userService

import (
	"../userEntitie"
	"../userModel"
)

func InitService(configFile string) {
	userModel.InitModel(configFile)
}

func FindById(id int) (userEntitie.User, error) {
	var user userEntitie.User
	user, err := userModel.FindById(id)
	if err != nil {
		return user, err
	}
	return user, err
}

func FindByUsernameAndPassword(username string, password string) (userEntitie.User, error) {
	var user userEntitie.User
	user, err := userModel.FindByUsernameAndPassword(username, password)
	if err != nil {
		return user, err
	}
	return user, err
}

func FindAll() (users []userEntitie.User, err error) {
	users, err = userModel.FindAll()
	if err != nil {
		return users, err
	}
	return users, err
}

func Create(user userEntitie.User) (id int, err error) {
	return userModel.Create(user)
}

func Update(user userEntitie.User) (err error) {
	_ = userModel.Update(user)
	return
}

func Delete(id int) (err error) {
	_ = userModel.Delete(id)
	return
}
