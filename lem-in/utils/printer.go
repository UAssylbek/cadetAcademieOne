package utils

import "fmt"

func PrintAntMovements(ants []AntInfo, minAntAllPath *AntAllPath) {

	for q := 0; q < minAntAllPath.FinnalTunnels-1; q++ {
		data := make(map[int]string)
		for i := 0; i < len(ants); i++ {
			if len(ants[i].Room) == 1 {
				if ants[i].Location < len(ants[i].Room) && IsValueNotInMap(data, ants[i].Room[ants[i].Location]) {
					data[ants[i].Id] = ants[i].Room[ants[i].Location]
					ants[i].Location++
				}
			} else {
				if ants[i].Location < len(ants[i].Room) && (ants[i].Location == len(ants[i].Room)-1 || IsValueNotInMap(data, ants[i].Room[ants[i].Location])) {
					data[ants[i].Id] = ants[i].Room[ants[i].Location]
					ants[i].Location++
				}
			}
		}
		for id, value := range data {
			fmt.Printf("L%d-%s ", id+1, value)
		}
		fmt.Println()
	}
}
