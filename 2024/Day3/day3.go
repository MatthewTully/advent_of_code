package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	mulPattern      = "mul\\(\\d{1,3},\\d{1,3}\\)"
	numsOnlyPattern = "(?<=mul\\()\\d{1,3},\\d{1,3}(?=\\))"
)

func getMultsFromData(data string) []string {
	re := regexp.MustCompile(mulPattern)
	return re.FindAllString(data, -1)

}

func splitAndParseNums(set string) []int {
	ss := strings.Split(set, ",")
	intSlice := make([]int, 2)

	cleanString1, _ := strings.CutPrefix(ss[0], "mul(")
	cleanString2, _ := strings.CutSuffix(ss[1], ")")
	log.Printf("post clean numbers: %v and %v", cleanString1, cleanString2)

	int1, err := strconv.Atoi(cleanString1)
	if err != nil {
		log.Fatal(err)
	}
	int2, err := strconv.Atoi(cleanString2)
	if err != nil {
		log.Fatal(err)
	}

	intSlice[0] = int1
	intSlice[1] = int2
	log.Printf("Returning numbers: %v and %v", int1, int2)
	return intSlice
}

func multPairs(pair []int) int {
	return pair[0] * pair[1]
}

func part2(data []byte) {

	actionableText := []string{}

	splitDonts := strings.Split(string(data), "don't()")
	actionableText = append(actionableText, splitDonts[0])

	for i, dontStr := range splitDonts {
		if i == 0 {
			continue
		}
		splitDos := strings.Split(dontStr, "do()")
		actionableText = append(actionableText, splitDos[1:]...)
	}
	if len(actionableText) == 0 {
		log.Fatal("Split data on dont't and do failed")
	}

	matches := getMultsFromData(strings.Join(actionableText, ""))
	if len(matches) == 0 {
		log.Fatal("No matches found in data")
	}
	total := 0
	for _, set := range matches {
		pair := splitAndParseNums(set)
		total += multPairs(pair)
	}

	log.Printf("Part 2: %v", total)
}

func part1(data []byte) {

	matches := getMultsFromData(string(data))
	if len(matches) == 0 {
		log.Fatal("No matches found in data")
	}

	total := 0
	for _, set := range matches {
		pair := splitAndParseNums(set)
		total += multPairs(pair)
	}

	log.Printf("Part 1: %v", total)

}

func main() {

	f, err := os.Open("./day3Input.txt")
	if err != nil {
		log.Fatal("err")
	}

	fReader := bufio.NewReader(f)

	data, err := io.ReadAll(fReader)
	if err != nil {
		log.Fatal(err)
	}
	part1(data)
	part2(data)

}
