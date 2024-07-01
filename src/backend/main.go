package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Struct to hold form data
type FormData struct {
	UserInput string
}

// Handler to serve the HTML form
func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../view/index.html", "../view/style.css")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// Handler to process form submission
func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	// Get user input
	userInput := r.FormValue("userInput")

	// Process the user input (for demonstration, we'll just print it)
	fmt.Fprintf(w, "You entered: %s\n", userInput)
}

func main() {
	http.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("../view"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../../assets"))))

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
