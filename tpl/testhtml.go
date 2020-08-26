package tpl

import (
	"html/template"
	"net/http"
)

type Person struct {
	Name string
	Age  int
	Header string
}

func tmpl(w http.ResponseWriter, r *http.Request) {
	t1, err := template.ParseFiles("./html/test.html")
	if err != nil {
		panic(err)
	}
	p := &Person{
		Name: "Tom",
		Age: 22,
		Header: "https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=3511831835,544094419&fm=26&gp=0.jpg",
	}
	t1.Execute(w, p)
}

func RunTplTest() {
	server := http.Server{
		Addr: "192.168.0.110:8081",
	}
	http.HandleFunc("/tmpl", tmpl)
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Welcome to Golang"))
	})
	server.ListenAndServe()
}
