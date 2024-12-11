package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type position struct {
	X int
	Y int
}

type Node struct {
	frequency string
	position  position
	antiNode  bool
}

func withinMapRange(x, y, gx, gy int) bool {
	return x >= 0 && x < gx && y >= 0 && y < gy
}

func sortFreq(node *Node, freqSlice []*Node, gridHeight int, gridLen int) []Node {
	antiNodeSlice := []Node{}
	for _, nextNode := range freqSlice {
		xDiff := node.position.X - nextNode.position.X
		yDiff := node.position.Y - nextNode.position.Y

		for i := range 2 {
			antiNode := position{}
			if i == 0 {
				antiNode.X = node.position.X + xDiff
				antiNode.Y = node.position.Y + yDiff
			} else {
				antiNode.X = nextNode.position.X - xDiff
				antiNode.Y = nextNode.position.Y - yDiff
			}
			log.Printf("antinode %v", antiNode)
			if withinMapRange(antiNode.X, antiNode.Y, gridLen, gridHeight) {
				antiNodeSlice = append(antiNodeSlice, Node{
					frequency: node.frequency,
					position:  antiNode,
				})
			}
		}

	}
	return antiNodeSlice
}

func updateGrid2(nodeMap [][]*Node, x, y int) {
	nodeMap[y][x].antiNode = true
	if nodeMap[y][x].frequency == "." {
		nodeMap[y][x].frequency = "#"
	}
}

func sortHarmonics(node *Node, freqSlice []*Node, nodeMap [][]*Node) {
	var x int
	var y int
	for _, nextNode := range freqSlice {
		xDiff := node.position.X - nextNode.position.X
		yDiff := node.position.Y - nextNode.position.Y

		for i := range 2 {
			if i == 0 {
				x = nextNode.position.X
				y = nextNode.position.Y
				for withinMapRange(x, y, len(nodeMap[0]), len(nodeMap)) {
					updateGrid2(nodeMap, x, y)
					x += xDiff
					y += yDiff
				}
			} else {
				x = node.position.X
				y = node.position.Y
				for withinMapRange(x, y, len(nodeMap[0]), len(nodeMap)) {
					updateGrid2(nodeMap, x, y)
					x -= xDiff
					y -= yDiff
				}
			}
		}
	}
}

func updateGrid(nodeMap [][]*Node, antinodes []Node) {
	for _, node := range antinodes {
		if nodeMap[node.position.Y][node.position.X].frequency == node.frequency {
			continue
		}
		nodeMap[node.position.Y][node.position.X].antiNode = true
		if nodeMap[node.position.Y][node.position.X].frequency == "." {
			nodeMap[node.position.Y][node.position.X].frequency = "#"
		}
	}
}

func printGrid(nodeMap [][]*Node) {
	var sb strings.Builder
	totalAntiNode := 0
	for y := 0; y < len(nodeMap); y++ {
		for x := 0; x < len(nodeMap[y]); x++ {
			sb.Write([]byte(nodeMap[y][x].frequency))
			if nodeMap[y][x].antiNode {
				log.Printf("anti node freq = %v", nodeMap[y][x])
				totalAntiNode++
			}
		}
		sb.Write([]byte("\n"))
	}
	log.Printf("Total Anti nodes: %v\n%v", totalAntiNode, sb.String())
}

func part1(nodeMap [][]*Node, freqMap map[string][]*Node) {
	antiNodes := []Node{}
	log.Printf("Nodes %v", freqMap)

	for k, v := range freqMap {
		if len(v) > 1 {
			log.Printf("Check Freq %v", k)
			for i := range v {
				antiNodes = append(antiNodes, sortFreq(v[i], v[i:], len(nodeMap), len(nodeMap[0]))...)
			}
		}
	}
	updateGrid(nodeMap, antiNodes)
	printGrid(nodeMap)
}

func part2(nodeMap [][]*Node, freqMap map[string][]*Node) {

	for _, v := range freqMap {
		if len(v) > 1 {
			for i := range v {
				if i+1 < len(v) {
					sortHarmonics(v[i], v[i+1:], nodeMap)
				}
			}
		}
	}
	printGrid(nodeMap)
}

func main() {
	//f, err := os.Open("./testInput.txt")
	f, err := os.Open("./day8Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := bufio.NewScanner(f)

	nodeMap := [][]*Node{}
	freqMap := map[string][]*Node{}
	currentLine := 0
	for fs.Scan() {
		line := fs.Text()
		if line == "EOF" {
			break
		}
		splitLine := strings.Split(line, "")
		nodeSlice := []*Node{}
		for i, s := range splitLine {
			newNode := Node{
				frequency: s,
				position: position{
					X: i,
					Y: currentLine,
				},
				antiNode: false,
			}
			if s != "." {
				ns, exists := freqMap[s]
				if !exists {
					freqMap[s] = []*Node{&newNode}
				} else {
					freqMap[s] = append(ns, &newNode)
				}

			}
			nodeSlice = append(nodeSlice, &newNode)
		}
		currentLine++
		nodeMap = append(nodeMap, nodeSlice)
	}
	part1(nodeMap, freqMap)
	part2(nodeMap, freqMap)

}
