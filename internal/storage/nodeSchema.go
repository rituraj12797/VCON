package storage

type DataType string 

const (
	Delta  DataType = "delta"
	Snapshot DataType = "snapshot"
)

type Node struct {
	number int;  // the node number
	version int;  // the version this node corrsponds to 
	lastSnapshotAncestor int; // ths last ancestor which was a snapshot node 

	nodeType DataType; // snapshot or delta node
	data interface{}; // holding the actual data which the node will store diffs or snapshot

	parentNumber int;
	childNodes []int;
}
