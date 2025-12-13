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
	rangestrs := strings.Split(lines[0], ",")
	ranges := [][2]int{}
	for _, r := range rangestrs {
		min, max := parseRange(r)
		ranges = append(ranges, [2]int{min, max})
	}

	fmt.Println(findInvalidIds1(ranges))
	fmt.Println(findInvalidIds2(ranges))
}

func findInvalidIds1(ranges [][2]int) int {
	answer := 0
	for _, r := range ranges {
		for i := r[0]; i <= r[1]; i++ {
			productIdStr := strconv.Itoa(i)
			if checkIsValid(i, len(productIdStr)/2) {
				answer += i
			}
		}
	}
	return answer
}

func findInvalidIds2(ranges [][2]int) int {
	answer := 0
	for _, r := range ranges {
		for i := r[0]; i <= r[1]; i++ {
			productIdStr := strconv.Itoa(i)
			for j := 1; j <= len(productIdStr)/2; j++ {
				if checkIsValid(i, j) {
					answer += i
					break
				}
			}
		}
	}
	return answer
}

func checkIsValid(productId int, sequenceLength int) bool {
	productIdStr := strconv.Itoa(productId)
	if sequenceLength == 0 || len(productIdStr)%sequenceLength != 0 {
		return false
	}
	match := productIdStr[:sequenceLength]
	for i := sequenceLength; i < len(productIdStr); i += sequenceLength {
		sequence := productIdStr[i : i+sequenceLength]
		if sequence != match {
			return false
		}
	}
	return true
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
