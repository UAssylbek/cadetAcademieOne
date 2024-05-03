package admin

import (
	"fmt"
	"lemin/path"
	"lemin/utils"
	"lemin/validators"
	"os"
)

func Run(filename string) {
	countAnts, start, end, tunnels, err := validators.ReadInput(filename)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}
	allPaths := path.FindAllPaths(start, end, tunnels)

	stringLists := make([][]string, 0)
	stringListLists := make([][][]string, 0)
	var one, two int
	allProcessed := utils.FindUniquePaths(one, allPaths, stringLists, stringListLists, two)
	utils.Sort3DArray(allProcessed)
	allProcessed = utils.RemoveDuplicateArrays(allProcessed)
	allProcessed = utils.RemoveDuplicates(allProcessed)
	allProcessed = utils.RemoveDuplicateArrays(allProcessed)
	antAllPaths := utils.AntAllPaths(allProcessed, countAnts)

	for _, antAllPath := range antAllPaths {
		if utils.MinAntAllPath == nil || antAllPath.FinnalTunnels < utils.MinAntAllPath.FinnalTunnels {
			utils.MinAntAllPath = antAllPath
		}
	}

	for i := 0; i < len(utils.MinAntAllPath.ArrayRooms); i++ {
		utils.MinAntAllPath.ArrayRooms[i] = utils.MinAntAllPath.ArrayRooms[i][1:]
	}
	ants := utils.RaspredelenyeAnt(utils.MinAntAllPath, countAnts)

	content, err := utils.ReadFile(filename)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}
	fmt.Println(content)
	fmt.Println()
	utils.PrintAntMovements(ants, utils.MinAntAllPath)
}
