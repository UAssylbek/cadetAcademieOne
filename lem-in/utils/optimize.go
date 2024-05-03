package utils

import (
	"os"
)

func RaspredelenyeAnt(minAntAllPath *AntAllPath, countAnts int) []AntInfo {
	antData := make([]AntInfo, countAnts)
	for i := 0; i < countAnts; i++ {
		antData[i].Ves = make([]int, len(minAntAllPath.ArrayRooms))
		for j := 0; j < len(minAntAllPath.ArrayRooms); j++ {
			antData[i].Ves[j] = len(minAntAllPath.ArrayRooms[j])
		}
	}

	for i := 0; i < countAnts; i++ {
		antData[i].Id = i
	}

	aa := antData[0].Ves
	for i := 0; i < countAnts; i++ {
		l := aa[0]
		id := 0
		for i1 := 0; i1 < len(aa); i1++ {
			if !(aa[i1] > l) {
				l = aa[i1]
				id = i1
			}
		}
		aa[id] = aa[id] + 1
		antData[i].Room = minAntAllPath.ArrayRooms[id]
	}
	return antData
}

func AntAllPaths(allProcessed [][][]string, n int) []*AntAllPath {
	antAllPaths := make([]*AntAllPath, len(allProcessed))

	for i, paths := range allProcessed {
		antAllPath := &AntAllPath{}
		for _, path := range paths {
			l := len(path)
			if l > antAllPath.MaxL {
				antAllPath.MaxL = l
			}
		}
		for _, path := range paths {
			antAllPath.MaxAnts += antAllPath.MaxL - len(path) + 1
		}

		antAllPath.ExcessAnts = n - antAllPath.MaxAnts
		antAllPath.AdditionTunels = (antAllPath.ExcessAnts + len(paths) - 1) / len(paths)
		if n <= antAllPath.MaxAnts {
			antAllPath.FinnalTunnels = antAllPath.MaxL
		} else {
			antAllPath.FinnalTunnels = antAllPath.AdditionTunels + antAllPath.MaxL
		}

		antAllPath.ArrayRooms = paths
		antAllPaths[i] = antAllPath
	}

	return antAllPaths
}

func FindUniquePaths(one int, allPaths [][][]string, stringLists [][]string, stringListLists [][][]string, two int) [][][]string {
	count := 0
	for i := 0; i < len(allPaths); i++ {
		num := 0
		for k := 0; k < len(allPaths[i]); k++ {
			if !HasDuplicates(allPaths[i][k], stringLists) {
				stringLists = append(stringLists, allPaths[i][k])
				num++
			}
		}
		if num == 0 {
			stringListLists = append(stringListLists, append([][]string{}, stringLists...))
			stringLists = make([][]string, 0)
			count++
		}
	}
	if two != 1000 {
		allPaths = MoveFirstToEnd(allPaths)
		return FindUniquePaths(one, allPaths, stringLists, stringListLists, two+1)
	}

	return stringListLists
}

func ReadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Выделение памяти для буфера, в котором будет храниться содержимое файла
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	// Чтение содержимого файла в буфер
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Преобразование буфера в строку
	fileContent := string(buffer)
	return fileContent, nil
}
