package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		return
	}

	data, _ := os.ReadFile(os.Args[1])
	outFileName := os.Args[2]

	str := SplitData(string(data))

	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == "" {
			str = append(str[:i], str[i+1:]...)
		}
	}
	pattern := regexp.MustCompile(`[0-9]\d{0,19}\)`)
	var stor string
	var soz string
	var koz string
	var boz string
	count := 0

Loop:
	for i := 0; i < len(str); i++ {
		matched := pattern.MatchString(str[i])
		if str[i] == "(cap)" {
			if str[i] != str[0] {
				count1 := 1
				for l := len(str[:i]) - 1; l >= 0; l-- {
					if !SortBonus(FirstRune(str[l])) && str[l] != "'" && str[l] != "\n" {
						str[l] = Capitalize(str[l])
						count++
					}
					if count1 == count {
						count = 0
						break
					}
				}
				str = append(str[:i], str[i+1:]...)
				goto Loop
			}

		} else if str[i] == "(up)" {
			if str[i] != str[0] {
				count1 := 1
				for l := len(str[:i]) - 1; l >= 0; l-- {
					if !SortBonus(FirstRune(str[l])) && str[l] != "'" && str[l] != "\n" {
						str[l] = ToUpper(str[l])
						count++
					}
					if count1 == count {
						count = 0
						break
					}
				}
				str = append(str[:i], str[i+1:]...)
				goto Loop
			}
		} else if str[i] == "(low)" {
			if str[i] != str[0] {
				count1 := 1
				for l := len(str[:i]) - 1; l >= 0; l-- {
					if !SortBonus(FirstRune(str[l])) && str[l] != "'" && str[l] != "\n" {
						str[l] = ToLower(str[l])
						count++
					}
					if count1 == count {
						count = 0
						break
					}
				}
				str = append(str[:i], str[i+1:]...)
				goto Loop
			}
		} else if str[i] == "(hex)" {
			if str[i] != str[0] {
				count1 := 1
				for l := len(str[:i]) - 1; l >= 0; l-- {
					if !SortBonus(FirstRune(str[l])) && str[l] != "'" && str[l] != "\n" {
						pattern2 := regexp.MustCompile(`^[1-9a-fA-F]+`)
						matched := pattern2.MatchString(str[l])
						if matched {
							num := str[l]
							bum := AtoiBase(num, "0123456789ABCDEF")
							str[l] = bum
						}
						count++
					}
					if count1 == count {
						count = 0
						break
					}
				}
				str = append(str[:i], str[i+1:]...)
				goto Loop
			}

		} else if str[i] == "(bin)" {
			if str[i] != str[0] {
				count1 := 1
				for l := len(str[:i]) - 1; l >= 0; l-- {
					if !SortBonus(FirstRune(str[l])) && str[l] != "'" && str[l] != "\n" {
						pattern3 := regexp.MustCompile(`[0-9]\d{0,19}`)
						matched := pattern3.MatchString(str[l])
						if matched {
							num := str[l]
							bum, err := strconv.ParseInt(num, 2, 64)
							if err != nil {
								break
							}
							str[l] = strconv.Itoa(int(bum))
						}
						count++
					}
					if count1 == count {
						count = 0
						break
					}
				}
				str = append(str[:i], str[i+1:]...)
				goto Loop
			}

		} else if Sort(FirstRune(str[i])) {
			if len(str[i]) > 1 {
				if str[i] != str[0] {
					if str[i-1] == "a" || str[i-1] == "an" {
						str[i-1] = "an"
					} else if str[i-1] == "A" || str[i-1] == "An" {
						str[i-1] = "An"
					}
				}
			}
		} else if AntiSort(FirstRune(str[i])) {
			if len(str[i]) > 1 {
				if str[i] != str[0] {
					if str[i-1] == "a" || str[i-1] == "an" {
						str[i-1] = "a"
					} else if str[i-1] == "A" || str[i-1] == "An" || str[i-1] == "AN" {
						str[i-1] = "A"
					}
				}
			}
		} else if matched {
			if str[i] != str[0] {
				stor = str[i-1]
				if stor == "(low," {
					koz = str[i]
					num, _ := strconv.Atoi(string(NumSort(koz)))
					if len(str[:i])-2 < num {
						for l := 0; l < len(str[:i]); l++ {
							str[l] = ToLower(str[l])
						}
					} else {
						for l := len(str[:i]) - 2; l >= 0; l-- {
							if !SortBonus(FirstRune(str[l])) && str[l] != "'" && str[l] != "\n" {
								str[l] = ToLower(str[l])
								count++
							}
							if num == count {
								count = 0
								break
							}
						}
						str = append(str[:i-1], str[i+1:]...)
						goto Loop
					}
				} else if stor == "(up," {
					boz = str[i]
					num, _ := strconv.Atoi(string(NumSort(boz)))
					if len(str[:i])-2 < num {
						for l := 0; l < len(str[:i]); l++ {
							str[l] = ToUpper(str[l])
						}
					} else {
						for l := len(str[:i]) - 2; l >= 0; l-- {
							if !SortBonus(FirstRune(str[l])) && str[l] != "'" && str[l] != "\n" {
								str[l] = ToUpper(str[l])
								count++
							}
							if num == count {
								count = 0
								break
							}
						}
						str = append(str[:i-1], str[i+1:]...)
						goto Loop
					}
				} else if stor == "(cap," {
					soz = str[i]
					num, _ := strconv.Atoi(string(NumSort(soz)))
					if len(str[:i])-2 < num {
						for l := 0; l < len(str[:i]); l++ {
							str[l] = Capitalize(str[l])
						}
						str = append(str[:i-1], str[i+1:]...)
						goto Loop
					} else {
						for l := len(str[:i]) - 2; l >= 0; l-- {
							if !SortBonus(FirstRune(str[l])) && str[l] != "'" && str[l] != "\n" {
								str[l] = Capitalize(str[l])
								count++
							}
							if num == count {
								count = 0
								break
							}
						}
						str = append(str[:i-1], str[i+1:]...)
						goto Loop
					}

				}
			}
		}
	}

	for i, c := range str {
		if SortBonus(FirstRune(c)) {
			if str[i] != str[0] {
				if SortBonus(SecondRune(c)) {
					str[i-1] = str[i-1] + c
					str[i] = ""
				} else {
					first := FirstRune(c)
					str[i-1] = str[i-1] + string(first)
					str[i] = AntiFirst(str[i])
				}
			}
		}
	}

	word := make([]string, 0)
	for i := 0; i < len(str); i++ {
		if str[i] != "" {
			word = append(word, str[i])
		}
	}

	var answer string
	if len(word) > 0 {
		for i := 0; i < len(word)-1; i++ {
			if SortBonus(FirstRune(word[i+1])) {
				answer += word[i]
			} else if word[i] == "\n" {
				answer += word[i]
			} else {
				answer += word[i] + " "
			}

		}
		answer += word[len(word)-1]
	}

	result := ""
	if len(answer) > 0 {
		for i := 0; i < len(answer)-1; i++ {
			if answer[i] == ' ' {
				if SortBonus(rune(answer[i+1])) {
					result = result + ""
				} else {
					result = result + string(answer[i])
				}
			} else if SortBonus(rune(answer[i])) {
				if SortBonus(rune(answer[i+1])) || answer[i+1] == ' ' {
					result = result + string(answer[i])
				} else {
					result = result + string(answer[i]) + " "
				}
			} else {
				result = result + string(answer[i])
			}
		}
		result += string(answer[len(answer)-1])
	}

	result = fixQuotes(result)
	result = strings.ReplaceAll(result, "( ", "(")
	fileName := outFileName
	err := os.WriteFile(fileName, []byte(result), 0644)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}

	fmt.Println("Строка успешно записана в файл:", fileName)
}

func SplitWhiteSpaces(s string) []string {
	var words []string
	wordStart := -1

	for i, r := range s {
		if r == ' ' || r == '\t' || r == '\n' {
			if wordStart != -1 {
				words = append(words, s[wordStart:i])
				wordStart = -1
			}
		} else if wordStart == -1 {
			wordStart = i
		}
	}

	if wordStart != -1 {
		words = append(words, s[wordStart:])
	}

	return words

}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	s = strings.ToLower(s)
	s = strings.ToUpper(s[0:1]) + s[1:]
	return s
}

func ToUpper(s string) string {
	k := []rune(s)
	for i := 0; i < len(k); i++ {
		if k[i] >= 'a' && k[i] <= 'z' {
			k[i] = k[i] - 32
		}
	}
	return string(k)
}

func ToLower(s string) string {
	k := []rune(s)
	for i := 0; i < len(k); i++ {
		if k[i] >= 'A' && k[i] <= 'Z' {
			k[i] = k[i] + 32
		}
	}
	return string(k)
}

func NumSort(s string) string {
	var a string

	for _, i := range s {
		if i != ')' {
			a += string(i)
		}
	}
	return a
}

func AtoiBase(s string, base string) string {
	res := 0
	s = string(ToUpper(s))
	n := len(base)
	m := 1
	for i := len(s) - 1; i >= 0; i-- {
		idx := find(base, rune(s[i]))
		res += m * idx
		m *= n
	}
	ans := strconv.Itoa(res)

	return ans
}

func find(base string, x rune) int {
	for i, c := range base {
		if c == x {
			return i
		}
	}
	return -1
}

func FirstRune(s string) rune {
	res := []rune(s)
	if len(s) == 0 {
		return -1
	}
	return res[0]
}

func SecondRune(s string) rune {
	res := []rune(s)
	if len(s) > 1 {
		return res[1]
	}
	return -1
}

func AntiFirst(s string) string {
	if len(s) == 1 {
		return ""
	}
	s = s[1:]
	return s
}

func Sort(s rune) bool {
	a := "aeiouhAEIOUH"
	for _, i := range a {
		if i == s {
			return true
		}
	}
	return false
}

func AntiSort(s rune) bool {
	a := "bcdfgjklmnpqrstvwxyzBCDFGJKLMNPQRSTVWXYZ"

	for _, i := range a {
		if i == s {
			return true
		}
	}
	return false
}

func SortBonus(s rune) bool {
	a := ",.!?:;)"

	for _, i := range a {
		if i == s {
			return true
		}
	}
	return false
}

func fixQuotes(text string) string {
	re := regexp.MustCompile(`'\s*([^']+?)\s*'`)
	fixedText := re.ReplaceAllString(text, "'$1'")
	return fixedText
}

func addSpaces(input string) string {
	punctuation := []string{".", ",", ":", "!", "?"}

	for _, p := range punctuation {
		input = strings.ReplaceAll(input, p, p+" ")
	}

	return input
}

func SplitData(data string) []string {
	regex := regexp.MustCompile(`\b\w*-\w*\w*\b|(\d+)\)|\b\w*'\w*\w*\b|\((\w+)\)|\(\w+\,|\(\w+,\s*\w+\)|\b\w+\b|[^\w\s]|\n|`)
	words := regex.FindAllString(data, -1)
	return words
}
