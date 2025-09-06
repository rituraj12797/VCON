package globalStore

/*

	CHANGE THE THING TO IDENTIFIER AS HASH INSTEAD OF A INTEGER

*/

import (
	"errors"
	"fmt"
	"sync"
	"vcon/internal/schema"

	"github.com/emirpasic/gods/maps/treemap"
)

// global store storign the string v identifier, identifier vs string maps
// usefull for storage optimisation ad retranslation at the time of rendering

type Store struct {
	stringToIdentifier treemap.Map      // stores string vs hash
	identifierToString treemap.Map      // stores hash vs String
	titleToDocument    treemap.Map      // stores title vs document
	currentDocument    *schema.Document // the current document on whihc the user is operating
	mutex              sync.RWMutex
}

// Global store defined here
var GlobalStore *Store

func InitializeStore() *Store {
	x := Store{
		stringToIdentifier: *treemap.NewWithStringComparator(),
		identifierToString: *treemap.NewWithStringComparator(),
		titleToDocument:    *treemap.NewWithStringComparator(),
		currentDocument:    nil, // currently it points to nothing
		mutex:              sync.RWMutex{},
	}

	return &x
}

func Initialize() {
	GlobalStore = InitializeStore()
}

func (t *Store) InternContentString(statement string) error {
	// check if the string already exits or not
	t.mutex.RLock()
	_, exist := t.stringToIdentifier.Get(statement)
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
	_, ext := t.stringToIdentifier.Get(statement)

	if ext {
		return nil
	}

	return nil
}

func (t *Store) GetStringFromIdentifier(identifier int) (string, error) {

	t.mutex.RLock()
	value, exist := t.identifierToString.Get(identifier)
	t.mutex.RUnlock()

	if !exist {
		return "", fmt.Errorf("identifier %d not found in store", identifier)
	}

	return value.(string), nil
}

// get identifier for a statemrnt
func (t *Store) GetIdentifier(statement string) (int, error) {
	t.mutex.RLock()
	id, exist := t.stringToIdentifier.Get(statement)
	t.mutex.RUnlock()

	if !exist {
		return 0, fmt.Errorf("statement '%s' not found in store", statement)
	}

	return id.(int), nil
}

func (t *Store) InsertNewDocument(title string, doc *schema.Document) error {
	if len(title) == 0 {
		return errors.New("empty title can't be inserted")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	//  if a document with this title already exists to prevent overwrites.
	if _, found := t.titleToDocument.Get(title); found {
		return fmt.Errorf("document with title '%s' already exists", title)
	}

	t.titleToDocument.Put(title, doc)

	return nil
}

func (t *Store) ChangeCurrent(doc *schema.Document) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.currentDocument = doc
}

func (t *Store) GetCurrentDoc() (*schema.Document, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.currentDocument == nil {
		return nil, errors.New("no current document is set")
	}

	return t.currentDocument, nil
}

func (t *Store) GetDocumentByTitle(title string) (*schema.Document, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	doc, found := t.titleToDocument.Get(title)
	if !found {
		return nil, fmt.Errorf("document with title '%s' not found", title)
	}

	return doc.(*schema.Document), nil
}
