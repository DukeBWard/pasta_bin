package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	component := Pasta_bin()

	http.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("."))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../../assets"))))

	// http.HandleFunc("/", formHandler)
	// http.HandleFunc("/submit", submitHandler)

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
