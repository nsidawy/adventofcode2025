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

type Machine struct {
	indicatorTarget []bool
	switches        [][]int
	joltageTarget   []int
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <filename>")
	}
	filename := os.Args[1]
	lines := readLinesFromFile(filename)
	machines := make([]Machine, len(lines))
	for i, line := range lines {
		machines[i] = parseMachine(line)
	}

	fmt.Println(getTotalIndicatorRoute(machines))
	fmt.Println(getTotalJoltageRoute(machines))
}

func getTotalIndicatorRoute(machines []Machine) int {
	total := 0
	for _, machine := range machines {
		paths := findIndicatorRoute(machine)
		sort.Slice(paths, func(i, j int) bool {
			return len(paths[i]) < len(paths[j])
		})
		count := len(paths[0])
		fmt.Println(count, paths[0], machine)
		total += count
	}
	return total
}

func getTotalJoltageRoute(machines []Machine) int {
	total := 0
	for _, machine := range machines {
		count := findJoltageRoute(machine.joltageTarget, machine.switches, []int{})
		fmt.Println(count, machine)
		total += count
	}
	return total
}

func findJoltageRoute(joltage []int, switches [][]int, path []int) int {
	indicator := getIndicatorFromJoltage(joltage)
	partialMachine := Machine{indicatorTarget: indicator, switches: switches, joltageTarget: joltage}
	partialPaths := findIndicatorRoute(partialMachine)

	minPath := math.MaxInt32
	for _, partialPath := range partialPaths {
		newJoltage := make([]int, len(joltage))
		copy(newJoltage, joltage)
		newPath := make([]int, len(path)+len(partialPath))
		copy(newPath, path)
		copy(newPath[len(path):], partialPath)
		for _, si := range partialPath {
			for _, sw := range switches[si] {
				newJoltage[sw]--
			}
		}

		is_negative := false
		for _, j := range newJoltage {
			if j < 0 {
				is_negative = true
				break
			}
		}
		if is_negative {
			continue
		}

		result := len(partialPath)
		if !isJoltageValid(newJoltage) {
			all_evens := true
			mult := 1
			for all_evens {
				mult *= 2
				for i := range newJoltage {
					newJoltage[i] = newJoltage[i] / 2
					if newJoltage[i]%2 == 1 {
						all_evens = false
					}
					all_evens = false
				}
			}

			pathScore := findJoltageRoute(newJoltage, switches, newPath)
			if pathScore == math.MaxInt32 {
				continue
			}
			result += pathScore * mult
		}
		if result < minPath {
			minPath = result
		}
	}
	return minPath
}

func isJoltageValid(joltage []int) bool {
	for _, j := range joltage {
		if j != 0 {
			return false
		}
	}
	return true
}

func getIndicatorFromJoltage(joltage []int) []bool {
	indicator := make([]bool, len(joltage))
	for i, j := range joltage {
		indicator[i] = j%2 == 1
	}
	return indicator
}

func findIndicatorRoute(machine Machine) [][]int {
	paths := make([][]int, 1)
	paths[0] = make([]int, 0)
	targetStateStr := boolSliceToString(machine.indicatorTarget)
	fullPaths := [][]int{}
	options := generateAllSubsets(len(machine.switches))
	for _, o := range options {
		indicators := make([]bool, len(machine.indicatorTarget))
		for _, sw := range o {
			indicators = applyIndicatorSwitches(indicators, machine.switches[sw])
		}
		if targetStateStr == boolSliceToString(indicators) {
			fullPaths = append(fullPaths, o)
		}
	}
	return fullPaths
}

func applyIndicatorSwitches(state []bool, switches []int) []bool {
	newState := append([]bool{}, state...)
	for _, sw := range switches {
		newState[sw] = !newState[sw]
	}
	return newState
}

func boolSliceToString(slice []bool) string {
	var sb strings.Builder
	for _, b := range slice {
		if b {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

func intSliceToString(slice []int) string {
	var sb strings.Builder
	for i, v := range slice {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func parseMachine(s string) Machine {
	parts := strings.Split(s, " ")
	targetStr := strings.Trim(parts[0], "[]")
	target := make([]bool, len(targetStr))
	for i, c := range targetStr {
		target[i] = c == '#'
	}
	switcheStrs := parts[1 : len(parts)-1]
	switches := make([][]int, len(switcheStrs))
	for i, s := range switcheStrs {
		intStrs := strings.Split(strings.Trim(s, "()"), ",")
		switches[i] = make([]int, len(intStrs))
		for j, intStr := range intStrs {
			switches[i][j], _ = strconv.Atoi(intStr)
		}
	}
	joltageStrs := strings.Split(strings.Trim(parts[len(parts)-1], "{}"), ",")
	joltageRequirements := make([]int, len(joltageStrs))
	for i, joltageStr := range joltageStrs {
		joltage, _ := strconv.Atoi(joltageStr)
		joltageRequirements[i] = joltage
	}
	return Machine{indicatorTarget: target, switches: switches, joltageTarget: joltageRequirements}
}

func readLinesFromFile(filename string) []string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}

// generateAllSubsets returns all possible subsets of numbers [0, 1, 2, ..., n]
// This is the power set - includes empty set and all combinations
func generateAllSubsets(n int) [][]int {
	var result [][]int
	var generate func([]int, int)

	generate = func(current []int, start int) {
		// Add current subset to result (make a copy)
		subset := make([]int, len(current))
		copy(subset, current)
		result = append(result, subset)

		// Generate all subsets that include numbers from start to n
		for i := start; i < n; i++ {
			current = append(current, i)
			generate(current, i+1)
			current = current[:len(current)-1] // backtrack
		}
	}

	generate([]int{}, 0)
	return result
}
