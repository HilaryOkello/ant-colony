package antfarm

import (
	"testing"

	"test/models"
)

func TestAntFarm_dfs(t *testing.T) {
	type fields struct {
		NumAnts int
		Ants    []*models.Ant
		Rooms   map[string]*models.Room
		Start   *models.Room
		End     *models.Room
	}
	type args struct {
		current *models.Room
		end     *models.Room
		visited map[*models.Room]bool
		path    []*models.Room
		paths   *[]models.Path
	}

	// Helper function to create rooms and connections
	createRooms := func(roomNames []string, connections map[string][]string) map[string]*models.Room {
		rooms := make(map[string]*models.Room)

		// Create rooms
		for _, name := range roomNames {
			rooms[name] = &models.Room{
				Name:      name,
				Connected: make([]*models.Room, 0),
			}
		}

		// Add connections
		for roomName, connectedRooms := range connections {
			for _, connectedRoom := range connectedRooms {
				rooms[roomName].Connected = append(rooms[roomName].Connected, rooms[connectedRoom])
			}
		}

		return rooms
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		expected int // Expected number of paths
	}{
		{
			name: "Single path",
			fields: fields{
				NumAnts: 1,
				Rooms: createRooms(
					[]string{"start", "middle", "end"},
					map[string][]string{
						"start":  {"middle"},
						"middle": {"end"},
					},
				),
			},
			args: args{
				visited: make(map[*models.Room]bool),
				path:    make([]*models.Room, 0),
				paths:   &[]models.Path{},
			},
			expected: 1,
		},
		{
			name: "Multiple paths",
			fields: fields{
				NumAnts: 1,
				Rooms: createRooms(
					[]string{"start", "a", "b", "end"},
					map[string][]string{
						"start": {"a", "b"},
						"a":     {"end"},
						"b":     {"end"},
					},
				),
			},
			args: args{
				visited: make(map[*models.Room]bool),
				path:    make([]*models.Room, 0),
				paths:   &[]models.Path{},
			},
			expected: 2,
		},
		{
			name: "Cycle in graph",
			fields: fields{
				NumAnts: 1,
				Rooms: createRooms(
					[]string{"start", "a", "b", "end"},
					map[string][]string{
						"start": {"a"},
						"a":     {"b", "end"},
						"b":     {"a", "end"},
					},
				),
			},
			args: args{
				visited: make(map[*models.Room]bool),
				path:    make([]*models.Room, 0),
				paths:   &[]models.Path{},
			},
			expected: 2,
		},
		{
			name: "No path",
			fields: fields{
				NumAnts: 1,
				Rooms: createRooms(
					[]string{"start", "a", "b", "end"},
					map[string][]string{
						"start": {"a"},
						"b":     {"end"},
					},
				),
			},
			args: args{
				visited: make(map[*models.Room]bool),
				path:    make([]*models.Room, 0),
				paths:   &[]models.Path{},
			},
			expected: 0,
		},
		{
			name: "Complex graph",
			fields: fields{
				NumAnts: 1,
				Rooms: createRooms(
					[]string{"start", "a", "b", "c", "d", "end"},
					map[string][]string{
						"start": {"a", "b"},
						"a":     {"c", "d"},
						"b":     {"c", "d"},
						"c":     {"end"},
						"d":     {"end"},
					},
				),
			},
			args: args{
				visited: make(map[*models.Room]bool),
				path:    make([]*models.Room, 0),
				paths:   &[]models.Path{},
			},
			expected: 4,
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

			// Initialize args with proper start and end rooms
			tt.args.current = af.Start
			tt.args.end = af.End

			// Run DFS
			af.dfs(tt.args.current, tt.args.end, tt.args.visited, tt.args.path, tt.args.paths)

			// Check number of paths found
			if len(*tt.args.paths) != tt.expected {
				t.Errorf("dfs() found %v paths, expected %v paths", len(*tt.args.paths), tt.expected)
			}

			// Verify that each path is valid
			for i, path := range *tt.args.paths {
				// Check if path starts at start room and ends at end room
				if path.Rooms[0] != af.Start {
					t.Errorf("Path %d does not start at start room", i)
				}
				if path.Rooms[len(path.Rooms)-1] != af.End {
					t.Errorf("Path %d does not end at end room", i)
				}

				// Check if path length is correct
				if path.Length != len(path.Rooms)-1 {
					t.Errorf("Path %d length %d does not match number of rooms %d", i, path.Length, len(path.Rooms)-1)
				}

				// Verify connections between consecutive rooms
				for j := 0; j < len(path.Rooms)-1; j++ {
					current := path.Rooms[j]
					next := path.Rooms[j+1]
					found := false
					for _, connected := range current.Connected {
						if connected == next {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Path %d: Room %s is not connected to %s", i, current.Name, next.Name)
					}
				}
			}
		})
	}
}
