package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	Id      int
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			time STRING,
			content STRING)`, tableName)

	//cmdの実行
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}
}

//データを一つ取得
func GetOne(r *http.Request) Todo {
	DbConnection, _ := sql.Open("sqlite3", "./todo.sql")
	defer DbConnection.Close()

	//編集したいデータのIDを取得
	idS := r.URL.Path[len("/edit/"):]
	idI, _ := strconv.Atoi(idS)
	cmd := fmt.Sprintf(`
		SELECT * FROM %s WHERE id = ? `, tableName)

	row := DbConnection.QueryRow(cmd, idI)

	var todo Todo
	err := row.Scan(&todo.Id, &todo.Time, &todo.Content)
	if err != nil {
		log.Println(err)
	}

	return todo
}

//DBの中身を全取得
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
		err := rows.Scan(&t.Id, &t.Time, &t.Content)
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
		INSERT INTO %s (id, time, content) VALUES (NULL, ?, ?)`, tableName)

	// Execの第2引数以降はDBに入れたいデータ
	_, err := DbConnection.Exec(cmd, time, content)
	if err != nil {
		log.Fatalln(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//Update
func Update(w http.ResponseWriter, r *http.Request) {
	DbConnection, _ := sql.Open("sqlite3", "./todo.sql")
	defer DbConnection.Close()

	//IDの取得
	idS := r.URL.Path[len("/update/"):]
	idI, _ := strconv.Atoi(idS)
	//情報の更新
	time := time.Now()
	content := r.FormValue("content")

	cmd := fmt.Sprintf(`
		UPDATE  %s SET (time, content) = (?, ?) WHERE id=? `, tableName)
	_, err := DbConnection.Exec(cmd, time, content, idI)

	if err != nil {
		log.Fatalln(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//Delete
func Delete(w http.ResponseWriter, r *http.Request) {
	DbConnection, _ := sql.Open("sqlite3", "./todo.sql")
	defer DbConnection.Close()

	idS := r.URL.Path[len("/delete/"):]
	idI, _ := strconv.Atoi(idS)
	cmd := fmt.Sprintf(`
		DELETE FROM  %s  WHERE id=? `, tableName)
	_, err := DbConnection.Exec(cmd, idI)
	if err != nil {
		log.Fatalln(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}