package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func MoveFirstToEnd(arr [][][]string) [][][]string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func HasDuplicates(a []string, b [][]string) bool {
	if len(b) == 0 {
		return false
	}
	for i := 0; i < len(b); i++ {
		for k := 1; k < len(b[i])-1; k++ {
			if !CheckForDuplicates(b[i][k], a) {
				return true
			}
		}
	}
	return false
}

func CheckForDuplicates(a string, b []string) bool {
	for i := 1; i < len(b)-1; i++ {
		if b[i] == a {
			return false
		}
	}
	return true
}

func IsFirstRoomUnique(a string, b []string) bool {
	for i := 0; i < len(b); i++ {
		if b[i] == a {
			return true
		}
	}
	return false
}

func SplitArray(a string, b [][]string) [][]string {
	allPaths := make([][]string, 0)
	for _, path := range b {
		if path[1] == a {
			allPaths = append(allPaths, path)
		}
	}
	return allPaths
}

func Sort3DArray(arr [][][]string) {
	for _, subArr := range arr {
		sort.Slice(subArr, func(i, j int) bool {
			return len(arr[i]) < len(arr[j])
		})
	}

	sort.Slice(arr, func(i, j int) bool {
		return CountElements(arr[i]) < CountElements(arr[j])
	})
}

func CountElements(arr [][]string) int {
	count := 0
	for _, subArr := range arr {
		for _, element := range subArr {
			if element != "" {
				count++
			}
		}
	}
	return count
}

func RemoveDuplicateArrays(allOnly2 [][][]string) [][][]string {
	uniqueArrays := make([][][]string, 0)
	encountered := make(map[string]bool)
	elementsCount := make(map[string]bool)

	for _, arr := range allOnly2 {
		key := ArrayKey(arr)

		if elementsCount[key] {
			continue
		}

		if !encountered[key] {
			encountered[key] = true
			uniqueArrays = append(uniqueArrays, arr)
		}

		elementsCount[key] = true
	}
	return uniqueArrays
}

func RemoveDuplicates(arr [][][]string) [][][]string {
	result := make([][][]string, len(arr))
	for i := 0; i < len(arr); i++ {
		uniqueMap := make(map[string]bool)
		var uniqueArr [][]string
		for j := 0; j < len(arr[i]); j++ {
			str := strings.Join(arr[i][j], "")
			if !uniqueMap[str] {
				uniqueMap[str] = true
				uniqueArr = append(uniqueArr, arr[i][j])
			}
		}
		result[i] = uniqueArr
	}

	return result
}

func ArrayKey(arr [][]string) string {
	var key strings.Builder
	key.WriteString(fmt.Sprintf("%d:", len(arr)))

	for _, subArr := range arr {
		for _, element := range subArr {
			key.WriteString(fmt.Sprintf("%s,", element))
		}
	}

	return key.String()
}

func IsValueNotInMap(data map[int]string, Location string) bool {
	for _, value := range data {
		if value == Location {
			return false
		}
	}
	return true
}

func СheckForEmptyStrings(slice []string) error {
	for _, str := range slice {
		if len(str) == 0 {
			return errors.New("найден пустой элемент в срезе")
		}
	}
	return nil
}
