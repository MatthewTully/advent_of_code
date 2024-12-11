package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Direction int

const (
	up Direction = iota
	right
	down
	left
)

type trailStep struct {
	trailHead bool
	left      *trailStep
	right     *trailStep
	up        *trailStep
	down      *trailStep
	val       int
	found     bool
}

type trail struct {
	ends []*trailStep
}

func (ts *trailStep) reset() {
	ts.found = false
}

func (ts *trailStep) availableDirections() []Direction {
	dirs := []Direction{}
	for i := up; i <= left; i++ {
		switch i {
		case up:
			if ts.up != nil && ts.val == ts.up.val-1 {
				dirs = append(dirs, up)
			}
		case right:
			if ts.right != nil && ts.val == ts.right.val-1 {
				dirs = append(dirs, right)
			}
		case down:
			if ts.down != nil && ts.val == ts.down.val-1 {
				dirs = append(dirs, down)
			}
		case left:
			if ts.left != nil && ts.val == ts.left.val-1 {
				dirs = append(dirs, left)
			}
		}
	}
	return dirs
}

func (ts *trailStep) walk_distinct_r(t *trail) *trail {
	log.Printf("cur val %v", ts.val)
	if ts.found {
		return t
	}

	if ts.val == 9 {
		ts.found = true
		t.ends = append(t.ends, ts)
		return t
	}
	dirsToWalk := ts.availableDirections()
	log.Printf("total avialable dirs: %v", len(dirsToWalk))
	for _, dir := range dirsToWalk {
		switch dir {
		case up:
			t = ts.up.walk_distinct_r(t)
		case right:
			t = ts.right.walk_distinct_r(t)
		case down:
			t = ts.down.walk_distinct_r(t)
		case left:
			t = ts.left.walk_distinct_r(t)
		}
	}
	return t

}

func (ts *trailStep) walk_r(score int) int {
	log.Printf("cur val %v", ts.val)
	if ts.val == 9 {
		score++
		return score
	}
	dirsToWalk := ts.availableDirections()
	log.Printf("total avialable dirs: %v", len(dirsToWalk))
	for _, dir := range dirsToWalk {
		switch dir {
		case up:
			score = ts.up.walk_r(score)
		case right:
			score = ts.right.walk_r(score)
		case down:
			score = ts.down.walk_r(score)
		case left:
			score = ts.left.walk_r(score)
		}
	}
	return score

}

func part2(start []*trailStep) {
	total := 0
	for _, trailS := range start {
		total += trailS.walk_r(0)
	}
	log.Printf("Part1, total score %v", total)
}

func part1(start []*trailStep) {
	total := 0
	for _, trailS := range start {
		t := trail{ends: []*trailStep{}}
		trailS.walk_distinct_r(&t)
		total += len(t.ends)
		for _, ts := range t.ends {
			ts.reset()
		}

		log.Printf("total so far %v", total)
	}
	log.Printf("Part1, total score %v", total)
}

func main() {
	//f, err := os.Open("./testInput.txt")
	f, err := os.Open("./day10Input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fs := bufio.NewScanner(f)

	trailStarts := []*trailStep{}
	prevRow := []*trailStep{}

	for fs.Scan() {
		line := fs.Text()
		if line == "EOF" {
			break
		}
		steps := strings.Split(line, "")
		currentRow := []*trailStep{}
		for i, n := range steps {
			val, err := strconv.Atoi(n)
			if err != nil {
				log.Fatal(err)
			}
			newStep := trailStep{
				val:       val,
				trailHead: val == 0,
			}
			if len(currentRow) > 0 {
				currentRow[i-1].right = &newStep
				newStep.left = currentRow[i-1]
			}
			if len(prevRow) > 0 {
				newStep.down = prevRow[i]
				prevRow[i].up = &newStep
			}

			if newStep.trailHead {
				trailStarts = append(trailStarts, &newStep)
			}
			currentRow = append(currentRow, &newStep)
		}
		prevRow = currentRow

	}

	if len(trailStarts) > 0 {
		part1(trailStarts)
		part2(trailStarts)
	}

}
