package antfarm

import (
	"fmt"
	"strings"
	"testing"

	"test/models"
)

func TestAntFarm_findAllPaths(t *testing.T) {
	// Helper function to create rooms
	createRoom := func(name string, x, y int, isStart, isEnd bool) *models.Room {
		return &models.Room{
			Name:      name,
			X:         x,
			Y:         y,
			IsStart:   isStart,
			IsEnd:     isEnd,
			Connected: make([]*models.Room, 0),
		}
	}

	// Helper function to connect rooms bidirectionally
	connectRooms := func(room1, room2 *models.Room) {
		room1.Connected = append(room1.Connected, room2)
		room2.Connected = append(room2.Connected, room1)
	}

	type fields struct {
		NumAnts int
		Ants    []*models.Ant
		Rooms   map[string]*models.Room
		Start   *models.Room
		End     *models.Room
	}

	tests := []struct {
		name   string
		fields fields
		want   []models.Path
	}{
		{
			name: "Single direct path",
			fields: fields{
				NumAnts: 1,
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", 0, 0, true, false)
					end := createRoom("end", 1, 1, false, true)
					connectRooms(start, end)
					return map[string]*models.Room{
						"start": start,
						"end":   end,
					}
				}(),
				Start: nil, // Will be set below
				End:   nil, // Will be set below
			},
			want: []models.Path{
				{
					Rooms:  []*models.Room{createRoom("start", 0, 0, true, false), createRoom("end", 1, 1, false, true)},
					Length: 1,
				},
			},
		},
		{
			name: "Two parallel paths",
			fields: fields{
				NumAnts: 2,
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", 0, 0, true, false)
					middle1 := createRoom("middle1", 1, 0, false, false)
					middle2 := createRoom("middle2", 1, 1, false, false)
					end := createRoom("end", 2, 0, false, true)

					connectRooms(start, middle1)
					connectRooms(start, middle2)
					connectRooms(middle1, end)
					connectRooms(middle2, end)

					return map[string]*models.Room{
						"start":   start,
						"middle1": middle1,
						"middle2": middle2,
						"end":     end,
					}
				}(),
				Start: nil, // Will be set below
				End:   nil, // Will be set below
			},
			want: []models.Path{
				{
					Rooms: []*models.Room{
						createRoom("start", 0, 0, true, false),
						createRoom("middle1", 1, 0, false, false),
						createRoom("end", 2, 0, false, true),
					},
					Length: 2,
				},
				{
					Rooms: []*models.Room{
						createRoom("start", 0, 0, true, false),
						createRoom("middle2", 1, 1, false, false),
						createRoom("end", 2, 0, false, true),
					},
					Length: 2,
				},
			},
		},
		{
			name: "No path available",
			fields: fields{
				NumAnts: 1,
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", 0, 0, true, false)
					end := createRoom("end", 1, 1, false, true)
					// Intentionally not connecting the rooms
					return map[string]*models.Room{
						"start": start,
						"end":   end,
					}
				}(),
				Start: nil, // Will be set below
				End:   nil, // Will be set below
			},
			want: []models.Path{}, // Empty path slice when no path exists
		},
		{
			name: "Complex path with multiple options",
			fields: fields{
				NumAnts: 3,
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", 0, 0, true, false)
					r1 := createRoom("r1", 1, 0, false, false)
					r2 := createRoom("r2", 1, 1, false, false)
					r3 := createRoom("r3", 2, 0, false, false)
					r4 := createRoom("r4", 2, 1, false, false)
					end := createRoom("end", 3, 0, false, true)

					connectRooms(start, r1)
					connectRooms(start, r2)
					connectRooms(r1, r3)
					connectRooms(r2, r4)
					connectRooms(r3, end)
					connectRooms(r4, end)
					connectRooms(r1, r4)
					connectRooms(r2, r3)

					return map[string]*models.Room{
						"start": start,
						"r1":    r1,
						"r2":    r2,
						"r3":    r3,
						"r4":    r4,
						"end":   end,
					}
				}(),
				Start: nil, // Will be set below
				End:   nil, // Will be set below
			},
			want: []models.Path{
				{
					Rooms: []*models.Room{
						createRoom("start", 0, 0, true, false),
						createRoom("r1", 1, 0, false, false),
						createRoom("r3", 2, 0, false, false),
						createRoom("end", 3, 0, false, true),
					},
					Length: 3,
				},
				{
					Rooms: []*models.Room{
						createRoom("start", 0, 0, true, false),
						createRoom("r2", 1, 1, false, false),
						createRoom("r4", 2, 1, false, false),
						createRoom("end", 3, 0, false, true),
					},
					Length: 3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set Start and End rooms from the map
			for _, room := range tt.fields.Rooms {
				if room.IsStart {
					tt.fields.Start = room
				}
				if room.IsEnd {
					tt.fields.End = room
				}
			}

			af := &AntFarm{
				NumAnts: tt.fields.NumAnts,
				Ants:    tt.fields.Ants,
				Rooms:   tt.fields.Rooms,
				Start:   tt.fields.Start,
				End:     tt.fields.End,
			}

			got := af.findAllPaths()

			// Check if the number of paths matches
			if len(got) != len(tt.want) {
				t.Errorf("AntFarm.findAllPaths() returned %v paths, want %v paths", len(got), len(tt.want))
				return
			}

			// For each found path, verify:
			// 1. It starts at the start room
			// 2. It ends at the end room
			// 3. All rooms in the path are connected
			// 4. Path length is correctly calculated
			for i, path := range got {
				// Check start room
				if !path.Rooms[0].IsStart {
					t.Errorf("Path %d does not start at start room", i)
				}

				// Check end room
				if !path.Rooms[len(path.Rooms)-1].IsEnd {
					t.Errorf("Path %d does not end at end room", i)
				}

				// Check path connectivity
				for j := 0; j < len(path.Rooms)-1; j++ {
					isConnected := false
					for _, connected := range path.Rooms[j].Connected {
						if connected == path.Rooms[j+1] {
							isConnected = true
							break
						}
					}
					if !isConnected {
						t.Errorf("Path %d has disconnected rooms at position %d", i, j)
					}
				}

				// Check path length
				if path.Length != len(path.Rooms)-1 {
					t.Errorf("Path %d length %d does not match actual length %d", i, path.Length, len(path.Rooms)-1)
				}
			}

			// Verify paths are sorted by length
			for i := 1; i < len(got); i++ {
				if got[i-1].Length > got[i].Length {
					t.Errorf("Paths are not sorted by length")
				}
			}
		})
	}
}

func TestAntFarm_filterNonOverlappingPaths(t *testing.T) {
	// Helper function to create rooms
	createRoom := func(name string, x, y int, isStart, isEnd bool) *models.Room {
		return &models.Room{
			Name:      name,
			X:         x,
			Y:         y,
			IsStart:   isStart,
			IsEnd:     isEnd,
			Connected: make([]*models.Room, 0),
		}
	}

	// Helper function to create a path
	createPath := func(rooms []*models.Room) models.Path {
		return models.Path{
			Rooms:  rooms,
			Length: len(rooms) - 1,
		}
	}

	// Common rooms for tests
	start := createRoom("start", 0, 0, true, false)
	end := createRoom("end", 5, 5, false, true)
	r1 := createRoom("r1", 1, 1, false, false)
	r2 := createRoom("r2", 2, 2, false, false)
	r3 := createRoom("r3", 3, 3, false, false)
	r4 := createRoom("r4", 4, 4, false, false)

	type fields struct {
		NumAnts int
		Ants    []*models.Ant
		Rooms   map[string]*models.Room
		Start   *models.Room
		End     *models.Room
	}
	type args struct {
		paths []models.Path
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []models.Path
	}{
		{
			name: "Empty paths",
			fields: fields{
				NumAnts: 1,
				Start:   start,
				End:     end,
			},
			args: args{
				paths: []models.Path{},
			},
			want: []models.Path{},
		},
		{
			name: "Single path",
			fields: fields{
				NumAnts: 1,
				Start:   start,
				End:     end,
			},
			args: args{
				paths: []models.Path{
					createPath([]*models.Room{start, r1, end}),
				},
			},
			want: []models.Path{
				createPath([]*models.Room{start, r1, end}),
			},
		},
		{
			name: "Two non-overlapping paths",
			fields: fields{
				NumAnts: 2,
				Start:   start,
				End:     end,
			},
			args: args{
				paths: []models.Path{
					createPath([]*models.Room{start, r1, end}),
					createPath([]*models.Room{start, r2, end}),
				},
			},
			want: []models.Path{
				createPath([]*models.Room{start, r1, end}),
				createPath([]*models.Room{start, r2, end}),
			},
		},
		{
			name: "Two overlapping paths",
			fields: fields{
				NumAnts: 2,
				Start:   start,
				End:     end,
			},
			args: args{
				paths: []models.Path{
					createPath([]*models.Room{start, r1, r2, end}),
					createPath([]*models.Room{start, r2, r3, end}),
				},
			},
			want: []models.Path{
				createPath([]*models.Room{start, r1, r2, end}),
			},
		},
		{
			name: "Complex overlapping paths",
			fields: fields{
				NumAnts: 3,
				Start:   start,
				End:     end,
			},
			args: args{
				paths: []models.Path{
					createPath([]*models.Room{start, r1, r2, end}), // Path 1
					createPath([]*models.Room{start, r2, r3, end}), // Path 2 (overlaps with 1)
					createPath([]*models.Room{start, r3, r4, end}), // Path 3
					createPath([]*models.Room{start, r1, r3, end}), // Path 4 (overlaps with 1 and 3)
				},
			},
			want: []models.Path{
				createPath([]*models.Room{start, r1, r2, end}),
				createPath([]*models.Room{start, r3, r4, end}),
			},
		},
		{
			name: "Multiple paths with shared start/end",
			fields: fields{
				NumAnts: 4,
				Start:   start,
				End:     end,
			},
			args: args{
				paths: []models.Path{
					createPath([]*models.Room{start, r1, end}),
					createPath([]*models.Room{start, r2, end}),
					createPath([]*models.Room{start, r3, end}),
					createPath([]*models.Room{start, r4, end}),
				},
			},
			want: []models.Path{
				createPath([]*models.Room{start, r1, end}),
				createPath([]*models.Room{start, r2, end}),
				createPath([]*models.Room{start, r3, end}),
				createPath([]*models.Room{start, r4, end}),
			},
		},
		{
			name: "Paths with different lengths",
			fields: fields{
				NumAnts: 2,
				Start:   start,
				End:     end,
			},
			args: args{
				paths: []models.Path{
					createPath([]*models.Room{start, r1, r2, r3, end}),
					createPath([]*models.Room{start, r4, end}),
				},
			},
			want: []models.Path{
				createPath([]*models.Room{start, r1, r2, r3, end}),
				createPath([]*models.Room{start, r4, end}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := &AntFarm{
				NumAnts: tt.fields.NumAnts,
				Ants:    tt.fields.Ants,
				Rooms:   tt.fields.Rooms,
				Start:   tt.fields.Start,
				End:     tt.fields.End,
			}

			got := af.filterNonOverlappingPaths(tt.args.paths)

			// Check if number of paths matches
			if len(got) != len(tt.want) {
				t.Errorf("filterNonOverlappingPaths() returned %v paths, want %v paths", len(got), len(tt.want))
				return
			}

			// Verify each path in the result:
			// 1. Has no overlapping rooms with other paths (except start/end)
			// 2. Matches expected path structure
			for i := 0; i < len(got); i++ {
				// Check for overlaps with other paths
				for j := i + 1; j < len(got); j++ {
					rooms1 := make(map[string]struct{})
					for k := 1; k < len(got[i].Rooms)-1; k++ {
						rooms1[got[i].Rooms[k].Name] = struct{}{}
					}

					for k := 1; k < len(got[j].Rooms)-1; k++ {
						if _, exists := rooms1[got[j].Rooms[k].Name]; exists {
							t.Errorf("Found overlapping rooms between paths %d and %d", i, j)
						}
					}
				}

				// Verify path structure
				if !got[i].Rooms[0].IsStart {
					t.Errorf("Path %d does not start with start room", i)
				}
				if !got[i].Rooms[len(got[i].Rooms)-1].IsEnd {
					t.Errorf("Path %d does not end with end room", i)
				}
				if got[i].Length != len(got[i].Rooms)-1 {
					t.Errorf("Path %d length %d does not match actual length %d", i, got[i].Length, len(got[i].Rooms)-1)
				}
			}
		})
	}
}

// Helper function to create a string representation of the ant-path map
func formatAntPathMap(m map[*models.Ant]models.Path) string {
	var result strings.Builder
	result.WriteString("map[\n")

	// Convert map to sorted slice for consistent output
	var antPaths []string
	for ant, path := range m {
		// Collect room names from the path
		var roomNames []string
		for _, room := range path.Rooms {
			roomNames = append(roomNames, room.Name)
		}

		antPath := fmt.Sprintf("    Ant(ID: %d): Path[%s]",
			ant.Id,
			strings.Join(roomNames, " -> "))
		antPaths = append(antPaths, antPath)
	}

	result.WriteString(strings.Join(antPaths, "\n"))
	result.WriteString("\n]")
	return result.String()
}

func TestAntFarm_assignAntsToPath(t *testing.T) {
	// Helper function to create ants
	createAnts := func(count int) []*models.Ant {
		ants := make([]*models.Ant, count)
		for i := 0; i < count; i++ {
			ants[i] = &models.Ant{
				Id:         i + 1,
				HasReached: false,
				PathIndex:  0,
			}
		}
		return ants
	}

	// Helper function to create a simple room
	createRoom := func(name string, x, y int, isStart, isEnd bool) *models.Room {
		return &models.Room{
			Name:      name,
			X:         x,
			Y:         y,
			IsStart:   isStart,
			IsEnd:     isEnd,
			Connected: make([]*models.Room, 0),
		}
	}

	// Helper function to connect rooms
	connectRooms := func(room1, room2 *models.Room) {
		room1.Connected = append(room1.Connected, room2)
		room2.Connected = append(room2.Connected, room1)
	}

	type fields struct {
		NumAnts int
		Ants    []*models.Ant
		Rooms   map[string]*models.Room
		Start   *models.Room
		End     *models.Room
	}

	tests := []struct {
		name   string
		fields fields
		want   map[*models.Ant]models.Path
	}{
		{
			name: "Empty farm - no ants",
			fields: fields{
				NumAnts: 0,
				Ants:    []*models.Ant{},
				Rooms:   nil,
				Start:   nil,
				End:     nil,
			},
			want: nil,
		},
		{
			name: "Single path - one ant",
			fields: fields{
				NumAnts: 1,
				Ants:    createAnts(1),
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", 0, 0, true, false)
					end := createRoom("end", 1, 1, false, true)
					connectRooms(start, end)

					return map[string]*models.Room{
						"start": start,
						"end":   end,
					}
				}(),
				Start: createRoom("start", 0, 0, true, false),
				End:   createRoom("end", 1, 1, false, true),
			},
			want: func() map[*models.Ant]models.Path {
				start := createRoom("start", 0, 0, true, false)
				end := createRoom("end", 1, 1, false, true)
				return map[*models.Ant]models.Path{
					{Id: 1}: {
						Rooms:  []*models.Room{start, end},
						Length: 1,
					},
				}
			}(),
		},
		{
			name: "Multiple paths - multiple ants",
			fields: fields{
				NumAnts: 3,
				Ants:    createAnts(3),
				Rooms: func() map[string]*models.Room {
					rooms := make(map[string]*models.Room)
					start := createRoom("start", 0, 0, true, false)
					middle1 := createRoom("middle1", 1, 0, false, false)
					middle2 := createRoom("middle2", 1, 1, false, false)
					end := createRoom("end", 2, 0, false, true)

					connectRooms(start, middle1)
					connectRooms(start, middle2)
					connectRooms(middle1, end)
					connectRooms(middle2, end)

					rooms["start"] = start
					rooms["middle1"] = middle1
					rooms["middle2"] = middle2
					rooms["end"] = end
					return rooms
				}(),
				Start: createRoom("start", 0, 0, true, false),
				End:   createRoom("end", 2, 0, false, true),
			},
			want: func() map[*models.Ant]models.Path {
				ants := createAnts(3)
				path1 := models.Path{
					Rooms: []*models.Room{
						createRoom("start", 0, 0, true, false),
						createRoom("middle1", 1, 0, false, false),
						createRoom("end", 2, 0, false, true),
					},
					Length: 2,
				}
				path2 := models.Path{
					Rooms: []*models.Room{
						createRoom("start", 0, 0, true, false),
						createRoom("middle1", 1, 0, false, false),
						createRoom("end", 2, 0, false, true),
					},
					Length: 2,
				}
				return map[*models.Ant]models.Path{
					ants[0]: path1,
					ants[1]: path2,
					ants[2]: path1,
				}
			}(),
		},
		{
			name: "No valid paths",
			fields: fields{
				NumAnts: 1,
				Ants:    createAnts(1),
				Rooms: func() map[string]*models.Room {
					rooms := make(map[string]*models.Room)
					start := createRoom("start", 0, 0, true, false)
					end := createRoom("end", 1, 0, false, true)
					// No connections between rooms
					rooms["start"] = start
					rooms["end"] = end
					return rooms
				}(),
				Start: createRoom("start", 0, 0, true, false),
				End:   createRoom("end", 1, 0, false, true),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := &AntFarm{
				NumAnts: tt.fields.NumAnts,
				Ants:    tt.fields.Ants,
				Rooms:   tt.fields.Rooms,
				Start:   tt.fields.Rooms["start"],
				End:     tt.fields.Rooms["end"],
			}

			// Call findAllPaths() to ensure paths are discovered
			// paths := af.findAllPaths()
			// if len(paths) == 0 {
			// 	t.Errorf("findAllPaths() returned 0 paths, expected at least 1 for test case %s", tt.name)
			// } else {
			// 	for _, path := range paths {
			// 		t.Logf("Discovered path: %v", path)
			// 	}
			// }

			// Call the assignAntsToPath() method
			got := af.assignAntsToPath()
			if len(got) != len(tt.want) {
				t.Errorf("\nAntFarm.assignAntsToPath()\nTest: %s\ngot = %v\nwant = %v",
					tt.name,
					formatAntPathMap(got),
					formatAntPathMap(tt.want))
			}
		})
	}
}
