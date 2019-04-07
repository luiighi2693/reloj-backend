package userModel

import (
	"../../util/config"
	"../userEntitie"
	"database/sql"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitModel(configFile string) {
	var err error

	Db, err = sql.Open("postgres",
		config.GenerateStringPostgresDbConnection(configFile))
	if err != nil {
		panic(err)
	}
}

func FindById(id int) (user userEntitie.User, err error) {
	user = userEntitie.User{}
	err = Db.QueryRow("select id, username, password from users where id = $1",
		id).Scan(&user.Id, &user.Username, &user.Password)
	return user, nil
}

func FindByUsernameAndPassword(username string, password string) (user userEntitie.User, err error) {
	user = userEntitie.User{}
	err = Db.QueryRow("select id, username, password from users where username = $1 and password = $2",
		username, password).Scan(&user.Id, &user.Username, &user.Password)
	return user, nil
}

func FindAll() (users []userEntitie.User, err error) {
	rows, err := Db.Query("select id, username, password from users")
	for rows.Next() {
		user := userEntitie.User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Password)
		users = append(users, user)
	}
	return users, nil
}

func Create(user userEntitie.User) (id int, err error) {
	statement := "insert into users (username, password) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(user.Username, user.Password).Scan(&user.Id)
	return user.Id, err
}

func Update(user userEntitie.User) (err error) {
	_, err = Db.Exec("update users set username = $2, password = $3 where id = $1",
		user.Id, user.Username, user.Password)
	return
}

func Delete(id int) (err error) {
	_, err = Db.Exec("delete from users where id = $1", id)
	return
}
