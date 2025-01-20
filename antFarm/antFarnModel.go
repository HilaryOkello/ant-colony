package antfarm

import (
	"test/models"
)

type AntFarm struct {
	NumAnts int
	Ants    []*models.Ant
	Rooms   map[string]*models.Room
	Start   *models.Room
	End     *models.Room
}
