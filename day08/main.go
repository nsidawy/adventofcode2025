package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
	z int
}

type Pair struct {
	points   [2]Point
	distance float64
}

type Circuit map[Point]bool

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <filename> <target>")
	}
	filename := os.Args[1]
	target, _ := strconv.Atoi(os.Args[2])
	lines := readLinesFromFile(filename)

	points := make([]Point, len(lines))
	for i, l := range lines {
		points[i] = parsePoint(l)
	}
	pairs := getPairs(points)

	circuits, _ := connectCircuit(points, pairs, target)
	fmt.Println(len(circuits[0]) * len(circuits[1]) * len(circuits[2]))

	// set target high enoguh to guarantee all points are in a single circuit
	_, pair := connectCircuit(points, pairs, len(pairs))
	fmt.Println(pair.points[0].x * pair.points[1].x)
}

func connectCircuit(points []Point, pairs []Pair, target int) ([]Circuit, Pair) {
	circuits := getCircuits(points)
	var pair Pair
	for i := 0; i < target; i++ {
		pair = pairs[i]
		circuit1Index := -1
		circuit2Index := -1
		for i, circuit := range circuits {
			if circuit[pair.points[0]] {
				circuit1Index = i
			}
			if circuit[pair.points[1]] {
				circuit2Index = i
			}
			if circuit1Index >= 0 && circuit2Index >= 0 {
				break
			}
		}
		if circuit1Index == circuit2Index {
			continue
		}

		for p, _ := range circuits[circuit2Index] {
			circuits[circuit1Index][p] = true
		}
		circuits[circuit2Index] = circuits[len(circuits)-1]
		circuits = circuits[:len(circuits)-1]

		if len(circuits) == 1 {
			break
		}
	}

	sort.Slice(circuits, func(i, j int) bool {
		return len(circuits[i]) > len(circuits[j])
	})
	//for _, circuit := range circuits {
	//	fmt.Println(len(circuit), circuit)
	//}
	return circuits, pair
}

func getPairs(points []Point) []Pair {
	pairs := make([]Pair, (len(points)*(len(points)-1))/2)
	pos := 0
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			distance := getDistance(points[i], points[j])
			pairs[pos] = Pair{
				points:   [2]Point{points[i], points[j]},
				distance: distance,
			}
			pos++
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].distance < pairs[j].distance
	})
	return pairs
}

func getCircuits(points []Point) []Circuit {
	circuits := make([]Circuit, len(points))
	for i, p := range points {
		circuits[i] = make(Circuit)
		circuits[i][p] = true
	}
	return circuits
}

func getDistance(a Point, b Point) float64 {
	return math.Sqrt(math.Pow(float64(b.x-a.x), 2) + math.Pow(float64(b.y-a.y), 2) + math.Pow(float64(b.z-a.z), 2))
}

func parsePoint(s string) Point {
	pointStrs := strings.Split(s, ",")
	x, _ := strconv.Atoi(pointStrs[0])
	y, _ := strconv.Atoi(pointStrs[1])
	z, _ := strconv.Atoi(pointStrs[2])
	return Point{x: x, y: y, z: z}
}

func readLinesFromFile(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}
