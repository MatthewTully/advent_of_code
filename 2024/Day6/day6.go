package main

import (
	"bufio"
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
)

type Node struct {
	leftEdge   bool
	rightEdge  bool
	topEdge    bool
	bottomEdge bool

	leftNode   *Node
	rightNode  *Node
	topNode    *Node
	bottomNode *Node

	containsBlocker bool
	guardHasBeen    bool
}

type Guard struct {
	startingNode       *Node
	currentNode        *Node
	uniqueNodesVisited uint
	facingDirection    Direction
	inArea             bool
	stepsTaken         uint
}

func (g *Guard) turn() {
	switch g.facingDirection {
	case up:
		g.facingDirection = right
	case right:
		g.facingDirection = down
	case down:
		g.facingDirection = left
	case left:
		g.facingDirection = up
	}
}

func (g *Guard) leftArea() {
	g.inArea = false
}

func (g *Guard) moveNode() {
	switch g.facingDirection {
	case up:
		if g.currentNode.topEdge {
			g.leftArea()
		} else if !g.currentNode.topNode.containsBlocker {
			g.currentNode = g.currentNode.topNode
		} else {
			g.turn()
			g.moveNode()
		}
	case right:
		if g.currentNode.rightEdge {
			g.leftArea()
		} else if !g.currentNode.rightNode.containsBlocker {
			g.currentNode = g.currentNode.rightNode
		} else {
			g.turn()
			g.moveNode()
		}
	case down:
		if g.currentNode.bottomEdge {
			g.leftArea()
		} else if !g.currentNode.bottomNode.containsBlocker {
			g.currentNode = g.currentNode.bottomNode

		} else {
			g.turn()
			g.moveNode()
		}
	case left:
		if g.currentNode.leftEdge {
			g.leftArea()
		} else if !g.currentNode.leftNode.containsBlocker {
			g.currentNode = g.currentNode.leftNode
		} else {
			g.turn()
			g.moveNode()
		}
	}

	if g.inArea {
		if !g.currentNode.guardHasBeen {
			g.currentNode.guardHasBeen = true
			g.uniqueNodesVisited++
		}
		g.stepsTaken++
	}
}

func (g *Guard) resetGuard() {
	g.currentNode = g.startingNode
	g.inArea = true
	g.stepsTaken = 0
	g.uniqueNodesVisited = 0
	g.facingDirection = up
}

func part1(g Guard) {
	g.moveNode()
}

func doesCreateLoop(node *Node, g *Guard, maxSteps uint) bool {
	if node.containsBlocker {
		return false
	}
	node.containsBlocker = true
	for g.stepsTaken < maxSteps && g.inArea {
		g.moveNode()
	}
	node.containsBlocker = false
	return g.inArea
}

func part2(grid *[][]*Node, guard *Guard, possibleSteps uint) {
	//loop over and add blocker, count steps if steps exceed 16083, then must be infinite loop
	possibleBlocker := 0
	for _, row := range *grid {
		for _, node := range row {
			guard.resetGuard()
			if doesCreateLoop(node, guard, possibleSteps) {
				possibleBlocker++
			}
		}
	}
	log.Printf("TotalPossible Blockers: %v", possibleBlocker)
}

func main() {

	f, err := os.Open("./day6Input.txt")
	//f, err := os.Open("./testInput.txt")
	if err != nil {
		log.Fatal("err")
	}

	fscanner := bufio.NewScanner(f)

	guard := Guard{
		facingDirection:    up,
		uniqueNodesVisited: 1,
		inArea:             true,
		stepsTaken:         0,
	}
	var grid [][]*Node
	var PrevRow []*Node
	noneBlockerNodes := 0
	for fscanner.Scan() {
		row := fscanner.Text()
		if row != "EOF" {
			nodes := strings.Split(row, "")
			currentRow := []*Node{}
			for i, n := range nodes {
				newNode := Node{
					topEdge:         len(PrevRow) == 0,
					leftEdge:        i == 0,
					rightEdge:       i == len(nodes)-1,
					bottomEdge:      false,
					containsBlocker: n == "#",
					guardHasBeen:    n == "^",
				}
				if n == "^" {
					guard.currentNode = &newNode
					guard.startingNode = &newNode
				}
				if len(currentRow) > 0 {
					currentRow[i-1].rightNode = &newNode
					newNode.leftNode = currentRow[i-1]
				}
				if len(PrevRow) > 0 {
					newNode.topNode = PrevRow[i]
					PrevRow[i].bottomNode = &newNode
				}
				if n != "#" && n != "^" {
					noneBlockerNodes++
				}
				currentRow = append(currentRow, &newNode)
			}
			grid = append(grid, currentRow)
			PrevRow = currentRow
		}
	}
	for _, n := range PrevRow {
		n.bottomEdge = true
	}

	part1(guard)
	part2(&grid, &guard, uint(noneBlockerNodes))
}
