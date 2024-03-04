package main

import (
	"fmt"
	"math"
	"os"
	"tetris/tetris"
)

func main() {
	Sample, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	if !tetris.ProSimvol(string(Sample)) {
		fmt.Println("ERROR")
		return
	}

	MassivStrok := tetris.Split(string(Sample))
	if !tetris.Kvadrat4(MassivStrok) {

		fmt.Println("ERROR")
		return
	}

	sosedy := tetris.ProSosedy(MassivStrok)

	if !tetris.ProKol(sosedy) {
		fmt.Println("ERROR")
		return
	}

	coordinates := tetris.Coordinates(MassivStrok)

	Tcor := tetris.LeftUp(coordinates)

	razmer := int(math.Ceil(math.Sqrt(float64(4 * len(Tcor)))))

	korobka := tetris.PusKorobka(razmer)

	result, err := tetris.PlaceTetrominoes(korobka, Tcor, 0)

	if err != nil {
		razmer++
		korobka = tetris.PusKorobka(razmer)
		result, err = tetris.PlaceTetrominoes(korobka, Tcor, 0)
	}

	for _, row := range result {
		fmt.Println(string(row))
	}

}
