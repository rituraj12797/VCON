// Add Document
// load document X
// Add verson Y to Curent Docuemnt
// get version X from current Document

// All thises update instruction ust update the SataBase as well as the in memory stores to avoid stale states

package services

import (
	"vcon/internal/globalStore"
	"vcon/internal/repository"
	"vcon/internal/schema"
	"vcon/internal/engine"
	"vcon/internal/hasher"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentService struct {
	documentRespository *repository.DocumentRespository
	contentStringRepository *repository.ContentStringRepository
}

func (t *DocumentService) NewDocumentRepository() *DocumentService {
	return &DocumentService{

	}
}

func (t *DocumentService) LoadTitleOfAllDocuments() ([]string, error) {

}

func (t *DocumentService) AddDocument(title string, stringArray []string, globalStore *globalStore.Store) (*schema.Document, error) {
	
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

	return nil, nil

}

func (t *DocumentService) AddVersionToDocument(docId primitive.ObjectID, docTitle string, parentNode int, stringArr []string, globalStore *globalStore.Store) error { // second arguement is the complete statement file 

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

	return nil 

}

func (t *DocumentService) GetDocumentByTitle(title string, globalStore *globalStore.Store) (schema.Document , error){ 
	// WHAT IT RECEIVES ? => 
	// title - title of the document
	//  globalStore - reference to the global store to Read/Write from 


	// WHAT IT DOES 
	// this will set the current document 
	// search if the title exists in the in memory map ....
	// if yes fetch it from there 
	// else  FetchDocumentFromDataBaseAndSetGlobalStore
	
	
	// WHAT IT RETURNS  ==> 
	//  A schema.Document object which has it's file aray populated and delta array emptied   

}

func (t *DocumentService) GetVersionFromDocument(versionNumber int, docTitle string, globalStore *globalStore.Store) ([]string, error) { 
	
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

func (t *DocumentService) FetchDocumentFromDataBaseAndSetGlobalStore(id primitive.ObjectID, globalStore *globalStore.Store) error {


	// WHAT IT DOES 
	// => fetch it from DB and set a title vs Document entry in globalStore  
	// => fetch the usable set of strings used to render any version of this and set them in globalstore string vs dentifier ( Hydration )

	return nil
}
 

