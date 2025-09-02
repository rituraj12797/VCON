package globalStore

import (
	"fmt"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
)

// global store storign the string v identifier, identifier vs string maps
// usefull for storage optimisation ad retranslation at the time of rendering

type Store struct {
	stringToIdentifier treemap.Map
	identifierToString treemap.Map

	nextAvailableIdentifier int // same as the next usable node 
	mutex                   sync.RWMutex
}

// Global store defined here
var GlobalStore *Store

func InitializeStore() *Store {
	x := Store{
		stringToIdentifier:      *treemap.NewWithStringComparator(),
		identifierToString:      *treemap.NewWithIntComparator(),
		nextAvailableIdentifier: 1,
		mutex:                   sync.RWMutex{},
	}

	return &x
}

func Initialize() {
	GlobalStore = InitializeStore()
}

func (t *Store) Intern(statement string) (int, error) {
	// check if the string already exits or not
	t.mutex.RLock()
	id, exist := t.stringToIdentifier.Get(statement)
	t.mutex.RUnlock()

	if exist {
		// this statement already exists, return its id
		return id.(int), nil
	}

	// not exits

	// lock - resource being written so enforced a write lock
	t.mutex.Lock()
	defer t.mutex.Unlock() // unlock happens even if there's a panic

	// check if some other go routine added it during this phase we may skip
	id, ext := t.stringToIdentifier.Get(statement)

	if ext {
		return id.(int), nil
	}

	t.stringToIdentifier.Put(statement, t.nextAvailableIdentifier)
	t.identifierToString.Put(t.nextAvailableIdentifier, statement)
	t.nextAvailableIdentifier = t.nextAvailableIdentifier + 1

	return t.nextAvailableIdentifier - 1, nil
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
