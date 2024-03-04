package main

import (
	"fmt"
	"os"
	"strings"
	"utils/utils"
)

func main() {
	if len(os.Args) == 0 {
		return
	}
	var arg, wordslice []string
	var colors [][]string
	var str, Filename, namefile, position, txtfile, kil string
	arg, Filename, txtfile, kil = utils.WriteArg(os.Args)
	if kil != "" {
		fmt.Println("Usage: go run . [OPTION]",
			"\n\nEX: go run . --reverse=<fileName>")
		return
	}
	Font, err := os.ReadFile("standard.txt")
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	slice := utils.StringTo2DArray(string(Font))
	slice[0] = utils.AntiFirst(slice[0])
	if txtfile != "" {
		Fon, err := os.ReadFile(txtfile)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		var qwer, qwer1 []string
		var qwer2 [][]string
		qwerty := ""
		upo := 0
		for i := 0; i < len(Fon); i++ {
			if Fon[i] == '\n' && i != 0 && Fon[i-1] != '\n' {
				qwer = append(qwer, qwerty)
				qwerty = ""
				upo++
			} else if Fon[i] == '\n' && i == 0 {
				qwer2 = append(qwer2, qwer1)
			} else if Fon[i] != '\n' {
				qwerty = qwerty + string(Fon[i])
			}
			if i != 0 && Fon[i] == '\n' && Fon[i-1] == '\n' {
				qwer2 = append(qwer2, qwer1)
			} else if upo == 8 {
				qwer2 = append(qwer2, qwer)
				qwer = []string{}
				upo = 0
			}
		}
		index1 := [][]int{}
		index := []int{}
		sttr := []string{}
		stttrr := [][]string{}

		for u := 0; u < len(qwer2); u++ {
			if len(qwer2[u]) != 0 {
				mino := len(qwer2[u][0])
				if mino == len(qwer2[u][1]) && mino == len(qwer2[u][2]) && mino == len(qwer2[u][3]) && mino == len(qwer2[u][4]) && mino == len(qwer2[u][5]) && mino == len(qwer2[u][6]) && mino == len(qwer2[u][7]) {
					if qwer2[u][0][len(qwer2[u][0])-1] == 32 && qwer2[u][1][len(qwer2[u][0])-1] == 32 && qwer2[u][2][len(qwer2[u][0])-1] == 32 && qwer2[u][3][len(qwer2[u][0])-1] == 32 && qwer2[u][4][len(qwer2[u][0])-1] == 32 && qwer2[u][5][len(qwer2[u][0])-1] == 32 && qwer2[u][6][len(qwer2[u][0])-1] == 32 && qwer2[u][7][len(qwer2[u][0])-1] == 32 {
						for i := 0; i < len(qwer2[u][0]); i++ {
							if qwer2[u][0][i] == 32 && qwer2[u][1][i] == 32 && qwer2[u][2][i] == 32 && qwer2[u][3][i] == 32 && qwer2[u][4][i] == 32 && qwer2[u][5][i] == 32 && qwer2[u][6][i] == 32 && qwer2[u][7][i] == 32 {
								index = append(index, i)
								if i+1 < len(qwer2[u][0]) {
									if qwer2[u][0][i+1] == 32 && qwer2[u][1][i+1] == 32 && qwer2[u][2][i+1] == 32 && qwer2[u][3][i+1] == 32 && qwer2[u][4][i+1] == 32 && qwer2[u][5][i+1] == 32 && qwer2[u][6][i+1] == 32 && qwer2[u][7][i+1] == 32 {
										i = i + 5
									}
								}
							}
						}
						index1 = append(index1, index)
						index = []int{}
					} else {
						fmt.Println("Usage: go run . [OPTION]",
							"\n\nEX: go run . --reverse=<fileName>")
						return
					}
				} else {
					fmt.Println("Usage: go run . [OPTION]",
						"\n\nEX: go run . --reverse=<fileName>")
					return
				}
			} else {
				indexer := []int{0}
				index1 = append(index1, indexer)
			}

		}
		p := 0
		s := 0
		var stttrr2 [][][]string
		ipo := [][]string{}
		for y := 0; y < len(index1); y++ {
			if index1[y][0] == 0 && len(index1[y]) == 1 {
				upi := []string{"\n"}
				ipo = append(ipo, upi)
				stttrr2 = append(stttrr2, ipo)
			} else {
				for i := 0; i < len(index1[y]); i++ {
					for j := 0; j < 8; j++ {
						sttr = append(sttr, qwer2[s][j][p:index1[y][i]+1])
					}
					stttrr = append(stttrr, sttr)
					sttr = []string{}
					p = index1[y][i] + 1
				}
				stttrr2 = append(stttrr2, stttrr)
				stttrr = [][]string{}
				sttr = []string{}
				p = 0
			}
			s++
		}
		for m := 0; m < len(stttrr2); m++ {
			for i := 0; i < len(stttrr2[m]); i++ {
				if stttrr2[m][i][0] == "\n" {
					continue
				} else {
					pol := 0
					for j := 0; j < len(slice); j++ {
						if slicesEqual(stttrr2[m][i], slice[j]) {
							pol++
						}
					}
					if pol == 0 {
						fmt.Println("Usage: go run . [OPTION]",
							"\n\nEX: go run . --reverse=<fileName>")
						return
					}
				}
			}
		}

		for m := 0; m < len(stttrr2); m++ {
			for i := 0; i < len(stttrr2[m]); i++ {
				if stttrr2[m][i][0] == "\n" {
					continue
				} else {
					pol := 0
					for j := 0; j < len(slice); j++ {
						if slicesEqual(stttrr2[m][i], slice[j]) {
							fmt.Print(string(rune(j + 32)))
							pol++
						}
					}
					if pol == 0 {
						fmt.Println("Usage: go run . [OPTION]",
							"\n\nEX: go run . --reverse=<fileName>")
						return
					}
				}
			}
			fmt.Println()
		}
	}

	if Filename != "" {
		Font, err := os.ReadFile(Filename)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}

		colors, arg, namefile, position, err = utils.FindOption(arg)
		if err != nil {
			fmt.Println("Usage: go run . [OPTION] [STRING]",
				"\n\nEX: go run . --color=<color> <letters to be colored> \"something\"")
			return
		}
		fmt.Println(arg)
		if len(arg) > 1 {
			fmt.Println("Usage: go run . [OPTION] [STRING]",
				"\n\nEX: go run . --color=<color> <letters to be colored> \"something\"")
			return
		}

		str = strings.Join(arg, " ")
		slice := utils.StringTo2DArray(string(Font))
		slice[0] = utils.AntiFirst(slice[0])
		if len(slice) != 95 {
			return
		}

		var cir interface{} = nil
		if position != "" {
			cleanedStr := strings.Join(strings.Fields(str), " ")
			cleanedStr = strings.Replace(cleanedStr, "\\n", "\n", -1)
			wordslice, cir = utils.CheckLatin(cleanedStr)
		} else {
			wordslice, cir = utils.CheckLatin(str)
		}
		if cir != nil {
			fmt.Println(cir)
			return
		}

		utils.PrintWords(namefile, wordslice, colors, slice, position)
	}
}

func slicesEqual(slice1, slice2 []string) bool {
	// Проверяем, имеют ли срезы одинаковую длину
	if len(slice1) != len(slice2) {
		return false
	}

	// Проверяем каждый элемент в срезах
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	// Срезы равны, если все элементы совпадают
	return true
}
