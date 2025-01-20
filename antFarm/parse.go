package antfarm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"test/models"
)

// parserState holds the parsing state and configuration
type parserState struct {
	scanner      *bufio.Scanner
	expectStart  bool
	expectEnd    bool
	parsingLinks bool
}

// ParseInput reads and parses the input file for the ant farm configuration.
// It returns an error if the file cannot be read or if the input format is invalid.
func (af *AntFarm) ParseInput(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	state := &parserState{
		scanner: bufio.NewScanner(file),
	}

	if err := af.parseNumAnts(state); err != nil {
		return fmt.Errorf("parsing number of ants: %w", err)
	}

	if err := af.parseRoomsAndLinks(state); err != nil {
		return fmt.Errorf("parsing rooms and links: %w", err)
	}

	if err := af.validate(); err != nil {
		return fmt.Errorf("validating ant farm: %w", err)
	}

	af.initializeAnts()
	return nil
}

// parseNumAnts reads and validates the number of ants from the first line.
func (af *AntFarm) parseNumAnts(state *parserState) error {
	const maxAnts = 10000

	if !state.scanner.Scan() {
		return &models.ParseError{Message: "empty file"}
	}

	numAnts, err := strconv.Atoi(state.scanner.Text())
	if err != nil {
		return &models.ParseError{Message: "invalid number of ants"}
	}

	if numAnts <= 0 {
		return &models.ParseError{Message: "number of ants must be positive"}
	}

	if numAnts > maxAnts {
		return &models.ParseError{Message: "number of ants exceeds maximum limit"}
	}

	af.NumAnts = numAnts
	return nil
}

// parseRoomsAndLinks processes the room definitions and link configurations.
func (af *AntFarm) parseRoomsAndLinks(state *parserState) error {
	for state.scanner.Scan() {
		line := state.scanner.Text()
		if line == "" {
			continue
		}

		if isComment := strings.HasPrefix(line, "#"); isComment {
			af.handleComment(line, state)
			continue
		}

		if err := af.parseLine(line, state); err != nil {
			return err
		}

		// Reset command flags after processing a room
		if !state.parsingLinks {
			state.expectStart, state.expectEnd = false, false
		}
	}

	return state.scanner.Err()
}

// handleComment processes comment lines and updates parser state accordingly.
func (af *AntFarm) handleComment(line string, state *parserState) {
	switch line {
	case "##start":
		state.expectStart = true
	case "##end":
		state.expectEnd = true
	}
}

// parseLine handles parsing either a room or link definition.
func (af *AntFarm) parseLine(line string, state *parserState) error {
	if strings.Contains(line, "-") {
		state.parsingLinks = true
		return af.parseLink(line)
	}

	if !state.parsingLinks {
		return af.parseRoom(line, state)
	}

	return nil
}

// validate ensures the ant farm configuration is complete and valid.
func (af *AntFarm) validate() error {
	if af.Start == nil {
		return &models.ParseError{Message: "no start room found"}
	}
	if af.End == nil {
		return &models.ParseError{Message: "no end room found"}
	}
	return nil
}

// parseRoom parses a room definition line and adds the room to the ant farm.
func (af *AntFarm) parseRoom(line string, state *parserState) error {
	room, err := af.parseRoomDefinition(line, state)
	if err != nil {
		return err
	}

	if err := af.addRoom(room); err != nil {
		return err
	}
	return nil
}

// parseRoomDefinition parses the room components from a line.
func (af *AntFarm) parseRoomDefinition(line string, state *parserState) (*models.Room, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, &models.ParseError{Message: "invalid room format"}
	}

	name := parts[0]
	if _, exists := af.Rooms[name]; exists {
		return nil, &models.ParseError{Message: "duplicate room name"}
	}

	x, err1 := strconv.Atoi(parts[1])
	y, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		return nil, &models.ParseError{Message: "invalid room coordinates"}
	}

	return &models.Room{
		Name:      name,
		X:         x,
		Y:         y,
		IsStart:   state.expectStart,
		IsEnd:     state.expectEnd,
		Connected: make([]*models.Room, 0),
	}, nil
}

// addRoom adds a new room to the ant farm based on the room definition.
func (af *AntFarm) addRoom(room *models.Room) error {
	if room.IsStart {
		if af.Start != nil {
			return &models.ParseError{Message: "multiple start rooms defined"}
		}
		af.Start = room
	}

	if room.IsEnd {
		if af.End != nil {
			return &models.ParseError{Message: "multiple end rooms defined"}
		}
		af.End = room
	}

	af.Rooms[room.Name] = room
	return nil
}


// parseLink parses a link definition line
func (af *AntFarm) parseLink(line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 || parts[0] == parts[1] {
		return &models.ParseError{Message: "invalid link format"}
	}
	room1, exists1 := af.Rooms[parts[0]]
	room2, exists2 := af.Rooms[parts[1]]
	if !exists1 || !exists2 {
		return &models.ParseError{Message: "link references nonexistent room"}
	}
	// Check if link already exists
	for _, connected := range room1.Connected {
		if connected.Name == room2.Name {
			return &models.ParseError{Message: "duplicate link"}
		}
	}
	room1.Connected = append(room1.Connected, room2)
	room2.Connected = append(room2.Connected, room1)
	return nil
}



