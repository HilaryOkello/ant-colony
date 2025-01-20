package models

import "fmt"

type Ant struct {
	Id          int
	CurrentRoom *Room
	Path        []*Room
	PathIndex   int
	HasReached  bool
}

// Room represents a single room in the ant farm
type Room struct {
	Name      string
	X, Y      int
	IsStart   bool
	IsEnd     bool
	Connected []*Room
	// ant       *Ant
}

// ParseError represents an error during parsing
type ParseError struct {
	Message string
}
type Path struct {
	Rooms  []*Room
	Length int
	InUse  bool
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("ERROR: invalid data format, %s", e.Message)
}
