package main

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func convReportToInt(report []string) []int {
	intSlice := make([]int, len(report))
	for i, x := range report {
		y, err := strconv.Atoi(x)
		if err != nil {
			log.Fatal(err)
		}
		intSlice[i] = y
	}
	return intSlice
}

func safeDistance(a, b int) bool {

	if (a-b) < 1 || (a-b) > 3 {
		return false
	}

	return true

}

func isLevelSafe(report []int) bool {
	asc := report[0] < report[1]
	switch asc {
	case true:
		for i := range report {
			if i+1 < len(report) {
				if report[i] > report[i+1] {
					return false
				}
				if !safeDistance(report[i+1], report[i]) {
					return false
				}
			}

		}
	default:
		for i := range report {
			if i+1 < len(report) {
				if report[i] < report[i+1] {
					return false
				}
				if !safeDistance(report[i], report[i+1]) {
					return false
				}
			}
		}
	}
	return true
}

func removeDuplicateLevels(report []int, skipUsed bool) (bool, []int, bool) {
	log.Printf("Checking level %v for duplicates", report)
	reportCP := make([]int, len(report))
	copy(reportCP, report[:])
	if !reflect.DeepEqual(reportCP, report) {
		log.Printf("%v", reportCP)
		log.Printf("%v", report)
		log.Fatal("copy error")
	}

	compact := slices.Compact(reportCP)
	log.Printf("compacted: %v", compact)
	lenDiff := len(report) - len(compact)
	if lenDiff == 0 {
		return true, report, skipUsed
	}
	if lenDiff == 1 && !skipUsed {
		return true, compact, true
	}
	return false, report, true

}

func isLevelOrdered(report []int, skipUsed bool) (bool, []int, bool, bool) {
	log.Printf("Checking level %v Order", report)
	asc := false

	reportCP := make([]int, len(report))
	copy(reportCP, report[:])
	if !reflect.DeepEqual(reportCP, report) {
		log.Printf("%v", reportCP)
		log.Printf("%v", report)
		log.Fatal("copy error")
	}

	sort.Slice(reportCP, func(i, j int) bool {
		return reportCP[i] < reportCP[j]
	})

	if reflect.DeepEqual(reportCP, report) {
		asc = true
		return true, report, asc, skipUsed
	}

	//not asc
	sort.Slice(reportCP, func(i, j int) bool {
		return reportCP[i] > reportCP[j]
	})

	if reflect.DeepEqual(reportCP, report) {
		asc = false
		return true, report, asc, skipUsed

	}
	if skipUsed {
		// skip already used, so can't fix
		return false, report, asc, skipUsed
	}

	// not ordered, need to check difference
	orderChange := 0
	log.Printf("Can't determine order.")
	for i := 0; i+1 < len(report); i++ {
		if report[i] > report[i+1] {
			if i == 0 {
				asc = false
				continue
			}
			if asc {
				orderChange += 1
				asc = false
			}

		}
		if report[i] < report[i+1] {
			if i == 0 {
				asc = true
				continue
			}
			if asc {
				orderChange += 1
				asc = true
			}
		}
	}
	log.Printf("orderChange = %v", orderChange)
	if orderChange <= 2 {
		log.Printf("order changes once, use skip and continue = %v", orderChange)
		return true, reportCP, false, true
	}
	log.Printf("order changes too often, not safe : change - %v, report - %v", orderChange, report)
	return false, report, false, true

}

func safeDifference(report []int, asc, skipUsed bool) bool {
	log.Printf("Checking level %v, order: %v", report, asc)
	for i := 0; i+1 < len(report); i++ {
		var safe bool

		if asc {
			safe = safeDistance(report[i+1], report[i])
		} else {
			safe = safeDistance(report[i], report[i+1])
		}

		if safe {
			continue
		}

		if skipUsed {
			log.Printf("%v distance from %v exceeds safe limit and skip already used (%v)", report[i], report[i+1], report)
			return false
		}

		if i+2 == len(report) {
			//i+1 is the issue and the last number, safe
			log.Printf("Last num is the issue, SAFE level %v", report)
			return true
		}

		if i+2 < len(report) {
			var iIsSafe bool
			var i1IsSafe bool
			if asc {
				iIsSafe = safeDistance(report[i+2], report[i])
				i1IsSafe = safeDistance(report[i+2], report[i+1])
			} else {
				iIsSafe = safeDistance(report[i], report[i+2])
				i1IsSafe = safeDistance(report[i+1], report[i+2])
			}

			if !iIsSafe {
				//i is the issue and needs removing
				log.Printf("%v distance from %v exceeds safe limit and %v (i) is the issue (%v)", report[i], report[i+1], report[i], report)
				skipUsed = true
				copy(report[i:], report[i+1:])
				report[len(report)-1] = 0
				report = report[:len(report)-1]
				i = -1
			} else if !i1IsSafe {
				//i+1 is the issue and needs removing
				log.Printf("%v distance from %v exceeds safe limit and %v (i+1) is the issue (%v)", report[i], report[i+1], report[i+1], report)
				skipUsed = true
				copy(report[i+1:], report[i+2:])
				report[len(report)-1] = 0
				report = report[:len(report)-1]
				i = -1
			} else {
				//i+1 is the issue and needs removing
				log.Printf("%v distance from %v exceeds safe limit and %v (i+1) is the issue (%v)", report[i], report[i+1], report[i+1], report)
				skipUsed = true
				copy(report[i+1:], report[i+2:])
				report[len(report)-1] = 0
				report = report[:len(report)-1]
				i = -1
			}
		}

	}
	log.Printf("SAFE level %v", report)
	return true
}

func part1(reports [][]int) {

	numSafe := 0

	for _, x := range reports {
		if isLevelSafe(x) {
			numSafe += 1
		}
	}

	log.Printf("Number of safe reports: %v", numSafe)

}

func part2(reports [][]int) {
	numSafe := 0
	for _, x := range reports {
		safe, ordered, asc, skipUsed := isLevelOrdered(x, false)
		if !safe {
			continue
		}
		safe, noDupes, skipUsed := removeDuplicateLevels(ordered, skipUsed)
		if !safe {
			continue
		}
		safeDif := safeDifference(noDupes, asc, skipUsed)
		if safeDif {
			numSafe += 1
		}

	}
	log.Printf("Part2 Number of safe reports: %v", numSafe)
}

func main() {

	f, err := os.Open("./day2Input.txt")
	//f, err := os.Open("./test_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fscanner := bufio.NewScanner(f)
	reports := [][]int{}

	for fscanner.Scan() {

		line := fscanner.Text()
		if line != "EOF" {
			report := convReportToInt(strings.Split(line, " "))
			reports = append(reports, report)
		}
	}

	part2(reports)
	part1(reports)

}
