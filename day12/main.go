package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Shape struct {
	areas [][3][3]bool
	size  int
	id    int
}

type Region struct {
	width  int
	height int
	shapes []int
}

type Score struct {
	id         int
	value      float32
	placements []Placement
}

type Placement struct {
	orientation int
	location    int
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <filename>")
	}
	filename := os.Args[1]
	lines := readLinesFromFile(filename)
	shapes, regions := parseLines(lines)
	for _, s := range shapes {
		fmt.Println("shape id", s.id)
		for i, o := range s.areas {
			fmt.Println("orientation", i)
			printShape(o)
		}
	}
	count := 0
	for _, region := range regions {
		if isRegionValid(region, shapes) {
			count++
		}
	}
	fmt.Println(count)
}

func isRegionValid(region Region, shapes []Shape) bool {
	grid := make([]bool, region.height*region.width)
	isValid := isRegionValidStep(grid, region, shapes)
	return isValid
}

func isRegionValidStep(grid []bool, region Region, shapes []Shape) bool {
	isComplete := true
	for _, r := range region.shapes {
		if r != 0 {
			isComplete = false
			break
		}
	}
	if isComplete {
		return true
	}

	// Check if any have less fit than remaining
	scores := calculateShapeScores(grid, region, shapes)
	fmt.Println("--------------\n", scores)
	fmt.Println(region.shapes)
	for _, score := range scores {
		if score.value < 1.0 {
			return false
		}
	}

	for _, s := range scores {
		for _, p := range s.placements {
			//fmt.Println("======New Grid=======", s.id, p)
			//fmt.Println(region.shapes)
			//printShape(shapes[s.id].areas[p.orientation])
			//printGridWithShape(grid, region.height, region.width, shapes[s.id].areas[p.orientation], p)
			applyShapePlacement(grid, region, shapes[s.id], p, true)
			region.shapes[s.id]--
			if isRegionValidStep(grid, region, shapes) {
				return true
			}
			applyShapePlacement(grid, region, shapes[s.id], p, false)
			//fmt.Println("======Backtrack Grid=======")
			//printGrid(grid, region.height, region.width)
			region.shapes[s.id]++
		}
	}

	return false
}

func applyShapePlacement(grid []bool, region Region, shape Shape, placement Placement, isAdd bool) {
	for l := 0; l < 3; l++ {
		for m := 0; m < 3; m++ {
			if shape.areas[placement.orientation][l][m] {
				if isAdd {
					grid[placement.location+l*region.width+m] = true
				} else {
					grid[placement.location+l*region.width+m] = false
				}
			}
		}
	}
}

func calculateShapeScores(grid []bool, region Region, shapes []Shape) []Score {
	scores := []Score{}
	for _, shape := range shapes {
		placements := []Placement{}
		if region.shapes[shape.id] == 0 {
			continue
		}
		for a, area := range shape.areas {
			for j := 0; j < region.height-2; j++ {
				for k := 0; k < region.width-2; k++ {
					if isShapeFits(grid, j, k, area, region) {
						//if a == 5 && shape.id == 4 && j*region.width+k == 24 && region.shapes[0] == 1 && region.shapes[2] == 0 {
						//	fmt.Println("FIND ME")
						//	printGrid(grid, region.height, region.width)
						//	printGridWithShape(grid, region.height, region.width, area, Placement{orientation: a, location: j*region.width + k})
						//}
						placements = append(placements, Placement{orientation: a, location: j*region.width + k})
					}
				}
			}
		}

		scores = append(scores, Score{id: shape.id, value: float32(len(placements)) / float32(region.shapes[shape.id]), placements: placements})
	}

	sort.Slice(scores, func(a int, b int) bool {
		return scores[a].value > scores[b].value
	})
	return scores
}

func isShapeFits(grid []bool, h int, w int, area [3][3]bool, region Region) bool {
	for l := 0; l < 3; l++ {
		for m := 0; m < 3; m++ {
			if area[l][m] && grid[(h+l)*region.width+w+m] {
				return false
			}
		}
	}
	return true
}

func printGrid(grid []bool, height int, width int) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if grid[i*width+j] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func printGridWithShape(grid []bool, height int, width int, area [3][3]bool, placement Placement) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			spot := i*width + j
			if grid[spot] {
				fmt.Print("#")
			} else if spot == placement.location && area[0][0] {
				fmt.Print("*")
			} else if spot == placement.location+1 && area[0][1] {
				fmt.Print("*")
			} else if spot == placement.location+2 && area[0][2] {
				fmt.Print("*")
			} else if spot == placement.location+width && area[1][0] {
				fmt.Print("*")
			} else if spot == placement.location+width+1 && area[1][1] {
				fmt.Print("*")
			} else if spot == placement.location+width+2 && area[1][2] {
				fmt.Print("*")
			} else if spot == placement.location+width*2 && area[2][0] {
				fmt.Print("*")
			} else if spot == placement.location+width*2+1 && area[2][1] {
				fmt.Print("*")
			} else if spot == placement.location+width*2+2 && area[2][2] {
				fmt.Print("*")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func printShape(area [3][3]bool) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if area[i][j] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func parseLines(lines []string) ([]Shape, []Region) {
	i := 0
	shapes := []Shape{}
	id := 0
	for true {
		matched, _ := regexp.MatchString(`^(\d+):$`, lines[i])
		if !matched {
			break
		}
		area := make([][3][3]bool, 8)
		for j := 0; j < 3; j++ {
			i++
			for k := 0; k < 3; k++ {
				area[0][j][k] = lines[i][k] == '#'
			}
		}

		for k := 1; k < 4; k++ {
			area[k][0][0] = area[k-1][2][0]
			area[k][1][0] = area[k-1][2][1]
			area[k][2][0] = area[k-1][2][2]
			area[k][0][1] = area[k-1][1][0]
			area[k][1][1] = area[k-1][1][1]
			area[k][2][1] = area[k-1][1][2]
			area[k][0][2] = area[k-1][0][0]
			area[k][1][2] = area[k-1][0][1]
			area[k][2][2] = area[k-1][0][2]
		}
		area[4][0][0] = area[0][0][2]
		area[4][1][0] = area[0][1][2]
		area[4][2][0] = area[0][2][2]
		area[4][0][1] = area[0][0][1]
		area[4][1][1] = area[0][1][1]
		area[4][2][1] = area[0][2][1]
		area[4][0][2] = area[0][0][0]
		area[4][1][2] = area[0][1][0]
		area[4][2][2] = area[0][2][0]
		for k := 5; k < 8; k++ {
			area[k][0][0] = area[k-1][2][0]
			area[k][1][0] = area[k-1][2][1]
			area[k][2][0] = area[k-1][2][2]
			area[k][0][1] = area[k-1][1][0]
			area[k][1][1] = area[k-1][1][1]
			area[k][2][1] = area[k-1][1][2]
			area[k][0][2] = area[k-1][0][0]
			area[k][1][2] = area[k-1][0][1]
			area[k][2][2] = area[k-1][0][2]
		}
		// Deduplicate areas
		area = deduplicateAreas(area)

		size := 0
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				if area[0][j][k] {
					size++
				}
			}
		}
		shape := Shape{areas: area, size: size, id: id}
		shapes = append(shapes, shape)
		i += 2
		id++
	}
	regions := []Region{}
	for _, l := range lines[i:] {
		parts := strings.Split(l, ": ")
		dimensionStrings := strings.Split(parts[0], "x")
		width, _ := strconv.Atoi(dimensionStrings[0])
		height, _ := strconv.Atoi(dimensionStrings[1])
		shapeStrings := strings.Split(parts[1], " ")
		shapes := make([]int, len(shapeStrings))
		for i, s := range shapeStrings {
			shapes[i], _ = strconv.Atoi(s)
		}
		region := Region{
			width:  width,
			height: height,
			shapes: shapes,
		}
		regions = append(regions, region)
	}

	return shapes, regions
}

func deduplicateAreas(areas [][3][3]bool) [][3][3]bool {
	unique := [][3][3]bool{}
	for _, area := range areas {
		isDuplicate := false
		for _, existing := range unique {
			if areasEqual(area, existing) {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			unique = append(unique, area)
		}
	}
	return unique
}

func areasEqual(a, b [3][3]bool) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func readLinesFromFile(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}
