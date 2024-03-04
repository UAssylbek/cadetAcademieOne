package tetris

import "fmt"

func PusKorobka(size int) [][]rune {
	array := make([][]rune, size)
	for i := range array {
		array[i] = make([]rune, size)
		for j := range array[i] {
			array[i][j] = '.'
		}
	}
	return array
}

func ProSimvol(s string) bool {
	for i, ch := range s {
		if s[i] != '#' && ch != '.' && s[i] != '\n' {
			return false
		}
	}
	return true
}

func Split(str string) [][]string {
	eki := [][]string{}
	bir := []string{}
	word := ""
	i := 0
	for i < len(str) {
		if str[i] == '\n' {
			if i+1 < len(str) {
				if i+1 < len(str) {
					if str[i+1] == '\n' {
						bir = append(bir, word)
						eki = append(eki, bir)
						bir = []string{}
						word = ""
						i = i + 2
						continue
					}
					if word != "\n" {
						bir = append(bir, word)
						word = ""
					}
				}
			}
		} else {
			word = word + string(str[i])
		}
		i++
	}
	if word != "" {
		bir = append(bir, word)
		eki = append(eki, bir)
	}
	return eki
}

func Kvadrat4(str [][]string) bool {
	for i := 0; i < len(str); i++ {

		if len(str[i]) != 4 {
			return false
		}
		for j := 0; j < len(str[i]); j++ {
			if len(str[i][j]) != 4 {
				return false
			}
		}
	}
	return true
}

func ProSosedy(s [][]string) []int {
	count := 0
	arr := []int{}
	for k := 0; k < len(s); k++ {
		for k1 := 0; k1 < len(s[k]); k1++ {
			val := s[k][k1]

			for k2 := 0; k2 < len(val); k2++ {
				if string(val[k2]) == "#" {
					count++
					if k2 != 0 && string(val[k2-1]) == "#" {
						count++
					}
					if k2 != 3 && string(val[k2+1]) == "#" {
						count++
					}
					if k1 != 3 && s[k][k1+1][k2] == '#' {
						count++
					}
					if k1 != 0 && s[k][k1-1][k2] == '#' {
						count++
					}
				}
			}
		}
		arr = append(arr, count)
		count = 0
	}
	return arr
}

func ProKol(arr []int) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] != 10 && arr[i] != 12 {
			return false
		}
	}
	return true
}

func Coordinates(massivtext [][]string) [][][]int {
	arr := []int{}
	arr2 := [][]int{}
	arr3 := [][][]int{}
	for g := 0; g < len(massivtext); g++ {
		for x := 0; x < len(massivtext[g]); x++ {
			for y := 0; y < len(massivtext[g][x]); y++ {
				if string(massivtext[g][x][y]) == "#" {
					arr = append(arr, x)
					arr = append(arr, y)
				}
				if len(arr) > 0 {
					arr2 = append(arr2, arr)
					arr = []int{}
				}
			}
		}
		if len(arr2) > 0 {
			arr3 = append(arr3, arr2)
			arr2 = [][]int{}
		}
	}
	return arr3
}

func LeftUp(bb [][][]int) [][][]int {
	transformedArray := make([][][]int, len(bb))
	for i, bb2 := range bb {
		minX, minY := findMinCoordinates(bb2)
		for _, bb1 := range bb2 {
			transformedArray[i] = append(transformedArray[i], []int{bb1[0] - minX, bb1[1] - minY})
		}
	}
	return transformedArray
}

func findMinCoordinates(bb2 [][]int) (minX, minY int) {
	minX = bb2[0][0]
	minY = bb2[0][1]
	for _, bb1 := range bb2 {
		if bb1[0] < minX {
			minX = bb1[0]
		}
		if bb1[1] < minY {
			minY = bb1[1]
		}
	}
	return minX, minY
}

func CanPlace(board [][]rune, tetromino [][]int, x, y int) bool {
	for _, coord := range tetromino {
		newX, newY := x+coord[0], y+coord[1]

		if newX < 0 || newX >= len(board) || newY < 0 || newY >= len(board[0]) {
			return false
		}

		if board[newX][newY] != '.' {
			return false
		}
	}

	return true
}

func PlaceTetromino(board [][]rune, tetromino [][]int, x, y, index int) [][]rune {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	newBoard := make([][]rune, len(board))
	for i := range newBoard {
		newBoard[i] = make([]rune, len(board[i]))
		copy(newBoard[i], board[i])
	}

	for _, coord := range tetromino {
		newX, newY := x+coord[0], y+coord[1]
		newBoard[newX][newY] = rune(alphabet[index])
	}
	return newBoard
}

func PlaceTetrominoes(board [][]rune, tetrominoes [][][]int, index int) ([][]rune, error) {
	if index == len(tetrominoes) {
		return board, nil
	}

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if CanPlace(board, tetrominoes[index], i, j) {
				newBoard := PlaceTetromino(board, tetrominoes[index], i, j, index)

				result, err := PlaceTetrominoes(newBoard, tetrominoes, index+1)
				if err == nil {
					return result, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("Cannot place tetromino %d", index)
}
