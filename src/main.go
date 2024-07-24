package src

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	userInput  string
	postId     string
	expiryTime time.Time
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	postID := r.FormValue("post_id")
	fmt.Print(postID)

	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
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

	iter := client.Collection("posts").Where("post_id", "==", postID).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Error finding document", http.StatusInternalServerError)
			return
		}

		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			http.Error(w, "Error deleting document", http.StatusInternalServerError)
			return
		}
	}

	component := pasta_deleted("http://localhost:8080/")
	component.Render(r.Context(), w)

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
			post.userInput = content
		} else {
			log.Fatalf("error parsing document data")
		}

		component := get_pasta_bin(post.userInput, urlParam)
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

	// get the expiry time in string minutes, convert to a time duration and add that to the now
	expiryTimeString := r.FormValue("expiryTime")
	expiryTimeDuration, err := time.ParseDuration(expiryTimeString + "m")
	expiryTime := time.Now().Add(expiryTimeDuration)

	if err != nil {
		fmt.Println(err)
	}

	post := FormData{
		userInput:  userInput,
		postId:     postID,
		expiryTime: expiryTime,
	}

	_, _, err = client.Collection("posts").Add(ctx, map[string]interface{}{
		"post_id":     post.postId,
		"body":        post.userInput,
		"expiry_time": post.expiryTime,
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
		component := create_pasta_bin("", "")
		templ.Handler(component).ServeHTTP(w, r)
	})

	r.Post("/submit", submitHandler)
	r.Get("/{url}", getHandler)
	r.Post("/delete", deleteHandler)

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
