package antfarm

import (
	"sort"

	"test/models"
)

// findAllPaths traverses the colony using a depth-first to return sorted non-overlapping paths
func (af *AntFarm) findAllPaths() []models.Path {
	paths := make([]models.Path, 0)
	visited := make(map[*models.Room]bool)
	currentPath := make([]*models.Room, 0)

	af.dfs(af.Start, af.End, visited, currentPath, &paths)
	// Sort paths by length
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].Length < paths[j].Length
	})

	// Filter out overlapping paths
	nonOverlappingPaths := af.filterNonOverlappingPaths(paths)

	return nonOverlappingPaths
}

// filterNonOverlappingPaths filters and returns the best combination of non-overlapping paths from a list of paths. 
// The function avoids overlaps by ensuring that no two paths in the final result share any "middle" rooms (rooms 
// that are not the start or end). 
func (af *AntFarm) filterNonOverlappingPaths(paths []models.Path) []models.Path {
	// Helper function to check if two paths overlap
	pathsOverlap := func(path1, path2 models.Path) bool {
		rooms1 := make(map[string]struct{})

		// Add middle rooms from path1 to set1 (excluding start and end)
		for i := 1; i < len(path1.Rooms)-1; i++ {
			rooms1[path1.Rooms[i].Name] = struct{}{}
		}

		// Check if any middle room from path2 exists in set1
		for i := 1; i < len(path2.Rooms)-1; i++ {
			if _, exists := rooms1[path2.Rooms[i].Name]; exists {
				return true
			}
		}
		return false
	}

	// Helper function to find best combination of non-overlapping paths
	findBestCombination := func(paths []models.Path) []models.Path {
		n := len(paths)
		if n == 0 {
			return []models.Path{}
		}

		bestCombination := []models.Path{paths[0]} // Initialize with first path
		maxPaths := 1

		// Function to check if a path is compatible with a combination
		isCompatible := func(path models.Path, combination []models.Path) bool {
			for _, existingPath := range combination {
				if pathsOverlap(path, existingPath) {
					return false
				}
			}
			return true
		}

		// Try each path as the starting path
		for startIdx := 0; startIdx < n; startIdx++ {
			currentCombination := []models.Path{paths[startIdx]}

			// Try to add other paths
			for i := 0; i < n; i++ {
				if i == startIdx {
					continue
				}

				// Check if current path can be added
				if isCompatible(paths[i], currentCombination) {
					currentCombination = append(currentCombination, paths[i])
				}
			}

			// Update best combination if current is better
			if len(currentCombination) > maxPaths {
				maxPaths = len(currentCombination)
				bestCombination = make([]models.Path, len(currentCombination))
				copy(bestCombination, currentCombination)
			}
		}

		return bestCombination
	}

	result := findBestCombination(paths)
	return result
}

// assignAntsToPath assigns ants to optimal paths
func (af *AntFarm) assignAntsToPath() map[*models.Ant]models.Path {
	paths := af.findAllPaths()
	if len(paths) == 0 {
		return nil
	}

	antPaths := make(map[*models.Ant]models.Path)
	pathAnts := make([]int, len(paths)) // Tracks the number of ants assigned to each path
	totalMoves := 0                     // Tracks the total moves required for all ants

	// Assign ants to paths by minimizing total moves
	for i := 0; i < len(af.Ants); i++ {
		bestTurns := int(^uint(0) >> 1) // Max int to find the path with the least moves
		bestPathIndex := 0

		// Find the best path for the current ant (with the fewest total moves)
		for j, path := range paths {
			antsOnPath := pathAnts[j]
			totalTurns := path.Length + antsOnPath

			// Choose the path with the fewest moves
			if totalTurns < bestTurns {
				bestTurns = totalTurns
				bestPathIndex = j
			}
		}

		// Assign the ant to the best path found
		antPaths[af.Ants[i]] = paths[bestPathIndex]
		pathAnts[bestPathIndex]++
		totalMoves += bestTurns // Add the moves for this assignment to the total moves
	}

	return antPaths
}

