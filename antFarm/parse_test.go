package antfarm

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"test/models"
)

func TestParseNumAnts(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"valid number of ants", "10", 10, false},
		{"zero ants", "0", 0, true},
		{"negative ants", "-5", 0, true},
		{"too many ants", "10001", 0, true},
		{"non-numeric input", "abc", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			af := &AntFarm{}
			state := &parserState{
				scanner: bufio.NewScanner(strings.NewReader(tc.input)),
			}

			err := af.parseNumAnts(state)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseNumAnts() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if af.NumAnts != tc.want {
				t.Errorf("ParseNumAnts() got = %d, want %d", af.NumAnts, tc.want)
			}
		})
	}
}

func TestParseRoomsAndLinks(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid rooms and links", `Room1 0 0
Room2 1 1
Room1-Room2`, false},
		{"duplicate room name", `Room1 0 0
Room1 1 1
Room1-Room2`, true},
		{"invalid room format", `Room1 0
Room2 1 1
Room1-Room2`, true},
		{"invalid link format", `Room1 0 0
Room2 1 1
Room1-Room2-Room3`, true},
		{"link references nonexistent room", `Room1 0 0
Room2 1 1
Room1-Room3`, true},
		{"duplicate link", `Room1 0 0
Room2 1 1
Room1-Room2
Room1-Room2`, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			af := &AntFarm{
				Rooms: make(map[string]*models.Room),
			}
			state := &parserState{
				scanner: bufio.NewScanner(strings.NewReader(tc.input)),
			}

			err := af.parseRoomsAndLinks(state)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseRoomsAndLinks() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestHandleComment(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *parserState
	}{
		{"##start comment", "##start", &parserState{expectStart: true, expectEnd: false, parsingLinks: false}},
		{"##end comment", "##end", &parserState{expectStart: false, expectEnd: true, parsingLinks: false}},
		{"other comment", "# this is a comment", &parserState{expectStart: false, expectEnd: false, parsingLinks: false}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			af := &AntFarm{}
			state := &parserState{}
			af.handleComment(tc.input, state)

			if state.expectStart != tc.want.expectStart || state.expectEnd != tc.want.expectEnd || state.parsingLinks != tc.want.parsingLinks {
				t.Errorf("HandleComment() got = %+v, want %+v", state, tc.want)
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
		setup   func(af *AntFarm) // Setup function to prepare AntFarm for test
	}{
		{
			name:    "valid room definition",
			input:   "Room1 0 0",
			wantErr: false,
			setup:   func(af *AntFarm) {}, // No setup needed
		},
		{
			name:    "valid link definition",
			input:   "Room1-Room2",
			wantErr: false,
			setup: func(af *AntFarm) { // Add Room1 and Room2 before testing the link
				af.Rooms["Room1"] = &models.Room{Name: "Room1"}
				af.Rooms["Room2"] = &models.Room{Name: "Room2"}
			},
		},
		{
			name:    "invalid room definition",
			input:   "Room1 0",
			wantErr: true,
			setup:   func(af *AntFarm) {}, // No setup needed
		},
		{
			name:    "invalid link definition",
			input:   "Room1-Room2-Room3",
			wantErr: true,
			setup:   func(af *AntFarm) {}, // No setup needed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			af := &AntFarm{
				Rooms: make(map[string]*models.Room),
			}
			state := &parserState{}

			// Apply the setup function to prepare the AntFarm as needed
			tc.setup(af)

			err := af.parseLine(tc.input, state)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name    string
		rooms   map[string]*models.Room
		start   *models.Room
		end     *models.Room
		wantErr bool
	}{
		{"valid configuration", map[string]*models.Room{
			"Start": {IsStart: true},
			"End":   {IsEnd: true},
		}, &models.Room{IsStart: true}, &models.Room{IsEnd: true}, false},
		{"missing start room", map[string]*models.Room{
			"End": {IsEnd: true},
		}, nil, &models.Room{IsEnd: true}, true},
		{"missing end room", map[string]*models.Room{
			"Start": {IsStart: true},
		}, &models.Room{IsStart: true}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			af := &AntFarm{
				Rooms: tc.rooms,
				Start: tc.start,
				End:   tc.end,
			}

			err := af.validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestParseRoom(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid room definition", "Room1 0 0", false},
		{"invalid room format", "Room1 0", true},
		{"duplicate room name", "Room1 0 0\nRoom1 1 1", true},
		{"invalid room coordinates", "Room1 abc 123", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			af := &AntFarm{
				Rooms: make(map[string]*models.Room),
			}
			state := &parserState{}

			err := af.parseRoom(tc.input, state)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseRoom() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestParseRoomDefinition(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    *models.Room
		wantErr bool
	}{
		{"valid room definition", "Room1 0 0", &models.Room{
			Name:      "Room1",
			X:         0,
			Y:         0,
			IsStart:   false,
			IsEnd:     false,
			Connected: make([]*models.Room, 0),
		}, false},
		{"invalid room format", "Room1 0", nil, true},
		{"invalid room coordinates", "Room1 abc 123", nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			af := &AntFarm{}
			state := &parserState{}

			got, err := af.parseRoomDefinition(tc.input, state)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseRoomDefinition() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseRoomDefinition() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAddRoom(t *testing.T) {
	testCases := []struct {
		name    string
		rooms   map[string]*models.Room
		start   *models.Room
		end     *models.Room
		room    *models.Room
		wantErr bool
	}{
		{"add valid room", map[string]*models.Room{}, nil, nil, &models.Room{Name: "Room1"}, false},
		{"add start room", map[string]*models.Room{}, nil, nil, &models.Room{Name: "Start", IsStart: true}, false},
		{"add end room", map[string]*models.Room{}, nil, nil, &models.Room{Name: "End", IsEnd: true}, false},
		{"multiple start rooms", map[string]*models.Room{"Start1": {IsStart: true}}, &models.Room{IsStart: true}, nil, &models.Room{Name: "Start2", IsStart: true}, true},
		{"multiple end rooms", map[string]*models.Room{"End1": {IsEnd: true}}, nil, &models.Room{IsEnd: true}, &models.Room{Name: "End2", IsEnd: true}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			af := &AntFarm{
				Rooms: tc.rooms,
				Start: tc.start,
				End:   tc.end,
			}

			err := af.addRoom(tc.room)
			if (err != nil) != tc.wantErr {
				t.Errorf("AddRoom() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				if tc.room.Name != "" {
					if _, ok := af.Rooms[tc.room.Name]; !ok {
						t.Errorf("AddRoom() did not add the room to the Rooms map")
					}
				}

				if tc.room.IsStart && af.Start != tc.room {
					t.Errorf("AddRoom() did not set the Start room correctly")
				}

				if tc.room.IsEnd && af.End != tc.room {
					t.Errorf("AddRoom() did not set the End room correctly")
				}
			}
		})
	}
}

func TestParseLink(t *testing.T) {
	testCases := []struct {
		name    string
		rooms   map[string]*models.Room
		input   string
		wantErr bool
	}{
		{"valid link", map[string]*models.Room{
			"Room1": {Name: "Room1", Connected: make([]*models.Room, 0)},
			"Room2": {Name: "Room2", Connected: make([]*models.Room, 0)},
		}, "Room1-Room2", false},
		{"invalid link format", map[string]*models.Room{
			"Room1": {Name: "Room1", Connected: make([]*models.Room, 0)},
		}, "Room1-Room2-Room3", true},
		{"link references nonexistent room", map[string]*models.Room{
			"Room1": {Name: "Room1", Connected: make([]*models.Room, 0)},
		}, "Room1-Room2", true},
		{"duplicate link", map[string]*models.Room{
			"Room1": {Name: "Room1", Connected: []*models.Room{{Name: "Room2"}}},
			"Room2": {Name: "Room2", Connected: []*models.Room{{Name: "Room1"}}},
		}, "Room1-Room2", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			af := &AntFarm{
				Rooms: tc.rooms,
			}
			err := af.parseLink(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseLink() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				room1 := af.Rooms["Room1"]
				room2 := af.Rooms["Room2"]
				if !contains(room1.Connected, room2) || !contains(room2.Connected, room1) {
					t.Errorf("ParseLink() did not correctly add the link between rooms")
				}
			}
		})
	}
}

func contains(rooms []*models.Room, room *models.Room) bool {
	for _, r := range rooms {
		if r == room {
			return true
		}
	}
	return false
}
