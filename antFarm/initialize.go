package antfarm

import "test/models"

// NewAntFarm creates a new ant farm instance
func NewAntFarm() *AntFarm {
	return &AntFarm{
		Rooms: make(map[string]*models.Room),
	}
}

// initializeAnts creates all ants in the start room
func (af *AntFarm) initializeAnts() {
	af.Ants = make([]*models.Ant, af.NumAnts)
	for i := 0; i < af.NumAnts; i++ {
		af.Ants[i] = &models.Ant{
			Id:          i + 1,
			CurrentRoom: af.Start,
			PathIndex:   0,
			HasReached:  false,
		}
	}
}
