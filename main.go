package main

import (
	"log"
	"net/http"
	"text/template"

	"goapp/todo/models"
)

var templates = template.Must(template.ParseFiles("templates/toppage.html", "templates/edit.html"))

func viewHandler(w http.ResponseWriter, r *http.Request) {
	todo, _ := models.GetAll()

	err := templates.ExecuteTemplate(w, "toppage.html", todo.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	todo := models.GetOne(r)

	err := templates.ExecuteTemplate(w, "edit.html", todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/new", models.Insert)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/update/", models.Update)
	http.HandleFunc("/delete/", models.Delete)

	//Webサーバー開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
