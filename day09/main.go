package main

import (
	"fmt"
	"log"
	"os"
	"sort"
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
	//getTilesMap(redTiles)
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
	maxArea := 2000000000
	emptyTiles := getTilesMap(redTiles)
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			area := getArea(redTiles[i], redTiles[j])
			if area > maxArea && isValidArea(redTiles[i], redTiles[j], emptyTiles) {
				maxArea = area
				fmt.Println(redTiles[i], redTiles[j], area)
			}
		}
	}
	return maxArea
}

func findMaxArea3(redTiles []Point) int {
	emptyTiles := getTilesMap(redTiles)
	pairs := make([]Pair, 0, len(redTiles)*(len(redTiles)-1)/2)
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			pairs = append(pairs, Pair{p1: redTiles[i], p2: redTiles[j], area: getArea(redTiles[i], redTiles[j])})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].area > pairs[j].area
	})
	fmt.Println(len(pairs))
	for i, pair := range pairs {
		if i%1000 == 0 {
			fmt.Println(i, pair.area)
		}
		if isValidArea(pair.p1, pair.p2, emptyTiles) {
			return pair.area
		}
	}
	return 0
}

func isValidArea(p1 Point, p2 Point, emptyTiles [][]bool) bool {
	lowerX := min(p1.x, p2.x)
	higherX := max(p1.x, p2.x)
	lowerY := min(p1.y, p2.y)
	higherY := max(p1.y, p2.y)
	for x := lowerX; x <= higherX; x++ {
		for y := lowerY; y <= higherY; y++ {
			if emptyTiles[x][y] {
				return false
			}
		}
	}
	return true
}

func isValidArea2(p1 Point, p2 Point, emptyTiles map[Point]bool) bool {
	lowerX := min(p1.x, p2.x)
	higherX := max(p1.x, p2.x)
	lowerY := min(p1.y, p2.y)
	higherY := max(p1.y, p2.y)

	for e := range emptyTiles {
		if e.x >= lowerX && e.x <= higherX && e.y >= lowerY && e.y <= higherY {
			return false
		}
	}

	return true
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

func getTilesMap(redTiles []Point) [][]bool {
	minX := 1000000000
	maxX := 0
	minY := 1000000000
	maxY := 0
	for _, tile := range redTiles {
		if tile.x > maxX {
			maxX = tile.x
		}
		if tile.x < minX {
			minX = tile.x
		}
		if tile.y > maxY {
			maxY = tile.y
		}
		if tile.y < minY {
			minY = tile.y
		}
	}
	maxY++
	maxX++
	minX--
	minY--

	borderTiles := getBorderTiles(redTiles)

	emptyTiles := make([][]bool, maxX+1)
	for i, _ := range emptyTiles {
		emptyTiles[i] = make([]bool, maxY+1)
	}
	queue := make([]Point, maxX*maxY/10)
	current := 0
	end := 0
	addToQueue := func(p Point) {
		if p.y < maxY {
			if !borderTiles[Point{x: p.x, y: p.y + 1}] && !emptyTiles[p.x][p.y+1] {
				queue[end] = Point{x: p.x, y: p.y + 1}
				end++
				emptyTiles[p.x][p.y+1] = true
			}
		}
		if p.y > minY {
			if !borderTiles[Point{x: p.x, y: p.y - 1}] && !emptyTiles[p.x][p.y-1] {
				queue[end] = Point{x: p.x, y: p.y - 1}
				end++
				emptyTiles[p.x][p.y-1] = true
			}
		}
		if p.x < maxX {
			if !borderTiles[Point{x: p.x + 1, y: p.y}] && !emptyTiles[p.x+1][p.y] {
				queue[end] = Point{x: p.x + 1, y: p.y}
				end++
				emptyTiles[p.x+1][p.y] = true
			}
		}
		if p.x > minX {
			if !borderTiles[Point{x: p.x - 1, y: p.y}] && !emptyTiles[p.x-1][p.y] {
				queue[end] = Point{x: p.x - 1, y: p.y}
				end++
				emptyTiles[p.x-1][p.y] = true
			}
		}
	}
	emptyTiles[minX][minY] = true
	addToQueue(Point{x: minX, y: minY})
	end++
	for current < end {
		p := queue[current]
		current++
		addToQueue(p)
		if current%1000000 == 0 {
			fmt.Println(current, end, end-current)
		}
	}

	fmt.Println("Empty Tiles", len(emptyTiles))
	return emptyTiles
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
