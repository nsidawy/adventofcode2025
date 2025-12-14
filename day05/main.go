package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	min int
	max int
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <filename>")
	}
	filename := os.Args[1]
	lines, err := readLinesFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	ranges, ingredients := parseInput(lines)
	freshIngredients := getFreshIngredients(ranges, ingredients)
	fmt.Println(len(freshIngredients))
	fmt.Println(getUniqueFreshIngredients(ranges))
}

func getUniqueFreshIngredients(ranges []Range) int {
	count := 0
	sort.Slice(ranges, func(i, j int) bool {
		gapI := ranges[i].max - ranges[i].min
		gapj := ranges[j].max - ranges[j].min
		return gapI < gapj
	})
	for i, r := range ranges {
		for _, r2 := range ranges[i+1:] {
			if r2.inRange(r.min) && r2.inRange(r.max) {
				r.max = 0
				r.min = 1
				break
			} else if r.min <= r2.max && r2.min <= r.max {
				if r.max < r2.max {
					r.max = r2.min - 1
				}
				if r.min > r2.min {
					r.min = r2.max + 1
				}
			}
			if r.max < r.min {
				break
			}
		}

		if r.max >= r.min {
			count += r.max - r.min + 1
		}
	}

	return count
}

func getFreshIngredients(ranges []Range, ingredients []int) []int {
	freshIngredients := []int{}
	for _, ingredient := range ingredients {
		isIngredientFresh := false
		for _, r := range ranges {
			if r.inRange(ingredient) {
				if !isIngredientFresh {
					freshIngredients = append(freshIngredients, ingredient)
					break
				}
			}
		}
	}
	return freshIngredients
}

func (r Range) inRange(i int) bool {
	return r.min <= i && i <= r.max
}

func parseInput(lines []string) ([]Range, []int) {
	ranges := []Range{}
	i := 0
	for lines[i] != "" {
		min, max := parseRange(lines[i])
		ranges = append(ranges, Range{min, max})
		i++
	}
	i++

	numbers := []int{}
	for i < len(lines) {
		num, _ := strconv.Atoi(lines[i])
		numbers = append(numbers, num)
		i++
	}
	return ranges, numbers
}

func parseRange(s string) (min, max int) {
	parts := strings.Split(s, "-")
	min, _ = strconv.Atoi(parts[0])
	max, _ = strconv.Atoi(parts[1])
	return
}

func readLinesFromFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(content)), "\n"), nil
}
