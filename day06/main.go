package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Operand rune

const (
	Add      Operand = '+'
	Multiply Operand = '*'
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

	fmt.Println(calculate1(lines))
	fmt.Println(calculate2(lines))
}

func calculate1(lines []string) int {
	answer := 0
	numbers, operands := parseInput(lines)
	for i := 0; i < len(numbers[0]); i++ {
		value := numbers[0][i]
		for j := 1; j < len(numbers); j++ {
			switch operands[i] {
			case Add:
				value += numbers[j][i]
			case Multiply:
				value *= numbers[j][i]
			}
		}
		answer += value
	}
	return answer
}

func calculate2(lines []string) int {
	totalAnswer := 0
	numbers := lines[:len(lines)-1]
	operands := lines[len(lines)-1]
	endPosition := len(operands)
	for endPosition > 0 {
		nextOperatorIndex := findNextOperatorIndex(operands, endPosition)
		operator := Operand(operands[nextOperatorIndex])

		answer := 0 // Add operator case
		if operator == Multiply {
			answer = 1
		}

		for i := endPosition - 1; i >= nextOperatorIndex; i-- {
			number := readNumberColumn(numbers, i)
			switch operator {
			case Add:
				answer += number
			case Multiply:
				answer *= number
			}
		}
		totalAnswer += answer
		endPosition = nextOperatorIndex - 1
	}
	return totalAnswer
}

func readNumberColumn(numbers []string, columnIndex int) int {
	digitCount := 0
	number := 0
	for i := len(numbers) - 1; i >= 0; i-- {
		if numbers[i][columnIndex] == ' ' {
			continue
		}

		digit, _ := strconv.Atoi(string(numbers[i][columnIndex]))
		number += digit * int(math.Pow10(digitCount))
		digitCount++
	}
	return number
}

func findNextOperatorIndex(operands string, position int) int {
	for i := position - 1; i >= 0; i-- {
		if operands[i] == '+' || operands[i] == '*' {
			return i
		}
	}
	return -1
}

func parseInput(lines []string) ([][]int, []Operand) {
	numbers := [][]int{}
	for _, line := range lines[:len(lines)-1] {
		numStrings := strings.Split(line, " ")
		n := []int{}
		for _, numStr := range numStrings {
			if numStr == "" {
				continue
			}
			val, _ := strconv.Atoi(numStr)
			n = append(n, val)
		}
		numbers = append(numbers, n)
	}

	operands := []Operand{}
	for _, s := range strings.Split(lines[len(lines)-1], " ") {
		if s == "" {
			continue
		}
		operands = append(operands, Operand(s[0]))
	}
	return numbers, operands
}

func readLinesFromFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}
