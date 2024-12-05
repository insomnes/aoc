package solution

import (
	"bufio"
	"log"
	"os"
)

func ReadInput(fp string) ([]string, error) {
	defer Track("Input")()
	log.Println("Reading", fp)
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
