package path

import (
	"lemin/utils"
	"sort"
)

func FindAllPaths(start, end string, tunnels [][]string) [][][]string {
	visited := make(map[string]bool)
	var paths [][]string
	var currentPath []string
	var dfs func(node string)

	dfs = func(node string) {
		currentPath = append(currentPath, node)
		visited[node] = true
		if node == end {
			path := make([]string, len(currentPath))
			copy(path, currentPath)
			paths = append(paths, path)
		} else {
			for _, connection := range tunnels {
				if connection[0] == node && !visited[connection[1]] {
					dfs(connection[1])
				} else if connection[1] == node && !visited[connection[0]] {
					dfs(connection[0])
				}
			}
		}
		currentPath = currentPath[:len(currentPath)-1]
		visited[node] = false
	}
	dfs(start)
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	allPaths3 := make([][][]string, 0)
	uniqueFirstRooms := make([]string, 0)
	for _, path := range paths {
		if !utils.IsFirstRoomUnique(path[1], uniqueFirstRooms) {
			uniqueFirstRooms = append(uniqueFirstRooms, path[1])
			allPaths3 = append(allPaths3, utils.SplitArray(path[1], paths))
		}
	}
	return allPaths3
}
