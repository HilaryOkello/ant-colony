package antfarm

import "test/models"

// dfs performs depth-first search to find all possible paths
func (af *AntFarm) dfs(current, end *models.Room, visited map[*models.Room]bool, path []*models.Room, paths *[]models.Path) {
	visited[current] = true
	path = append(path, current)

	if current == end {
		// Create a new path
		newPath := models.Path{
			Rooms:  make([]*models.Room, len(path)),
			Length: len(path) - 1,
			InUse:  false,
		}
		copy(newPath.Rooms, path)
		*paths = append(*paths, newPath)
	} else {
		for _, next := range current.Connected {
			if !visited[next] {
				af.dfs(next, end, visited, path, paths)
			}
		}
	}

	// Backtrack
	visited[current] = false
}
