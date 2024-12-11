package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func numOfOccurs(a int, listB []int) int {

	occurs := 0

	for _, b := range listB {
		if a == b {
			occurs += 1
		}
	}
	return occurs
}

func part1(listA, listB []int) {
	sort.Slice(listA, func(i, j int) bool {
		return listA[i] < listA[j]
	})
	sort.Slice(listB, func(i, j int) bool {
		return listB[i] < listB[j]
	})

	totalDist := 0
	for i := range listA {
		var dif int
		if listA[i] > listB[i] {
			dif = listA[i] - listB[i]
		} else {
			dif = listB[i] - listA[i]
		}
		totalDist += dif
	}

	log.Printf("Part1 Total Distance: %v", totalDist)
}

func part2(listA, listB []int) {

	simScore := 0
	for _, x := range listA {
		occurs := numOfOccurs(x, listB)
		simScore += (x * occurs)
	}
	log.Printf("Part2 Total Sim Score: %v", simScore)
}

func main() {

	f, err := os.Open("./day1Input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fscanner := bufio.NewScanner(f)

	listA := []int{}
	listB := []int{}

	for fscanner.Scan() {
		line := fscanner.Text()

		if line != "EOF" {
			splitLine := strings.Split(fscanner.Text(), "   ")

			numA, err := strconv.Atoi(strings.Trim(splitLine[0], " "))
			if err != nil {
				log.Fatal(err)
			}

			numB, err := strconv.Atoi(strings.Trim(splitLine[1], " "))
			if err != nil {
				log.Fatal(err)
			}

			listA = append(listA, numA)
			listB = append(listB, numB)
		}
	}

	part1(listA, listB)
	part2(listA, listB)

}
