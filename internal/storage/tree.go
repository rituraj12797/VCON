package storage

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
)

type Tree struct {
	tree       []Node
	versionMap *treemap.Map // version vs node number
}




// ONE thign to note is that the DB wont store the childArray in DB to save space
//  But while building the tree here in RAM we will store the child array as it will be used t o iterate through our path in get version X 

// 
func NewTree() *Tree { // returns a reference to a new tree
	var newtre Tree
	newtre.tree = append(newtre.tree, Node{}) // empty node this will be the 0th node we are doing 1 based indexing so root will be 1st index hence starter is 0th node to avoid any confusion
	newtre.versionMap = treemap.NewWithStringComparator()
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

	t.versionMap.Put(nodeVersion, nodeNumber) // this version corrosponds to this node number

	t.tree = append(t.tree, Node{}) // add a empty child and then modify values for the child

	parent := &t.tree[parentNode]
	childDepth := parent.depth + 1

	thresholdDepth := 3

	childNode := Node{

		number:  nodeNumber,
		version: nodeVersion,

		lastSnapshotAncestor: func() int {
			if childDepth%thresholdDepth == 0 || nodeNumber == 1 {
				return nodeNumber
			}
			return parent.lastSnapshotAncestor
		}(),

		depth: childDepth,

		nodeType: func() DataType {
			if childDepth%thresholdDepth == 0 || nodeNumber == 1 {
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

func (t *Tree) GetVersionX(version string) {
	// this will return a interface type data

	// main logic revolves around doing a dfs from ancestor to child

	// now finding this path from ancestor to child will be dificult as number of child could be > 1 so high number of paths how to find the one whihc connects ancestor snapshto ndoe to this one ??

	// answer => from child traverse it's parent if parent = LSA good else go to it's parent
	// this way we will find the path from chiild to ancestor correclty without exploring much ( finding path in constant time )

	// for now this function will pritn the path from the LSA to this version

	// checks
	targetNodenumber, found := t.versionMap.Get(version)

	if !found {
		fmt.Println(" version not found ")
		return
	}

	// if exists means it is valid
	nodeNum := targetNodenumber.(int) // since the map is of type interface {} vs interface {} type convert before using thi value
	lsa := t.tree[nodeNum].lastSnapshotAncestor

	curNode := nodeNum
	// path must have both lsa and curNode

	var path []int
	path = append(path, nodeNum)

	for j := curNode; j != lsa; {
		path = append(path, t.tree[j].parentNumber)
		j = t.tree[j].parentNumber
	}

	fmt.Println(" nodeNum : ", nodeNum)
	fmt.Println(" lsa : ", lsa)

	fmt.Println(" path : ", path)
	fmt.Println("=============")

	// path reversal 
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
    path[i], path[j] = path[j], path[i]
}

	var data string;
	
	for j := 0; j< len(path); j++ {
		data += t.tree[path[j]].data.(string)
	}
	
	fmt.Println(" The Data is : ",data);

	// if path length = 1 means this node it self is the LSA
	// if path length > 1 means someone else ( one etreme is lsa and other exreme is the nodeNum )

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
