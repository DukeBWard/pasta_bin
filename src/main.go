package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FormData struct {
	UserInput string
}

// func formHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseFiles("../view/index.html", "../view/style.css")
// 	if err != nil {
// 		http.Error(w, "Could not load template", http.StatusInternalServerError)
// 		return
// 	}

// 	tmpl.Execute(w, nil)
// }

func getHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../view/index.html", "../view/style.css")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

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

	defer client.Close()

	iter := client.Collection("posts").Select()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
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

	defer client.Close()

	post_id := uuid.New().String()

	userInput := r.FormValue("userInputHidden")

	_, _, err = client.Collection("posts").Add(ctx, map[string]interface{}{
		"post_id": post_id,
		"body":    userInput,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Fprintf(w, "You entered: %s\n", post_id)
}

func main() {

	component := Pasta_bin()

	// need to use /view/ like how I do in the index.html for the style sheet.  basically virtual link
	http.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("."))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))

	//http.HandleFunc("/", formHandler)
	http.Handle("/", templ.Handler(component))
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc()

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
