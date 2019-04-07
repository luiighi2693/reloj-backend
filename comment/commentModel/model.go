package commentModel

import (
	"../../util/config"
	"../commentEntitie"
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

func FindById(id int) (comment commentEntitie.Comment, err error) {
	comment = commentEntitie.Comment{}
	err = Db.QueryRow("select id, content, author from comments where id = $1",
		id).Scan(&comment.Id, &comment.Content, &comment.Author)
	return comment, nil
}

func FindAll() (comments []commentEntitie.Comment, err error) {
	rows, err := Db.Query("select id, content, author from comments")
	for rows.Next() {
		comment := commentEntitie.Comment{}
		err = rows.Scan(&comment.Id, &comment.Content, &comment.Author)
		comments = append(comments, comment)
	}
	return comments, nil
}

func Create(comment commentEntitie.Comment) (id int, err error) {
	statement := "insert into comments (content, author) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(comment.Content, comment.Author).Scan(&comment.Id)
	return comment.Id, err
}

func Update(comment commentEntitie.Comment) (err error) {
	_, err = Db.Exec("update comments set content = $2, author = $3 where id = $1",
		comment.Id, comment.Content, comment.Author)
	return
}

func Delete(id int) (err error) {
	_, err = Db.Exec("delete from comments where id = $1", id)
	return
}
