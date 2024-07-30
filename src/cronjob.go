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
	log.Print("hit it here1")

	godotenv.Load()
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "pasta-bin-2aae4", option.WithCredentialsFile(os.Getenv("CRED")))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	c := cron.New()
	err = c.AddFunc("@every 10m", func() { deleteExpiredDocuments(ctx, client, "posts") })
	if err != nil {
		log.Fatalf("Failed to schedule task: %v", err)
	}
	c.Start()
	select {}
}

// func addDocumentWithExpiry(ctx context.Context, client *firestore.Client, collection string, data map[string]interface{}) error {
// 	data["expiry"] = time.Now().Add(24 * time.Hour) // Set expiry to 24 hours from now
// 	_, _, err := client.Collection(collection).Add(ctx, data)
// 	return err
// }

func deleteExpiredDocuments(ctx context.Context, client *firestore.Client, collection string) {
	now := time.Now()
	log.Print("hit it here")
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
	// iter := client.Collection(collection).Where("expiry", "<=", now).Documents(ctx)

	// writer := client.BulkWriter(ctx)
	// count := 0

	// for {
	// 	doc, err := iter.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Printf("Failed to iterate documents: %v", err)
	// 		return
	// 	}
	// 	doc.Ref.Delete(ctx)
	// 	// writer.Delete(doc.Ref)
	// 	count++
	// }

	// writer.Flush()
}
