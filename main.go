package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"goapp/todo/models"
)

var templates = template.Must(template.ParseFiles("templates/toppage.html"))

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
	//
	//	todo, _ := models.GetAll()
	//	err := templates.ExecuteTemplate(w, "toppage.html", todo.Data)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	todo, _ := models.GetAll()

	err := templates.ExecuteTemplate(w, "toppage.html", todo.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	//	http.Handle("/", &templateHandler{filename: "toppage.html"})
	http.Handle("/newpost", &templateHandler{filename: "newpost.html"})
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/new", models.Insert)

	//Webサーバー開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
