package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error")
		return
	}

	datastring, err := os.ReadFile(os.Args[1])
	if len(datastring) == 0 {
		fmt.Println("ERROR")
		return
	}
	if err != nil {
		fmt.Println("ERROR")
		return
	}
	var dataint []int
	data := strings.Split(string(datastring), "\n")
	if data[len(data)-1] == "" {
		data = data[:len(data)-1]
	}
	for _, i := range data {
		j, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println("ERROR")
			return
		} else {
			dataint = append(dataint, j)
		}
	}
	aa := Average(dataint)
	vra, dev := Variance(int(aa), dataint)
	fmt.Println("Avarege:", Average(dataint))
	fmt.Println("Median:", Median(dataint))
	fmt.Println("Variance:", int(vra))
	fmt.Println("Standard Deviation:", int(dev))
}

func Average(s []int) float64 {
	var num float64
	for _, i := range s {
		num += float64(i)
	}
	q := math.Round(num / float64(len(s)))

	return q
}

func Median(s []int) float64 {
	sort.Ints(s)
	var j float64
	if len(s)%2 == 0 {
		j = float64(s[len(s)/2] + s[(len(s)/2)-1])
		j = j / 2
		k := math.Round(j)
		return float64(k)
	}

	return float64(s[len(s)/2])
}

func Variance(num int, s []int) (float64, float64) {
	var nums float64
	for i := 0; i < len(s); i++ {
		nums += math.Pow((float64(s[i]) - float64(num)), 2)
	}
	nums = nums / float64(len(s))
	w := math.Round(nums)
	Q := math.Sqrt(float64(nums))
	e := math.Round(Q)
	return float64(w), float64(e)
}
