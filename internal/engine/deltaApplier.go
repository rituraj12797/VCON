package engine

import (
	"vcon/internal/schema"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
)

// applier transformative Delta to version x to et version x + 1
// return the new array of identifiers which is content of version x + 1
func ApplyDelta(parent []string, delta []schema.DeltaInstruction) []string {

	var addMap treemap.Map = *treemap.NewWithIntComparator() // int vs vector ( of insertions )
	var deleteSet treeset.Set = *treeset.NewWithIntComparator()

	for j := 0; j <= len(parent); j++ {
		addMap.Put(j, []string{}) // initialize with empty slices
	}

	for j := range delta {
		if delta[j].DeltaType == 1 { // 1 means add
			line := delta[j].Line
			vec, found := addMap.Get(line)
			var arr []string
			if found {
				arr = vec.([]string)
			}
			arr = append(arr, delta[j].Val)
			addMap.Put(line, arr) // map of this line ==> but wouldn;t copying and then putting cost extra time ====> ay optimisation for this ??
		} else {
			deleteSet.Add(delta[j].Line) // this line needs to be deleted
		}
	}

	var result []string

	for j := 0; j <= len(parent); j++ {
		// add insetion before accessing this index
		if additions, found := addMap.Get(j); found { // if found then go
			insertions := additions.([]string)
			result = append(result, insertions...) // append insertions in result
		}

		// ignore instead of wrting and then deleting
		if j < len(parent) && deleteSet.Contains(j) {
			continue // simple ignore as ritu raj says
		}

		// add original line as it was not deleted
		if j < len(parent) {
			result = append(result, parent[j])
		}
	}

	return result

}
