package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
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
	err = c.AddFunc("@every 10m", func() { deleteExpiredDocuments(ctx, client, "posts") })
	if err != nil {
		log.Fatalf("Failed to schedule task: %v", err)
	}
	c.Start()

}

// func addDocumentWithExpiry(ctx context.Context, client *firestore.Client, collection string, data map[string]interface{}) error {
// 	data["expiry"] = time.Now().Add(24 * time.Hour) // Set expiry to 24 hours from now
// 	_, _, err := client.Collection(collection).Add(ctx, data)
// 	return err
// }

func deleteExpiredDocuments(ctx context.Context, client *firestore.Client, collection string) {
	now := time.Now()
	iter := client.Collection(collection).Where("expiry", "<=", now).Documents(ctx)

	writer := client.BulkWriter(ctx)
	count := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate documents: %v", err)
			return
		}
		writer.Delete(doc.Ref)
		count++
	}

	writer.Flush()
}
