package schema

import (
	"time"
	"vcon/internal/engine"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeType int

const (
	NodeTypeDelta    NodeType = 2
	NodeTypeSnapshot NodeType = 1
)

type Node struct { // this is not a Document which will be stored in MongoDB

	ParrentNode          int      `bson:"parent_id"` // parrent of this node
	LastSnapshotAncestor int      `bson:"lsa_id"`    // LSA of this node
	NodeNumber           int      `bson:"node_id"`   // number of this node in version tree
	Depth 				 int 	  `bson:"depth"`
	NodeType             NodeType `bson:"node_type"` // snapshot or delta node
	
	VersionString        string   `bson:"version_name"` // say uset gives this name as " version1.1 " will be used in mapping with ndoe number

	// ChildrenArray []int // this array contains al the children of this node
	// Not needed as tree could be constructed at run time and this array could be made using parrent pointer to create 

	/*
	You fetch the flat list of nodes from the database, and then your Go application quickly loops through them at runtime to build the tree structure in memory by linking children to their parents via the ParrentNode ID.
	*/ 

	DeltaInstructions []engine.DeltaInstruction `bson:"delta_instructions,omitempty"` // for delta nodes [ {0/1 ( Add/ DEL) , X ( Line number ), YYYY ( identifier of val) },....]
	FileArray         []string   `bson:"file_content,omitempty"`                 // this file array is for only Snapshot nodes and thic contains identifiers (SHa 56 of this statement ) of 

}

type Document struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Title         string             `bson:"title"`
	NumberOfNodes int                `bson:"numberOfNodes"`
	NodeArray     []Node             `bson:"nodes"` // the node array 0 represents pseudo root  1 represents actual root

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
