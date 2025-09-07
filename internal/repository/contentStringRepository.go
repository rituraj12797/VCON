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
	for _, cs := range contentString {
		fmt.Println(" hash : ",cs.Hash, "  string : ",cs.Content)
		model := mongo.NewInsertOneModel().SetDocument(cs)
		models = append(models, model)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := r.collection.BulkWrite(ctx, models,opts)

	if err != nil {
		// check what are the erros if error is for duplicacy it is acceptable and should nto be considered as n error 
		if bwe, ok := err.(mongo.BulkWriteException); ok {
			// this peace of code check if any bulkwrite exceptios occoured ot not ?? 
			// this ok = true means err is of type BulkWriteException means some write resulted in error 
			// ok => true
			
			for _,e := range bwe.WriteErrors {
				if e.Code != 11000 { // 11000 means duplicacy error 
					// key is not duplicate means error is something else 
					return err
				}
			}
			return nil
		}

		return err; // it is not Bulk write exception means errors must be comsehitn else 

	}

	return nil
}

func (r *ContentStringRepository) BulkReader(ctx context.Context, hashArray []string) ([]schema.ContentString, error) {

	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": hashArray}})

	if err != nil {
		return nil, err
	}

	// for bulk queries the Database fecthes and stores the result on DB server and does not sends direclty but sends a cursor which is a remote controll for tat result in the DB, using it you can fetch the result 

	// after fetching the result the cleanup of this cursor is required as the cursor mkes the DataBase server hold result and thereby hold resources in a stateful manner 

	defer cursor.Close(ctx)

	var result []schema.ContentString

	// fetch all the stored result and store in result array 
	if x := cursor.All(ctx,&result) ; x != nil {
		return nil, x
	}

	return result, nil
}
