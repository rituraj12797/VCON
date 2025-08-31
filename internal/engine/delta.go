package engine

type DeltaType string

const (
	A DeltaType = "A"
	D DeltaType = "D"
)

type DeltaInstruction struct {
	DeltaType DeltaType
	Line      int // the line of parrent array which wil be affected by it
	Val       int // identifier of data being added or removed
}
