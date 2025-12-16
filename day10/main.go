package main

import (
	"fmt"
	"log"
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

	fmt.Println(getTotalIndicatorRoute(machines, findIndicatorRoute))
	fmt.Println(getTotalIndicatorRoute(machines, findJoltageRoute))
}

func getTotalIndicatorRoute(machines []Machine, routeFunc func(Machine) int) int {
	total := 0
	for _, machine := range machines {
		count := routeFunc(machine)
		fmt.Println(count, machine)
		total += count
	}
	return total
}

func findJoltageRoute(machine Machine) int {
	return findJoltageRouteStep(machine, make([]int, len(machine.joltageTarget)), make(map[string]bool), intSliceToString(machine.joltageTarget))
}

func findJoltageRouteStep(machine Machine, state []int, seen map[string]bool, targetStateStr string) int {
	if targetStateStr == intSliceToString(state) {
		return 0
	}
	sortedSwitches := make([][]int, len(machine.switches))
	copy(sortedSwitches, machine.switches)
	sort.Slice(sortedSwitches, func(i, j int) bool {
		val1, val2 := 0, 0
		for _, sw := range sortedSwitches[i] {
			val1 += machine.joltageTarget[sw] - state[sw]
		}
		for _, sw := range sortedSwitches[j] {
			val2 += machine.joltageTarget[sw] - state[sw]
		}
		return val1 > val2
	})
	//fmt.Println(sortedSwitches)
	best := 10000
	for _, sw := range sortedSwitches {
		nextState := make([]int, len(machine.joltageTarget))
		copy(nextState, state)
		tooHigh := false
		for _, s := range sw {
			nextState[s] += 1
			if nextState[s] > machine.joltageTarget[s] {
				tooHigh = true
				break
			}
		}
		if tooHigh {
			continue
		}
		nextStateStr := intSliceToString(nextState)
		if seen[nextStateStr] {
			continue
		}
		seen[nextStateStr] = true
		count := findJoltageRouteStep(machine, nextState, seen, targetStateStr)
		if count != -1 && count < best {
			return count + 1
		}
	}
	return -1
}

func findIndicatorRoute(machine Machine) int {
	seen := make(map[string]bool)
	states := [][]bool{make([]bool, len(machine.indicatorTarget))}
	targetStateStr := boolSliceToString(machine.indicatorTarget)
	count := 0
	for true {
		count++
		nextStates := [][]bool{}
		for _, state := range states {
			for _, sw := range machine.switches {
				nextState := applyIndicatorSwitches(state, sw)
				nextStateStr := boolSliceToString(nextState)
				if seen[nextStateStr] {
					continue
				}
				if nextStateStr == targetStateStr {
					return count
				}
				seen[nextStateStr] = true
				nextStates = append(nextStates, nextState)
			}
		}

		states = nextStates
	}
	return count
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
