package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Direction int
type Letter struct {
	letter string
	N      *Letter
	NE     *Letter
	E      *Letter
	SE     *Letter
	S      *Letter
	SW     *Letter
	W      *Letter
	NW     *Letter
}

const (
	North Direction = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

func (l *Letter) Search(word string) int {
	numFound := 0
	letterSplice := strings.Split(word, "")
	if l.letter == letterSplice[0] {
		//search all directions for next letter
		found := false
		for i := North; i <= NorthWest; i++ {
			switch i {
			case North:
				if l.N != nil {
					found = l.N.SearchDirection(North, letterSplice, 1)
				}
			case NorthEast:
				if l.NE != nil {
					found = l.NE.SearchDirection(NorthEast, letterSplice, 1)
				}
			case East:
				if l.E != nil {
					found = l.E.SearchDirection(East, letterSplice, 1)
				}
			case SouthEast:
				if l.SE != nil {
					found = l.SE.SearchDirection(SouthEast, letterSplice, 1)
				}
			case South:
				if l.S != nil {
					found = l.S.SearchDirection(South, letterSplice, 1)
				}
			case SouthWest:
				if l.SW != nil {
					found = l.SW.SearchDirection(SouthWest, letterSplice, 1)
				}
			case West:
				if l.W != nil {
					found = l.W.SearchDirection(West, letterSplice, 1)
				}
			case NorthWest:
				if l.NW != nil {
					found = l.NW.SearchDirection(NorthWest, letterSplice, 1)
				}
			}
			if found {
				found = false
				numFound += 1
			}
		}
	}

	return numFound
}

func (l *Letter) SearchDirection(dir Direction, wordSplice []string, index int) bool {
	if l.letter != wordSplice[index] {
		return false
	}
	if index == len(wordSplice)-1 {
		return true
	}

	switch dir {
	case North:
		if l.N != nil {
			return l.N.SearchDirection(North, wordSplice, index+1)
		}
	case NorthEast:
		if l.NE != nil {
			return l.NE.SearchDirection(NorthEast, wordSplice, index+1)
		}
	case East:
		if l.E != nil {
			return l.E.SearchDirection(East, wordSplice, index+1)
		}
	case SouthEast:
		if l.SE != nil {
			return l.SE.SearchDirection(SouthEast, wordSplice, index+1)
		}
	case South:
		if l.S != nil {
			return l.S.SearchDirection(South, wordSplice, index+1)
		}
	case SouthWest:
		if l.SW != nil {
			return l.SW.SearchDirection(SouthWest, wordSplice, index+1)
		}
	case West:
		if l.W != nil {
			return l.W.SearchDirection(West, wordSplice, index+1)
		}
	case NorthWest:
		if l.NW != nil {
			return l.NW.SearchDirection(NorthWest, wordSplice, index+1)
		}
	}
	return false
}

func (l *Letter) searchXmas() bool {
	if l.letter != "A" {
		return false
	}
	masFound := 0
	if l.NE != nil && l.NW != nil && l.SE != nil && l.SW != nil {
		if (l.NW.letter == "M" && l.SE.letter == "S") || (l.NW.letter == "S" && l.SE.letter == "M") {
			masFound += 1
		}

		if (l.NE.letter == "M" && l.SW.letter == "S") || l.NE.letter == "S" && l.SW.letter == "M" {
			masFound += 1
		}
	}
	if masFound > 1 {
		return true
	}
	return false
}

func part2(grid [][]*Letter) {
	totalFound := 0
	for _, row := range grid {
		for _, col := range row {
			found := col.searchXmas()
			if found {
				totalFound += 1
			}
		}
	}
	log.Printf("Part 2: %v", totalFound)
}

func part1(grid [][]*Letter, word string) {
	totalFound := 0
	for _, row := range grid {
		for _, col := range row {
			totalFound += col.Search(word)

		}

	}
	log.Printf("Part 1: %v", totalFound)
}

func main() {

	f, err := os.Open("./day4Input.txt")
	//f, err := os.Open("./testInput.txt")
	if err != nil {
		log.Fatal(err)
	}

	fscanner := bufio.NewScanner(f)

	grid := [][]*Letter{}
	var PrevRow []*Letter
	for fscanner.Scan() {
		row := fscanner.Text()
		if row != "EOF" {
			letterSplice := strings.Split(row, "")
			currentRow := []*Letter{}
			for i, letter := range letterSplice {
				newLetter := Letter{
					letter: letter,
				}
				if len(currentRow) > 0 {
					currentRow[i-1].E = &newLetter
					newLetter.W = currentRow[i-1]
				}
				if len(PrevRow) > 0 {
					newLetter.N = PrevRow[i]
					newLetter.NW = PrevRow[i].W
					newLetter.NE = PrevRow[i].E

					PrevRow[i].S = &newLetter
					PrevRow[i].SW = newLetter.W
					if i > 0 {
						PrevRow[i-1].SE = &newLetter
					}

				}
				currentRow = append(currentRow, &newLetter)
			}
			PrevRow = currentRow
			grid = append(grid, currentRow)
		}
	}
	part1(grid, "XMAS")
	part2(grid)
}
