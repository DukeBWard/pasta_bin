package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func cronjob() {

	godotenv.Load()
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "pasta-bin-2aae4", option.WithCredentialsFile(os.Getenv("CRED")))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	c := cron.New()
	err = c.AddFunc("@every 5s", func() { deleteExpiredDocuments(ctx, client, "posts") })
	if err != nil {
		log.Fatalf("Failed to schedule task: %v", err)
	}
	c.Start()
}

func deleteExpiredDocuments(ctx context.Context, client *firestore.Client, collection string) {
	now := time.Now()
	godotenv.Load()

	// Use a service account
	//ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv("CRED"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	iter := client.Collection(collection).Where("expiry_time", "<=", now).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Print("Error finding document", http.StatusInternalServerError)
			return
		}

		doc.Ref.Delete(ctx)
	}
}
