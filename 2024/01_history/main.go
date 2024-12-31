package main

import (
	"fmt"
	"log"
	"os"

	"history/solution"
)

const (
	testPath   = "test.txt"
	normalPath = "input.txt"
)

func main() {
	defer solution.Track("main")()

	inpPath := testPath
	if len(os.Args) > 1 {
		inpPath = normalPath
	}

	inpLines, err := solution.ReadInput(inpPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to read input %v\n", err))
	}

	inp, err := solution.ParseInput(inpLines)
	if err != nil {
		log.Fatal(fmt.Sprintf("Bad input: %v\n", err))
	}

	resOne := solution.PartOne(inp)
	log.Println("Result one", resOne)

	resTwo := solution.PartTwo(inp)
	log.Println("Result two", resTwo)
}
