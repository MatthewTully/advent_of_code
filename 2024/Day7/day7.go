package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Equation struct {
	testValue       uint64
	calibrationNums []uint64
}

func convertStringSliceToInt(strSlice []string) []uint64 {
	intSlice := []uint64{}
	for _, str := range strSlice {
		intVal, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			log.Fatal(err)
		}
		intSlice = append(intSlice, uint64(intVal))
	}
	if len(intSlice) != len(strSlice) {
		log.Fatal("Converted slice is missing values")
	}
	return intSlice
}

func checkOperation_r(testVal uint64, total uint64, depth int, op []uint64) bool {
	totalCpy := total
	sum := op[depth]

	totalCpy += sum
	if depth+1 == len(op) {
		log.Printf("Last op (add), does %v = expected? %v", totalCpy, totalCpy == testVal)
		if totalCpy == testVal {
			return true
		}
	} else {
		if totalCpy > testVal {
			return false
		}

		valid := checkOperation_r(testVal, totalCpy, depth+1, op)
		if valid {
			return valid
		}
	}

	totalCpy = total
	if totalCpy > 0 {
		totalCpy *= sum
	} else {
		totalCpy = sum
	}
	if depth+1 == len(op) {
		log.Printf("Last op (mult), does %v = expected ? %v", totalCpy, totalCpy == testVal)
		return totalCpy == testVal
	}

	if totalCpy > testVal {
		return false
	}
	return checkOperation_r(testVal, totalCpy, depth+1, op)

}

func calibrateEquation(eq Equation) bool {
	log.Printf("Checking Value %v, with numbers: %v", eq.testValue, eq.calibrationNums)
	return checkOperation_r(eq.testValue, eq.calibrationNums[0], 1, eq.calibrationNums)
}

func part1(eqSlice []Equation) {
	var total uint64
	for _, eq := range eqSlice {
		if eq.testValue > eq.calibrationNums[0] && calibrateEquation(eq) {
			log.Printf("Valid equation: %v", eq.testValue)
			total += eq.testValue
		}
	}
	log.Printf("Total Eq: %v", len(eqSlice))
	log.Printf("Part 1 Total: %v", total)
}

func main() {

	f, err := os.Open("./testInput.txt")
	//f, err := os.Open("./day7Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := bufio.NewScanner(f)

	eqSlice := []Equation{}

	for fs.Scan() {
		line := fs.Text()
		if line == "EOF" {
			break
		}

		newEq := Equation{}

		firstSplit := strings.Split(line, ":")
		testVal, err := strconv.Atoi(firstSplit[0])
		if err != nil {
			log.Fatal(err)
		}
		newEq.testValue = uint64(testVal)

		secondSplit := strings.Split(strings.TrimSpace(firstSplit[1]), " ")
		newEq.calibrationNums = convertStringSliceToInt(secondSplit)
		eqSlice = append(eqSlice, newEq)
	}

	part1(eqSlice)
}
