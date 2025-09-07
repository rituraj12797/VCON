// Add Document
// load document X
// Add verson Y to Curent Docuemnt
// get version X from current Document

// All thises update instruction ust update the SataBase as well as the in memory stores to avoid stale states

package services

import (
	"context"
	"fmt"
	"time"
	"vcon/internal/engine"
	"vcon/internal/globalStore"
	"vcon/internal/hasher"
	"vcon/internal/repository"
	"vcon/internal/schema"

	"github.com/emirpasic/gods/sets/treeset"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocumentService struct {
	documentRespository     *repository.DocumentRespository
	contentStringRepository *repository.ContentStringRepository
	globalStore             *globalStore.Store
}

func (t *DocumentService) NewDocumentRepository(db *mongo.Database,
	globalStore *globalStore.Store) *DocumentService {

	return &DocumentService{
		documentRespository:     repository.NewDocumentRepository(db),
		contentStringRepository: repository.NewContentStringRepository(db),
		globalStore:             globalStore,
	}
}

func (t *DocumentService) LoadTitleOfAllDocuments(ctx context.Context) ([]string, error) {

	// fill all the entries into the globalStore to shwo that these exists and for now the valeus keep it as an empty object and will load the object as required with lazy architecture

	stringArr, err := t.documentRespository.FindTitleOfAllDocument(ctx)

	if err != nil {
		return []string{}, fmt.Errorf("failed to load document titles from DB: %w", err)
	}

	for _, title := range stringArr {
		t.globalStore.InsertNewDocument(title, nil) // all initially pointing to null document they wil be loaded lazily
	}

	return stringArr, nil
}

func (t *DocumentService) AddDocument(ctx context.Context, title string, stringArray []string) (*schema.Document, error) {

	// WHAT IT RECEIVES
	// title - title for the document
	// stringArray - array of strings the original human readable format
	// globalStore - referene to the global store

	// WHAT IT DOES
	// using hasher hashes the strign array and converts into a string of hashes
	// send this (title, array of hash ) to docRepo.createDocument to entry in DB
	// once succesfull we ger return the document that is saved in the DB
	// fills the globalStore's cntent string part with new hash , string pairs
	// fills the globalStore's title vs document store with this new title and document

	// What it returns => ??
	// when done it returns reference to the created Document with a nil error

	// heck if this file already exists in DB or nor checking in cache ( global store ) wont work here since it may not have that document since it was not needed and hence it may create duplictae titles document

	// but we can potimise it via multilevel check
	// first check our globalStore if the title exists or nto - (How does this enure we have data of all document user has created ?? => because at start we run LoadTitleOfAllDocuments)
	// this pre warms the globalStore with name of all documents that exists and since we will be  updating our globalStore on each file addition we are making our globalStore reliable
	// this not only provides spedd but high cache hit rates
	// for cases where it is not found ( false negative it eixted in Db but not in cache ) we do a Db call and there by ensuring 100% correctness

	if _, found := t.globalStore.GetDocumentByTitle(title); found {
		return nil, fmt.Errorf(" Same Title document found with title :", title)
	}

	doc, err := t.documentRespository.FindByTitle(ctx, title)

	if err == nil {
		// a document with same title was found
		return nil, fmt.Errorf(" Same Title document found with title :", title)
	}
	// this checks that the erro is mongoDocumentNotFound else if the other s somethign else we wont proceed
	if err != mongo.ErrNoDocuments {
		// problem like a network issue. We must stop and return this error.
		return nil, fmt.Errorf("failed to check for document existence: %w", err)
	}

	// NWO WE ARE SURE THAT THIS DOCUMENT DOES NTO EXIST

	// call hasher to give a hash array
	hasedArray := hasher.Hasher(stringArray) //

	var contentStringArray []schema.ContentString

	for i, str := range stringArray {
		contentStringArray = append(contentStringArray, schema.ContentString{
			Hash:    hasedArray[i],
			Content: str,
		})
	}

	err = t.contentStringRepository.AddBulk(ctx, contentStringArray)
	// save the hash, string pairs in the DB also
	if err != nil {
		return nil, err //
	}

	// save the documet in DB
	// send title, hash array to docRepo.CreateDocument
	doc, err = t.documentRespository.CreateDocument(ctx, title, hasedArray)

	if err != nil {
		return nil, err
	}

	// stored in the DB
	// update the global store
	// string vs identifier, identifier vs string, title vs document

	// title vs document
	t.globalStore.InsertNewDocument(title, doc)

	// string vs identifier
	for index, hash := range hasedArray {
		t.globalStore.InternContentString(hash, stringArray[index])
	}

	return doc, nil

}

func (t *DocumentService) AddVersionToDocument(ctx context.Context, docId primitive.ObjectID, docTitle string, parentNode int, stringArr []string) error { // second arguement is the complete statement file

	// WHAT it receives ??
	// docId - the doument id using whihc we are going to find it in the Data Base
	// docTitle - the title of the document
	// parentNode - the parrent node where this new node is going to attach
	// stringArray - the new version raw content in human readable form
	// globalStore  - a reference to the global store to Read/Write from

	// WHAT IT DOES
	// create hash array from hasher for this array of string
	// fetches the data regarding parrent node from the in memory global store
	// based on parrent depth and new depth identify if this is a new Delta or Snapshot node
	// SNAPSHOT ??  if it is a  just make a Node schem object with it and call AddNode from Document Repo
	// DELTA ?? then load the parrentNode version hash array and now
	// 			=> generate LCS, generate Delta
	// 			=> Create a Node and and fill the delta instructions in it

	// WHAT IT RETURNS => an error ( nil if succesfully otheriwse a complete error )

	// if document is not in globalStore find and hydrate it

	if _, found := t.globalStore.GetDocumentByTitle(docTitle); found == false {
		// document not in global store
		err := t.FetchDocumentFromDataBaseAndSetGlobalStore(ctx, docTitle)
		if err != nil {
			return fmt.Errorf("Can't create version since document not found ")
		}
	}

	// reference to document from globalStore
	doc, _ := t.globalStore.GetDocumentByTitle(docTitle)

	// document now found
	// verify if parrentNode exists
	// parrentNode = parrentVersion

	if len(doc.NodeArray) <= parentNode {
		return fmt.Errorf(" parrent version does not exists ")
	}

	// parrent also exists now
	// efrecnes to parrent node
	childNodeNumber := len(doc.NodeArray)
	pNode := doc.NodeArray[parentNode]
	pDepth := pNode.Depth
	// if depth % 10 == 0 | depth == 1 =====> Snapshot
	// else ================================> Delta
	childDepth := pDepth + 1

	var lastAncestor int
	var nodeType schema.NodeType
	var isSnapShot bool

	// by default asume delta and then if it matches criteria for becoming snapshot make it

	nodeType = schema.NodeTypeDelta
	lastAncestor = pNode.LastSnapshotAncestor
	isSnapShot = false

	if childDepth%10 == 0 || childDepth == 1 {
		isSnapShot = true
		lastAncestor = childNodeNumber // self LSA
		nodeType = schema.NodeTypeSnapshot
	}

	childHash := hasher.Hasher(stringArr)
	//add the new hashed statement to database and globalStore
	var newContentStrings []schema.ContentString
	for i, hash := range childHash {

		t.globalStore.InternContentString(hash, stringArr[i])
		// for bulk write
		newContentStrings = append(newContentStrings, schema.ContentString{
			Hash:    hash,
			Content: stringArr[i],
		})
	}

	// perform bulk write
	if len(newContentStrings) > 0 {
		err := t.contentStringRepository.AddBulk(ctx, newContentStrings)
		if err != nil {
			return fmt.Errorf("failed to save new content strings to DB: %w", err)
		}
	}

	var childNode schema.Node = schema.Node{
		ParrentNode:          pNode.NodeNumber,
		LastSnapshotAncestor: lastAncestor,
		NodeNumber:           childNodeNumber,
		Depth:                childDepth,
		NodeType:             nodeType,
		VersionString:        "nil", // we wont be allowing naming versions for now
		DeltaInstructions:    []schema.DeltaInstruction{},
		FileArray:            []string{},
	}

	if isSnapShot {
		childNode.FileArray = childHash
	} else {
		// operate on the engine and fill the delta array

		parrentHash, err := t.GetVersionFromDocument(ctx, pNode.NodeNumber, docTitle)

		if err != nil {
			return fmt.Errorf(" Unable to retriev parrent version to compute delta : %w ", err)
		}

		// now we have parrent hash array  and child hash array use it to

		// 1 find the lcs
		// find the delta using parrent hash and child hash and lcs
		// store the delta in child's DeltaInstructions

		// first task find the lcs
		lcs := engine.LCS(&parrentHash, &childHash)

		// find the delta array
		var deltaArray []schema.DeltaInstruction
		deltaArray = engine.GenerateDelta(&parrentHash, &childHash, &lcs)

		// store the delta array
		childNode.DeltaInstructions = deltaArray
	}

	// save the node
	err := t.documentRespository.AddNode(ctx, docId, childNode)

	if err != nil {
		return err
	}

	// UPDATED SO NOW UPDATE THE GLOBAL STORE

	// cache updated
	doc.NodeArray = append(doc.NodeArray, childNode)
	doc.NumberOfNodes++
	doc.UpdatedAt = time.Now()

	return nil

}

func (t *DocumentService) GetDocumentByTitle(ctx context.Context, title string) (*schema.Document, error) {
	// WHAT IT RECEIVES ? =>
	// title - title of the document
	// globalStore - reference to the global store to Read/Write from

	// WHAT IT DOES
	// search if the title exists in the in memory map ....

	doc, found := t.globalStore.GetDocumentByTitle(title)

	// if not found do FetchDocumentFromDataBaseAndSetGlobalStore and store document in globalStor
	if found == false {
		// try to fetch from Db and hydrate the global store now
		ner := t.FetchDocumentFromDataBaseAndSetGlobalStore(ctx, title)

		if ner != nil {
			return nil, ner
		}
	}

	// now by here it is surely inside the globalStorage
	doc, _ = t.globalStore.GetDocumentByTitle(title)
	// set current Document
	t.globalStore.ChangeCurrent(doc)

	return doc, nil
	// this will set the current document
	// WHAT IT RETURNS  ==>
	//  A schema.Document object

}

func (t *DocumentService) GetVersionFromDocument(ctx context.Context, versionNumber int, docTitle string) ([]string, error) {

	// WHAT IT RECEIVES ??
	// versionNumber - the node number = version number which teh user is asking
	// doctitle - the title of the document from whihc teh user wants this version
	//  globalStore - the reference to the global store to Read/Write from

	// WHAT IT DOES
	// finds this document from the global store
	// if it is present okay else FetchDocumentFromDataBaseAndSetGlobalStore
	// fetch the reference to the document
	// using engine find the path to LSA node frm this node
	// perform delta chain application and get the hash array for this version
	// using rendered convert this hash array into the humn redale string array

	// WHAT IT RETURNS
	// return the hashed array and from there, content renderer will take charge and rendere the complete human readable format





}



func (t *DocumentService) FetchDocumentFromDataBaseAndSetGlobalStore(ctx context.Context, title string) error {

	// WHAT IT DOES
	// => fetch it from DB and set a title vs Document entry in globalStore
	// => fetch the usable set of strings used to render any version of this and set them in globalstore string vs dentifier ( Hydration )

	doc, err := t.documentRespository.FindByTitle(ctx, title)

	if err != nil {
		// fmt.Println(" Error While fetching and saving in Global Store")
		return err
	}

	// document fetched soccesfully
	// hydrated the store
	t.globalStore.InsertNewDocument(title, doc)

	// PART 2

	// hydrate the contentString store now
	// for this we will consider all the hashes from the snapshot nodes
	// and
	// consider all the hashes that were related to add query in delta nodes and add them too
	// this final hash node list make it unique
	// then return a cursor pointer to this result

	set := treeset.NewWithStringComparator()

	for j, node := range doc.NodeArray {
		// node is doc.NodeArray[j] now

		if j > 0 { // skip the pseudo root node
			if node.NodeType == schema.NodeTypeSnapshot {
				for _, hash := range node.FileArray {
					set.Add(hash)
				}
			} else {
				for _, deltaInstruction := range node.DeltaInstructions {
					if deltaInstruction.DeltaType == schema.A { // Add type instruction
						set.Add(deltaInstruction.Val) // added hash of dela instruction
					}
				}
			}
		}
	}

	// set is the list of hashes required to render any version of this file
	// / /get them into an array
	var hashArr []string
	for _, v := range set.Values() {
		hashArr = append(hashArr, v.(string))
	}

	var resultArr []schema.ContentString // the result from the dataBase ( contains objects with hash vs string )

	// then using this hash array we would run the bulk read operation from contentStringrepo and hydrate the globalStore.stringToIdentifier and globalStore.identifiertoString
	resultArr, err = t.contentStringRepository.BulkReader(ctx, hashArr)

	if err != nil {
		return err
	}

	// we got result aray hydrate the global store now
	for _, contentString := range resultArr {
		t.globalStore.InternContentString(contentString.Hash, contentString.Content) // pased hash , content
	}

	// find all hashes whihc are used in this

	return nil
}
