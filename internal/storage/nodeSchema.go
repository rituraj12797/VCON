package storage

type DataType string 

const (
	Delta  DataType = "delta"
	Snapshot DataType = "snapshot"
)

type Node struct {
	number int;  // the node number
	version string;  // the version this node corrsponds to 
	lastSnapshotAncestor int; // ths last ancestor which was a snapshot node 

	depth int; // depth from root
	nodeType DataType; // snapshot or delta node
	data interface{}; // holding the actual data which the node will store diffs or snapshot

	// snapshots will be array of identifiers 
	// deltas will be object data containing infor about change 

	parentNumber int;
	childNodes []int;
}
