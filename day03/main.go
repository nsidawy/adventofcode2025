package main

import (
	"fmt"
	"log"
	"math"
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

	banks := make([][]int, len(lines))
	for i, line := range lines {
		banks[i] = parseBank(line)
	}

	fmt.Println(findTotalJoltage(banks, 2))
	fmt.Println(findTotalJoltage(banks, 12))
}

func findTotalJoltage(banks [][]int, num int) int {
	answer := 0
	for _, bank := range banks {
		answer += findBankMaxJoltage(bank, num)
	}
	return answer
}

func findBankMaxJoltage(bank []int, num int) int {
	positions := make([]int, num)
	for i := 0; i < num; i++ {
		if i == 0 {
			positions[0] = 0
		} else {
			positions[i] = positions[i-1] + 1
		}
		offset := positions[i] + 1
		for j, battery := range bank[offset : len(bank)-(num-(i+1))] {
			if bank[positions[i]] < battery {
				positions[i] = j + offset
			}
		}
	}
	value := 0
	for i := 0; i < num; i++ {
		value += bank[positions[i]] * int(math.Pow10(num-i-1))
	}
	fmt.Println(positions, value)
	return value
}

func parseBank(bank string) []int {
	batteries := []int{}
	for _, battery := range bank {
		num, _ := strconv.Atoi(string(battery))
		batteries = append(batteries, num)
	}
	return batteries
}

func readLinesFromFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(content)), "\n"), nil
}
