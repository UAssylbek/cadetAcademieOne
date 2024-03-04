package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func StringTo2DArray(str string) [][]string {
	str = strings.ReplaceAll(str, "\r", "")
	test1 := [][]string{}
	test := strings.Split(str, "\n\n")
	for _, i := range test {
		w := strings.Split(i, "\n")
		test1 = append(test1, w)
	}
	return test1
}

func CheckLatin(str string) ([]string, error) {
	//check for printable characters
	for i := 0; i < len(str); i++ {
		if (str[i] < 32 || str[i] > 126) && str[i] != '\n' {
			return nil, errors.New("wrong sintax in input data")
		}
	}
	wordslice := strings.Split(strings.ReplaceAll(str, "\\n", "\n"), "\n")
	return wordslice, nil
}

func PrintWords(namefile string, wordslice []string, colors [][]string, slice [][]string, position string) {
	_, size, _ := GetTerminalWidth()
	char2slice := [][]string{}
	file, err := os.Create(namefile)
	O := 0
	centered := 0
	Probel := 0
	Probel1 := 0
	q := 0
	i := 0
	for i < len(wordslice) {
		slovo := ""
		count := 0
		if wordslice[i] == "" {
			if namefile == "" {
				fmt.Println()
			} else {
				_, err := file.WriteString("\n")
				if err != nil {
					fmt.Println("Error writing to file:", err)
				}
			}
		} else {
			if len(wordslice[i]) != 0 {
				for w, l := range wordslice[i] {
					if l != '\n' {
						char2slice = append(char2slice, slice[(l-32)])
						if l == 32 {
							if w != len(wordslice[i])-1 && w != 0 {
								Probel++
							}
							Probel1++
						}
					}
				}
				RemoveColors(char2slice)
				for k := 0; k < 8; k++ {
					for j := 0; j < len(char2slice); j++ {
						if k == 0 {
							count = count + len(char2slice[j][k])
						}

					}
				}
			Loop:
				if len(colors) != 0 {
					if len(colors[O]) == 1 {
						RemoveColors(char2slice)
						for y := 0; y < 8; y++ {
							for j := 0; j < len(char2slice); j++ {
								char2slice[j][y] = colors[O][0] + char2slice[j][y] + FindColor("reset")
							}
						}
					} else if len(colors[O]) == 2 {
						for k := 0; k < 8; k++ {
							for j := 0; j < len(char2slice); j++ {
								if IsLetterExist(string(wordslice[i][j]), colors[O][1]) {
									char2slice[j][k] = colors[O][0] + RemoveColorsFromString(char2slice[j][k]) + FindColor("reset")
								}
							}
						}
					}
					if O+1 < len(colors) {
						O++
						goto Loop
					}
				}
				if position == "left" {
					centered = 0
				} else if position == "right" {
					centered = size - count
				} else if position == "center" {
					centered = ((size - count) / 2)
				} else if position == "justify" {
					if Probel == 0 {
						centered = 0
					} else {
						centered = ((size - (count - (Probel1 * 6))) / Probel)
					}
				}
				for l := 0; l < centered; l++ {
					slovo = slovo + " "
				}

				if namefile == "" {
					if position == "right" || position == "left" || position == "center" {
						for k := 0; k < 8; k++ {
							fmt.Print(slovo)
							for j := 0; j < len(char2slice); j++ {
								fmt.Print(char2slice[j][k])
							}
							fmt.Println()
						}
					} else if position == "justify" {
						for k := 0; k < 8; k++ {
							for j := 0; j < len(char2slice); j++ {
								if wordslice[i][q] == 32 {
									if j != 0 {
										if j != len(wordslice[i])-1 {
											fmt.Print(slovo)
										}
									}
								} else {
									fmt.Print(char2slice[j][k])
								}
								if q+1 < len(wordslice[i]) {
									q++
								}
							}
							q = 0
							fmt.Println()
						}
					} else {
						for k := 0; k < 8; k++ {
							for j := 0; j < len(char2slice); j++ {
								fmt.Print(char2slice[j][k])
							}
							fmt.Println()
						}
					}
				} else {
					if position == "justify" {
						if err != nil {
							fmt.Println("Unable to create file:", err)
						}
						defer file.Close()
						if wordslice[i] == "\n" {
							_, err := file.WriteString("\n")
							if err != nil {
								fmt.Println("Error writing to file:", err)
							}
						} else {
							if len(wordslice[i]) != 0 {
								for k := 0; k < 8; k++ {
									for j := 0; j < len(char2slice); j++ {
										if wordslice[i][q] == 32 {
											_, err := file.WriteString(slovo)
											if err != nil {
												fmt.Println("Error writing to file:", err)
											}
										}
										if q+1 < len(wordslice[i]) {
											q++
										}
										_, err := file.WriteString(char2slice[j][k])
										if err != nil {
											fmt.Println("Error writing to file:", err)
										}
									}
									q = 0
									_, err = file.WriteString("\n")
									if err != nil {
										fmt.Println("Error writing to file:", err)
									}
								}
							}
						}
						file.WriteString("\n")
					} else {
						if err != nil {
							fmt.Println("Unable to create file:", err)
						}
						defer file.Close()
						if wordslice[i] == "\n" {
							_, err := file.WriteString("\n")
							if err != nil {
								fmt.Println("Error writing to file:", err)
							}
						} else {
							if len(wordslice[i]) != 0 {
								for k := 0; k < 8; k++ {
									_, err := file.WriteString(slovo)
									if err != nil {
										fmt.Println("Error writing to file:", err)
									}
									for j := 0; j < len(char2slice); j++ {
										_, err := file.WriteString(char2slice[j][k])
										if err != nil {
											fmt.Println("Error writing to file:", err)
										}
									}
									_, err = file.WriteString("\n")
									if err != nil {
										fmt.Println("Error writing to file:", err)
									}
								}
							}
						}
						// file.WriteString("\n")
					}
				}
			}
		}
		Probel = 0
		Probel1 = 0
		char2slice = [][]string{}
		i++
	}
}

func AntiFirst(s []string) []string {
	s = s[1:]
	return s
}

func FindColor(letter string) string {
	switch letter {
	case "red":
		return "\033[31m"
	case "green":
		return "\033[32m"
	case "yellow":
		return "\033[33m"
	case "blue":
		return "\033[34m"
	case "magenta":
		return "\033[35m"
	case "cyan":
		return "\033[36m"
	case "white":
		return "\033[37m"
	case "reset":
		return "\033[0m"
	case "black":
		return "\033[30m"
	case "orange":
		return "\033[38;5;208m"
	}
	return "I"
}

func IsLetterExist(str string, letter string) bool {

	for _, i := range str {
		for _, j := range letter {
			if i == j {
				return true
			}
		}
	}
	return false

}

func WriteArg(args []string) ([]string, string, string, string) {
	var arg []string
	Filename := ""
	txtfile := ""
	kil := ""
	if len(os.Args) > 2 {
		for i := 1; i < len(os.Args); i++ {
			if len(os.Args[i]) > 10 {
				if os.Args[i][:10] == "--reverse=" {
					kil = "Hello"
					return arg, Filename, txtfile, kil
				}
			}
		}
	}
	if len(os.Args[1]) > 10 {

		if os.Args[1][:10] == "--reverse=" {
			txtfile = os.Args[1][10:]
			arg = os.Args[1 : len(os.Args)-1]
			return arg, Filename, txtfile, kil
		} else if CheckReverse(os.Args[1][:10]) {
			kil = "Hello"
			return arg, Filename, txtfile, kil
		}
	}

	if len(os.Args) == 2 {
		if len(os.Args[1]) > 10 {

			if os.Args[1][:10] == "--reverse=" {
				txtfile = os.Args[1][10:]
				arg = os.Args[1 : len(os.Args)-1]
				return arg, Filename, txtfile, kil
			} else if CheckReverse(os.Args[1][:10]) {
				kil = "Hello"
				return arg, Filename, txtfile, kil
			}
		}
	}

	// if len(args) == 2 {
	// 	if strings.HasPrefix(args[1], "--reverse=") {
	// 		txtfile = strings.TrimPrefix(args[1], "--reverse=")
	// 	}
	// 	arg = os.Args[1 : len(os.Args)-1]
	// 	return arg, Filename, txtfile
	// }

	if len(os.Args) > 2 {
		if os.Args[len(os.Args)-1] == "standard" {
			Filename = "standard.txt"
			if len(os.Args[1]) > 9 {
				if len(os.Args) == 3 && (os.Args[1][:8] == "--color=" || os.Args[1][:9] == "--output=" || os.Args[1][:8] == "--align=") {
					arg = os.Args[1:]
					Filename = "standard.txt"
				} else if len(os.Args) > 3 && (os.Args[1][:8] == "--color=" || os.Args[1][:9] == "--output=" || os.Args[1][:8] == "--align=") {
					arg = os.Args[1 : len(os.Args)-1]
				} else {
					arg = os.Args[1:]
				}
			} else {
				arg = os.Args[1 : len(os.Args)-1]
			}
		} else if os.Args[len(os.Args)-1] == "shadow" {
			Filename = "shadow.txt"
			if len(os.Args[1]) > 9 {
				if len(os.Args) == 3 && (os.Args[1][:8] == "--color=" || os.Args[1][:9] == "--output=" || os.Args[1][:8] == "--align=") {
					arg = os.Args[1:]
					Filename = "standard.txt"
				} else if len(os.Args) > 3 && (os.Args[1][:8] == "--color=" || os.Args[1][:9] == "--output=" || os.Args[1][:8] == "--align=") {
					arg = os.Args[1 : len(os.Args)-1]
				} else {
					arg = os.Args[1:]
				}
			} else {
				arg = os.Args[1 : len(os.Args)-1]
			}
		} else if os.Args[len(os.Args)-1] == "thinkertoy" {
			Filename = "thinkertoy.txt"
			if len(os.Args[1]) > 9 {
				if len(os.Args) == 3 && (os.Args[1][:8] == "--color=" || os.Args[1][:9] == "--output=" || os.Args[1][:8] == "--align=") {
					arg = os.Args[1:]
					Filename = "standard.txt"
				} else if len(os.Args) > 3 && (os.Args[1][:8] == "--color=" || os.Args[1][:9] == "--output=" || os.Args[1][:8] == "--align=") {
					arg = os.Args[1 : len(os.Args)-1]
				} else {
					arg = os.Args[1:]
				}
			} else {
				arg = os.Args[1 : len(os.Args)-1]
			}
		} else {
			Filename = "standard.txt"
			arg = os.Args[1:]
		}
	} else {
		Filename = "standard.txt"
		arg = os.Args[1:]

	}
	return arg, Filename, txtfile, kil
}

func RemoveColors(char2slice [][]string) {
	for y := 0; y < 8; y++ {
		for j := 0; j < len(char2slice); j++ {
			re := regexp.MustCompile(`\x1b\[[0-9;]+m`)
			char2slice[j][y] = re.ReplaceAllString(char2slice[j][y], "")
		}
	}
}

func RemoveColorsFromString(input string) string {
	re := regexp.MustCompile(`\033\[[0-9;]+m`)
	return re.ReplaceAllString(input, "")
}

func FindOption(arg []string) ([][]string, []string, string, string, error) {
	var color, char, namefile, position string
	var colour []string
	var colors [][]string

	v := ""
Loop:
	if len(arg) > 1 {
		for i := 0; i < len(arg); i++ {
			if len(arg[i]) > 8 {
				if arg[i][:8] == "--color=" {
					color = FindColor(arg[i][8:])
					if color == "I" {
						return nil, nil, v, position, errors.New("Error")
					} else {
						if len(arg) > 2 {
							if len(arg[i+1]) > 8 {
								if arg[i+1][:8] != "--color=" && arg[i+1][:9] != "--output=" && arg[i+1][:8] != "--align=" {
									char = arg[i+1]
									colour = append(colour, color, char)
									colors = append(colors, colour)
									colour = []string{}
									arg = append(arg[:i+1], arg[i+2:]...)
									arg = append(arg[:i], arg[i+1:]...)
									goto Loop
								} else {
									colour = append(colour, color)
									colors = append(colors, colour)
									colour = []string{}
								}
							} else {
								char = arg[i+1]
								colour = append(colour, color, char)
								colors = append(colors, colour)
								colour = []string{}
								arg = append(arg[:i+1], arg[i+2:]...)
								arg = append(arg[:i], arg[i+1:]...)
								goto Loop
							}
						} else {
							colour = append(colour, color)
							colors = append(colors, colour)
							colour = []string{}
						}
					}
					if len(arg) > 1 {
						arg = append(arg[:i], arg[i+1:]...)
						goto Loop
					}
				} else if arg[i][:9] == "--output=" {
					namefile = arg[i][9:]
					if namefile == "" {
						return nil, nil, v, position, errors.New("Error")
					}
					if len(arg) > 1 {
						arg = append(arg[:i], arg[i+1:]...)
						goto Loop
					}
				} else if arg[i][:8] == "--align=" {
					if arg[i][8:] == "right" || arg[i][8:] == "left" || arg[i][8:] == "center" || arg[i][8:] == "justify" {
						position = arg[i][8:]
						arg = append(arg[:i], arg[i+1:]...)
						goto Loop
					} else {
						return nil, nil, v, position, errors.New("Error")
					}
				}
			}
		}
	}
	return colors, arg, namefile, position, nil
}

func GetTerminalWidth() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	size := strings.Split(strings.TrimSpace(string(output)), " ")
	if len(size) != 2 {
		return 0, 0, fmt.Errorf("unexpected output format from 'stty size'")
	}
	rows, err := strconv.Atoi(size[0])
	if err != nil {
		return 0, 0, err
	}
	columns, err := strconv.Atoi(size[1])
	if err != nil {
		return 0, 0, err
	}
	return rows, columns, nil
}

func CheckReverse(s string) bool {
	res := "--reverse="
	count := 0

	for _, i := range s {
		for _, j := range res {
			if i == j {
				count++
			}
		}
	}
	if count > 17 && count < 23 {
		return true
	}
	return false
}
