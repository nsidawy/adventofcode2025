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
	neighbors := make(map[string][]string)
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
		neighbors[start] = ends
	}
	//fmt.Println(getPaths(neighbors, "you", "out"))
	fmt.Println(getPathsReverse(reverse, "you", "out"))
	c1 := getPathsReverse(reverse, "svr", "fft")
	fmt.Println(c1)
	c2 := getPathsReverse(reverse, "fft", "dac")
	fmt.Println(c2)
	c3 := getPathsReverse(reverse, "dac", "out")
	fmt.Println(c3)
	fmt.Println(c1 * c2 * c3)
}

func getPaths(neighbors map[string][]string, start string, end string) int {
	count := 0
	current := []string{start}
	for len(current) > 0 {
		next := []string{}
		for _, c := range current {
			if c == end {
				count++
				continue
			}
			next = append(next, neighbors[c]...)
		}
		current = next
	}
	return count
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

type Path struct {
	lastNode string
	is_fft   bool
	is_dac   bool
}

func getPaths2(neighbors map[string][]string) int {
	count := 0
	current := []Path{{lastNode: "svr", is_fft: false, is_dac: false}}
	loop := 0
	for len(current) > 0 {
		fmt.Println(loop, len(current))
		next := []Path{}
		for _, p := range current {
			if p.lastNode == "out" {
				if p.is_fft && p.is_dac {
					count++
				}
				continue
			}
			for _, n := range neighbors[p.lastNode] {
				newPath := Path{lastNode: n, is_fft: p.is_fft || n == "fft", is_dac: p.is_dac || n == "dac"}
				next = append(next, newPath)
			}
		}
		current = next
		loop++
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
