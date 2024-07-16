package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FormData struct {
	UserInput string
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	urlParam := chi.URLParam(r, "url")

	r.ParseForm()

	if r.Form.Has("delete") {

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

		iter := client.Collection("posts").Where("post_id", "==", urlParam).Documents(ctx)

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return
			}

			doc.Ref.Delete(ctx)
		}
	}

}

func getHandler(w http.ResponseWriter, r *http.Request) {

	urlParam := chi.URLParam(r, "url")

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

	iter := client.Collection("posts").Where("post_id", "==", urlParam).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}

		var post FormData
		docData := doc.Data()

		if content, ok := docData["body"].(string); ok {
			post.UserInput = content
		} else {
			log.Fatalf("error parsing document data")
		}

		component := pasta_bin(post.UserInput)
		component.Render(r.Context(), w)

	}

	// tmpl.Execute(w, nil)
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

	postID := uuid.New().String()
	userInput := r.FormValue("userInputHidden")

	_, _, err = client.Collection("posts").Add(ctx, map[string]interface{}{
		"post_id": postID,
		"body":    userInput,
	})
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	component := submit("http://localhost:8080/" + postID)
	component.Render(r.Context(), w)
}

func main() {
	// Initialize Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// need to use /view/ like how I do in the index.html for the style sheet.  basically virtual link
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))

	// Define routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		component := pasta_bin("")
		templ.Handler(component).ServeHTTP(w, r)
	})

	r.Post("/submit", submitHandler)
	r.Get("/{url}", getHandler)
	r.Delete("/{url}", deleteHandler)

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
