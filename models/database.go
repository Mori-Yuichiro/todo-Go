package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//コネクションプールの作成
var DbConnection *sql.DB

type Data struct {
	Data []Todo
}

//データ格納用のストラクト
type Todo struct {
	Time    string
	Content string
}

const (
	tableName = "todo"
)

func init() {
	//データベースに関する処理
	DbConnection, _ := sql.Open("sqlite3", "./todo.sql")
	//Connectionをクローズ
	defer DbConnection.Close()

	//todoテーブルの作成
	cmd := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			Time STRING,
			Content STRING)`, tableName)

	//cmdの実行
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}
}

//DBの中身を取得
func GetAll() (todo *Data, err error) {
	DbConnection, _ := sql.Open("sqlite3", "./todo.sql")
	defer DbConnection.Close()

	cmd := fmt.Sprintf(`
		SELECT * FROM %s `, tableName)
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return
	}
	defer rows.Close()

	//データの保存領域を確保
	todo = &Data{}
	for rows.Next() {
		var t Todo
		// Scanにて、structのアドレスにデータを入れる
		err := rows.Scan(&t.Time, &t.Content)
		if err != nil {
			log.Println(err)
		}
		//データの取得
		todo.Data = append(todo.Data, t)
	}
	return todo, nil
}

//データの挿入
func Insert(w http.ResponseWriter, r *http.Request) {
	time := time.Now()
	content := r.FormValue("content")
	DbConnection, _ := sql.Open("sqlite3", "./todo.sql")
	defer DbConnection.Close()

	cmd := fmt.Sprintf(`
		INSERT INTO %s (Time, Content) VALUES (?, ?)`, tableName)

	// Execの第2引数以降はDBに入れたいデータ
	_, err := DbConnection.Exec(cmd, time, content)
	if err != nil {
		log.Fatalln(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
