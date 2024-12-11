package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Rule struct {
	isBefore []string
	isAfter  []string
}

func sortRules(ruleSlice []string) map[string]Rule {
	sortedMap := map[string]Rule{}
	for _, rule := range ruleSlice {
		splitRule := strings.Split(rule, "|")
		data, exists := sortedMap[splitRule[0]]
		if !exists {
			sortedMap[splitRule[0]] = Rule{
				isBefore: []string{splitRule[1]},
				isAfter:  []string{},
			}
		} else {
			data.isBefore = append(data.isBefore, splitRule[1])
			sortedMap[splitRule[0]] = data
		}

		data, exists = sortedMap[splitRule[1]]
		if !exists {
			sortedMap[splitRule[1]] = Rule{
				isBefore: []string{},
				isAfter:  []string{splitRule[0]},
			}
		} else {
			data.isAfter = append(data.isAfter, splitRule[0])
			sortedMap[splitRule[1]] = data
		}
	}
	return sortedMap
}

func checkPreviousEntriesValid(curRule Rule, prevEntries []string) (bool, int) {
	for i, num := range prevEntries {
		if slices.Contains(curRule.isBefore, num) {
			return false, i
		}
	}
	return true, -1
}

func isCorrectOrder(row []string, ruleMap map[string]Rule) bool {
	for i, num := range row {
		rule, exists := ruleMap[num]

		if !exists || i == 0 {
			continue
		}
		valid, _ := checkPreviousEntriesValid(rule, row[:i])
		if !valid {
			return false
		}
	}
	return true
}

func getMiddleNumber(row []string) int {
	if len(row)%2 == 0 {
		return 0
	}

	middleIndex := int(((len(row) - 1) / 2))
	midNum, err := strconv.Atoi(row[middleIndex])
	if err != nil {
		log.Fatal(err)
	}
	return midNum
}

func part1(ruleSlice, rowSlice []string) {
	sortedRules := sortRules(ruleSlice)

	validRows := [][]string{}
	for _, row := range rowSlice {
		splitRow := strings.Split(row, ",")
		valid := isCorrectOrder(splitRow, sortedRules)
		if valid {
			validRows = append(validRows, splitRow)
		}
	}

	totalMidNum := 0

	for _, validR := range validRows {
		totalMidNum += getMiddleNumber(validR)
	}
	log.Printf("Part1: %v", totalMidNum)
}

func orderInvalidRows(row []string, ruleMap map[string]Rule) []string {
	valid := true
	invalidIndex := -1
	indexToMove := -1
	for i, num := range row {
		rule, exists := ruleMap[num]

		if !exists || i == 0 {
			continue
		}
		valid, invalidIndex = checkPreviousEntriesValid(rule, row[:i])
		if !valid {
			indexToMove = i
			break
		}
	}
	if valid {
		log.Printf("row is now valid")
		return row
	}

	newRow := []string{}
	newRow = append(newRow, row[:invalidIndex]...)
	newRow = append(newRow, row[indexToMove])
	newRow = append(newRow, row[invalidIndex])
	if indexToMove != invalidIndex+1 {
		newRow = append(newRow, row[invalidIndex+1:indexToMove]...)
	}
	newRow = append(newRow, row[indexToMove+1:]...)
	return orderInvalidRows(newRow, ruleMap)
}

func part2(ruleSlice, rowSlice []string) {
	sortedRules := sortRules(ruleSlice)

	invalidRows := [][]string{}
	validRows := [][]string{}
	for _, row := range rowSlice {
		splitRow := strings.Split(row, ",")
		valid := isCorrectOrder(splitRow, sortedRules)
		if !valid {
			invalidRows = append(invalidRows, splitRow)
		}
	}

	for _, row := range invalidRows {
		validRows = append(validRows, orderInvalidRows(row, sortedRules))
	}

	totalMidNum := 0

	for _, validR := range validRows {
		totalMidNum += getMiddleNumber(validR)
	}
	log.Printf("Part2: %v", totalMidNum)
}

func main() {
	//f, err := os.Open("./testInput.txt")
	f, err := os.Open("./day5Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := bufio.NewScanner(f)

	rulesSlice := []string{}
	rowsSlice := []string{}

	swapOver := false

	for fs.Scan() {
		line := fs.Text()
		if line != "EOF" {
			if line == "" {
				swapOver = true
				continue
			}
			if !swapOver {
				rulesSlice = append(rulesSlice, line)
				continue
			}

			rowsSlice = append(rowsSlice, line)
		}
	}

	if len(rulesSlice) == 0 || len(rowsSlice) == 0 {
		log.Fatal("File read error")
	}
	part1(rulesSlice, rowsSlice)
	part2(rulesSlice, rowsSlice)
}
