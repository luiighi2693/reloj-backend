package relojModel

import (
	"../../util/config"
	"../relojEntitie"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func InitModel(configFile string) {
	var err error

	Db, err = sql.Open("mysql",
		config.GenerateStringMysqlDbConnection(configFile))
	if err != nil {
		panic(err)
	}

	err = Db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

//func FindById(id int) (user relojEntitie.User, err error) {
//	user = relojEntitie.User{}
//	err = Db.QueryRow("select id, username, password from users where id = $1",
//		id).Scan(&user.Id, &user.Username, &user.Password)
//	return user, nil
//}
//
//func FindByUsernameAndPassword(username string, password string) (user relojEntitie.User, err error) {
//	user = relojEntitie.User{}
//	err = Db.QueryRow("select id, username, password from users where username = $1 and password = $2",
//		username, password).Scan(&user.Id, &user.Username, &user.Password)
//	return user, nil
//}

func FindAll() (relojs []relojEntitie.Reloj, err error) {
	rows, err := Db.Query("select cod_us, ing, sale, estado, tiempo from reloj limit 20")
	for rows.Next() {
		reloj := relojEntitie.Reloj{}
		err = rows.Scan(&reloj.CodUs, &reloj.Ing, &reloj.Sale, &reloj.Estado, &reloj.Tiempo)
		relojs = append(relojs, reloj)
	}
	return relojs, nil
}

func VerifyDni(dni string) (hasError bool, err error) {
	user := relojEntitie.User{}
	err = Db.QueryRow("SELECT nro_doc, psw FROM usuario WHERE nro_doc = ?", dni).Scan(&user.NroDoc, &user.Psw)
	if err != nil {
		return true, err
	}

	return false, err
}

func VerifyDniAndPassword(user relojEntitie.User) (hasError bool, err error) {
	h := md5.New()
	h.Write([]byte(user.Psw))
	bs := h.Sum(nil)
	pswToString := hex.EncodeToString(bs)

	err = Db.QueryRow("SELECT nro_doc, psw FROM usuario WHERE nro_doc = ? AND psw = ?", user.NroDoc, pswToString).Scan(&user.NroDoc, &user.Psw)
	if err != nil {
		return true, err
	}

	return false, err
}

func VerifyMarkByOperation(nroDoc string, operation string) (hasError bool, err error, message string) {
	user := relojEntitie.User{}
	err = Db.QueryRow("SELECT nro_doc, cod_us, nombres, apellido FROM usuario WHERE nro_doc = ?", nroDoc).Scan(&user.NroDoc, &user.CodUs, &user.Nombres, &user.Apellido)
	if err != nil {
		return true, err, "Error en validacion de tarjeta!"
	}

	return false, err, "OK"
}

//func Create(user relojEntitie.User) (id int, err error) {
//	statement := "insert into users (username, password) values ($1, $2) returning id"
//	stmt, err := Db.Prepare(statement)
//	if err != nil {
//		return
//	}
//	defer stmt.Close()
//	err = stmt.QueryRow(user.Username, user.Password).Scan(&user.Id)
//	return user.Id, err
//}
//
//func Update(user relojEntitie.User) (err error) {
//	_, err = Db.Exec("update users set username = $2, password = $3 where id = $1",
//		user.Id, user.Username, user.Password)
//	return
//}
//
//func Delete(id int) (err error) {
//	_, err = Db.Exec("delete from users where id = $1", id)
//	return
//}
