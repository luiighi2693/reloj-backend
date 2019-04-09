package relojService

import (
	"../relojEntitie"
	"../relojModel"
)

func InitService(configFile string) {
	relojModel.InitModel(configFile)
}

//func FindById(id int) (reloj relojEntitie.Reloj, error) {
//	var user relojEntitie.User
//	user, err := relojModel.FindById(id)
//	if err != nil {
//		return user, err
//	}
//	return user, err
//}
//
//func FindByUsernameAndPassword(username string, password string) (relojEntitie.User, error) {
//	var user relojEntitie.User
//	user, err := relojModel.FindByUsernameAndPassword(username, password)
//	if err != nil {
//		return user, err
//	}
//	return user, err
//}

func FindAll() (relojs []relojEntitie.Reloj, err error) {
	relojs, err = relojModel.FindAll()
	if err != nil {
		return relojs, err
	}
	return relojs, err
}

func VerifyDni(dni string) (hasError bool, err error) {
	return relojModel.VerifyDni(dni)
}

func VerifyDniAndPassword(user relojEntitie.User) (hasError bool, err error) {
	return relojModel.VerifyDniAndPassword(user)
}

func VerifyMarkByOperation(nroDoc string, operation string) (hasError bool, err error, message string) {
	return relojModel.VerifyMarkByOperation(nroDoc, operation)
}

//func Create(user relojEntitie.User) (id int, err error) {
//	return relojModel.Create(user)
//}
//
//func Update(user relojEntitie.User) (err error) {
//	_ = relojModel.Update(user)
//	return
//}
//
//func Delete(id int) (err error) {
//	_ = relojModel.Delete(id)
//	return
//}
