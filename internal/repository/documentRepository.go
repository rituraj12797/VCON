package repository

import (
	"context"
	"time"
	"vcon/internal/schema"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This is the abstraction to define the interface through which we will tak with the DataBase since Go doesnot have a standard ODm we define it ourselvs the way we talk with DB.

type DocumentRespository struct {
	collection *mongo.Collection //
}

func NewDocumentRepository(db *mongo.Database) *DocumentRespository { // takes arguemtn of databae and return the interface to communicate with the DocumentCollection ( Document table ) inside that Db

	return &DocumentRespository{
		collection: db.Collection("document"), // return the collection as this "document" collection
	}
}

func (r *DocumentRespository) CreateDocument(context context.Context, title string, hashes []string) (*schema.Document, error) {

	// the context the database to communicate with the title string and th base file ( version 1)
	// CRITICAL -> verify if a document with same title already exists or not ?? if exists dont create this file

	var existingDoc schema.Document

	err := r.collection.FindOne(context, bson.M{"title": title}).Decode(&existingDoc) // if doc found it will write into ecistingDoc

	if err == nil { // if error is nil means a document is found
		return nil, errors.New(" A document with same title already exists")
	}

	// say if file is not found is no found then what is the reason ??

	// reason 1: the file was actually not found and ErrNoDocument was returned
	// else somethign was wrong and error will be done
	if err != mongo.ErrNoDocuments {
		// means the erros not indicated to the same titled document being esisting error is eomthign else
		return nil, err
	}

	// Proceed to create a new document
	now := time.Now()

	// crea a new document with a pseudo root node the 0th node

	// the 0th node
	pseudoNode := schema.Node{
		ParrentNode:          -1,
		LastSnapshotAncestor: -1,
		NodeNumber:           0,
		Depth:                0,
		NodeType:             schema.NodeTypeSnapshot,

		VersionString: "pseudoRoot",
	}

	// the 1st Node
	rootNode := schema.Node{
		ParrentNode:          0,
		LastSnapshotAncestor: 1,
		NodeNumber:           1,
		Depth:                1,
		NodeType:             schema.NodeTypeSnapshot,

		VersionString: "base", // this is fixed the version string for first version is forced to be base
		FileArray:     hashes,
	}

	doc := schema.Document{
		Title:         title,
		NumberOfNodes: 2, // pseudo Root and root
		NodeArray:     []schema.Node{pseudoNode, rootNode},

		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := r.collection.InsertOne(context, doc) // result = return the saved document from DB

	if err != nil {
		// panic(err)
		return nil, err
	}

	doc.ID = result.InsertedID.(primitive.ObjectID)

	return &doc, nil
}

// Simply adds the new node to the corrosponding Document document
func (r *DocumentRespository) AddNode(ctx context.Context, docId primitive.ObjectID, newNode schema.Node) error {

	// add this node to this document

	// define the update query
	// push newNode into the nodes
	// increment numberOfNodes by 1
	// set  updated_at at time.Now()

	update := bson.M{
		"$push": bson.M{"nodes": newNode},
		"$inc":  bson.M{"numberOfNodes": 1},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateByID(ctx, docId, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("document not found or modified")
	}

	return nil

}

func (r *DocumentRespository) FindByID(ctx context.Context, docId primitive.ObjectID) (*schema.Document, error) {

	var doc schema.Document
	err := r.collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&doc)

	if err != nil {
		// means the document was not found
		return nil, err
	}

	return &doc, nil
}

func (r *DocumentRespository) FindByTitle(ctx context.Context, title string) (*schema.Document, error) {

	var doc schema.Document
	err := r.collection.FindOne(ctx, bson.M{"title": title}).Decode(&doc)

	if err != nil {
		// means the document was not found
		return nil, err
	}

	return &doc, nil
}

func (r *DocumentRespository) FindTitleOfAllDocument(ctx context.Context) ([]string, error) {
	// Use projection to only fetch the 'title' field

	// what are projections and what do they do ??

	// Projections are query options whihc optimised our query by twlling whihc files i required and which i do not

	// with this projection we enabled title and disabled _id
	// why disable only id and not any other things ??
	// ==> in mongoDB id is by default included and rest are by default exluded so need to diable it manually

	opts := options.Find().SetProjection(bson.M{"title": 1, "_id": 0})

	// what below querry does is it finds all the document with matching query bson.M{} means all will be matched and opts is our Projection query whihc tells we only aggregate the title and nothing else

	// at last the cursor is the remote aess or the pointer provided to us
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var titles []string
	for cursor.Next(ctx) {
		var result struct {
			Title string `bson:"title"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		titles = append(titles, result.Title)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return titles, nil
}

func (r *DocumentRespository) DeleteByTitle(ctx context.Context, title string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"title": title})
	return err
}
