package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

type Pair struct {
	p1   Point
	p2   Point
	area int
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <filename>")
	}
	filename := os.Args[1]
	lines := readLinesFromFile(filename)
	redTiles := make([]Point, len(lines))
	for i, l := range lines {
		redTiles[i] = parsePoint(l)
	}

	fmt.Println(findMaxArea1(redTiles))
	fmt.Println(findMaxArea2(redTiles))
	fmt.Println("done")
}

func findMaxArea1(redTiles []Point) int {
	maxArea := 0
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			area := getArea(redTiles[i], redTiles[j])
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

func findMaxArea2(redTiles []Point) int {
	maxArea := 0
	borderTiles := getBorderTiles(redTiles)
	for i := 0; i < len(redTiles); i++ {
		fmt.Println(i, maxArea)
		for j := i + 1; j < len(redTiles); j++ {
			area := getArea(redTiles[i], redTiles[j])
			if area > maxArea {
				if isValidArea(redTiles[i], redTiles[j], redTiles, borderTiles) {
					maxArea = area
					fmt.Println("Valid area found:", redTiles[i], redTiles[j], area)
				}
			}
		}
	}
	return maxArea
}

func isValidArea(p1 Point, p2 Point, redTiles []Point, borderTiles map[Point]bool) bool {
	lowerX := min(p1.x, p2.x)
	higherX := max(p1.x, p2.x)
	lowerY := min(p1.y, p2.y)
	higherY := max(p1.y, p2.y)
	if !isTileInBorder(Point{x: p1.x, y: p2.y}, borderTiles) || !isTileInBorder(Point{x: p2.x, y: p1.y}, borderTiles) {
		return false
	}
	for i, _ := range redTiles {
		point1 := redTiles[i]
		point2 := redTiles[(i+1)%len(redTiles)]

		// Calculate bounding box for the segment (point1, point2)
		segLowerX := min(point1.x, point2.x)
		segHigherX := max(point1.x, point2.x)
		segLowerY := min(point1.y, point2.y)
		segHigherY := max(point1.y, point2.y)

		// Check for intersection with bounding box defined by lowerX, higherX, lowerY, higherY
		if segHigherX <= lowerX || segLowerX >= higherX || segHigherY <= lowerY || segLowerY >= higherY {
			continue
		} else {
			return false
		}
	}
	return true
}

func isTileInBorder(p Point, borderTiles map[Point]bool) bool {
	current := p
	count := 0
	if borderTiles[p] {
		return true
	}
	prevBorder := false
	for current.x >= 0 {
		if borderTiles[current] {
			if !prevBorder {
				count++
			}
			prevBorder = true
		} else {
			prevBorder = false
		}
		current.x = current.x - 1
	}
	return count%2 == 1
}

func getBorderTiles(redTiles []Point) map[Point]bool {
	borderTiles := make(map[Point]bool)
	connectTiles := func(p1 Point, p2 Point) {
		if p1.x == p2.x {
			lowerY := min(p1.y, p2.y)
			higherY := max(p1.y, p2.y)
			for y := lowerY; y <= higherY; y++ {
				borderTiles[Point{x: p1.x, y: y}] = true
			}
		} else {
			lowerX := min(p1.x, p2.x)
			higherX := max(p1.x, p2.x)
			for x := lowerX; x <= higherX; x++ {
				borderTiles[Point{x: x, y: p1.y}] = true
			}
		}
	}

	// create tile border
	for i, point := range redTiles {
		if i == 0 {
			continue
		}
		prevPoint := redTiles[i-1]
		connectTiles(point, prevPoint)
	}
	connectTiles(redTiles[0], redTiles[len(redTiles)-1])

	return borderTiles
}

func getArea(p1 Point, p2 Point) int {
	dx := p1.x - p2.x
	if dx < 0 {
		dx = -dx
	}
	dx++
	dy := p1.y - p2.y
	if dy < 0 {
		dy = -dy
	}
	dy++

	return dx * dy
}

func parsePoint(s string) Point {
	parts := strings.Split(s, ",")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	return Point{x: x, y: y}
}

func readLinesFromFile(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}
