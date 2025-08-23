package storage

import "fmt"

type Tree struct {
	tree []Node
}

func NewTree() *Tree { // returns a reference to a new tree
	var newtre Tree
	newtre.tree = append(newtre.tree, Node{}) // empty node this will be the 0th node we are doing 1 based indexing so root will be 1st index hence starter is 0th node to avoid any confusion
	return &newtre
}

func (t *Tree) AddNode(parentNode int, nodeNumber int, nodeVersion string, childData interface{}) error {

	// considering sequential linesr adding of node i.e 26 will be aded if 1 to 25 all are aded
	// adding a node

	// verify is tree exists or not
	if len(t.tree) == 0 {
		return fmt.Errorf(" tree does not exists ")
	}

	// verify if parent exists or not ??

	if parentNode >= len(t.tree) || parentNode < 0 { // parent node =0 means first node is being aded so thats cmpletely fine
		return fmt.Errorf(" invalid parent ")

	}

	// verify the size of tree it must be = nodenumber
	// say i am adding 26th node
	// means already 25 nodes should be there +  a 0th node whihc s our starter( not the root )

	if len(t.tree) != nodeNumber {

		if len(t.tree) > nodeNumber {
			return fmt.Errorf(" node already exists ")
		} else {
			return fmt.Errorf(" smaller nodes are not added add them first ")
		}
	}

	t.tree = append(t.tree, Node{}) // add a empty child and then modify values for the child

	parent := &t.tree[parentNode]
	childDepth := parent.depth + 1

	childNode := Node{

		number:  nodeNumber,
		version: nodeVersion,

		lastSnapshotAncestor: func() int {
			if childDepth%10 == 0 {
				return nodeNumber
			}
			return parent.lastSnapshotAncestor
		}(),

		depth: childDepth,

		nodeType: func() DataType {
			if childDepth%10 == 0 {
				return Snapshot
			}
			return Delta
		}(),

		data: childData,

		parentNumber: parentNode,

		childNodes: make([]int, 0), // no children as of now
	}

	t.tree[nodeNumber] = childNode
	t.tree[parentNode].childNodes = append(t.tree[parentNode].childNodes, nodeNumber)

	fmt.Println(" Child Added to Parrent in tree ")
	fmt.Println(" Parent : ", parentNode, " Child list : ", parent.childNodes)
	fmt.Println(" ======= CHILDREN ======= ")
	fmt.Println(childNode)
	fmt.Println("========================= ")

	return nil
}

func (t *Tree) ShowTree() {

	for i := 1; i < len(t.tree); i++ {
		fmt.Println(" Node : ", i)
		fmt.Println(" version : ", t.tree[i].version)
		fmt.Println(" depth : ", t.tree[i].depth)
		fmt.Println(" nodeType : ", t.tree[i].nodeType)
		fmt.Println(" data : ", t.tree[i].data)
		fmt.Println(" Parrent : ", t.tree[i].parentNumber)
		fmt.Println(" Children : ", t.tree[i].childNodes)
		fmt.Println("")
	}

}
