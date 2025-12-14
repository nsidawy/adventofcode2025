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
	lines, err := readLinesFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	grid := makeGrid(lines)

	fmt.Println(len(getAccessible(grid)))
	fmt.Println(countTotalRemoved(grid))
}

func countTotalRemoved(grid [][]bool) int {
	totalRemoved := 0
	accessible := getAccessible(grid)
	for len(accessible) > 0 {
		totalRemoved += len(accessible)
		for _, pos := range accessible {
			grid[pos[0]][pos[1]] = false
		}
		accessible = getAccessible(grid)
	}
	return totalRemoved
}

func getAccessible(grid [][]bool) [][2]int {
	accessible := [][2]int{}
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] && countAdjacents(grid, i, j) < 4 {
				accessible = append(accessible, [2]int{i, j})
			}
		}
	}
	return accessible
}

func makeGrid(lines []string) [][]bool {
	width := len(lines[0])
	height := len(lines)
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}
	for i, line := range lines {
		for j, char := range line {
			grid[i][j] = char == '@'
		}
	}
	return grid
}

func countAdjacents(grid [][]bool, x int, y int) int {
	count := 0
	if x > 0 {
		if grid[x-1][y] {
			count++
		}
		if y > 0 && grid[x-1][y-1] {
			count++
		}
		if y < len(grid[0])-1 && grid[x-1][y+1] {
			count++
		}
	}
	if x < len(grid)-1 {
		if grid[x+1][y] {
			count++
		}
		if y > 0 && grid[x+1][y-1] {
			count++
		}
		if y < len(grid[0])-1 && grid[x+1][y+1] {
			count++
		}
	}
	if y > 0 && grid[x][y-1] {
		count++
	}
	if y < len(grid[0])-1 && grid[x][y+1] {
		count++
	}
	return count
}

func readLinesFromFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(content)), "\n"), nil
}
