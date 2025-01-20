package antfarm

import (
	"errors"
	"fmt"
	"strings"

	"test/models"
)

// simulateMovement simulates the movement of all ants using multiple paths
func (af *AntFarm) SimulateMovement() (string, error) {
	antPaths := af.assignAntsToPath()
	if antPaths == nil {
		return "", errors.New("ERROR: no valid path found between start and end")
	}

	if len(af.Ants) == 0 {
		return "", errors.New("no ants available")
	}

	// Initialize ant positions
	for ant := range antPaths {
		path := antPaths[ant]
		ant.Path = path.Rooms
		ant.PathIndex = 0
		ant.CurrentRoom = path.Rooms[0]
		ant.HasReached = false
	}

	// Track room occupancy
	occupiedRooms := make(map[*models.Room]*models.Ant)
	allMoves := ""

	// Simulate movements
	for {
		moves := make([]string, 0)
		allReached := true
		move := 0

		// Clear non-start/end room occupancy at start of turn
		for room := range occupiedRooms {
			if !room.IsStart && !room.IsEnd {
				occupiedRooms[room] = nil
			}
		}

		// Try to move each ant
		for _, ant := range af.Ants {
			move++
			if ant == nil {
				return "", errors.New("ant is nil")
			}

			if len(ant.Path) == 0 {
				return "", fmt.Errorf("ant %d has no valid path", ant.Id)
			}
			if ant.CurrentRoom == nil {
				return "", fmt.Errorf("ant %d has no current room set", ant.Id)
			}

			if ant.HasReached {
				move = 0
				continue
			}
			if len(ant.Path) == 2 && move >= 2 {
				continue
			}

			allReached = false

			// Check if ant can move
			if ant.PathIndex < len(ant.Path)-1 {
				nextRoom := ant.Path[ant.PathIndex+1]

				// Check if next room is available
				if occupiedRooms[nextRoom] == nil || nextRoom.IsEnd {
					// Move ant
					if !ant.CurrentRoom.IsStart && !ant.CurrentRoom.IsEnd {
						occupiedRooms[ant.CurrentRoom] = nil
					}

					ant.CurrentRoom = nextRoom
					if !nextRoom.IsStart && !nextRoom.IsEnd {
						occupiedRooms[nextRoom] = ant
					}

					ant.PathIndex++
					moves = append(moves, fmt.Sprintf("L%d-%s", ant.Id, nextRoom.Name))

					if nextRoom.IsEnd {
						ant.HasReached = true
					}
				}
			}

		}

		move = 0
		if len(moves) > 0 {
			allMoves += strings.Join(moves, " ") + "\n"
		}

		if allReached {
			break
		}
	}
	return allMoves, nil
}
