package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <filename>")
	}
	filename := os.Args[1]
	lines, err := readLinesFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	calculatePassword1(lines)
	calculatePassword2(lines)
}

func calculatePassword1(lines []string) {
	position := 50
	password := 0
	for _, line := range lines {
		distance, _ := strconv.Atoi(line[1:])
		if line[0] == 'L' {
			position -= distance
		} else {
			position += distance
		}

		if position%100 == 0 {
			password++
		}
	}
	fmt.Println(password)
}

func calculatePassword2(lines []string) {
	position := 50
	password := 0
	for _, line := range lines {
		distance, _ := strconv.Atoi(line[1:])
		direction := line[0]
		for range distance {
			if direction == 'L' {
				position--
			} else {
				position++
			}
			switch position {
			case 100:
				position = 0
			case -1:
				position = 99
			}

			if position == 0 {
				password++
			}
		}
	}
	fmt.Println(password)
}

func readLinesFromFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(content)), "\n"), nil
}
