package antfarm

import (
	"reflect"
	"testing"

	"test/models"
)

func TestNewAntFarm(t *testing.T) {
	tests := []struct {
		name string
		want *AntFarm
	}{
		{
			name: "New empty ant farm initialization",
			want: &AntFarm{
				NumAnts: 0,
				Ants:    nil,
				Rooms:   make(map[string]*models.Room),
				Start:   nil,
				End:     nil,
			},
		},
		{
			name: "Verify room map is initialized but empty",
			want: &AntFarm{
				Rooms: make(map[string]*models.Room),
			},
		},
		{
			name: "Verify default values of all fields",
			want: &AntFarm{
				NumAnts: 0,
				Ants:    []*models.Ant(nil),
				Rooms:   make(map[string]*models.Room),
				Start:   (*models.Room)(nil),
				End:     (*models.Room)(nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAntFarm()

			// Test overall structure equality
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAntFarm() = %v, want %v", got, tt.want)
			}

			// Specific tests for the Rooms map
			if got.Rooms == nil {
				t.Error("NewAntFarm() Rooms map is nil, want initialized empty map")
			}

			if len(got.Rooms) != 0 {
				t.Errorf("NewAntFarm() Rooms map has %d elements, want empty map", len(got.Rooms))
			}

			// Test map functionality
			got.Rooms["test"] = &models.Room{Name: "test"}
			if len(got.Rooms) != 1 {
				t.Error("NewAntFarm() Rooms map is not properly initialized for adding elements")
			}

			// Verify zero values for other fields
			if got.NumAnts != 0 {
				t.Errorf("NewAntFarm() NumAnts = %d, want 0", got.NumAnts)
			}

			if got.Ants != nil {
				t.Errorf("NewAntFarm() Ants = %v, want nil", got.Ants)
			}

			if got.Start != nil {
				t.Errorf("NewAntFarm() Start = %v, want nil", got.Start)
			}

			if got.End != nil {
				t.Errorf("NewAntFarm() End = %v, want nil", got.End)
			}
		})
	}

	// Test concurrent initialization
	t.Run("Concurrent initialization", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func() {
				farm := NewAntFarm()
				if farm.Rooms == nil {
					t.Error("NewAntFarm() Rooms map is nil in concurrent initialization")
				}
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	// Test memory allocation
	t.Run("Memory allocation", func(t *testing.T) {
		farms := make([]*AntFarm, 1000)
		for i := 0; i < 1000; i++ {
			farms[i] = NewAntFarm()
			if farms[i].Rooms == nil {
				t.Error("NewAntFarm() failed to allocate memory for Rooms map")
			}
		}
	})
}

func TestAntFarm_initializeAnts(t *testing.T) {
	// Define a sample room for testing
	startRoom := &models.Room{
		Name:      "start",
		X:         0,
		Y:         0,
		IsStart:   true,
		IsEnd:     false,
		Connected: make([]*models.Room, 0),
	}

	endRoom := &models.Room{
		Name:      "end",
		X:         1,
		Y:         1,
		IsStart:   false,
		IsEnd:     true,
		Connected: make([]*models.Room, 0),
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
		want   []*models.Ant
	}{
		{
			name: "Initialize zero ants",
			fields: fields{
				NumAnts: 0,
				Ants:    nil,
				Rooms: map[string]*models.Room{
					"start": startRoom,
					"end":   endRoom,
				},
				Start: startRoom,
				End:   endRoom,
			},
			want: []*models.Ant{},
		},
		{
			name: "Initialize single ant",
			fields: fields{
				NumAnts: 1,
				Ants:    nil,
				Rooms: map[string]*models.Room{
					"start": startRoom,
					"end":   endRoom,
				},
				Start: startRoom,
				End:   endRoom,
			},
			want: []*models.Ant{
				{
					Id:          1,
					CurrentRoom: startRoom,
					PathIndex:   0,
					HasReached:  false,
				},
			},
		},
		{
			name: "Initialize multiple ants",
			fields: fields{
				NumAnts: 3,
				Ants:    nil,
				Rooms: map[string]*models.Room{
					"start": startRoom,
					"end":   endRoom,
				},
				Start: startRoom,
				End:   endRoom,
			},
			want: []*models.Ant{
				{
					Id:          1,
					CurrentRoom: startRoom,
					PathIndex:   0,
					HasReached:  false,
				},
				{
					Id:          2,
					CurrentRoom: startRoom,
					PathIndex:   0,
					HasReached:  false,
				},
				{
					Id:          3,
					CurrentRoom: startRoom,
					PathIndex:   0,
					HasReached:  false,
				},
			},
		},
		{
			name: "Initialize with existing ants array",
			fields: fields{
				NumAnts: 2,
				Ants:    make([]*models.Ant, 1), // Existing array with different size
				Rooms: map[string]*models.Room{
					"start": startRoom,
					"end":   endRoom,
				},
				Start: startRoom,
				End:   endRoom,
			},
			want: []*models.Ant{
				{
					Id:          1,
					CurrentRoom: startRoom,
					PathIndex:   0,
					HasReached:  false,
				},
				{
					Id:          2,
					CurrentRoom: startRoom,
					PathIndex:   0,
					HasReached:  false,
				},
			},
		},
		{
			name: "Verify start room properties",
			fields: fields{
				NumAnts: 1,
				Ants:    nil,
				Rooms: map[string]*models.Room{
					"start": startRoom,
					"end":   endRoom,
				},
				Start: startRoom,
				End:   endRoom,
			},
			want: []*models.Ant{
				{
					Id:          1,
					CurrentRoom: startRoom,
					PathIndex:   0,
					HasReached:  false,
				},
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

			// Call the method being tested
			af.initializeAnts()

			// Verify the length of the ants array
			if len(af.Ants) != tt.fields.NumAnts {
				t.Errorf("initializeAnts() created %v ants, want %v", len(af.Ants), tt.fields.NumAnts)
			}

			// Verify each ant's properties
			for i := 0; i < tt.fields.NumAnts; i++ {
				if af.Ants[i] == nil {
					t.Errorf("initializeAnts() ant at index %v is nil", i)
					continue
				}

				// Check ID
				if af.Ants[i].Id != i+1 {
					t.Errorf("Ant[%d].Id = %v, want %v", i, af.Ants[i].Id, i+1)
				}

				// Check CurrentRoom
				if af.Ants[i].CurrentRoom != af.Start {
					t.Errorf("Ant[%d].CurrentRoom = %v, want %v", i, af.Ants[i].CurrentRoom, af.Start)
				}

				// Check if CurrentRoom is actually the start room
				if !af.Ants[i].CurrentRoom.IsStart {
					t.Errorf("Ant[%d].CurrentRoom is not marked as start room", i)
				}

				// Check PathIndex
				if af.Ants[i].PathIndex != 0 {
					t.Errorf("Ant[%d].PathIndex = %v, want 0", i, af.Ants[i].PathIndex)
				}

				// Check HasReached
				if af.Ants[i].HasReached {
					t.Errorf("Ant[%d].HasReached = true, want false", i)
				}
			}
		})
	}
}
