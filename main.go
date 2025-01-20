package main

import (
	"fmt"
	"log"
	"os"

	antfarm "test/antFarm"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: go run . <filename>")
	}

	farm := antfarm.NewAntFarm()
	if err := farm.ParseInput(os.Args[1]); err != nil {
		log.Fatalln(err)
	}

	moves, err := farm.SimulateMovement()
	if err != nil {
		log.Fatalln(err)
	}

	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(input) + "\n")
	fmt.Print(moves)
}
