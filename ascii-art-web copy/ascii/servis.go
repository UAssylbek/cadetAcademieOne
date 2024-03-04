package ascii

import (
	"fmt"
	"os"
	"strings"
)

func AsciiImpl(text string, style string) (string, string) {
	teextt, latin := CheckLatin(text)
	if latin != "" {
		return "", "no Latin"
	}
	Banner, err := os.ReadFile(style)
	if err != nil {
		fmt.Println("ERROR:", err)
		return "", "no banner"
	}
	Shifr := StringTo2DArray(string(Banner))

	Shifr[0] = AntiFirst(Shifr[0])
	soz := ""
	char2slice := [][]string{}
	i := 0
	for i < len(teextt) {
		if teextt[i] == "" {
			soz += "\n"
		} else {
			if len(teextt[i]) != 0 {
				for _, l := range teextt[i] {
					if l != '\n' {
						char2slice = append(char2slice, Shifr[(l-32)])
					}
				}
			}
			for k := 0; k < 8; k++ {
				for j := 0; j < len(char2slice); j++ {
					soz += char2slice[j][k]
				}
				soz += "\n"
			}
		}
		char2slice = [][]string{}
		i++
	}
	return soz, ""
}

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

func CheckLatin(str string) ([]string, string) {
	str = strings.ReplaceAll(str, "\r", "")
	//check for printable characters
	for i := 0; i < len(str); i++ {
		if (str[i] < 32 || str[i] > 126) && str[i] != 10 {
			return nil, "not Latin"
		}
	}

	wordslice := strings.Split(str, "\n")
	return wordslice, ""
}

func AntiFirst(s []string) []string {
	s = s[1:]
	return s
}
