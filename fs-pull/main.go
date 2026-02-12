package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/avast/retry-go"
)

const (
	projectID  = "lh-prod-cust-lucem-test"
	collection = "workflows"
	pageSize   = 1000
	outputFile = "workflows.ndjson"
)

func main() {
	ctx := context.Background()

	fmt.Println("Creating fs client...")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	fmt.Println("success: fs client created")

	fmt.Println("Creating output file...")
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Println("success: output file created")

	query := client.Collection(collection).
		Where("name", ">=", "wrath-").
		OrderBy("name", firestore.Asc)

	var lastDoc *firestore.DocumentSnapshot
	exported := 0

	fmt.Printf("Querying fs collection %s...", collection)

	err = retry.Do(
		func() error {
			for {
				q := query.Limit(pageSize)
				if lastDoc != nil {
					q = q.StartAfter(lastDoc.Data()["name"])
				}

				docs, err := q.Documents(ctx).GetAll()
				if err != nil {
					return err
				}
				if len(docs) == 0 {
					break
				}

				for _, doc := range docs {
					record := map[string]interface{}{
						"id": doc.Ref.ID,
					}
					for k, v := range doc.Data() {
						record[k] = v
					}

					b, _ := json.Marshal(record)
					file.Write(append(b, '\n'))
					exported++
				}

				lastDoc = docs[len(docs)-1]
			}
			return nil
		})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Exported %d documents to %s\n", exported, outputFile)
}
