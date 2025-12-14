package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <filename>")
	}
	filename := os.Args[1]
	lines := readLinesFromFile(filename)

	startPosition := strings.IndexRune(lines[0], 'S')
	fmt.Println(countSplits(lines, startPosition))
	fmt.Println(countTimelines(lines, startPosition, 0, make(map[[2]int]int)))
}

func countSplits(lines []string, startPosition int) int {
	count := 0
	positions := make([]bool, len(lines[0]))
	positions[startPosition] = true
	for i := 0; i < len(lines); i++ {
		nextPositions := make([]bool, len(lines[0]))
		for j := 0; j < len(lines[0]); j++ {
			if !positions[j] {
				continue
			}
			if lines[i][j] == '^' {
				nextPositions[j-1] = true
				nextPositions[j+1] = true
				count++
			} else {
				nextPositions[j] = true
			}
		}
		positions = nextPositions
	}
	return count
}

func countTimelines(lines []string, beamPosition int, rowIndex int, memo map[[2]int]int) int {
	if rowIndex == len(lines) {
		return 1
	}
	if memo[[2]int{rowIndex, beamPosition}] != 0 {
		return memo[[2]int{rowIndex, beamPosition}]
	}
	count := 0
	if lines[rowIndex][beamPosition] == '^' {
		count += countTimelines(lines, beamPosition-1, rowIndex+1, memo)
		count += countTimelines(lines, beamPosition+1, rowIndex+1, memo)
	} else {
		count = countTimelines(lines, beamPosition, rowIndex+1, memo)
	}
	memo[[2]int{rowIndex, beamPosition}] = count
	return count
}

func readLinesFromFile(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}
