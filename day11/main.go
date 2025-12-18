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
	reverse := make(map[string][]string)
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		start := parts[0]
		ends := strings.Split(parts[1], " ")
		if reverse[start] == nil {
			reverse[start] = []string{}
		}
		for _, e := range ends {
			if reverse[e] == nil {
				reverse[e] = []string{}
			}
			reverse[e] = append(reverse[e], start)
		}
	}

	calculatePaths := func(path []string) int {
		total := 1
		for i := 0; i < len(path)-1; i++ {
			p := getPathsReverse(reverse, path[i], path[i+1])
			total = p * total
		}
		return total
	}

	path1 := calculatePaths([]string{"you", "out"})
	fmt.Println(path1)

	path2a := calculatePaths([]string{"svr", "fft", "dac", "out"})
	path2b := calculatePaths([]string{"svr", "dac", "fft", "out"})
	fmt.Println(path2a, path2b)
}

func getPathsReverse(reverse map[string][]string, start string, end string) int {
	count := 0
	current := map[string]int{end: 1}
	for len(current) > 0 {
		next := map[string]int{}
		for k, v := range current {
			if k == start {
				count += v
				continue
			}
			for _, c := range reverse[k] {
				next[c] += v
			}
		}
		current = next
	}
	return count
}

func readLinesFromFile(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}
