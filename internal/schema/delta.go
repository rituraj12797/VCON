package schema

type DeltaType int

const (
	A DeltaType = 1
	D DeltaType = 0
)

type DeltaInstruction struct {
	DeltaType DeltaType
	Line      int    // the line of parrent array which wil be affected by it
	Val       string // SHA-256 identifier of data being added or removed
}
