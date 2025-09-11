package repository

import (
	"context"
	"fmt"
	"vcon/internal/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContentStringRepository struct {
	collection *mongo.Collection
}

func NewContentStringRepository(db *mongo.Database) *ContentStringRepository {
	return &ContentStringRepository{
		collection: db.Collection("contentString"),
	}
}

func (r *ContentStringRepository) AddBulk(ctx context.Context, contentString []schema.ContentString) error {

	if len(contentString) == 0 {
		return nil // do nothing
	}

	var models []mongo.WriteModel

	// fil the values into a mongo.WriteModel to do a bulk write
	// We use UpdateOne with Upsert to avoid inserting duplicates.\\

	/*
	*	What this does
	*	we use a upsert what it does is that it insers if something based on on filter does not exist
	*	but if it eists based on the filter it will do the update mentioned in the SetUpdate part
	*	what we are doing is we are finding document on basis of hash and if found we update it to itsel else we insert ( this update to it something or insert is activated by setting .SetUpsert to true )
	* this  removed redundency by skiping the insetin of already existing document
	 */
	for _, cs := range contentString {
		// fmt.Println(" hash : ",cs.Hash, "  string : ",cs.Content)

		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"hash": cs.Hash}).
			SetUpdate(bson.M{"$setOnInsert": cs}). // $setOnInsert only applies fields during an insert (upsert)
			SetUpsert(true)

		models = append(models, model)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := r.collection.BulkWrite(ctx, models, opts)

	if err != nil {
		return fmt.Errorf("bulk write to contentString failed: %w", err)
	}

	return nil
}

func (r *ContentStringRepository) BulkReader(ctx context.Context, hashArray []string) ([]schema.ContentString, error) {

	cursor, err := r.collection.Find(ctx, bson.M{"hash": bson.M{"$in": hashArray}})

	// fmt.Println(" caled Bulkreader ")

	if err != nil {
		return nil, err
	}

	// for bulk queries the Database fecthes and stores the result on DB server and does not sends direclty but sends a cursor which is a remote controll for tat result in the DB, using it you can fetch the result

	// after fetching the result the cleanup of this cursor is required as the cursor mkes the DataBase server hold result and thereby hold resources in a stateful manner

	defer cursor.Close(ctx)
	// fmt.Println(" ======= ARARAR ==========")
	var result []schema.ContentString

	// fetch all the stored result and store in result array
	if x := cursor.All(ctx, &result); x != nil {
		return nil, x
	}

	return result, nil
}
