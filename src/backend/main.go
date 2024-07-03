package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type FormData struct {
	UserInput string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../view/index.html", "../view/style.css")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	userInput := r.FormValue("userInputHidden")

	fmt.Fprintf(w, "You entered: %s\n", userInput)
}

func main() {
	godotenv.Load()

	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv("CRED"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	_, _, err = client.Collection("posts").Add(ctx, map[string]interface{}{
		"post_id": uuid.New().String(),
	})
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	defer client.Close()

	http.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("../view"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../../assets"))))

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
