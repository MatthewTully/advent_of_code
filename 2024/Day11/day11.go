package main

import (
	"bufio"
	"io"
	"log"
	"maps"
	"os"
	"strconv"
	"strings"
)

type stone struct {
	stoneMap map[int]int
}

var TOTAL_BLINKS = 75

func blink(stones map[int]int) map[int]int {
	MapToChange := map[int]int{}
	maps.Copy(MapToChange, stones)
	for k, v := range MapToChange {
		if v == 0 {
			continue
		}
		if k == 0 {
			stones[1] += v
			stones[0] -= v
		} else if len(strconv.Itoa(k))%2 == 0 {
			strVal := strconv.Itoa(k)
			val, err := strconv.Atoi(strVal[:int(len(strVal)/2)])

			if err != nil {
				log.Fatal(err)
			}
			stones[val] += v
			val, err = strconv.Atoi(strVal[int(len(strVal)/2):])
			if err != nil {
				log.Fatal(err)
			}
			stones[val] += v
			stones[k] -= v
		} else {
			stones[k*2024] += v
			stones[k] -= v
		}

	}
	return stones
}

func part1(s stone) {
	for TOTAL_BLINKS != 0 {
		s.stoneMap = blink(s.stoneMap)
		TOTAL_BLINKS--
	}
	log.Printf("total stones: %v", totalStones(s.stoneMap))

}

func totalStones(stoneMap map[int]int) int {
	total := 0
	for _, val := range stoneMap {
		total += val
	}
	return total
}

func parseInput() []string {
	//f, err := os.Open("./testInput.txt")
	f, err := os.Open("./day11Input.txt")

	if err != nil {
		log.Fatal(err)
	}
	fReader := bufio.NewReader(f)

	data, err := io.ReadAll(fReader)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(data), " ")
}

func createStones(splitS []string) stone {

	stone := stone{
		stoneMap: map[int]int{},
	}

	for _, s := range splitS {
		val, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}

		stone.stoneMap[val]++
	}
	return stone
}

func main() {
	splitS := parseInput()
	s := createStones(splitS)
	part1(s)
}
