package antfarm

import (
	"strings"
	"testing"

	"test/models"
)

func TestAntFarm_SimulateMovement(t *testing.T) {
	// Helper functions
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

	createRoom := func(name string, isStart, isEnd bool) *models.Room {
		return &models.Room{
			Name:      name,
			IsStart:   isStart,
			IsEnd:     isEnd,
			Connected: make([]*models.Room, 0),
		}
	}

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
	// Test cases
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Single ant with direct path",
			fields: fields{
				NumAnts: 1,
				Ants:    createAnts(1),
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", true, false)
					end := createRoom("end", false, true)
					connectRooms(start, end)

					return map[string]*models.Room{
						"start": start,
						"end":   end,
					}
				}(),
				Start: createRoom("start", true, false),
				End:   createRoom("end", false, true),
			},
			want:    "L1-end", // Path assignments should be logged as L1-end
			wantErr: false,
		},
		{
			name: "Multiple ants with one path",
			fields: fields{
				NumAnts: 3,
				Ants:    createAnts(3),
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", true, false)
					mid := createRoom("mid", false, false)
					end := createRoom("end", false, true)

					connectRooms(start, mid)
					connectRooms(mid, end)

					return map[string]*models.Room{
						"start": start,
						"mid":   mid,
						"end":   end,
					}
				}(),
				Start: createRoom("start", true, false),
				End:   createRoom("end", false, true),
			},
			want:    "L1-mid\nL1-end L2-mid\nL2-end L3-mid\nL3-end", // Ants should move across the path in the order
			wantErr: false,
		},
		{
			name: "No path between start and end",
			fields: fields{
				NumAnts: 2,
				Ants:    createAnts(2),
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", true, false)
					end := createRoom("end", false, true)

					return map[string]*models.Room{
						"start": start,
						"end":   end,
					}
				}(),
				Start: createRoom("start", true, false),
				End:   createRoom("end", false, true),
			},
			want:    "", // No paths should be assigned due to disconnected rooms
			wantErr: true,
		},
		{
			name: "Complex path with multiple routes",
			fields: fields{
				NumAnts: 2,
				Ants:    createAnts(2),
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", true, false)
					mid1 := createRoom("mid1", false, false)
					mid2 := createRoom("mid2", false, false)
					end := createRoom("end", false, true)

					connectRooms(start, mid1)
					connectRooms(start, mid2)
					connectRooms(mid1, end)
					connectRooms(mid2, end)

					return map[string]*models.Room{
						"start": start,
						"mid1":  mid1,
						"mid2":  mid2,
						"end":   end,
					}
				}(),
				Start: createRoom("start", true, false),
				End:   createRoom("end", false, true),
			},
			want:    "L1-mid1 L2-mid2\nL1-end L2-end", // Ants should follow paths to mid1 and mid2, then end
			wantErr: false,
		},
		{
			name: "Edge case: No ants",
			fields: fields{
				NumAnts: 0,
				Ants:    []*models.Ant{},
				Rooms: func() map[string]*models.Room {
					start := createRoom("start", true, false)
					end := createRoom("end", false, true)
					connectRooms(start, end)

					return map[string]*models.Room{
						"start": start,
						"end":   end,
					}
				}(),
				Start: createRoom("start", true, false),
				End:   createRoom("end", false, true),
			},
			want:    "", // No ants to simulate, no paths expected
			wantErr: true,
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

			// Call SimulateMovement and capture the result and error
			got, err := af.SimulateMovement()

			// Check for errors
			if (err != nil) != tt.wantErr {
				t.Errorf("\nTest: %s\nExpected error: %v, got: %v",
					tt.name,
					tt.wantErr,
					err != nil)
				return
			}

			// Compare the resulting path with the expected path
			if len(strings.Split(strings.TrimSpace(got), " ")) != len(strings.Split(strings.TrimSpace(tt.want), " ")) {
				t.Errorf("\nTest: %s\ngot = %v\nwant = %v",
					tt.name,
					got,
					tt.want)
			}
		})
	}
}
