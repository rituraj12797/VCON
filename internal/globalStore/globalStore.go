package globalStore

/*

	CHANGE THE THING TO IDENTIFIER AS HASH INSTEAD OF A INTEGER

*/

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"vcon/internal/schema"

	"github.com/emirpasic/gods/maps/treemap"
)

// global store storign the string v identifier, identifier vs string maps
// usefull for storage optimisation ad retranslation at the time of rendering

type Store struct {
	StringToIdentifier treemap.Map      // stores string vs hash
	IdentifierToString treemap.Map      // stores hash vs String
	TitleToDocument    treemap.Map      // stores title vs document
	CurrentDocument    *schema.Document // the current document on whihc the user is operating
	mutex              sync.RWMutex
}

// Global store defined here
var GlobalStore *Store

func InitializeStore() *Store {
	x := Store{
		StringToIdentifier: *treemap.NewWithStringComparator(),
		IdentifierToString: *treemap.NewWithStringComparator(),
		TitleToDocument:    *treemap.NewWithStringComparator(),
		CurrentDocument:    nil, // currently it points to nothing
		mutex:              sync.RWMutex{},
	}

	return &x
}

func Initialize() {
	GlobalStore = InitializeStore()
}

func (t *Store) InternContentString(hash string, statement string) error {
	// check if the string already exits or not
	t.mutex.RLock()
	_, exist := t.StringToIdentifier.Get(statement)
	t.mutex.RUnlock()

	if exist {
		// this statement already exists, return its id
		return nil
	}

	// not exits

	// lock - resource being written so enforced a write lock
	t.mutex.Lock()
	defer t.mutex.Unlock() // unlock happens even if there's a panic

	// check if some other go routine added it during this phase we may skip
	_, ext := t.StringToIdentifier.Get(statement)

	if ext {
		return nil
	}

	// nope no one has inserted till now
	t.StringToIdentifier.Put(statement, hash) //
	t.IdentifierToString.Put(hash, statement) //

	return nil
}

func (t *Store) GetStringFromIdentifier(identifier string) (string, error) {

	t.mutex.RLock()
	value, exist := t.IdentifierToString.Get(identifier)
	t.mutex.RUnlock()

	if !exist {
		return "", fmt.Errorf("identifier %s not found in store", identifier)
	}

	return value.(string), nil
}

// get identifier for a statemrnt
func (t *Store) GetIdentifier(statement string) (string, error) {
	t.mutex.RLock()
	id, exist := t.StringToIdentifier.Get(statement)
	t.mutex.RUnlock()

	if !exist {
		return "", fmt.Errorf("statement '%s' not found in store", statement)
	}

	return id.(string), nil
}

func (t *Store) InsertNewDocument(title string, doc *schema.Document) error {
	if len(title) == 0 {
		return errors.New("empty title can't be inserted")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	//  if a document with this title already exists to prevent overwrites.
	if _, found := t.TitleToDocument.Get(title); found {
		return fmt.Errorf("document with title '%s' already exists", title)
	}

	t.TitleToDocument.Put(title, doc)

	return nil
}

func (t *Store) ChangeCurrent(doc *schema.Document) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.CurrentDocument = doc
}

func (t *Store) GetCurrentDoc() (*schema.Document, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.CurrentDocument == nil {
		return nil, errors.New("no current document is set")
	}

	return t.CurrentDocument, nil
}

func (t *Store) GetDocumentByTitle(title string) (*schema.Document, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	doc, found := t.TitleToDocument.Get(title)
	if !found {
		return nil, false
	}

	return doc.(*schema.Document), true
}

func (t *Store) GetStringArray(hashArray []string) []string {

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	var stringArr []string

	for _, hash := range hashArray {
		str, _ := t.IdentifierToString.Get(hash)
		stringArr = append(stringArr, str.(string))
	}

	return stringArr
}

func (t *Store) AddNodeToDocument(title string, node schema.Node) {

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// update the document and d this node to it
	// this is the thread safe way to updae the global store
	// no function in the service or any other layer should be able to update the globalStore it must only be the function defined inside glbalStore that could update the gobalStore

	if doc, found := t.GetDocumentByTitle(title); found == true && doc != nil {
		doc.NodeArray = append(doc.NodeArray, node)
		doc.NumberOfNodes++
		doc.UpdatedAt = time.Now()
	}

}
