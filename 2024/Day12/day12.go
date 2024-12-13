package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Direction int

const (
	up Direction = iota
	right
	down
	left
	NotSet
)

type field struct {
	regionId      int
	plant         string
	parimeter     int
	top           *field
	right         *field
	bottom        *field
	left          *field
	outterBoarder Direction
	visitedOutter bool
	visitedInner  bool
	x             int
	y             int
}

type plot struct {
	id        int
	area      int
	parimeter []*int
}

type group struct {
	id            int
	numberOfSides int
	area          int
	startField    *field
	boarderFields []*field
}

func addPlot(plots map[int]*plot, fields *field) {
	_, exists := plots[fields.regionId]
	if !exists {
		newPlot := plot{
			id:        fields.regionId,
			parimeter: []*int{&fields.parimeter},
			area:      1,
		}
		plots[fields.regionId] = &newPlot
		return
	}
	plots[fields.regionId].area++
	plots[fields.regionId].parimeter = append(plots[fields.regionId].parimeter, &fields.parimeter)

}

func isBoarderParim(f *field) bool {
	if f.bottom == nil && f.top == nil {
		return true
	} else if f.left == nil && f.right == nil {
		return true
	}
	return false
}

func addGroup(groups map[int]*group, fields *field) {
	_, exists := groups[fields.regionId]
	if !exists {
		newGroup := group{
			id:         fields.regionId,
			startField: fields,
			area:       1,
		}
		groups[fields.regionId] = &newGroup
		return
	}
	groups[fields.regionId].area++
	if (fields.parimeter == 2 && isBoarderParim(fields)) || fields.parimeter == 1 {
		groups[fields.regionId].boarderFields = append(groups[fields.regionId].boarderFields, fields)
	}
}

func (p *plot) sumParimeter() int {
	total := 0
	for _, parim := range p.parimeter {
		total += *parim
	}
	return total
}

func (p *plot) calcCost() int {
	return p.area * p.sumParimeter()
}

func (g *group) calcCost() int {
	return g.area * g.numberOfSides
}

func Part1(grid [][]*field) {
	plots := map[int]*plot{}
	groupFields(grid)
	for _, row := range grid {
		for _, i := range row {
			addPlot(plots, i)
		}
	}
	total := 0
	for _, region := range plots {
		total += region.calcCost()
	}
	log.Printf("Part1: %v", total)
}

func part2(grid [][]*field) {
	grouped := map[int]*group{}
	groupFields(grid)
	for _, row := range grid {
		for _, field := range row {
			addGroup(grouped, field)
		}
	}
	total := 0
	for _, group := range grouped {
		group.countSides()
		total += group.calcCost()
	}
	log.Printf("Part2: %v", total)
}

func (f *field) group(id int) {
	if f.regionId != -1 {
		return
	}

	f.regionId = id
	if f.parimeter == 4 {
		return
	}

	if f.right != nil {
		f.right.group(id)
	}
	if f.top != nil {
		f.top.group(id)
	}
	if f.bottom != nil {
		f.bottom.group(id)
	}
	if f.left != nil {
		f.left.group(id)
	}
}

func (g *group) countSides() {
	f := g.startField
	if f.left == nil && f.right == nil && f.bottom == nil && f.top == nil {
		g.numberOfSides = 4
		return
	}
	g.numberOfSides = g.crawlBoarder(f, g.numberOfSides)
	if len(g.boarderFields) > 1 {
		for _, f := range g.boarderFields {
			if f.visitedInner {
				continue
			}
			g.numberOfSides += g.crawlBoarderInverse(f, 0)
		}
	}

}

func determineNextWall(comeFrom, lastwall Direction) (Direction, error) {
	switch lastwall {
	case left:
		if comeFrom == right {
			return up, nil
		}
		return down, nil
	case down:
		if comeFrom == up {
			return left, nil
		}
		return right, nil
	case right:
		if comeFrom == left {
			return down, nil
		}
		return up, nil
	case up:
		if comeFrom == down {
			return right, nil
		}
		return left, nil
	}
	log.Fatal("could not determine direction")
	return 0, fmt.Errorf("could not determine direction")
}

func determineNextWallInverse(comeFrom, lastwall Direction) (Direction, error) {
	switch lastwall {
	case left:
		if comeFrom == right {
			return down, nil
		}
		return up, nil
	case down:
		if comeFrom == up {
			return right, nil
		}
		return left, nil
	case right:
		if comeFrom == left {
			return up, nil
		}
		return down, nil
	case up:
		if comeFrom == down {
			return left, nil
		}
		return right, nil
	}
	log.Fatal("could not determine direction")
	return 0, fmt.Errorf("could not determine direction")
}

func determineFirstWallInverse(f *field) Direction {
	if f.parimeter > 2 {
		log.Fatal("shouldn't be here without 2 parim")
	}
	if f.parimeter == 2 {
		if f.outterBoarder == up && f.bottom == nil {
			return down
		}
		if f.outterBoarder == down && f.top == nil {
			return up
		}
		if f.outterBoarder == right && f.left == nil {
			return left
		}
		if f.outterBoarder == left && f.right == nil {
			return right
		}
	}
	if f.parimeter == 1 && !f.visitedOutter {
		if f.bottom == nil {
			return down
		}
		if f.top == nil {
			return up
		}
		if f.left == nil {
			return left
		}
		if f.right == nil {
			return right
		}
	}
	return NotSet
}

func (g *group) crawlBoarderInverse(f *field, count int) int {
	startX := f.x
	startY := f.y

	dir := determineFirstWallInverse(f)
	if dir == NotSet {
		return count
	}
	validInner := true
	comeFrom := up
	lastWall := left

	switch dir {
	case up:
		if f.outterBoarder == up {
			return 0
		}
		f, count = g.crawlTopWallInverse(f, count)
		comeFrom = left
		lastWall = up
		if f.top != nil {
			f = f.top
			f.visitedInner = true
			comeFrom = down
		}
	case right:
		if f.outterBoarder == right {
			return 0
		}
		f, count = g.crawlRightWallInverse(f, count)
		comeFrom = up
		lastWall = right
		if f.right != nil {
			f = f.right
			f.visitedInner = true
			comeFrom = left
		}
	case left:
		if f.outterBoarder == left {
			return 0
		}
		f, count = g.crawlLeftWallInverse(f, count)
		comeFrom = down
		lastWall = left
		if f.left != nil {
			f = f.left
			f.visitedInner = true
			comeFrom = right
		}
	case down:
		if f.outterBoarder == down {
			return 0
		}
		f, count = g.crawlBottomWallInverse(f, count)
		comeFrom = right
		lastWall = down
		if f.bottom != nil {
			f = f.bottom
			f.visitedInner = true
			comeFrom = up
		}
	}

	nextDir, err := determineNextWallInverse(comeFrom, lastWall)
	if err != nil {
		log.Fatal(err)
	}

	safetyX := f.x
	safetyY := f.y
	loop := 0

	for !(startX == f.x && startY == f.y) && !(safetyX == f.x && safetyY == f.y && loop == 1) {
		loop = 1
		if !validInner {
			break
		}
		if f.parimeter == 0 {
			log.Fatal("shouldn't have got here")
		}
		switch nextDir {
		case left:
			if f.outterBoarder == left {
				return 0
			}
			f, count = g.crawlLeftWallInverse(f, count)
			lastWall = left
			comeFrom = down
			if f.left != nil {
				f = f.left
				f.visitedInner = true
				comeFrom = right
			}
		case right:
			if f.outterBoarder == right {
				return 0

			}
			f, count = g.crawlRightWallInverse(f, count)
			lastWall = right
			comeFrom = up
			if f.right != nil {
				f = f.right
				f.visitedInner = true
				comeFrom = left
			}
		case up:
			if f.outterBoarder == up {
				return 0
			}
			f, count = g.crawlTopWallInverse(f, count)
			lastWall = up
			comeFrom = left
			if f.top != nil {
				f = f.top
				f.visitedInner = true
				comeFrom = down

			}
		case down:
			if f.outterBoarder == down {
				return 0
			}
			f, count = g.crawlBottomWallInverse(f, count)
			lastWall = down
			comeFrom = right
			if f.bottom != nil {
				f = f.bottom
				f.visitedInner = true
				comeFrom = up
			}
		}
		nextDir, err = determineNextWallInverse(comeFrom, lastWall)
		if err != nil {
			log.Fatal(err)
		}
	}

	if validInner {
		if count%2 == 0 {
			return count
		}
		return count - 1
	}
	return 0
}

func (g *group) crawlBoarder(f *field, count int) int {
	f, count = g.crawlLeftWall(f, count)
	comeFrom := up
	lastWall := left
	if f.left != nil {
		f = f.left
		f.visitedOutter = true
		if f.outterBoarder == NotSet {
			f.outterBoarder = up
		}
		comeFrom = right
	}

	nextDir, err := determineNextWall(comeFrom, lastWall)
	if err != nil {
		log.Fatal(err)
	}
	notMovedYet := false
	if f == g.startField {
		notMovedYet = true
	}
	pushedToLast := false

	for notMovedYet || f != g.startField {
		if f.parimeter == 0 {
			log.Fatal("shouldn't have got here")
		}
		switch nextDir {
		case left:
			f, count = g.crawlLeftWall(f, count)
			lastWall = left
			comeFrom = up
			if f.left != nil {
				f = f.left
				f.visitedOutter = true
				if f.outterBoarder == NotSet {
					f.outterBoarder = up
				}
				if f == g.startField {
					pushedToLast = true
				}
				comeFrom = right
			}
		case right:
			f, count = g.crawlRightWall(f, count)
			lastWall = right
			comeFrom = down
			if f.right != nil {
				f = f.right
				f.visitedOutter = true
				if f.outterBoarder == NotSet {
					f.outterBoarder = down
				}
				if f == g.startField {
					pushedToLast = true
				}
				comeFrom = left
			}
		case up:
			f, count = g.crawlTopWall(f, count)
			lastWall = up
			comeFrom = right
			if f.top != nil {
				f = f.top
				f.visitedOutter = true
				if f.outterBoarder == NotSet {
					f.outterBoarder = right
				}
				if f == g.startField {
					pushedToLast = true
				}
				comeFrom = down

			}
		case down:
			f, count = g.crawlBottomWall(f, count)
			lastWall = down
			comeFrom = left
			if f.bottom != nil {
				f = f.bottom
				f.visitedOutter = true
				if f.outterBoarder == NotSet {
					f.outterBoarder = left
				}
				if f == g.startField {
					pushedToLast = true
				}
				comeFrom = up
			}
		}
		if f != g.startField {
			notMovedYet = false
		}
		nextDir, err = determineNextWall(comeFrom, lastWall)
		if err != nil {
			log.Fatal(err)
		}
	}
	if pushedToLast {
		count += 2
	} else if lastWall != up {
		count++
	}
	return count
}

func (g *group) crawlLeftWall(f *field, count int) (*field, int) {
	for f.left == nil && f.bottom != nil {
		if f.outterBoarder == right && f.visitedOutter {
			f.visitedInner = true
		}
		f = f.bottom
		f.visitedOutter = true
		if f.outterBoarder == NotSet {
			f.outterBoarder = left
		}
	}
	count++
	return f, count
}

func (g *group) crawlBottomWall(f *field, count int) (*field, int) {
	for f.bottom == nil && f.right != nil {
		if f.outterBoarder == up && f.visitedOutter {
			f.visitedInner = true
		}
		f = f.right
		f.visitedOutter = true
		if f.outterBoarder == NotSet {
			f.outterBoarder = down
		}
	}
	count++
	return f, count
}

func (g *group) crawlRightWall(f *field, count int) (*field, int) {
	for f.right == nil && f.top != nil {
		if f.outterBoarder == left && f.visitedOutter {
			f.visitedInner = true
		}
		f = f.top
		f.visitedOutter = true
		if f.outterBoarder == NotSet {
			f.outterBoarder = right
		}
	}
	count++
	return f, count
}

func (g *group) crawlTopWall(f *field, count int) (*field, int) {
	for f.top == nil && f.left != nil {
		if f.outterBoarder == down && f.visitedOutter {
			f.visitedInner = true
		}
		f = f.left
		f.visitedOutter = true
		if f.outterBoarder == NotSet {
			f.outterBoarder = up
		}
	}
	count++
	return f, count
}

func (g *group) crawlLeftWallInverse(f *field, count int) (*field, int) {
	for f.left == nil && f.top != nil {
		f = f.top
		f.visitedInner = true
	}
	count++
	return f, count
}

func (g *group) crawlBottomWallInverse(f *field, count int) (*field, int) {
	for f.bottom == nil && f.left != nil {
		f = f.left
		f.visitedInner = true
	}
	count++
	return f, count
}

func (g *group) crawlRightWallInverse(f *field, count int) (*field, int) {
	for f.right == nil && f.bottom != nil {
		f = f.bottom
		f.visitedInner = true
	}
	count++
	return f, count
}

func (g *group) crawlTopWallInverse(f *field, count int) (*field, int) {
	for f.top == nil && f.right != nil {
		f = f.right
		f.visitedInner = true
	}
	count++
	return f, count
}

func groupFields(grid [][]*field) {
	id := 0
	for _, row := range grid {
		for _, field := range row {
			if field.regionId == -1 {
				field.group(id)
				id++
			}
		}
	}
}

func main() {
	//f, err := os.Open("./testInput.txt")
	f, err := os.Open("./day12Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := bufio.NewScanner(f)
	grid := [][]*field{}
	lineNum := 0
	for fs.Scan() {
		line := fs.Text()
		if line == "EOF" {
			break
		}
		splitL := strings.Split(line, "")
		currentRow := []*field{}
		for i, f := range splitL {
			newField := field{
				regionId:      -1,
				plant:         f,
				parimeter:     4,
				visitedOutter: false,
				visitedInner:  false,
				outterBoarder: NotSet,
				x:             i,
				y:             lineNum,
			}
			if len(grid) > 0 && grid[lineNum-1][i].plant == f {
				grid[lineNum-1][i].parimeter -= 1
				newField.parimeter -= 1
				grid[lineNum-1][i].bottom = &newField
				newField.top = grid[lineNum-1][i]
			}
			if i > 0 && currentRow[i-1].plant == f {
				currentRow[i-1].parimeter -= 1
				newField.parimeter -= 1
				currentRow[i-1].right = &newField
				newField.left = currentRow[i-1]
			}
			currentRow = append(currentRow, &newField)
		}
		grid = append(grid, currentRow)
		lineNum++
	}

	Part1(grid)
	part2(grid)
}
