package src

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/robfig/cron/v3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func cronjob() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "your-project-id", option.WithCredentialsFile("path/to/serviceAccountKey.json"))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	c := cron.New()
	_, err = c.AddFunc("@hourly", func() { deleteExpiredDocuments(ctx, client, "your-collection-name") })
	if err != nil {
		log.Fatalf("Failed to schedule task: %v", err)
	}
	c.Start()

	select {} // Keep the program running
}

// func addDocumentWithExpiry(ctx context.Context, client *firestore.Client, collection string, data map[string]interface{}) error {
// 	data["expiry"] = time.Now().Add(24 * time.Hour) // Set expiry to 24 hours from now
// 	_, _, err := client.Collection(collection).Add(ctx, data)
// 	return err
// }

func deleteExpiredDocuments(ctx context.Context, client *firestore.Client, collection string) {
	now := time.Now()
	iter := client.Collection(collection).Where("expiry", "<=", now).Documents(ctx)

	batch := client.BulkWriter(ctx)
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
		batch.Delete(doc.Ref)
		count++
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		log.Printf("Failed to delete documents: %v", err)
	} else {
		log.Printf("Deleted %d expired documents", count)
	}
}
